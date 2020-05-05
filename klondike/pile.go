package klondike

import (
	"errors"
)

type Pile struct {
	cards        []Card
	flippedCount int
}

func (p Pile) bytes() []byte {
	return cardPile(p.cards).bytes()
}

func (p Pile) Len() int {
	return len(p.cards)
}

func (p Pile) FlippedCount() int {
	return p.flippedCount
}

func (p Pile) Card(idx int) Card {
	if idx < 0 || idx >= len(p.cards) {
		panic(errors.New("invalid index"))
	}

	c := p.cards[idx]

	// Hide value of unflipped card
	if !c.Flipped {
		return Card{Flipped: false}
	}

	return c
}

func (p Pile) pop(n int) (Pile, []Card, error) {
	if n > len(p.cards) {
		panic(errors.New("invalid pop"))
	}

	popped := p.cards[:n]
	if !popped[n-1].Flipped {
		return Pile{}, []Card{}, ErrInvalidMove
	}

	p.cards = p.cards[n:]
	if len(p.cards) > 0 && !p.cards[0].Flipped {
		p.cards[0].Flipped = true
	}

	// Update flippedCount
	p.flippedCount = 0
	for ; p.flippedCount < len(p.cards); p.flippedCount++ {
		if !p.cards[p.flippedCount].Flipped {
			break
		}
	}

	return p, popped, nil
}

func (p Pile) place(cards []Card) (Pile, error) {
	if len(cards) == 0 {
		panic(errors.New("cards cannot be empty"))
	}

	bottom := p.cards[0]
	top := cards[len(cards)-1]
	if !top.stackablePile(bottom) {
		return Pile{}, ErrInvalidMove
	}

	p.cards = append(cards, p.cards...)
	p.flippedCount += len(cards)

	return p, nil
}
