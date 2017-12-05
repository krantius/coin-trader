package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Server struct {
	Brokers []*Broker `json:"brokers"`
}

type trackerConfig struct {
	Sustain    int     `json:"sustain"`
	DataSize   int     `json:"dataSize"`
	BuyPercent float64 `json:"buyPercent"`
	Name       string  `json:"name"`
}

func NewServer() *Server {
	configs := []trackerConfig{
		{12, 60, 10, "Two minute sustain, 10 minute history, 10% buy percent"},
		{6, 60, 10, "One minute sustain, 10 minute history, 10% buy percent"},
		{12, 60, 15, "Two minute sustain, 10 minute history, 15% buy percent"},
		{3, 60, 20, "30 second sustain, 10 minute history, 20% buy percent"},
		{30, 180, 8, "5 minute sustain, 30 minute history, 8% buy percent"},
	}

	s := &Server{}
	for _, c := range configs {
		s.Brokers = append(s.Brokers, NewFakeBroker(c))
	}

	return s
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

func (s *Server) Run() {
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/status", s.statusHandler)

	for _, b := range s.Brokers {
		go b.Work()
	}

	port := os.Getenv("PORT")
	fmt.Printf("Starting server on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Printf("error yo: %v\n", err)
}
