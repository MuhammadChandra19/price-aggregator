package config

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func KafkaDialer(topic string) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, 0)

	return conn, err
}
