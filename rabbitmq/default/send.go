package main

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Send(msg string) error {
	// 连接rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@123.57.167.109:5672/")
	if err != nil {
		fmt.Println("connect error:", err)
		return err
	}
	defer conn.Close()
	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel error:", err)
		return err
	}
	defer ch.Close()
	// 创建队列,使用默认的交换机
	q, err := ch.QueueDeclare(
		"lp_default", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		fmt.Println("queue declare error:", err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println(q.Name)
	// body := "Hello World!"
	for i := 0; i < 100; i++ {
		err = ch.PublishWithContext(ctx,
			"",     // exchange,默认交换机
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(msg + fmt.Sprintf("-%d", i+1)),
				DeliveryMode: amqp.Persistent,
			})
		if err != nil {
			fmt.Println("publish error:", err)
			return err
		}
	}
	return nil
}

func main() {
	Send("Hello world")
}
