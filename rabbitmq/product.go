package main

import (
	"fmt"
	"time"

	"github.com/mark-0609/go-examples/rabbitmq"
)

func main() {

	ch, conn := rabbitmq.CreateConnAndChannel()
	defer conn.Close()
	defer ch.Close()

	ch.DeclareQueueAndExchange()

	// 生产者
	go func() {
		for {
			body := []byte(fmt.Sprintf("Hello, RabbitMQ! - %d", time.Now().Unix()))
			rabbitmq.PublishMessageWithConfirm(ch, body)
			time.Sleep(1 * time.Second)
		}
	}()
}

// func fanIn([]chan interface{}) {}

// func fanOut([]chan interface{}) []chan interface{} {

// }
