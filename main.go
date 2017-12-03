package main

import "github.com/mkrant/coin-trader/server"

func main() {
	s := server.NewServer()
	s.Run()
}
