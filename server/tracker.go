package server

import (
	"fmt"
	"time"
)

const (
	dataSize = 120
)

type Tracker struct {
	market        string
	arr           [dataSize]float64
	last          int
	base          int
	percentChange float64
	wrap          bool
	buyChannel    chan string
	hasAlerted    bool
}

func NewTracker(market string, c chan string) *Tracker {
	t := &Tracker{market: market, buyChannel: c}
	return t
}

func (t *Tracker) Start() {
	fmt.Printf("Tracking market %s\n", t.market)
	tick := time.NewTicker(15 * time.Second)
	for range tick.C {
		v, _ := GetTicker(t.market)
		t.update(v.Last)
	}
}

func (t *Tracker) calculateChange() {
	base := t.getBase()
	current := t.getCurrent()
	t.percentChange = ((current - base) / base) * 100

	if t.percentChange > 15.0 {
		if !t.hasAlerted {
			t.buyChannel <- t.market
			t.hasAlerted = true
		}
	} else if t.percentChange < -10.0 {
		fmt.Printf("%s DUMPING: %f from %f to %f\n", t.market, t.percentChange, base, current)
		t.hasAlerted = false
	}
}

func (t *Tracker) getBase() float64 {
	return t.arr[t.base]
}

func (t *Tracker) getCurrent() float64 {
	return t.arr[t.last]
}

func (t *Tracker) update(last float64) {
	t.arr[t.last] = last
	t.calculateChange()
	t.last++

	if t.last == dataSize {
		t.last = 0
		t.wrap = true
	}

	if !t.wrap {
		return
	}

	t.base++
	if t.base == dataSize {
		t.base = 0
	}
}

/*
func (t *Tracker) update(last float64) {
	t.arr[t.last] = last
	t.calculateChange()
	t.last++

	if t.last == dataSize {
		t.last = 0
		t.wrap = true
	}

	if !t.wrap {
		return
	}

	t.base++
	if t.base == dataSize {
		t.base = 0
	}
}
*/