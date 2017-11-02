package server

import "fmt"

type TickerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Ticker  Ticker `json:"result"`
}

type Ticker struct {
	Bid  float64 `json:"Bid"`
	Ask  float64 `json:"Ask"`
	Last float64 `json:"Last"`
}

func (t Ticker) String() string {
	return fmt.Sprintf("Last: %f", t.Last)
}

type MarketResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Markets []Market `json:"result"`
}

type Market struct {
	MarketName     string  `json:"MarketName"`
	High           float64 `json:"High"`
	Low            float64 `json:"Low"`
	Volume         float64 `json:"Volume"`
	Last           float64 `json:"Last"`
	BaseVolume     float64 `json:"BaseVolume"`
	Bid            float64 `json:"Bid"`
	Ask            float64 `json:"Ask"`
	OpenBuyOrders  uint64  `json:"OpenBuyOrders"`
	OpenSellOrders uint64  `json:"OpenSellOrders"`
}

type Coins struct {
	Coins []Coin `json:"coins"`
}

type Coin struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}
