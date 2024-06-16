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

// Aggregator struct to hold trades and publisher instance
type Aggregator struct {
	trades    map[string]*module.Trade
	publisher publisher.MarketPublisher
}

// NewAggregator creates a new Aggregator instance with the given publisher.
// It initializes the trades map and assigns the provided publisher to the Aggregator.
func NewAggregator(publisher publisher.MarketPublisher) *Aggregator {
	return &Aggregator{
		trades:    map[string]*module.Trade{},
		publisher: publisher,
	}
}

// Process starts the aggregation process by managing WebSocket feeds, collecting trades,
// and periodically publishing aggregated trade data. It takes a context for cancellation,
// a cancel function, and a slice of WebSocket feeds. It spawns goroutines to listen
// to each WebSocket feed and collects trades into a channel. It also sets up a ticker
// to periodically publish the aggregated trades and ensures proper shutdown by waiting
// for all goroutines to complete.
func (a *Aggregator) Process(ctx context.Context, cancelFunc context.CancelFunc, wsfs []ws.WebSocketFeed) {
	trades := make(chan module.Trade, len(wsfs)) // Channel to receive trades
	wg := sync.WaitGroup{}
	wg.Add(len(wsfs)) // WaitGroup to wait for all goroutines to finish

	for _, wsf := range wsfs {
		go func(wsf ws.WebSocketFeed) {
			defer wg.Done()
			a.Listen(ctx, cancelFunc, wsf, trades)
		}(wsf)
	}

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
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

// aggregatePrices processes a single trade, updating the trades map with the
// aggregated prices. It locks the trades map to ensure thread safety, parses
// the trade prices to float, and updates both the USD pair and the original
// trading pair with the new average prices.
func (a *Aggregator) aggregatePrices(trade module.Trade) {
	mtx := sync.Mutex{} // Mutex to ensure thread safety
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
		originalPairTrade.AskFloat = a.getAvg(originalPairTrade.AskFloat, trade.AskFloat)
		originalPairTrade.BidFloat = a.getAvg(originalPairTrade.BidFloat, trade.BidFloat)
		a.trades[trade.Pair] = originalPairTrade
	}

	mtx.Unlock()
}

// getAvg calculates the average of two float64 values and returns it.
// This is used to compute the new average prices for the trades.
func (a *Aggregator) getAvg(existing float64, new float64) float64 {
	return (existing + new) / float64(2)
}

// Listen subscribes to a WebSocket feed and continuously reads from it,
// sending received trades to the provided trade channel. It also handles
// errors by cancelling the context and turning off the WebSocket feed
// if necessary. It stops reading when the context is done.
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
