package market

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type MarketHandlerUsecase interface {
	GetMarketPrice(ctx context.Context, market string) ([]MarketPair, error)
}

type MarketHandler struct {
	store redis.Client
}

func NewMarketHandler(store redis.Client) MarketHandlerUsecase {
	return &MarketHandler{store: store}
}

func (m *MarketHandler) GetMarketPrice(ctx context.Context, market string) ([]MarketPair, error) {
	pattern := fmt.Sprintf("price:*:%s*", market)
	keys, err := m.store.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("Failed to find keys in Redis: %v", err)
		return nil, err
	}

	marketPairs := []MarketPair{}

	for _, key := range keys {
		val, err := m.store.Get(ctx, key).Result()
		if err != nil {
			log.Printf("Failed to get market data from Redis: %v", err)
			return nil, err
		}

		var market Market
		err = json.Unmarshal([]byte(val), &market)
		if err != nil {
			log.Printf("Failed to unmarshal market data: %v", err)
			return nil, err
		}

		marketPairs = append(marketPairs, map[string]Market{market.Pair: market})

	}
	return marketPairs, nil
}
