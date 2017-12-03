package server

import (
	"context"
	"fmt"
	"math"
	"time"
)

const (
	dataSize = 120
)

type Tracker struct {
	Market        string            `json:"Market"`
	Arr           [dataSize]float64 `json:"Data"`
	current       int
	base          int
	PercentChange float64 `json:"PercentChange"`
	wrap          bool
	buyChannel    chan string
	hasAlerted    bool
}

func NewTracker(market string, c chan string) *Tracker {
	t := &Tracker{
		Market:        market,
		PercentChange: 0.0,
		buyChannel:    c,
	}
	return t
}

func (t *Tracker) Start(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			v, _ := GetTicker(t.Market)
			t.update(v.Last)
		case <-ctx.Done():
			return
		}
	}
}

func (t *Tracker) calculateChange() {
	base := t.getBase()
	current := t.getCurrent()
	t.PercentChange = ((current - base) / base) * 100
	if math.IsNaN(t.PercentChange) {
		t.PercentChange = 0
		return
	}

	if t.PercentChange > 15.0 {
		if !t.hasAlerted {
			t.buyChannel <- t.Market
			t.hasAlerted = true
		}
	} else if t.PercentChange < -10.0 {
		fmt.Printf("%s DUMPING: %f from %f to %f\n", t.Market, t.PercentChange, base, current)
		t.hasAlerted = false
	}
}

func (t *Tracker) getBase() float64 {
	return t.Arr[t.base]
}

func (t *Tracker) getCurrent() float64 {
	return t.Arr[t.current]
}

func (t *Tracker) update(last float64) {
	t.Arr[t.current] = last
	t.calculateChange()
	t.current++

	if t.current == dataSize {
		t.current = 0
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
