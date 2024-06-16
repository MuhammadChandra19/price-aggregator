package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/publisher"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws"
)

type Aggregator struct {
	trades    map[string]*module.Trade
	publisher publisher.MarketPublisher
}

func NewAggregator(publisher publisher.MarketPublisher) *Aggregator {
	return &Aggregator{
		trades:    map[string]*module.Trade{},
		publisher: publisher,
	}
}

func (a *Aggregator) Process(ctx context.Context, cancelFunc context.CancelFunc, wsfs []ws.WebSocketFeed) {
	trades := make(chan module.Trade, len(wsfs))
	wg := sync.WaitGroup{}
	wg.Add(len(wsfs))
	for _, wsf := range wsfs {

		go func(wsf ws.WebSocketFeed) {
			defer wg.Done()
			a.Listen(ctx, cancelFunc, wsf, trades)
		}(wsf)
	}

	ticker := time.NewTicker(1 * time.Minute) // Create a ticker that ticks every 1 minute
	go func() {
		// wg.Done()
		for range ticker.C {
			p, err := json.Marshal(a.trades)
			if err != nil {
				fmt.Println(err)
			}
			a.publisher.PublishEvent(publisher.PublisherRequest{
				Msg: p,
			})
		}
	}()

	go func() {
		wg.Wait()
		close(trades)
		ticker.Stop()
		a.publisher.CloseEvent()
	}()

	for trade := range trades {
		a.aggregatePrices(trade)
	}

}

func (a *Aggregator) aggregatePrices(trade module.Trade) {
	mtx := sync.Mutex{}
	mtx.Lock()

	err := trade.ParsePriceToFloat()
	if err == nil {
		usdPair := strings.ReplaceAll(trade.Pair, "USDT", "USD")
		usdPair = strings.ReplaceAll(usdPair, "USDC", "USD")
		usdPairTrade, ok := a.trades[usdPair]
		if !ok {
			usdPairTrade = &trade
		}
		usdPairTrade.AskFloat = a.getAvg(usdPairTrade.AskFloat, trade.AskFloat)
		usdPairTrade.BidFloat = a.getAvg(usdPairTrade.BidFloat, trade.BidFloat)
		a.trades[usdPair] = usdPairTrade

		originalPairTrade, ok := a.trades[trade.Pair]
		if !ok {
			originalPairTrade = &trade
		}
		// fmt.Println(originalPairTrade)
		originalPairTrade.AskFloat = a.getAvg(originalPairTrade.AskFloat, trade.AskFloat)
		originalPairTrade.BidFloat = a.getAvg(originalPairTrade.BidFloat, trade.BidFloat)
		a.trades[trade.Pair] = originalPairTrade
	}

	mtx.Unlock()
}

func (a *Aggregator) getAvg(existing float64, new float64) float64 {
	return (existing + new) / float64(2)
}

func (a *Aggregator) Listen(ctx context.Context, cancelFunc context.CancelFunc, wsf ws.WebSocketFeed, tradeChan chan module.Trade) {
	err := wsf.Subscribe()
	if err != nil {
		cancelFunc()
		wsf.TurnOff()
		log.Fatal(err)
	}

	for {
		resp, err := wsf.Read()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if resp == nil {
			continue
		}

		select {
		case tradeChan <- *resp:

		case <-ctx.Done():
			cancelFunc()
			wsf.TurnOff()
			return
		}
	}
}
