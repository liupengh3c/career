package bootstrap

import (
	"career-server/resource"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/time/rate"
)

func Init() {
	// 这里可以初始化一些全局变量或者执行一些初始化的操作
	// ElasticInit()
	LimitInit()
	// RabbitMQInit()
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

func LimitInit() {
	// 这里可以初始化一些全局变量或者执行一些初始化的操作
	resource.IpLimiter = make(map[string]*rate.Limiter)
	resource.GlobalLimiter = rate.NewLimiter(rate.Limit(resource.GlobalLimiterCnt), resource.GlobalLimiterMax)
}

func RedisInit() {
	resource.RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func RabbitMQInit() {
	conn, err := amqp.Dial("amqp://guest:guest@192.168.3.15:5672/")
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	// defer conn.Close()
	resource.RabbitMQChannel, err = conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel")
		return
	}
	fmt.Println("exchange declare start")
	err = resource.RabbitMQChannel.ExchangeDeclare(
		"tag_engine", // exchange name
		"direct",     // exchange type
		true,         // durable
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("Failed to declare an exchange:", err)
		return
	}
	fmt.Println("RabbitMQ init success")

	// 监听连接关闭事件，实现自动重连
	go func() {
		notifyClose := conn.NotifyClose(make(chan *amqp.Error))
		for err := range notifyClose {
			if err != nil {
				fmt.Printf("RabbitMQ connection closed: %v, attempting to reconnect...\n", err)
				// 这里可以实现重连逻辑
				// 为了简单，暂时只打印日志
			}
		}
	}()
}
