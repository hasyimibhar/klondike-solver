package klondike

import (
	"errors"
)

type CardType int

var (
	ErrInvalidMove = errors.New("invalid move")
)

const (
	Unknown CardType = iota
	Heart
	Spade
	Diamond
	Club
)

type Card struct {
	Type    CardType
	Number  int
	Flipped bool
}

type cardPile []Card

func (c Card) bytes() []byte {
	var flipped byte
	if c.Flipped {
		flipped = 1
	}

	return []byte{byte(c.Type), byte(c.Number), flipped}
}

func (p cardPile) bytes() []byte {
	b := []byte{}
	for _, c := range p {
		b = append(b, c.bytes()...)
	}
	return b
}

func (c Card) stackablePile(cc Card) bool {
	// Must be different color
	if c.Type%2 == cc.Type%2 {
		return false
	}
	// Top card number must be bottom card number-1
	if cc.Number-c.Number != 1 {
		return false
	}
	return true
}

func (c Card) stackableFoundation(cc Card) bool {
	// Must be same type
	if c.Type != cc.Type {
		return false
	}
	// Top card number must be bottom card number+1
	if c.Number-cc.Number != 1 {
		return false
	}
	return true
}
