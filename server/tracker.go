package server

import (
	"context"
	"fmt"
	"math"
	"time"
)

const (
//dataSize = 60
)

type Tracker struct {
	Market        string `json:"Market"`
	arr           []float64
	current       int
	base          int
	PercentChange float64 `json:"PercentChange"`
	wrap          bool
	buyChannel    chan string
	hasAlerted    bool
	sustain       int
	buyPercent    float64
	ds            int
}

func NewTracker(market string, c chan string, config trackerConfig) *Tracker {
	t := &Tracker{
		Market:        market,
		PercentChange: 0.0,
		buyChannel:    c,
		sustain:       config.Sustain,
		buyPercent:    config.BuyPercent,
		ds:            config.DataSize,
		arr:           make([]float64, config.DataSize),
	}
	return t
}

func (t *Tracker) Start(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
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

	if t.PercentChange > 10.0 {
		t.sustain++
		if !t.hasAlerted && t.sustain >= 12 {
			t.buyChannel <- t.Market
			t.hasAlerted = true
		}
		return

	}

	if t.PercentChange < -15.0 {
		fmt.Printf("%s DUMPING: %f from %f to %f\n", t.Market, t.PercentChange, base, current)
		t.hasAlerted = false
		return
	}

	t.sustain = 0
}

func (t *Tracker) getBase() float64 {
	return t.arr[t.base]
}

func (t *Tracker) getCurrent() float64 {
	return t.arr[t.current]
}

func (t *Tracker) update(last float64) {
	t.arr[t.current] = last
	t.calculateChange()
	t.current++

	if t.current == t.ds {
		t.current = 0
		t.wrap = true
	}

	if !t.wrap {
		return
	}

	t.base++
	if t.base == t.ds {
		t.base = 0
	}
}
