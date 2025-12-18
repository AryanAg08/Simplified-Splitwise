package queue

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var Connection *amqp091.Connection
var Channel *amqp091.Channel

func ConnectRabbitMQ() {
	var err error

	Connection, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatal("Rabbit MQ connection failed!!", err)
	}

	Channel, err = Connection.Channel()

	if err != nil {
		log.Fatal("Rabbit MQ Channel failed!!", err)
	}

	_, err = Channel.QueueDeclare(
		"expense_added",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Queue Creation Failed!!", err)
	}

	log.Println("Connected to RabbitMQ!!")
}
