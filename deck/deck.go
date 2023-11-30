package deck

import (
	"math/rand"
	"fmt"
	"strconv"
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
func NewCard(s Suit, v int) Card {
	if v > 13 {
		panic("the value of the card cannot be higher than 13")
	}

	return Card {
		suit: s,
		value: v,
	}
}

func (c Card) String() string {
	value := strconv.Itoa(c.value)
	if value == "1" {
		value = "ACE"
	}
	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

type Deck [52]Card 

func New() Deck {
	nSuits := 4
	nCards := 13
	d := [52]Card{}
	x := 0

	for i := 0; i < nSuits; i++ {
		for j:= 0; j < nCards; j++ {
			d[x] = NewCard(Suit(i), j + 1)
			x++
		}
	}
	return Shuffle(d)
}

func Shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		r := rand.Intn(i + 1)
		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}
	return d
}

func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("Invalid card suit")
	}
}
