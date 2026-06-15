package services

import (
	"encoding/json"
	"log"

	"cinema/config"
	"cinema/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "cinema.events"
	queueName    = "booking.events"
)

var mqConn *amqp.Connection
var mqChannel *amqp.Channel

func InitMQ(cfg *config.Config) error {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	if err := ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return err
	}
	if err := ch.QueueBind(queueName, "booking.success", exchangeName, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return err
	}

	mqConn = conn
	mqChannel = ch
	return nil
}

func PublishBookingSuccess(event models.BookingEvent) error {
	if mqChannel == nil {
		return nil
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return mqChannel.Publish(exchangeName, "booking.success", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func StartBookingConsumer(onSuccess func(models.BookingEvent)) {
	if mqChannel == nil {
		return
	}

	msgs, err := mqChannel.Consume(queueName, "cinema-consumer", true, false, false, false, nil)
	if err != nil {
		log.Printf("mq consumer error: %v", err)
		return
	}

	go func() {
		for msg := range msgs {
			var event models.BookingEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				continue
			}
			onSuccess(event)
		}
	}()
}
