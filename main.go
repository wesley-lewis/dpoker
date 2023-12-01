package main

import (
	"time"

	"github.com/wesley-lewis/dpoker/p2p"
)

func main() {
	cfg := p2p.ServerConfig {
		ListenAddr: ":3000",
		Version: "DPOKER V0.1-alpha",
	}
	server := p2p.NewServer(cfg)
	go server.Start()

	remoteCfg := p2p.ServerConfig {
		ListenAddr: ":4000",
		Version: "DPOKER V0.1-alpha",
		GameVariant: p2p.TexasHoldem,
	}

	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()

	time.Sleep(time.Second * 1)
	if err := remoteServer.Connect(server.ServerConfig.ListenAddr); err != nil {
		panic(err)
	}

	select {
	}
	// for j := 0; j < 10; j++ {
	// 	d := deck.New()
	// 	fmt.Println(d)
	// }
}
