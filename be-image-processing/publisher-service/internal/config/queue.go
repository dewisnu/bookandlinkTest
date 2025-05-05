package config

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
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

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Setup exchanges and queues
	err = setupExchangesAndQueues(channel)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %w", err)
	}

	// Inspect the already-declared queue (avoid re-declaration with wrong args)
	queue, err := channel.QueueInspect("image_jobs")
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to inspect main queue: %w", err)
	}

	log.Println("✅ Successfully connected to RabbitMQ with retry & DLQ setup")
	return &RabbitMQ{conn: conn, channel: channel, queue: queue}, nil
}

func setupExchangesAndQueues(ch *amqp.Channel) error {
	// Dead Letter Exchange
	if err := ch.ExchangeDeclare("dead_letter_exchange", "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare DLX: %w", err)
	}

	// Retry Exchange
	if err := ch.ExchangeDeclare("retry_exchange", "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare retry exchange: %w", err)
	}

	// Retry Queue
	_, err := ch.QueueDeclare(
		"retry_queue",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl":             int32(5000),  // 5s
			"x-dead-letter-exchange":    "",           // default exchange
			"x-dead-letter-routing-key": "image_jobs", // back to main
		},
	)
	if err != nil {
		return fmt.Errorf("declare retry_queue: %w", err)
	}
	if err := ch.QueueBind("retry_queue", "retry", "retry_exchange", false, nil); err != nil {
		return fmt.Errorf("bind retry_queue: %w", err)
	}

	// Failed Jobs Queue
	_, err = ch.QueueDeclare(
		"failed_jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare failed_jobs queue: %w", err)
	}
	if err := ch.QueueBind("failed_jobs", "failed_jobs", "dead_letter_exchange", false, nil); err != nil {
		return fmt.Errorf("bind failed_jobs queue: %w", err)
	}

	// Main Image Jobs Queue
	_, err = ch.QueueDeclare(
		"image_jobs",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    "retry_exchange",
			"x-dead-letter-routing-key": "retry",
		},
	)
	if err != nil {
		return fmt.Errorf("declare image_jobs queue: %w", err)
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
	message := JobMessage{ID: jobID, Filename: filename}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	err = r.channel.Publish(
		"",           // default exchange
		r.queue.Name, // routing key (image_jobs)
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish job: %w", err)
	}

	log.Printf("📤 Published job ID %d to queue '%s'", jobID, r.queue.Name)
	return nil
}
