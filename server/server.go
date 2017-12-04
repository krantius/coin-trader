package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Server struct {
	b        *Broker
	Trackers []*Tracker `json:"trackers"`
}

func NewServer() *Server {
	return &Server{
		b: NewFakeBroker(),
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
	j, err := json.Marshal(s.b)
	if err != nil {
		fmt.Printf("%v", err)
	}

	w.Write(j)
}

func (s *Server) topHandler(w http.ResponseWriter, r *http.Request) {
	sort.Slice(s.Trackers, func(i, j int) bool {
		return s.Trackers[i].PercentChange < s.Trackers[j].PercentChange
	})

	j, err := json.Marshal(s.Trackers)
	if err != nil {
		fmt.Printf("%v", err)
	}

	w.Write(j)
}

func (s *Server) Run() {
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/status", s.statusHandler)
	http.HandleFunc("/orders", s.orderHandler)
	http.HandleFunc("/top", s.topHandler)

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

	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Printf("error yo: %v\n", err)
}
