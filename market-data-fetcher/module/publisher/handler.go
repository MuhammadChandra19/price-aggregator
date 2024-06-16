package publisher

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

type PublisherRequest struct {
	Msg   []byte
	Topic string
}

//go:generate mockgen -source handler.go -destination mock/handler_mock.go -package=publisher
type MarketPublisher interface {
	PublishEvent(payload PublisherRequest) error
	CloseEvent()
}

type MarketPublisherConnection struct {
	conn *kafka.Conn
}

func NewMarketPublisher(conn *kafka.Conn) MarketPublisher {
	return &MarketPublisherConnection{
		conn: conn,
	}
}

func (s *MarketPublisherConnection) PublishEvent(payload PublisherRequest) error {
	message := kafka.Message{
		Value: payload.Msg,
	}
	_, err := s.conn.WriteMessages(message)
	if err != nil {
		fmt.Println("Failed to write message to Kafka:", err)
		return err
	}

	return nil
}

func (s *MarketPublisherConnection) CloseEvent() {
	if err := s.conn.Close(); err != nil {
		fmt.Println("Error when closing:", err)
	}
}
