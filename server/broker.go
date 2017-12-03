package server

import (
	"fmt"
	"sync"
)

const (
	BUY     = 0
	BUYING  = 1
	SELL    = 2
	SELLING = 3
)

type Broker struct {
	buyCh     chan string
	trackerCh chan string
	state     int
	stateM    sync.Mutex
	exchange  Exchange
}

func NewBroker() *Broker {
	return &Broker{
		buyCh:     make(chan string),
		trackerCh: make(chan string),
		state:     BUY,
		exchange:  NewRealExchange(),
	}
}

func NewMockBroker() *Broker {
	return &Broker{
		buyCh:     make(chan string),
		trackerCh: make(chan string),
		state:     BUY,
		exchange:  NewMockExchange(),
	}
}

func (b *Broker) listen() {
	for {
		market := <-b.trackerCh

		fmt.Printf("Received buy order for %s\n", market)
		if b.getState() == BUY {
			b.buyCh <- market
		}
	}
}

func (b *Broker) Work() {
	go b.listen()

	for {
		switch b.state {
		case BUY:
			fmt.Printf("Wait for buy with %f bitcoin\n", b.exchange.GetValue())
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
