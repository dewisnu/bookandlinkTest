package main

import (
	"log"

	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// 👇 Ganti QueueDeclare dengan QueueInspect (queue sudah dibuat oleh publisher)
	q, err := ch.QueueInspect("image_jobs")
	if err != nil {
		log.Fatalf("Failed to inspect the queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("📥 Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			processJob(db, d) // Handle job
		}
	}()

	<-forever
}
