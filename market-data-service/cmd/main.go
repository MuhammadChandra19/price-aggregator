package main

import (
	"log"

	"github.com/muhammadchandra19/price-aggregator/market-data-service/api"
	"github.com/muhammadchandra19/price-aggregator/market-data-service/config"
)

func main() {
	store := config.RedisClient()
	server := api.NewServer(*store)

	err := server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
