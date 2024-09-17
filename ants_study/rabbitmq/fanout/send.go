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
	err = ch.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		fmt.Println("Failed to declare an exchange")
		return
	}
	body := "Hello World!"
	err = ch.Publish("logs", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		fmt.Println("Failed to publish a message")
		return
	}
}
