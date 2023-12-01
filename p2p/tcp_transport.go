package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type Peer struct {
	conn net.Conn
}

type Message struct {
	Payload					io.Reader	
	From						net.Addr
}

func(p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)

	return err
}

func(p *Peer) ReadLoop(msgCh chan *Message) {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		msgCh <- &Message{
			From:			p.conn.RemoteAddr(),
			Payload:	bytes.NewReader(buf[:n]),
		}

	}

	// TODO: unregister the peer
	p.conn.Close()
}

type TCPTransport struct {
	version						string
	listenAddr				string 
	listener					net.Listener
	AddPeer						chan *Peer
	DelPeer						chan *Peer
}

func NewTcpTransport(addr string) *TCPTransport {
	return &TCPTransport{
		version: "DPOKER v0.1-alpha",
		listenAddr: addr,
	}
}

func(t *TCPTransport) ListenAndAccept() error {
	lis, err := net.Listen("tcp",t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = lis

	for {
		conn, err := lis.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer {
			conn: conn,
		}

		t.AddPeer <- peer

		// peer.ReadLoop()
	}
	
	return fmt.Errorf("TCP Transport stopped TODO: ?")
}
