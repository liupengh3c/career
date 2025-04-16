package main

import (
	"es_batch_go/elastic"
	"es_batch_go/rabbitmq"
	"fmt"
	"time"
)

func main() {
	host := "amqp://guest:guest@192.168.3.47:5672/"
	exchangeName := "dw_exchange"
	exchangeType := "direct"
	queue := "dw_queue"
	batchSize := 100
	maxDelay := 10
	fmt.Println("connecting...")
	mq, err := rabbitmq.NewRabbitMq(host, exchangeName, exchangeType, queue)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("rabbitmq connect ok")
	mq.BeginConsume()

	es, err := elastic.NewEsBatchInsert("https://localhost:9200", "elastic", "xpE4DQGWE9bCkoj7WXYE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("elastic connect ok")
	now := time.Now().Unix()
	indexId := 0
	for {
		// 打batch，等待足够数据到达batchSize，如果超过最大等待时间没有达到batchsize，则直接处理
		if time.Now().Unix()-now < int64(maxDelay) && mq.GetBufferLen() < batchSize {
			continue
		}
		cnt := batchSize
		if mq.GetBufferLen() < batchSize {
			cnt = mq.GetBufferLen()
		}
		msgs, _ := mq.GetMessages(cnt)
		docs := make([]elastic.Document, 0)
		for _, v := range msgs {
			indexId++
			doc := elastic.Document{
				ID:     fmt.Sprintf("%v", indexId),
				Index:  "blog_202504", // 根据业务逻辑灵活设计索引
				Source: map[string]any{"data": v},
			}
			docs = append(docs, doc)
		}
		es.BatchInSert(docs)
		fmt.Println("batch insert ok,len=", cnt)
		now = time.Now().Unix()
	}
}
