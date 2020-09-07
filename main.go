package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

const (
	bootstrapServers = "pkc-4r297.europe-west1.gcp.confluent.cloud:9092"
	ccloudAPIKey     = ConfluentApiKey
	ccloudAPISecret  = ConfluentSecret
)

var topics = []string{"alert-topic"}

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     ccloudAPIKey,
		"sasl.password":     ccloudAPISecret,
		"group.id":          "alert-consumer",
		"auto.offset.reset": "earliest"})
	if err != nil {
		log.Fatalf("Failed to connect to Kafka %v", err)
	}

	// Read from Kafka
	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic: %v\n", err)
	}
	defer consumer.Close()

	for {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Fatalf("%% Error: %v\n", e)
			//default:
			//	log.Printf("Ignored %v\n", e)
		}
	}

	// Persist to Postgres

	// Publish message to Kafka (notification queue)
}
