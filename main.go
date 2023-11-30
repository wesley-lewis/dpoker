package main

import (

	"github.com/wesley-lewis/dpoker/p2p"
)

func main() {
	cfg := p2p.ServerConfig {
		ListenAddr: ":3000",
	}
	s := p2p.NewServer(cfg)
	s.Start()
	// for j := 0; j < 10; j++ {
	// 	d := deck.New()
	// 	fmt.Println(d)
	// }
}
// EP1: 35:14
