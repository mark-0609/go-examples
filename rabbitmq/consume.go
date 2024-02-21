package rabbitmq

import (
	"fmt"
)

func Consume() {
	fmt.Println("Consume Start....")

	ch, conn := CreateConnAndChannel()
	defer conn.Close()
	defer ch.Close()

	// DeclareQueueAndExchange(ch)
	ConsumeMessagesWithAck(ch)
	select {}
}
