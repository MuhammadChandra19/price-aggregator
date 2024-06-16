package degate

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws"
)

type DegateConfig struct {
	Endpoint string
	Pairs    []string
}

type DegateFeed struct {
	conf *DegateConfig
	ws   ws.WebsocketHandler
}

func NewDegateFeed(handler ws.WebsocketHandler) ws.WebSocketFeed {
	return &DegateFeed{
		ws: handler,
		conf: &DegateConfig{
			Endpoint: "wss://v1-mainnet-ws.degate.com/ws",
			Pairs: []string{
				"ethusdt",
				"ethusdc",
				"aaveusdc",
				"solusdc",
				"rndrusdc",
				"linkusdc",
				"arbusdc",
				"maticusdc",
				"pepeusdc",
				"enausdc",
			},
		},
	}
}

func (d *DegateFeed) Subscribe() error {
	if err := d.ws.Connect(d.conf.Endpoint); err != nil {
		fmt.Printf("fail to connect into websocket=%s err=%s", d.conf.Endpoint, err)
		return err
	}

	err := d.ws.WriteJSON(map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": d.parsePairs(),
		"id":     1,
	})

	if err != nil {
		log.Printf("fail to subscribe into websocket=%s err=%s", d.conf.Endpoint, err)
		return err
	}

	return nil
}

func (d *DegateFeed) Read() (*module.Trade, error) {
	resp := &DegateResponse{}
	err := d.ws.ReadJSON(&resp)

	if errors.Is(err, ws.ErrStopped) {
		log.Println("websocket stopped")
		return nil, ws.ErrStopped
	}

	if err != nil {
		// d, _ := json.MarshalIndent(resp, "", " ")
		// fmt.Println("error", string(d))
		return nil, fmt.Errorf("websocket read fail: %s", err)
	}

	// if resp.Pair0 == 0 || resp.Pair1 == 0 {
	// 	return nil, nil
	// }

	return &module.Trade{
		Source:    "DeGate",
		Pair:      strings.ToUpper(resp.GetSymbol()),
		Ask:       resp.AskPrice,
		Bid:       resp.BidPrice,
		UpdatedAt: time.Now(),
	}, nil

}

func (d *DegateFeed) TurnOff() error {
	err := d.ws.WriteJSON(map[string]interface{}{
		"method": "UNSUBSCRIBE",
		"params": []string{"btcusdt@ticker", "ethusdt@ticker"},
		"id":     1,
	})

	if err != nil {
		log.Printf("error to unsubscribe DeGate channel=%s err=%s", d.conf.Endpoint, err)
		return err
	}

	if err := d.ws.Disconnect(); err != nil {
		log.Printf("fail to disconnect into DeGate websocket=%s err=%s", d.conf.Endpoint, err)
		return err
	}

	fmt.Println("Connection turned off")
	return nil
}

func (d *DegateFeed) parsePairs() []string {
	pairs := []string{}
	for _, pair := range d.conf.Pairs {
		v, ok := mapPairs[pair]
		if ok {
			pairs = append(pairs, v+"@bookTicker")
		}
	}
	return pairs
}
