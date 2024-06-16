package binance

import (
	"errors"
	"fmt"
	"log"
	"time"

	module "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws"
)

// Response represents Binance Channel WebSocket Response
type Response struct {
	Symbol          string `json:"s"`
	BidPrice        string `json:"b"`
	BidBestQuantity string `json:"B"`
	AskPrice        string `json:"a"`
	AskBestQuantity string `json:"A"`
}

type BinanceConfig struct {
	Endpoint string
	Pairs    []string
}

type BinanceFeed struct {
	conf *BinanceConfig
	ws   ws.WebsocketHandler
}

func NewBinanceFeed(handler ws.WebsocketHandler) ws.WebSocketFeed {
	return &BinanceFeed{
		ws: handler,
		conf: &BinanceConfig{
			Endpoint: "wss://stream.binance.com:9443/ws",
			Pairs: []string{
				"btcusdt@bookTicker",
				"ethusdt@bookTicker",
				"aaveusdt@bookTicker",
				"enausdt@bookTicker",
				"solusdt@bookTicker",
				"linkusdt@bookTicker",
				"arbusdt@bookTicker",
				"maticusdt@bookTicker",
				"pepeusdt@bookTicker",
				"rndrusdt@bookTicker",
			},
		},
	}
}

func (b *BinanceFeed) Subscribe() error {
	if err := b.ws.Connect(b.conf.Endpoint); err != nil {
		fmt.Printf("fail to connect into websocket=%s err=%s", b.conf.Endpoint, err)
		return err
	}

	err := b.ws.WriteJSON(map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": b.conf.Pairs,
		"id":     1,
	})

	if err != nil {
		log.Printf("fail to subscribe into websocket=%s err=%s", b.conf.Endpoint, err)
		return err
	}

	return nil

}

func (b *BinanceFeed) Read() (*module.Trade, error) {
	resp := &Response{}
	err := b.ws.ReadJSON(&resp)

	if errors.Is(err, ws.ErrStopped) {
		log.Println("websocket stopped")
		return nil, ws.ErrStopped
	}

	if err != nil {
		return nil, fmt.Errorf("websocket read fail: %s", err)
	}

	return &module.Trade{
		Pair:      resp.Symbol,
		Ask:       resp.AskPrice,
		Bid:       resp.BidPrice,
		Source:    "Binance",
		UpdatedAt: time.Now(),
	}, nil

}

func (b *BinanceFeed) TurnOff() error {
	err := b.ws.WriteJSON(map[string]interface{}{
		"method": "UNSUBSCRIBE",
		"params": b.conf.Pairs,
		"id":     1,
	})

	if err != nil {
		log.Printf("error to unsubscribe channel=%s err=%s", b.conf.Endpoint, err)
		return err
	}

	if err := b.ws.Disconnect(); err != nil {
		log.Printf("fail to disconnect into websocket=%s err=%s", b.conf.Endpoint, err)
		return err
	}

	fmt.Println("Connection turned off")
	return nil
}
