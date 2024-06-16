package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/muhammadchandra19/price-aggregator/market-data-service/module/market"
)

type Server struct {
	store  redis.Client
	router *gin.Engine
}

func NewServer(store redis.Client) *Server {
	marketHandler := market.NewMarketHandler(store)
	marketService := NewMarket(marketHandler)

	server := &Server{store: store}
	router := gin.Default()
	router.GET("/market/:market", marketService.GetMarketPrice)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func SuccessResponse(code int, data interface{}) gin.H {
	var statusCode string
	if code == http.StatusOK {
		statusCode = "F00000"
	}
	return gin.H{
		"status":         1,
		"status_code":    statusCode,
		"status_message": "success",
		"data":           data,
	}
}
