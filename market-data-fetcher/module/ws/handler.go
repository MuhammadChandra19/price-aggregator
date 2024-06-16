package ws

import (
	"encoding/json"
	"errors"

	gws "github.com/gorilla/websocket"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
)

var (
	ErrStopped = errors.New("websocket service stopped")
)

// WebsocketHandler interface defines the essential methods for handling WebSocket connections.
// It includes methods for writing and reading JSON messages, connecting to an endpoint, and disconnecting.
//
//go:generate mockgen -source handler.go -destination mock/handler_mock.go -package=ws
type WebsocketHandler interface {
	// WriteJSON sends a JSON-encoded message over the WebSocket connection.
	// It returns an error if the WebSocket service is stopped or if there is a problem writing the message.
	WriteJSON(message interface{}) error

	// ReadJSON reads a JSON-encoded message from the WebSocket connection.
	// It returns an error if the WebSocket service is stopped or if there is a problem reading the message.
	ReadJSON(message interface{}) error

	// Connect establishes a WebSocket connection to the specified endpoint.
	// It updates the connection state and returns an error if there is a problem connecting to the endpoint.
	Connect(endpoint string) error

	// Disconnect closes the WebSocket connection and updates the connection state.
	// It returns an error if there is a problem closing the connection.
	Disconnect() error
}

// WebSocketFeed interface defines minimal WebSocket functions for subscribing to trade feeds.
// It includes methods for subscribing, turning off the feed, and reading trades.
type WebSocketFeed interface {
	Subscribe() error
	TurnOff() error
	Read() (*module.Trade, error)
}

// WebSocket struct implements the WebsocketHandler interface.
// It manages the WebSocket connection state and provides methods to interact with the WebSocket.
type WebSocket struct {
	IsConnected bool
	Conn        *gws.Conn
}

// NewWebSocket creates a new WebSocket instance with an initial disconnected state.
// It returns an instance of WebsocketHandler.
func NewWebSocket() WebsocketHandler {
	return &WebSocket{
		IsConnected: false,
	}
}

func (c *WebSocket) WriteJSON(message interface{}) error {
	if !c.IsConnected {
		return ErrStopped
	}
	if err := c.Conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (c *WebSocket) ReadJSON(message interface{}) error {
	if !c.IsConnected {
		return ErrStopped
	}
	_, r, err := c.Conn.NextReader()
	if err != nil {
		return err
	}
	d := json.NewDecoder(r)

	return d.Decode(message)
}

func (c *WebSocket) Connect(endpoint string) error {
	if c.Conn != nil {
		c.IsConnected = true
		return nil
	}

	conn, _, err := gws.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return err
	}

	c.Conn = conn
	c.IsConnected = true
	return nil
}

func (c *WebSocket) Disconnect() error {
	c.IsConnected = false
	if err := c.Conn.Close(); err != nil {
		return err
	}
	return nil
}
