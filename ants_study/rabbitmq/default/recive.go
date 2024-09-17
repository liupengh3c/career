package main

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@123.57.167.109:5672/")
	if err != nil {
		fmt.Println("connect error:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Channel error:", err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"lp_default", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		fmt.Println("Queue Declare error:", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println("Consume error:", err)
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)
			// time.Sleep(time.Second * 1)
			d.Ack(false)
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
