package main

import (
	"database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
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
			alert := createAlert(e.Value)
			persistAlert(alert)
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Fatalf("%% Error: %v\n", e)
		}
	}
	// Persist to Postgres

	// Publish message to Kafka (notification queue)
}

// createAlert creates an alert from bytes.
func createAlert(msg []byte) *Alert {
	alert := &Alert{}
	if err := proto.Unmarshal(msg, alert); err != nil {
		log.Fatalf("failed to unmarshal kafka message into alert: %v\n", err)
	}
	return alert
}

// persistAlert saves the alert to a PostgreSQL store
func persistAlert(alert *Alert) {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		INSTANCE_CONNECTION_NAME,
		DATABASE_NAME,
		DATABASE_USER,
		PASSWORD)
	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v\n", err)
	}
	rows, err := db.Query("SELECT id, model_name FROM face_detection_model")
	if err != nil {
		log.Fatalf("failed to query DB: %v\n", err)

	}
	var id string
	var name string
	for rows.Next() {
		if err = rows.Scan(&id, &name); err != nil {
			log.Fatalf("failed to scan row: %v\n", err)
		}
		fmt.Printf("Face detection model ID# %s name: `%s`\n", id, name)

	}
}
