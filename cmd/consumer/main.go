package main

import (
	"Shortxn/internal/domain"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"analytics_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"analytics",
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var event domain.ClickEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}
			// TODO: Store analytics data in database
			log.Printf("Processed click event for URL ID: %s", event.URLId)
		}
	}()

	log.Printf(" [*] Waiting for analytics events")
	<-forever
}
