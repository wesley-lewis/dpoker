package main

import (
	"fmt"

	"github.com/wesley-lewis/dpoker/deck"
)

func main() {
	card := deck.NewCard(deck.Spades, 1)

	fmt.Println(card)
}
