package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

func main() {

	// подключение к серверу NATS Streaming
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	bytes, err := os.ReadFile("publisher/model.json")
	if err != nil {
		log.Fatalf("Failed open file: %v", err)
	}

	err = nc.Publish("purchases", bytes)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Println("Message published successfully")
}
