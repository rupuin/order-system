package async

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishEvent(key string, headers map[string]string, payload any) error
	Close() error
}

type kafkaProducer struct {
	brokers []string
	topic   string
	writer  *kafka.Writer
}

func NewProducer(brokers []string, topic string) Producer {
	return &kafkaProducer{
		brokers: brokers,
		topic:   topic,
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *kafkaProducer) PublishEvent(key string, headers map[string]string, payload any) error {
	var kafkaHeaders []kafka.Header
	for key, value := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{
			Key:   key,
			Value: []byte(value),
		})
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshalling event: %v", err)
	}

	err = p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:     []byte(key),
			Value:   payloadJSON,
			Headers: kafkaHeaders,
		},
	)

	if err != nil {
		return fmt.Errorf("Error sending Kafka message: %v", err)
	}

	return nil
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
