package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	rabbitMQURL  = "amqp://guest:guest@114.132.210.241:5672/"
	queueName    = "example_queue"
	exchangeName = "example_exchange"
	routingKey   = "example_key"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: `%s`", msg, err)
	}
}

func CreateConnAndChannel() (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch, conn
}

func DeclareQueueAndExchange(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,   //amqp.Table{"alternate-exchange": "my-backup-exchange"}, // 设置备份交换机
	)
	failOnError(err, "Failed to declare exchange")

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	failOnError(err, "Failed to bind queue to exchange")
}

func PublishMessageWithConfirm(ch *amqp.Channel, body []byte, retry, timeout func()) {
	err := ch.Confirm(false)
	failOnError(err, "Failed to enable publisher confirms")

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 10000))

	err = ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	failOnError(err, "Failed to publish a message")
	select {
	case confirm := <-confirms:
		if !confirm.Ack {
			log.Printf("Failed delivery of message with body %s", body)
			// 可以在这里实现重试逻辑
			retry()
		}
	case <-time.After(3 * time.Second):
		log.Println("Timed out waiting for confirmation", err)
		// 可以在这里实现超时后的处理逻辑
		timeout()
	}
	log.Println("end...", string(body))
}

func ConsumeMessagesWithAck(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool, 4)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			// 在这里处理消息，确保没有发生错误，否则消息可能会被丢失
			err := processMessage(d.Body)
			if err == nil {
				d.Ack(false) // 确认消息已经被处理
			} else {
				log.Printf("Error processing message: %v", err)
				// 可以在这里实现错误处理和重试逻辑
				// 处理消息失败时，重新发布消息到队列
				d.Nack(false, true)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processMessage(body []byte) error {
	return nil
}
