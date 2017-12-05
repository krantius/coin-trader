package server

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	BUY     = 0
	BUYING  = 1
	SELL    = 2
	SELLING = 3
)

type Broker struct {
	buyCh        chan string
	trackerCh    chan string
	state        int
	stateM       sync.Mutex
	exchange     Exchange
	multipleBuys bool
	Success      int     `json:"success"`
	Fail         int     `json:"fail"`
	Orders       []Order `json:"orders"`
	trackers     []*Tracker
	Config       trackerConfig `json:"trackerConfig"`
}

func NewBroker() *Broker {
	return &Broker{
		buyCh:     make(chan string),
		trackerCh: make(chan string),
		state:     BUY,
		exchange:  NewRealExchange(),
	}
}

func NewFakeBroker(config trackerConfig) *Broker {
	return &Broker{
		buyCh:     make(chan string, 25),
		trackerCh: make(chan string),
		state:     BUY,
		exchange:  NewFakeExchange(),
		Orders:    []Order{},
		Config:    config,
	}
}

func (b *Broker) listen() {
	for market := range b.trackerCh {
		fmt.Printf("Received buy order for %s\n", market)

		go b.HandleTrade(market)
	}
}

func (b *Broker) Work() {
	go b.listen()

	markets, _ := GetMarkets()
	for _, m := range markets {
		if !strings.HasPrefix(m.MarketName, "BTC") {
			continue
		}
		t := NewTracker(m.MarketName, b.trackerCh, b.Config)
		b.trackers = append(b.trackers, t)
		go t.Start(context.Background())
	}
	/*
		for {
			switch b.state {
			case BUY:
				fmt.Printf("Wait for buy with %v bitcoin\n", b.exchange.GetValue())
				market := <-b.buyCh
				fmt.Printf("Acting on buy order with Market: %s\n", market)
				if err := b.exchange.Buy(market); err != nil {
					fmt.Printf("Got error when buying: %v", err.Error())

					continue
				}

				b.state = BUYING

			case BUYING:
				// wait for order to fullfil etc
				fmt.Println("Buying")
				b.state = SELL
			case SELL:
				fmt.Println("Sell")
				b.exchange.Sell()
				b.state = SELLING
				// sell stuff
			case SELLING:
				// wait for sell order to fullfil
				fmt.Println("Selling")
				b.state = BUY
			}
		}*/
}

func (b *Broker) HandleTrade(currency string) {
	bo, err := b.exchange.Buy(currency)
	if err != nil {
		fmt.Printf("Got error when buying: %v", err.Error())
		return
	}

	b.Orders = append(b.Orders, bo)

	target := bo.Rate * TargetGainPercent
	stopLoss := bo.Rate * StopLossPercent

	WaitForTargetPrice(currency, target, stopLoss)
	so, err := b.exchange.Sell(bo.Currency, bo.Units)
	if err != nil {
		fmt.Printf("got error during sell: %v", err)
	}

	b.Orders = append(b.Orders, so)

	if bo.Rate < so.Rate {
		b.Success++
	} else {
		b.Fail++
	}
}

// Blocking function to wait for a target price
func WaitForTargetPrice(m string, target, stopLoss float64) float64 {

	var sellPrice float64
	t := time.NewTicker(15 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Printf("Waiting for target price %v for Market %s, or stop loss at %v\n", target, m, stopLoss)
	for {
		select {
		case <-t.C:
			ticker, _ := GetTicker(m)

			if ticker.Last >= target {
				newTarget := target * TargetGainPercent
				stopLoss = target * .97
				fmt.Printf("Hit target %v let's try upping the target to %v and stoploss to %v\n", target, newTarget, stopLoss)
				target = newTarget
			}

			if ticker.Last <= stopLoss {
				cancel()
				t.Stop()
				sellPrice = ticker.Last
				fmt.Printf("Hit sell price of %v for coin %s\n", sellPrice, m)
			}

		case <-ctx.Done():
			if sellPrice < target {
				fmt.Printf("Selling shitcoin %s for loss at Value %v\n", m, sellPrice)
			} else {
				fmt.Printf("Selling shitcoin %s for gain at Value %v\n", m, sellPrice)
			}
			return sellPrice
		}
	}

}

func (b *Broker) GetOrders() []Order {
	return b.exchange.GetOrders()
}

func (b *Broker) getState() int {
	b.stateM.Lock()
	defer b.stateM.Unlock()
	return b.state
}

func (b *Broker) Exchange() Exchange {
	return b.exchange
}
