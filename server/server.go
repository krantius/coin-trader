package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Server struct {
	b        *Broker
	Trackers []*Tracker `json:"trackers"`
}

func NewServer() *Server {
	return &Server{
		b: NewMockBroker(),
	}
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	m, _ := json.Marshal("hello")
	w.Write(m)
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("%v", err)
	}

	w.Write(j)
}

func (s *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(s.b.Exchange())
	if err != nil {
		fmt.Printf("%v", err)
	}

	w.Write(j)
}

func (s *Server) Run() {
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/status", s.statusHandler)
	http.HandleFunc("/exchange", s.orderHandler)

	go s.b.Work()

	markets, _ := GetMarkets()
	for _, m := range markets {
		if !strings.HasPrefix(m.MarketName, "BTC") {
			continue
		}
		t := NewTracker(m.MarketName, s.b.trackerCh)
		s.Trackers = append(s.Trackers, t)
		go t.Start(context.Background())
	}

	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Printf("error yo: %v\n", err)
}
