package deck 

import (
	"fmt"
)
type Suit int

func(s Suit) String() string{
	switch s {
	case Spades:
		return "SPADES"
	case Hearts:
		return "HEARTS"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	default:
		panic("Invalid card suit")
	}
}

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// Card represents a single card in pack of 52 playing cards.
// Value will not be higher than 13
type Card struct {
	suit Suit 
	value int
}

func (c Card) String() string {
	return fmt.Sprintf("%d of %s", c.value, c.suit)
}

func NewCard(s Suit, v int) Card {
	if v > 12 {
		panic("the value of the card cannot be higher than 13")
	}

	return Card {
		suit: s,
		value: v,
	}
}
