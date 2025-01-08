package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func ImageSearch(client *elasticsearch.Client, query string) []map[string]any {
	results := make([]map[string]any, 0)
	docs := Documents{}
	search := esapi.SearchRequest{
		Index: []string{"vector_search_202412"},
		Body:  strings.NewReader(query),
	}
	resp, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Println("search err:", err.Error())
		return results
	}
	err = json.NewDecoder(resp.Body).Decode(&docs)
	if err != nil {
		fmt.Println("decode err:", err.Error())
		return results
	}
	for _, v := range docs.Hits.Hits {
		result := make(map[string]any)
		result["name"] = v.Source["name"]
		result["image_url"] = v.Source["path"]
		results = append(results, result)
	}
	return results
}
