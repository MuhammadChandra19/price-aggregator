package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
)

//go:generate mockgen -source handler.go -destination mock/handler_mock.go -package=ingestor
type IngestorHandler interface {
	WriteMarketFeed(ctx context.Context, feeds map[string]*module.Trade)
}

type Ingestor struct {
	redis *redis.Client
}

func NewIngestor(redis *redis.Client) IngestorHandler {
	return &Ingestor{
		redis: redis,
	}
}

func (i *Ingestor) WriteMarketFeed(ctx context.Context, feeds map[string]*module.Trade) {
	for key, feed := range feeds {
		pair0, pair1 := splitPair(key)
		market := map[string]interface{}{
			"pair":       key, // BTCUSDT
			"pair0":      pair0,
			"pair1":      pair1,
			"ask":        feed.AskFloat,
			"bid":        feed.BidFloat,
			"spread":     feed.Spread,
			"updated_at": feed.UpdatedAt.Local().Format("2024-06-16T03:05.30Z"),
		}
		tradeJSON, err := json.Marshal(market)
		if err != nil {
			log.Printf("Failed to marshal trade data: %v", err)
		}

		redisKey := fmt.Sprintf("price:%s:%s", feed.Source, key)

		err = i.redis.Set(ctx, redisKey, tradeJSON, time.Hour).Err()
		if err != nil {
			log.Fatalf("Failed to save trade data to Redis: %v", err)
		}
	}

	fmt.Println("Data saved to Redis successfully with TTL of 1 hour!")
}

func splitPair(pair string) (string, string) {
	for _, token := range module.SupportedTokens {
		if strings.HasSuffix(pair, strings.ToUpper(token)) {
			return pair[:len(pair)-len(token)], pair[len(pair)-len(token):]
		}
	}
	// Default behavior for unknown pairs
	if len(pair) >= 4 {
		return pair[:3], pair[3:]
	}
	return pair, ""
}
