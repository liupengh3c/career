package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/liupengh3c/esbuilder"
)

// 最外层数据结构
type Documents struct {
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

func main() {
	SearchFromSize()
}

func SearchFromSize() {
	st := time.Now()
	defer func() {
		fmt.Println("cost:", time.Since(st).Milliseconds(), "ms")
	}()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	docs := Documents{}
	cert, _ := os.ReadFile("/Users/liupeng/Documents/study/elasticsearch-8.17.0/config/certs/http_ca.crt")
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  "elastic",
		Password:  "XBS=adqa799j_Aoz=A+h",
		Addresses: []string{"https://127.0.0.1:9200"},
		CACert:    cert,
	})

	if err != nil {
		fmt.Println("create client err:", err.Error())
		return
	}

	dslQuery := esbuilder.NewDsl()
	boolQuery := esbuilder.NewBoolQuery()
	boolQuery.Filter(esbuilder.NewRangeQuery("doc_id").Gte(1))
	dslQuery.SetQuery(boolQuery)
	dslQuery.SetFrom(0)
	dslQuery.SetSize(10)
	dslQuery.SetOrder(esbuilder.NewSortQuery("doc_id", "asc"))
	dsl := dslQuery.BuildJson()
	search := esapi.SearchRequest{
		Index: []string{"new_tag_202411"},
		Body:  strings.NewReader(dsl),
	}
	resp, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		fmt.Println("decode err:", err.Error())
		return
	}
	fmt.Println(docs.Hits.Hits[len(docs.Hits.Hits)-1].Sort)
	dslQuery.SetSearchAfter(docs.Hits.Hits[len(docs.Hits.Hits)-1].Sort)
	for {
		search := esapi.SearchRequest{
			Index: []string{"new_tag_202411"},
			Body:  strings.NewReader(dslQuery.BuildJson()),
		}
		fmt.Println(dslQuery.BuildJson())
		resp, err = search.Do(context.Background(), client)
		if err != nil {
			fmt.Println("search err:", err.Error())
			return
		}
		err = json.NewDecoder(resp.Body).Decode(&docs)
		if err != nil {
			fmt.Println("decode err:", err.Error())
			return
		}
		if len(docs.Hits.Hits) == 0 {
			fmt.Println("no more data")
			break
		}

		fmt.Println(len(docs.Hits.Hits), docs.Hits.Hits[len(docs.Hits.Hits)-1].Source["doc_id"])
		dslQuery.SetSearchAfter(docs.Hits.Hits[len(docs.Hits.Hits)-1].Sort)
	}
}
