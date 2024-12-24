package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@123.57.167.109:5672/")
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ")
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel")
		return
	}
	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		fmt.Println("Failed to declare an exchange,err:", err)
		return
	}
	body := "Hello World by topic exchange"
	err = ch.Publish(
		"logs_topic",       // exchange
		"quick.orange.fox", // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		fmt.Println("Failed to publish a message")
		return
	}
	fmt.Println("send message " + body + " success")
}
