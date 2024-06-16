package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
	"github.com/muhammadchandra19/price-aggregator/market-data-ingestor/config"
	"github.com/muhammadchandra19/price-aggregator/market-data-ingestor/module/ingestor"
	"github.com/segmentio/kafka-go"
)

func main() {

	ctx := context.Background()
	signalCtx, cancelFn := cancelSignal(ctx)
	// set up a Kafka consumer
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "market-ingestor",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  20e6, // 10MB
	})

	rdb := config.RedisClient()
	ingestor := ingestor.NewIngestor(rdb)

	messageChan := make(chan kafka.Message)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			m, err := consumer.ReadMessage(signalCtx)
			if err != nil {
				fmt.Println("error while reading message", err)
			}

			select {
			case messageChan <- m:
			case <-signalCtx.Done():
				cancelFn()
				wg.Done()

				consumer.Close()
			}

		}
	}()

	go func() {
		wg.Wait()
		close(messageChan)
	}()

	for msg := range messageChan {
		var trades map[string]*module.Trade
		err := json.Unmarshal(msg.Value, &trades)
		if err == nil {
			ingestor.WriteMarketFeed(signalCtx, trades)
		} else {
			fmt.Println(err)
		}
	}
}

// cancelSignal Listen OS signal to stop execution close network connection properly
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
