package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	jsoniter "github.com/json-iterator/go"
)

// 最外层数据结构
type Documents struct {
	PitId    string      `json:"pit_id"`
	Shards   Shards      `json:"_shards"`
	Hits     HitOutLayer `json:"hits"`
	TimedOut bool        `json:"timed_out"`
	Took     int         `json:"took"`
}
type Shards struct {
	Failed     int `json:"failed"`
	Skipped    int `json:"skipped"`
	Successful int `json:"successful"`
	Total      int `json:"total"`
}
type HitOutLayer struct {
	Hits     []Hits  `json:"hits"`
	MaxScore float64 `json:"max_score"`
	Total    Total   `json:"total"`
}
type Hits struct {
	ID     string         `json:"_id"`
	Index  string         `json:"_index"`
	Score  float64        `json:"_score"`
	Source map[string]any `json:"_source"`
	Type   string         `json:"_type"`
	Sort   []any          `json:"sort"`
}
type Total struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}

func NewEsClient() (*elasticsearch.Client, error) {
	cert, _ := os.ReadFile("/Users/liupeng/Documents/study/elasticsearch-8.17.0/config/certs/http_ca.crt")
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  "elastic",
		Password:  "XBS=adqa799j_Aoz=A+h",
		Addresses: []string{"https://127.0.0.1:9200"},
		CACert:    cert,
	})

	if err != nil {
		fmt.Println("create client err:", err.Error())
		return client, err
	}
	return client, nil
}
func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	pitIds := map[string]int{}
	client, _ := NewEsClient()
	pit := esapi.OpenPointInTimeRequest{
		Index:     []string{"new_tag_202411"},
		KeepAlive: "60m",
	}
	resp, _ := pit.Do(context.Background(), client)
	pitResp := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&pitResp)
	id := pitResp["id"].(string)
	// id := "yvaYBAEObmV3X3RhZ18yMDI0MTEWd3oxYVR3N2dTTktDbWhwblhCVHdYUQABFjhvQ0tlaU9tUVFHUkRrb3NKQXo0aHcAAQAAAAAAAJouFkxDY1lLamNLU3NHQ25DNjlpTmRYZVEAARZ3ejFhVHc3Z1NOS0NtaHBuWEJUd1hRAAA="
	fmt.Println("id", id)
	docs := Documents{}
	mapDsl := map[string]any{
		"size": 1,
		"sort": []map[string]any{
			{
				"doc_id": "asc",
			},
		},
		"pit": map[string]any{
			"id":         id,
			"keep_alive": "60m",
		},
	}
	strDsl, _ := json.MarshalToString(mapDsl)
	search := esapi.SearchRequest{
		Body: strings.NewReader(strDsl),
	}
	searchResp, _ := search.Do(context.Background(), client)
	json.NewDecoder(searchResp.Body).Decode(&docs)
	id = docs.PitId
	if _, ok := pitIds[id]; !ok {
		pitIds[id] = 1
	}
	sort := docs.Hits.Hits[len(docs.Hits.Hits)-1].Sort
	for {
		docs := Documents{}
		mapDsl := map[string]any{
			"size": 1,
			"sort": []map[string]any{
				{
					"doc_id": "asc",
				},
			},
			"pit": map[string]any{
				"id":         id,
				"keep_alive": "60m",
			},
			"search_after": sort,
		}
		strDsl, _ := json.MarshalToString(mapDsl)
		search := esapi.SearchRequest{
			Body: strings.NewReader(strDsl),
		}
		searchResp, _ := search.Do(context.Background(), client)
		json.NewDecoder(searchResp.Body).Decode(&docs)
		if len(docs.Hits.Hits) == 0 {
			break
		}
		id = docs.PitId
		if _, ok := pitIds[id]; !ok {
			pitIds[id] = 1
		}
		sort = docs.Hits.Hits[len(docs.Hits.Hits)-1].Sort
		fmt.Println(docs.Hits.Hits[0].Source["doc_id"])
	}
	fmt.Println(pitIds)
}
