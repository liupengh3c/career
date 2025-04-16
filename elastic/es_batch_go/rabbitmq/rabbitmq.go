package rabbitmq

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	host         string
	connection   *amqp.Connection
	exchangeName string
	exchangeType string
	queueName    string
	channel      *amqp.Channel
	buffer       chan amqp.Delivery
}

func NewRabbitMq(host string, exchangeName, exchangType, queueName string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(host)
	if err != nil {
		return nil, err
	}
	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel")
		return nil, err
	}
	// 创建交换机
	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangType,  // type
		true,         // durable
		false,        // delete when unused
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    "dlk:" + exchangeName,
		"x-dead-letter-routing-key": "dlk:" + queueName,
	}
	// 创建队列
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		return nil, err
	}
	ch.QueueBind(queueName, queueName, exchangeName, false, nil)
	return &RabbitMQ{
		host:         host,
		connection:   conn,
		exchangeName: exchangeName,
		exchangeType: exchangType,
		queueName:    queueName,
		channel:      ch,
		buffer:       make(chan amqp.Delivery, 1500),
	}, nil
}

func (r *RabbitMQ) Publish(msg string) error {
	err := r.channel.Publish(
		r.exchangeName, // exchange
		r.queueName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		},
	)
	return err
}

func (r *RabbitMQ) BeginConsume() error {
	msgs, err := r.channel.Consume(
		r.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			r.buffer <- d
		}
	}()
	return nil
}

func (r *RabbitMQ) GetBufferLen() int {
	return len(r.buffer)
}

func (r *RabbitMQ) GetMessages(cnt int) ([]string, error) {
	msgs := make([]string, 0)
	if cnt > 1500 {
		return msgs, errors.New("too many messages")
	}
	for i := 0; i < cnt; i++ {
		if len(r.buffer) == 0 {
			break
		}
		msg, ok := <-r.buffer
		if !ok {
			return msgs, errors.New("channel closed")
		}
		msgs = append(msgs, string(msg.Body))
		msg.Ack(true)
	}
	return msgs, nil
}
