package async

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	ProcessMessages(ctx context.Context, handler MessageHandler)
	Close() error
}

type MessageHandler func(Message) error

type kafkaConsumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MaxBytes: 10e6,
	})

	return &kafkaConsumer{reader: reader}
}

func (c *kafkaConsumer) ProcessMessages(ctx context.Context, handler MessageHandler) {
	log.Printf("Starting consumer...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Consumer stopped")
			return
		default:
			msgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			kafkaMsg, err := c.reader.ReadMessage(msgCtx)
			log.Printf("Received message: key=%s, offset=%d", string(kafkaMsg.Key), kafkaMsg.Offset)
			cancel()
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "deadline exceeded") {
					continue
				}
				log.Printf("Error reading message: %v", err)
				continue
			}

			message := &kafkaMessage{
				key:   string(kafkaMsg.Key),
				value: kafkaMsg.Value,
				topic: kafkaMsg.Topic,
			}

			if err := handler(message); err != nil {
				log.Printf("Error proccessing message: %v", err)
			}
		}
	}
}
func (c *kafkaConsumer) Close() error {
	return c.reader.Close()
}
