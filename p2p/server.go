package server

import (
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

type ServerConfig struct {
	listenAddr string
}

type Server struct {
	ServerConfig		ServerConfig

	mu							sync.RWMutex
	peers						map[net.Addr]*Peer
	addPeer					chan *Peer
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers: make(map[net.Addr]*Peer),
		addPeer: make(chan *Peer),
	}
}

func (s *Server) Start() {
	go s.loop()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <- s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected %s", peer.conn.RemoteAddr())
		}
	}
}
