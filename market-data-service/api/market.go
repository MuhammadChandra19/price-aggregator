package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammadchandra19/price-aggregator/market-data-service/module/market"
)

type MarketHandler interface {
	GetMarketPrice(ctx *gin.Context)
}

type MarketConfig struct {
	marketHandler market.MarketHandlerUsecase
}

func NewMarket(marketHandler market.MarketHandlerUsecase) MarketHandler {
	return &MarketConfig{
		marketHandler: marketHandler,
	}
}

type GetMarketRequest struct {
	Market string `uri:"market" binding:"required,min=2"`
}

func (m *MarketConfig) GetMarketPrice(ctx *gin.Context) {
	var request GetMarketRequest
	err := ctx.ShouldBindUri(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	marketPair, err := m.marketHandler.GetMarketPrice(ctx, request.Market)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, SuccessResponse(http.StatusOK, marketPair))

}
