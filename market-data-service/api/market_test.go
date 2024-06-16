package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/muhammadchandra19/price-aggregator/market-data-service/module/market"
	"github.com/stretchr/testify/assert"
)

func TestMarket_GetMarketPrice(t *testing.T) {
	client, mock := redismock.NewClientMock()

	server := NewServer(client)

	testCases := []struct {
		name     string
		param    string
		mockFn   func(mock redismock.ClientMock)
		assertFn func(t *testing.T, res []market.MarketPair, err error)
	}{
		{
			name:  "should return successfully",
			param: "BTC",
			mockFn: func(mock redismock.ClientMock) {
				mock.ExpectKeys("price:*:BTC*").SetVal([]string{"BTCUSDT"})
				mock.ExpectGet("BTCUSDT").SetVal(`{
					"pair":"BTCUSDT",
					"pair0":"BTC",
					"pair1":"USDT",
					"ask":70000.98,
					"bid":70000,
					"spread":0.98,
					"updated_at":"datetime"
					}`,
				)
			},
			assertFn: func(t *testing.T, res []market.MarketPair, err error) {
				assert.Equal(t, []market.MarketPair{{
					"BTCUSDT": {
						Pair:      "BTCUSDT",
						Pair0:     "BTC",
						Pair1:     "USDT",
						Ask:       70000.98,
						Bid:       70000,
						Spread:    0.98,
						UpdatedAt: "datetime",
					},
				}}, res)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// ctrl := gomock.NewController(t)
			recorder := httptest.NewRecorder()

			// marketModuleMock := mockMarket.NewMockMarketHandlerUsecase(ctrl)
			test.mockFn(mock)

			url := fmt.Sprintf("/market/%s", test.param)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			server.router.ServeHTTP(recorder, request)
			res, err := getDataResponse(recorder.Body)
			test.assertFn(t, res, err)
		})
	}
}

func getDataResponse(body *bytes.Buffer) ([]market.MarketPair, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	val, ok := resp["data"].([]interface{})
	fmt.Println(ok, val)
	if !ok {
		return nil, errors.New("unexpected response")
	}

	var pairs []market.MarketPair

	for _, pair := range val {
		for key, value := range pair.(map[string]interface{}) {
			if data, ok := value.(map[string]interface{}); ok {
				fmt.Println(data)
				ask := data["ask"].(float64)
				bid := data["bid"].(float64)
				pair := data["pair"].(string)
				pair0 := data["pair0"].(string)
				pair1 := data["pair1"].(string)
				spread := data["spread"].(float64)

				pairs = append(pairs, market.MarketPair{key: market.Market{
					Ask:       ask,
					Bid:       bid,
					Pair:      pair,
					Pair0:     pair0,
					Pair1:     pair1,
					Spread:    spread,
					UpdatedAt: data["updated_at"].(string),
				},
				})
			}
		}
	}

	return pairs, nil
}
