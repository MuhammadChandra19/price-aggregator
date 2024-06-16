package market

type MarketPair map[string]Market
type Market struct {
	Pair      string  `json:"pair"`
	Pair0     string  `json:"pair0"`
	Pair1     string  `json:"pair1"`
	Ask       float64 `json:"ask"`
	Bid       float64 `json:"bid"`
	Spread    float64 `json:"spread"`
	UpdatedAt string  `json:"updated_at"`
}
