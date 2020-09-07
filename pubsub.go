package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
)

const (
	bootstrapServers = ConfluentServer
	ccloudAPIKey     = ConfluentApiKey
	ccloudAPISecret  = ConfluentSecret
)

// publish send a message to the given topic.
func publish(msg []byte, topic string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname: %v\n", err)
	}
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         hostname,
		"sasl.mechanisms":   "PLAIN",
		"security.protocol": "SASL_SSL",
		"sasl.username":     ccloudAPIKey,
		"sasl.password":     ccloudAPISecret})

	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg},
		nil)

	// Wait for delivery report
	e := <-producer.Events()

	message := e.(*kafka.Message)
	if message.TopicPartition.Error != nil {
		log.Printf("failed to deliver message: %v\n",
			message.TopicPartition)
	} else {
		log.Printf("delivered to topic %s [%d] at offset %v\n",
			*message.TopicPartition.Topic,
			message.TopicPartition.Partition,
			message.TopicPartition.Offset)
	}

	producer.Close()
}

// getConsumer returns a Kafka consumer for the given topic.
func getConsumer(topic string) *kafka.Consumer {
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
	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic: %v\n", err)
	}
	return consumer
}
