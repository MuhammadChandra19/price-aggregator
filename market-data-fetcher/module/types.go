package module

import (
	"fmt"
	"strconv"
	"time"
)

var SupportedTokens = []string{"btc", "eth", "aave", "ena", "sol", "link", "arb", "matic", "rndr"}

// Trade represents the main trading operation
type Trade struct {
	Pair      string    `json:"pair"`
	Ask       string    `json:"ask"`
	Bid       string    `json:"bid"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updated_at"`

	AskFloat float64
	BidFloat float64
	Spread   float64
}

func (t *Trade) ParsePriceToFloat() error {
	ask, err := strconv.ParseFloat(t.Ask, 64)
	if err != nil {
		fmt.Println("unexpected error when parsing ask", t, err)
		return err
	}

	t.AskFloat = ask

	bid, err := strconv.ParseFloat(t.Bid, 64)
	if err != nil {
		fmt.Println("unexpected error when parsing bid", t, err)
		return err
	}

	t.BidFloat = bid
	t.Spread = ask - bid

	return nil
}
