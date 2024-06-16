package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/aggregator"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/binance"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/config"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/degate"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/publisher"
	wsh "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws"
)

// main function is the entry point of the application.
// It sets up WebSocket feeds for Binance and Degate, establishes a connection to Kafka
// for publishing aggregated data, and starts the aggregation process.
// It also listens for OS signals to gracefully stop the application.
func main() {
	ws1 := wsh.NewWebSocket()
	ws2 := wsh.NewWebSocket()

	binance := binance.NewBinanceFeed(ws1)
	dg := degate.NewDegateFeed(ws2)

	ctx := context.Background()
	signalCtx, cancelFn := cancelSignal(ctx)

	kafkaConn, err := config.KafkaDialer("market-ingestor")
	if err != nil {
		fmt.Println("Failed to dial Kafka leader:", err)
	}

	publisher := publisher.NewMarketPublisher(kafkaConn)

	// rdb := config.RedisClient()
	// ingestor := ingestor.NewIngestor(rdb)
	engine := aggregator.NewAggregator(publisher)
	engine.Process(signalCtx, cancelFn, []wsh.WebSocketFeed{binance, dg})

	log.Println("aggregator engine stopped")

}

// cancelSignal listens for OS signals to stop execution and close network connections properly.
// It returns a context that is canceled when an interrupt signal is received.
func cancelSignal(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	ctxCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-osSignal
	}()
	return ctxCancel, cancel
}
