package server

import (
	"fmt"
	"net"
	"net/http"
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

	listener				net.Listener
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
	if err := s.listen(); err != nil {
		panic(err)
	}

	conn, err := s.listener.Accept()
	if err != nil {
		panic(err)
	}

	go s.handleConn(conn)
}

func (s *Server) handleConn(conn net.Conn) {

}

func(s* Server) listen() error {
	lis, err := net.Listen("tcp",s.ServerConfig.listenAddr)
	if err != nil {
		return err
	}

	s.listener = lis

	return nil
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
