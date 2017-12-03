package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Bittrex                     = "https://bittrex.com/api/v1.1"
	TargetGainPercent           = 1.1
	StopLossPercent             = .9
	ORDER_BUY         OrderType = "buy"
	ORDER_SELL        OrderType = "sell"
)

type Exchange interface {
	Buy(market string) error
	Sell() error
	GetValue() float64
	GetOrders() []Order
}

type RealExchange struct {
	url string
}

type MockExchange struct {
	Value    float64 `json:"Value"`
	Currency string  `json:"Currency"`
	BuyPrice float64 `json:"BuyPrice"`
	Orders   []Order `json:"Orders"`
}

func NewRealExchange() *RealExchange {
	return &RealExchange{
		url: Bittrex,
	}
}

func NewMockExchange() *MockExchange {
	return &MockExchange{
		Value:    1.0,
		Currency: "btc",
		Orders:   []Order{},
	}
}

func (*RealExchange) Buy(market string) error {
	return nil
}

func (e *MockExchange) Buy(market string) error {
	ticker, err := GetTicker(market)

	if err != nil {
		return fmt.Errorf("Error buying Market %s %v", market, err)
	}

	if e.Currency != "btc" {
		return fmt.Errorf("expecting bitcoin, got %s", e.Currency)
	}
	e.Currency = market
	btc := e.Value
	e.Value = btc / ticker.Ask
	e.BuyPrice = ticker.Ask

	order := &Order{
		Type:     ORDER_BUY,
		Currency: e.Currency,
		Rate:     e.BuyPrice,
		Units:    e.Value,
	}

	e.Orders = append(e.Orders, *order)
	fmt.Printf("Buy: %v\n", order)

	return nil
}

func (*RealExchange) Sell() error {
	return nil
}

func (e *MockExchange) Sell() error {
	sellPrice := WaitForTargetPrice(e.Currency, e.BuyPrice)

	shitCoin := e.Value
	e.Value = sellPrice * shitCoin

	order := &Order{
		Type:     ORDER_SELL,
		Currency: e.Currency,
		Rate:     sellPrice,
		Units:    e.Value,
	}

	e.Orders = append(e.Orders, *order)
	fmt.Printf("Sold: %v\n", order)
	e.Currency = "btc"

	return nil
}

func (e *RealExchange) GetValue() float64 {
	return 0.0
}

func (e *MockExchange) GetValue() float64 {
	return e.Value
}

func (e *RealExchange) GetOrders() []Order {
	return nil
}

func (e *MockExchange) GetOrders() []Order {
	return e.Orders
}

func GetTicker(currency string) (*Ticker, error) {
	resp := get(Bittrex + "/public/getticker?Market=" + currency)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error during parsing json", err)
	}

	ticker := TickerResponse{}
	json.Unmarshal([]byte(body), &ticker)

	if !ticker.Success {
		fmt.Printf("Ticker failed: %v", ticker)
	}
	return &ticker.Ticker, nil
}

func GetMarkets() ([]Market, error) {
	resp := get(Bittrex + "/public/getmarkets")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error during parsing json", err)
	}

	mr := MarketResponse{}
	json.Unmarshal([]byte(body), &mr)

	return mr.Markets, nil
}

// Blocking function to wait for a target price
func WaitForTargetPrice(m string, bp float64) float64 {
	target := bp * TargetGainPercent
	stopLoss := bp * StopLossPercent

	var sellPrice float64
	t := time.NewTicker(15 * time.Second)

	fmt.Printf("Waiting for target price %f for Market %s\n", target, m)
	for range t.C {
		ticker, _ := GetTicker(m)

		if ticker.Last >= target || ticker.Last <= stopLoss {
			t.Stop()
			sellPrice = ticker.Last
		}
	}

	if sellPrice < target {
		fmt.Printf("Selling shitcoin %s for loss at Value %f\n", m, sellPrice)
	} else {
		fmt.Printf("Selling shitcoin %s for gain at Value %f\n", m, sellPrice)
	}
	return sellPrice
}

func get(url string) *http.Response {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Errorf("Error getting http request", err)
	}

	return resp
}
