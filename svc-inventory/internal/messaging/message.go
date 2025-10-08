package messaging

type Message interface {
	Key() string
	Value() []byte
	Topic() string
}

type kafkaMessage struct {
	key   string
	value []byte
	topic string
}

func (m *kafkaMessage) Key() string   { return m.key }
func (m *kafkaMessage) Value() []byte { return m.value }
func (m *kafkaMessage) Topic() string { return m.topic }
