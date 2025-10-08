package messaging

import (
	"os"
	"strings"
)

func GetBrokers() []string {
	return strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
}
