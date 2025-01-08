package bootstrap

import (
	"career-server/resource"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

func Init() {
	// 这里可以初始化一些全局变量或者执行一些初始化的操作
	ElasticInit()
}
func ElasticInit() {
	cert, _ := os.ReadFile("/Users/liupeng/Documents/study/elasticsearch-8.17.0/config/certs/http_ca.crt")
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  "elastic",
		Password:  "xpE4DQGWE9bCkoj7WXYE",
		Addresses: []string{"https://127.0.0.1:9200"},
		CACert:    cert,
	})

	if err != nil {
		fmt.Println("create client err:", err.Error())
		return
	}
	resource.ElasticClient = client
}
