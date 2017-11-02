package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func tickerHandler(w http.ResponseWriter, r *http.Request) {

	c := Coins{}
	for _, currency := range r.URL.Query()["currency"] {
		t, _ := GetTicker(currency)
		c.Coins = append(c.Coins, Coin{currency, t.Last})
	}

	b, _ := json.Marshal(c)
	w.Write(b)
}

func marketHandler(w http.ResponseWriter, r *http.Request) {
	markets, _ := GetMarkets()

	for _, m := range markets {
		j, _ := json.Marshal(m)
		w.Write(j)
	}
}

func Run() {
	http.HandleFunc("/ticker", tickerHandler)
	http.HandleFunc("/markets", marketHandler)

	fmt.Println("Starting server on port 8080")

	c := make(chan string, 1)
	markets, _ := GetMarkets()

	for _, m := range markets {
		if !strings.HasPrefix(m.MarketName, "BTC") {
			continue
		}

		t := NewTracker(m.MarketName, c)
		go t.Start()
	}

	http.ListenAndServe(":8080", nil)
}
