package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/liupengh3c/esbuilder"
)

// 最外层数据结构
type Documents struct {
	ScrollID string      `json:"_scroll_id"`
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
}
type Total struct {
	Relation string `json:"relation"`
	Value    int    `json:"value"`
}

func main() {
	client, err := NewEsClient()
	if err != nil {
		fmt.Println("create client err:", err.Error())
		return
	}
	fmt.Println("connect success")
	for i := 0; i < 1; i++ {
		ScrollSearch(client)
	}
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
		// fmt.Println("create client err:", err.Error())
		return client, err
	}
	return client, nil
}

func ScrollSearch(client *elasticsearch.Client) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	docs := Documents{}
	dslQuery := esbuilder.NewDsl()
	boolQuery := esbuilder.NewBoolQuery()

	dslQuery.SetOrder(esbuilder.NewSortQuery("doc_id", "asc"))
	dslQuery.SetQuery(boolQuery)
	dslQuery.SetSize(1000)

	res, err := client.Search(
		client.Search.WithIndex("new_tag_202411"),
		client.Search.WithBody(strings.NewReader(dslQuery.BuildJson())),
		client.Search.WithScroll(time.Minute*50),
	)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return
	}
	err = json.NewDecoder(res.Body).Decode(&docs)
	if err != nil {
		fmt.Println("decode err:", err.Error())
		return
	}
	fmt.Println("search count:", len(docs.Hits.Hits))
	scrollId := docs.ScrollID
	for {
		docs = Documents{}
		res, err = client.Scroll(
			client.Scroll.WithScrollID(scrollId),
			client.Scroll.WithScroll(time.Minute),
		)
		if err != nil {
			fmt.Println("scroll err:", err.Error())
			return
		}
		fmt.Println("==========", err)
		json.NewDecoder(res.Body).Decode(&docs)

		if len(docs.Hits.Hits) == 0 {
			break
		}
		fmt.Println("search count:", len(docs.Hits.Hits))
		scrollId = docs.ScrollID
	}
}
