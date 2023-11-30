package p2p 

import (
	"fmt"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

type ServerConfig struct {
	ListenAddr string
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
	
	fmt.Printf("starting the server\n")
	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("%s", string(buf[:n]))
	}
}

func(s* Server) listen() error {
	lis, err := net.Listen("tcp",s.ServerConfig.ListenAddr)
	if err != nil {
		return err
	}

	s.listener = lis
	fmt.Println("Server listening on port:", s.ServerConfig.ListenAddr)
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

// 54:59
