package p2p

import (
	"io"
	"fmt"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct {
	Version			string
}

func NewDefaultHandler() Handler {
	return &DefaultHandler{
		Version: "v.0.1-alpha",
	}
}

func (h *DefaultHandler) HandleMessage(msg *Message) error {
	b, err := io.ReadAll(msg.Payload)
	fmt.Printf("handling the msg from %s: %s\n", msg.From, string(b))

	return err
}
