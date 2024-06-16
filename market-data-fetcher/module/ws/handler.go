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

//go:generate mockgen -source handler.go -destination mock/handler_mock.go -package=ws
type WebsocketHandler interface {
	WriteJSON(message interface{}) error
	ReadJSON(message interface{}) error
	Connect(endpoint string) error
	Disconnect() error
}

// WebSocketFeed defines minimal websocket functions to *Trade feed
type WebSocketFeed interface {
	Subscribe() error
	TurnOff() error
	Read() (*module.Trade, error)
}

type WebSocket struct {
	IsConnected bool
	Conn        *gws.Conn
}

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
