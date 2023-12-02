package p2p

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func init() {
	gob.Register(GameVariant(0))
	gob.Register(GameVariant(1))
	gob.Register(Handshake{})
}

func (gv GameVariant) String() string {
	switch gv {
		case TexasHoldem:
			return "TEXAS HOLDEM"
		case Other:
			return "OTHER"
		default:
			return "UNKNOWN"
	}
}

type ServerConfig struct {
	ListenAddr			string
	Version					string	
	GameVariant			GameVariant
}

type Server struct {
	ServerConfig

	transport *TCPTransport
	listener				net.Listener
	peers						map[net.Addr]*Peer
	addPeer					chan *Peer
	delPeer					chan *Peer
	msgCh						chan *Message 

	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {
	s :=  &Server{
		ServerConfig: cfg,
		peers: make(map[net.Addr]*Peer),
		addPeer: make(chan *Peer),
		delPeer: make(chan *Peer),
		msgCh: make(chan *Message),
		gameState: NewGameState(),
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

func(s *Server) SendHandshake(peer *Peer) error {
	hs := &Handshake {
		GameVariant: s.ServerConfig.GameVariant,
		Version: s.ServerConfig.Version,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(&hs); err != nil {
		return err
	}
	// if err := hs.Encode(buf); err != nil {
	// 	return err
	// }

	return peer.Send(buf.Bytes())
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

			// if a new peer connects to the server we send our handhskake message and wait for his reply.
		case peer := <- s.addPeer:
			s.SendHandshake(peer)

			if err := s.handshake(peer); err != nil {
				logrus.Errorf("handshake with incoming player failed: %s", err.Error())
				continue
			}
			// TODO: check max players and other game state logic
			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields {
				"addr": peer.conn.RemoteAddr(),
				}).Info("handshake successful: new player connected")
			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <- s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

type Handshake struct {
	Version							string
	GameVariant					GameVariant
}

func(hs *Handshake) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &hs.Version); err != nil {
		return err
	}
	fmt.Println("Decoding:",hs.Version)

	if err := binary.Read(r, binary.LittleEndian, &hs.GameVariant); err != nil {
		return err
	}
	
	return nil
}

func(hs *Handshake) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, []byte(hs.Version)); err != nil {
		return err
	}
	if err :=  binary.Write(w, binary.LittleEndian, uint8(hs.GameVariant)); err != nil {
		return err
	}
	return nil
}

func(s *Server) handshake(peer *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(peer.conn).Decode(hs); err != nil {
		return fmt.Errorf("Can't connect: %s", err)
	}
	// if err := hs.Decode(peer.conn); err != nil {
	// 	return err
	// }

	if s.ServerConfig.GameVariant != hs.GameVariant {
		return fmt.Errorf("Invalid gamevariant %s", hs.GameVariant)
	}

	if string(s.ServerConfig.Version) != string(hs.Version) {
		return fmt.Errorf("invalid version %s", hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer": peer.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Info("received handshake")
	return nil
}

func(s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}

