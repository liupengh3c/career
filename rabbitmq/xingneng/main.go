package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TestConfig struct {
	Mode         string
	EngineName   string
	QueueName    string
	QueueDurable bool
	DeliveryMode uint8
	UseConfirm   bool
	MessageSize  int
	MessageCount int
	Concurrency  int
}

func runTest(cfg TestConfig) {
	fmt.Printf("\n===== Running test: %s =====\n", cfg.Mode)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("connect failed: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel failed: %v", err)
	}
	defer ch.Close()

	if cfg.UseConfirm {
		if err := ch.Confirm(false); err != nil {
			log.Fatalf("confirm mode failed: %v", err)
		}
	}

	// q, err := ch.QueueDeclare(
	// 	"bench_queue",
	// 	cfg.QueueDurable,
	// 	true,  // auto-delete
	// 	false, // exclusive
	// 	false, // no-wait
	// 	nil,
	// )
	if err != nil {
		log.Fatalf("declare queue failed: %v", err)
	}

	msgBody := make([]byte, cfg.MessageSize)
	for i := range msgBody {
		msgBody[i] = 'A'
	}

	var wg sync.WaitGroup
	msgPerGoroutine := cfg.MessageCount / cfg.Concurrency

	startTime := time.Now()
	var totalDuration time.Duration
	var mu sync.Mutex

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ch2, err := conn.Channel()
			if err != nil {
				log.Fatalf("subchannel failed: %v", err)
			}
			defer ch2.Close()

			if cfg.UseConfirm {
				if err := ch2.Confirm(false); err != nil {
					log.Fatalf("confirm mode failed: %v", err)
				}
			}

			startSend := time.Now()
			for j := 0; j < msgPerGoroutine; j++ {
				err = ch2.Publish(
					cfg.EngineName, // exchange
					cfg.QueueName,  // routing key
					false,          // mandatory
					false,          // immediate
					amqp.Publishing{
						ContentType:  "text/plain",
						Body:         msgBody,
						DeliveryMode: cfg.DeliveryMode, // 1=非持久化, 2=持久化
					})
				if err != nil {
					log.Printf("publish failed: %v", err)
				}

				if cfg.UseConfirm {
					// // 等待单个确认
					// select {
					// case confirm := <-ch2.NotifyPublish(make(chan amqp.Confirmation, 1)):
					// 	if !confirm.Ack {
					// 		log.Printf("message not confirmed")
					// 	} else {
					// 		log.Printf("message confirmed")
					// 	}
					// case <-time.After(50 * time.Second):
					// 	log.Printf("confirm timeout")
					// }
				}
			}
			if cfg.UseConfirm {
				confirms := ch2.NotifyPublish(make(chan amqp.Confirmation, msgPerGoroutine))
				var wg1 sync.WaitGroup
				wg1.Add(1)
				cnt := 0
				go func() {
					defer wg1.Done()
					// for {
					// 	if len(confirms) < msgPerGoroutine {
					// 		time.Sleep(1 * time.Second)
					// 		fmt.Println("len", len(confirms))
					// 		continue
					// 	}
					// 	break
					// }
					for cnt < msgPerGoroutine {
						for confirm := range confirms {
							cnt = int(confirm.DeliveryTag)
							if !confirm.Ack {
								log.Printf("message not confirmed")
							} else {
								// log.Printf("message confirmed")
								// fmt.Println("confirm.DeliveryTag", confirm.DeliveryTag)
							}
							if len(confirms) == 0 {
								break
							}
						}
					}
				}()
				wg1.Wait()
			}
			dur := time.Since(startSend)
			mu.Lock()
			totalDuration += dur
			mu.Unlock()
		}()

	}

	wg.Wait()
	elapsed := time.Since(startTime)
	// avgLatency := float64(totalDuration.Milliseconds()) / float64(cfg.MessageCount)
	rate := float64(cfg.MessageCount) / elapsed.Seconds()

	fmt.Printf("Test: %s\n", cfg.Mode)
	fmt.Printf("Messages: %d, Concurrency: %d\n", cfg.MessageCount, cfg.Concurrency)
	fmt.Printf("Total time: %.5fs\n", elapsed.Seconds())
	fmt.Printf("Throughput: %.2f msg/s\n", rate)
	// fmt.Printf("Avg Latency: %.2f ms/msg\n", avgLatency)
}

func main() {
	tests := []TestConfig{
		// {
		// 	Mode:         "NonDurable_NoConfirm",
		// 	EngineName:   "tag_engine_1",
		// 	QueueName:    "tag_insert_1",
		// 	QueueDurable: false,
		// 	DeliveryMode: 1,
		// 	UseConfirm:   false,
		// 	MessageSize:  1024,
		// 	MessageCount: 200000,
		// 	Concurrency:  20,
		// },
		// {
		// 	Mode:         "Durable_NoConfirm",
		// 	EngineName:   "tag_engine_2",
		// 	QueueName:    "tag_insert_2",
		// 	QueueDurable: true,
		// 	DeliveryMode: 2,
		// 	UseConfirm:   false,
		// 	MessageSize:  1024,
		// 	MessageCount: 200000,
		// 	Concurrency:  20,
		// },
		{
			Mode:         "Durable_WithConfirm",
			EngineName:   "tag_engine_3",
			QueueName:    "tag_insert_3",
			QueueDurable: true,
			DeliveryMode: 2,
			UseConfirm:   true,
			MessageSize:  1024,
			MessageCount: 200000,
			Concurrency:  20,
		},
	}
	for i := 0; i < 20; i++ {
		for _, cfg := range tests {
			runTest(cfg)
		}
	}
}
