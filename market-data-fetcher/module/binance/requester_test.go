package binance

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	wsMock "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module/ws/mock"
	"github.com/stretchr/testify/assert"
)

func TestBinanceFeed_Subscribe(t *testing.T) {
	testCases := []struct {
		name     string
		mockFn   func(mockWs *wsMock.MockWebsocketHandler)
		assertFn func(t *testing.T, err error)
	}{
		{
			name: "should handle error when connect",
			mockFn: func(mockWs *wsMock.MockWebsocketHandler) {
				mockWs.EXPECT().Connect("wss://stream.binance.com:9443/ws").Return(errors.New("fail to connect"))
			},
			assertFn: func(t *testing.T, err error) {
				assert.Equal(t, err, errors.New("fail to connect"))
			},
		},
		{
			name: "should handle error when write",
			mockFn: func(mockWs *wsMock.MockWebsocketHandler) {
				mockWs.EXPECT().Connect("wss://stream.binance.com:9443/ws").Return(nil)
				mockWs.EXPECT().WriteJSON(map[string]interface{}{
					"method": "SUBSCRIBE",
					"params": []string{"btcusdt@ticker", "ethusdt@ticker"},
					"id":     1,
				}).Return(errors.New("fail to write"))
			},
			assertFn: func(t *testing.T, err error) {
				assert.Equal(t, err, errors.New("fail to write"))
			},
		},
		{
			name: "success",
			mockFn: func(mockWs *wsMock.MockWebsocketHandler) {
				mockWs.EXPECT().Connect("wss://stream.binance.com:9443/ws").Return(nil)
				mockWs.EXPECT().WriteJSON(map[string]interface{}{
					"method": "SUBSCRIBE",
					"params": []string{"btcusdt@ticker", "ethusdt@ticker"},
					"id":     1,
				}).Return(nil)
			},
			assertFn: func(t *testing.T, err error) {
				assert.Equal(t, err, nil)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := wsMock.NewMockWebsocketHandler(ctrl)
			binance := NewBinanceFeed(mock)

			test.mockFn(mock)

			err := binance.Subscribe()
			if err != nil {
				test.assertFn(t, err)
			}

		})
	}
}
