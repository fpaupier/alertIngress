package main

import (
	"database/sql"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/protobuf/proto"
	"log"
	"time"
)

const (
	IngressTopic       = "alert-topic"
	AlertToNotifyTopic = "to-notify-topic"
)

func main() {
	consumer := getConsumer(IngressTopic)
	defer consumer.Close()

	for {
		ev := consumer.Poll(0)
		switch e := ev.(type) {
		case *kafka.Message:
			alert := createAlert(e.Value)
			alertId := persistAlert(alert)
			pingNotificationService(alertId)
		case kafka.PartitionEOF:
			log.Printf("%% Reached %v\n", e)
		case kafka.Error:
			log.Fatalf("%% Error: %v\n", e)
		}
	}
}

// pingNotificationService publishes the id to a Kafka topic for it to be processed by a notification service.
func pingNotificationService(id int) {
	publish([]byte{byte(id)}, AlertToNotifyTopic)
}

// createAlert creates an alert from bytes.
func createAlert(msg []byte) *Alert {
	alert := &Alert{}
	if err := proto.Unmarshal(msg, alert); err != nil {
		log.Fatalf("failed to unmarshal kafka message into alert: %v\n", err)
	}
	return alert
}

// persistAlert saves the alert to a PostgreSQL store return the id of the alert created
func persistAlert(alert *Alert) int {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		InstanceConnectionName,
		DatabaseName,
		DatabaseUser,
		Password)
	db, err := sql.Open("cloudsqlpostgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v\n", err)
	}
	// Insert image
	//goland:noinspection SqlResolve
	rows, err := db.Query("INSERT INTO image (format, width, height, data) VALUES ($1, $2, $3, $4) RETURNING id;",
		alert.Image.Format,
		alert.Image.Size.Width,
		alert.Image.Size.Height,
		alert.Image.Data,
	)
	if err != nil {
		log.Fatalf("failed to insert image: %v\n", err)
	}
	var imageId int
	for rows.Next() {
		if err = rows.Scan(&imageId); err != nil {
			log.Fatalf("failed to recover last image inserted id: %v\n", err)
		}
	}
	_ = rows.Close()
	log.Printf("Saved image of type %s (%dH x %dW) with id %d\n", alert.Image.Format, alert.Image.Size.Height, alert.Image.Size.Width, imageId)

	// Insert alert record
	receivedAt := time.Now()
	//goland:noinspection ALL
	query := "INSERT INTO alert (event_time, received_at, device_id, face_model_id, mask_model_id, image_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	rows, err = db.Query(query, alert.EventTime, receivedAt, alert.CreatedBy.Guid, alert.FaceDetectionModel.Guid, alert.MaskClassifierModel.Guid, imageId)
	if err != nil {
		log.Fatalf("failed to execute query: %v\n", err)
	}
	var alertId int
	for rows.Next() {
		if err = rows.Scan(&alertId); err != nil {
			log.Fatalf("failed to recover last alert inserted id: %v\n", err)
		}
	}
	_ = rows.Close()
	log.Printf("Saved alert with id #%d\n", alertId)
	return alertId
}
