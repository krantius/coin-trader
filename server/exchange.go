package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Bittrex                     = "https://bittrex.com/api/v1.1"
	TargetGainPercent           = 1.1
	StopLossPercent             = .9
	ORDER_BUY         OrderType = "buy"
	ORDER_SELL        OrderType = "sell"
)

type Exchange interface {
	Buy(market string) (Order, error)
	Sell(market string, units float64) (Order, error)
	//GetValue() float64
	GetOrders() []Order
}

type RealExchange struct {
	url string
}

type FakeExchange struct {
	Orders     []Order          `json:"Orders"`
	RecentBuys map[string]Order `json:"RecentBuys"`
}

func NewRealExchange() *RealExchange {
	return &RealExchange{
		url: Bittrex,
	}
}

func NewFakeExchange() *FakeExchange {
	return &FakeExchange{
		Orders: []Order{},
	}
}

func (*RealExchange) Buy(market string) (Order, error) {
	return Order{}, nil
}

// Each buy starts with 1 btc
func (e *FakeExchange) Buy(market string) (Order, error) {
	ticker, err := GetTicker(market)

	if err != nil {
		return Order{}, fmt.Errorf("error buying Market %s %v", market, err)
	}

	units := 1 / ticker.Last

	order := &Order{
		Type:     ORDER_BUY,
		Currency: market,
		Rate:     ticker.Last,
		Units:    units,
		Bitcoin:  1.0,
	}

	e.Orders = append(e.Orders, *order)
	fmt.Printf("Buy: %v\n", order)

	return *order, nil
}

func (*RealExchange) Sell(market string, units float64) (Order, error) {
	return Order{}, nil
}

func (e *FakeExchange) Sell(market string, units float64) (Order, error) {
	last, _ := GetTicker(market)

	btc := last.Last * units

	order := &Order{
		Type:     ORDER_SELL,
		Currency: market,
		Rate:     last.Last,
		Units:    units,
		Bitcoin:  btc,
	}

	e.Orders = append(e.Orders, *order)
	fmt.Printf("Sold: %v\n", order)

	return *order, nil
}

func (e *RealExchange) GetValue() float64 {
	return 0.0
}

func (e *RealExchange) GetOrders() []Order {
	return nil
}

func (e *FakeExchange) GetOrders() []Order {
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

func get(url string) *http.Response {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Errorf("Error getting http request", err)
	}

	return resp
}
