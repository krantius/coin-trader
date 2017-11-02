package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Bittrex = "https://bittrex.com/api/v1.1"
)

func GetTicker(currency string) (*Ticker, error) {
	resp := get(Bittrex + "/public/getticker?market=" + currency)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error during parsing json", err)
	}

	ticker := TickerResponse{}
	json.Unmarshal([]byte(body), &ticker)

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

func Memory() *string {
	i := "sdf"
	return &i
}
