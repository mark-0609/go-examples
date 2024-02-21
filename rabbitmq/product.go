package rabbitmq

import (
	"fmt"
)

func Product() {

	ch, _ := CreateConnAndChannel()
	// defer conn.Close()
	defer ch.Close()

	DeclareQueueAndExchange(ch)
	go func() {
		for i := 0; i < 4000; i++ {
			body := []byte(fmt.Sprintf("RabbitMQ! - %d", i))
			PublishMessageWithConfirm(ch, body, func() {}, func() {})
			// time.Sleep(600 * time.Millisecond)
		}
	}()
}

// func fanIn([]chan interface{}) {}

// func fanOut([]chan interface{}) []chan interface{} {

// }
