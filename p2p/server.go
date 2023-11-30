package p2p

import (
	"bytes"
	"fmt"
	"log"
	"io"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
}

func(p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)

	return err
}

type ServerConfig struct {
	ListenAddr			string
	Version					string
}

type Message struct {
	Payload					io.Reader	
	From						net.Addr
}

type Server struct {
	ServerConfig		ServerConfig

	handler					Handler
	listener				net.Listener
	mu							sync.RWMutex
	peers						map[net.Addr]*Peer
	addPeer					chan *Peer
	delPeer					chan *Peer
	msgCh						chan *Message 
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		handler: NewDefaultHandler(),
		peers: make(map[net.Addr]*Peer),
		addPeer: make(chan *Peer),
		delPeer: make(chan *Peer),
		msgCh: make(chan *Message),
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

// TODO: redundant code to add new peers
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	peer := &Peer {
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(s.ServerConfig.Version))
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		peer := &Peer {
			conn: conn, 
		}

		s.addPeer <- peer

		peer.Send([]byte(s.ServerConfig.Version + "\n"))

		go s.handleConn(peer)
	}
}

func (s *Server) handleConn(peer *Peer) {
	buf := make([]byte, 1024)

	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			break
		}

		s.msgCh <- &Message{
			From:			peer.conn.RemoteAddr(),
			Payload:	bytes.NewReader(buf[:n]),
		}
		
		fmt.Printf("%s", string(buf[:n]))
	}

	s.delPeer <- peer
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
		case peer := <- s.delPeer:
			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)
			fmt.Printf("player disconnected %s\n", addr)

		case peer := <- s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())
		case msg := <- s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
