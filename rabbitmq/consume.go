package main

import (
	"fmt"

	"github.com/mark-0609/go-examples/rabbitmq"
)

func main() {
	fmt.Println("Consume Start....")

	ch, conn := rabbitmq.CreateConnAndChannel()
	defer conn.Close()
	defer ch.Close()

	ch.DeclareQueueAndExchange()
	go rabbitmq.ConsumeMessagesWithAck(ch)
}
