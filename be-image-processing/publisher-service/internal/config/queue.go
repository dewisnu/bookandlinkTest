package config

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

// RabbitMQ represents a connection to RabbitMQ
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type JobMessage struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
}

type RabbitMQConfig struct {
	RabbitMQUrl string `json:"rabbitMQUrl"`
}

type InitRabbitMQParams struct {
	Conf *RabbitMQConfig
}

func InitRabbitMQ(param *InitRabbitMQParams) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(param.Conf.RabbitMQUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare a queue
	queue, err := channel.QueueDeclare(
		"image_jobs", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Create dead letter exchange and queue for failed jobs
	err = setupDeadLetterQueue(channel)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup dead letter queue: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")
	return &RabbitMQ{conn: conn, channel: channel, queue: queue}, nil
}

// setupDeadLetterQueue sets up a dead letter exchange and queue for failed jobs
func setupDeadLetterQueue(ch *amqp.Channel) error {
	// Declare dead letter exchange
	err := ch.ExchangeDeclare(
		"dead_letter_exchange", // name
		"direct",               // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter exchange: %w", err)
	}

	// Declare dead letter queue
	dlq, err := ch.QueueDeclare(
		"failed_jobs", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter queue: %w", err)
	}

	// Bind the dead letter queue to the exchange
	err = ch.QueueBind(
		dlq.Name,               // queue name
		"failed_jobs",          // routing key
		"dead_letter_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind dead letter queue to exchange: %w", err)
	}

	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// PublishJob sends a job message to the queue
func (r *RabbitMQ) PublishJob(jobID int64, filename string) error {
	// Create message
	message := JobMessage{
		ID:       jobID,
		Filename: filename,
	}

	// Marshal to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling job message: %w", err)
	}

	// Publish message
	err = r.channel.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("error publishing job message: %w", err)
	}

	log.Printf("Published job %d to queue", jobID)
	return nil
}
