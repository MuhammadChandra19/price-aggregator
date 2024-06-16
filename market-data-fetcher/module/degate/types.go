package degate

import "fmt"

type DegateResponse struct {
	Pair0           int    `json:"B"`
	Pair1           int    `json:"U"`
	BidPrice        string `json:"b"`
	BidBestQuantity string `json:"I"`
	AskPrice        string `json:"a"`
	AskBestQuantity string `json:"A"`
}

func (d *DegateResponse) GetSymbol() string {
	symbolCode := fmt.Sprintf("%d.%d", d.Pair0, d.Pair1)
	symbol := ""
	for key, v := range mapPairs {
		if v == symbolCode {
			symbol = key
			return symbol
		}
	}

	return symbol
}

var mapPairs = map[string]string{
	"ethusdt":   "0.3",
	"ethusdc":   "0.2",
	"aaveusdc":  "36.2",
	"solusdc":   "143.2",
	"rndrusdc":  "42.2",
	"linkusdc":  "50.2",
	"arbusdc":   "90.2",
	"maticusdc": "68.2",
	"pepeusdc":  "94.2",
	"enausdc":   "147.2",
}
