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
	err = ch.ExchangeDeclare("logs", "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Println("Failed to declare an exchange")
		return
	}
	q, err := ch.QueueDeclare(
		"logs_topic", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	err = ch.QueueBind(
		q.Name,       // queue name
		"*.orange.*", // routing key(binding key)
		"logs_topic", // exchange
		false,
		nil,
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	var forever chan struct{}

	go func() {
		for d := range msgs {
			fmt.Printf(" [x] %s\n", d.Body)
		}
	}()

	fmt.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
