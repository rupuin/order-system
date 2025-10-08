package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishEvent(topic, key string, event any) error
}

func NewProducer(brokers []string) Producer {
	return &kafkaProducer{brokers: brokers}
}

type kafkaProducer struct {
	brokers []string
}

func (p *kafkaProducer) PublishEvent(topic, key string, event any) error {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(p.brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("Error marshalling event: %v", err)
	}

	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: eventJSON,
		},
	)

	if err != nil {
		return fmt.Errorf("Error sending Kafka message: %v", err)
	}

	return nil
}
