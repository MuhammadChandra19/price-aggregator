package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
	publisher "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/publisher/mock"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws"
	wsMock "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws/mock"
)

func TestProcess(t *testing.T) {
	// Create a new aggregator
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	publisher := publisher.NewMockMarketPublisher(ctrl)
	agg := NewAggregator(publisher)

	// Create a context with cancel function
	stopCtx, cancel := context.WithCancel(context.Background())

	// Create mock feeds
	mockFeed1 := wsMock.NewMockFeed(cancel)
	// mockFeed2 := wsMock.NewMockFeed()

	trades := []module.Trade{
		{
			Pair: "BTCUSDT",
			Ask:  "70000",
			Bid:  "6999.8",
		},
		{
			Pair: "ETHUSDT",
			Ask:  "3600",
			Bid:  "3598",
		},
	}

	mockFeed1.SetTrades(trades)
	// mockFeed2.SetTrades(trades)

	// Add mock feeds to a slice
	feeds := []ws.WebSocketFeed{mockFeed1}

	// Run the aggregator process
	agg.Process(stopCtx, cancel, feeds)
	d, _ := json.MarshalIndent(agg.trades, "", " ")
	fmt.Println(string(d))

	// Check the aggregated trades
	if len(agg.trades) != 4 {
		t.Errorf("expected 4 trades, got %d", len(agg.trades))
	}

	if trade, ok := agg.trades["BTCUSD"]; !ok || trade.AskFloat != 70000 || trade.BidFloat != 6999.8 {
		t.Errorf("unexpected BTCUSD trade: %+v", trade)
	}

	if trade, ok := agg.trades["ETHUSD"]; !ok || trade.AskFloat != 3600 || trade.BidFloat != 3598 {
		t.Errorf("unexpected ETHUSD trade: %+v", trade)
	}
}
