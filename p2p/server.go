package p2p

import (
	"fmt"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func (gv GameVariant) String() string {
	switch gv {
		case TexasHoldem:
			return "TEXAS HOLDEM"
		case Other:
			return "OTHER"
		default:
			return "-"
	}
}

type ServerConfig struct {
	ListenAddr			string
	Version					string
	GameVariant			GameVariant
}

type Server struct {
	ServerConfig		ServerConfig

	transport *TCPTransport
	handler					Handler
	listener				net.Listener
	peers						map[net.Addr]*Peer
	addPeer					chan *Peer
	delPeer					chan *Peer
	msgCh						chan *Message 
}

func NewServer(cfg ServerConfig) *Server {
	s :=  &Server{
		ServerConfig: cfg,
		handler: NewDefaultHandler(),
		peers: make(map[net.Addr]*Peer),
		addPeer: make(chan *Peer),
		delPeer: make(chan *Peer),
		msgCh: make(chan *Message),
	}

	tr := NewTcpTransport(s.ServerConfig.ListenAddr)
	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer
	s.transport = tr

	return s
}

func (s *Server) Start() {
	go s.loop()

	fmt.Printf("game server running on port %s\n", s.transport.listenAddr)
	logrus.WithFields(logrus.Fields{
		"port": s.transport.listenAddr,
		"variant": s.ServerConfig.GameVariant,
	}).Info("started new game server")

	s.transport.ListenAndAccept()
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


func (s *Server) loop() {
	for {
		select {
		case peer := <- s.delPeer:
			logrus.WithFields(logrus.Fields {
				"addr": peer.conn.RemoteAddr(),
			}).Info("player disconnected")
			
			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)
			fmt.Printf("player disconnected %s\n", addr)

		case peer := <- s.addPeer:
			// TODO: check max players and other game state logic
			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields {
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")
			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <- s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
