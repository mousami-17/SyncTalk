package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"realtime-chat/src/config"
	"realtime-chat/src/models"
)

type KafkaConfig struct {
	Host       string
	Port       string
	kafkaTopic string
	GroupID    string
}

var kafkaConfig KafkaConfig

func init() {
	kafkaConfig = KafkaConfig{
		Host:       config.Config.KafkaHost,
		Port:       config.Config.KafkaPort,
		kafkaTopic: config.Config.KafkaTopic,
		GroupID:    config.Config.KafkaGroupID,
	}
}

func PublishMessage(message models.Message) {
	log.Printf("[Kafka] Publishing message to topic %s at %s:%s", kafkaConfig.kafkaTopic, kafkaConfig.Host, kafkaConfig.Port)

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaConfig.Host + ":" + kafkaConfig.Port},
		Topic:   kafkaConfig.kafkaTopic,
	})
	defer w.Close()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("[Kafka] Error marshalling message:", err)
		return
	}

	// Publish message to kafka
	err = w.WriteMessages(context.Background(), kafka.Message{
		Value: messageBytes,
	})

	if err != nil {
		log.Println("[Kafka] Error publishing message to kafka:", err)
		return
	}

	log.Printf("[Kafka] Message published successfully: %s", message.Content)
}