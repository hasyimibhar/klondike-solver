package klondike

import (
	"errors"
)

type Foundation struct {
	cardType CardType
	cards    []Card
}

func (f Foundation) bytes() []byte {
	return cardPile(f.cards).bytes()
}

func (f Foundation) Len() int {
	return len(f.cards)
}

func (f Foundation) CardType() CardType {
	return f.cardType
}

func (f Foundation) Card() Card {
	if len(f.cards) == 0 {
		return Card{Type: CardTypeUnknown}
	}

	return f.cards[0]
}

func (f Foundation) pop(n int) (Foundation, []Card, error) {
	if len(f.cards) == 0 {
		panic(errors.New("cannot pop empty foundation"))
	}
	if n != 1 {
		panic(errors.New("can only pop 1 card from foundation"))
	}

	popped := f.cards[:1]
	f.cards = f.cards[1:]

	return f, popped, nil
}

func (f Foundation) place(cards []Card) (Foundation, error) {
	if len(cards) != 1 {
		panic(errors.New("can only place 1 card on foundation"))
	}

	c := cards[0]
	if len(f.cards) > 0 && !c.stackableFoundation(f.cards[0]) {
		return Foundation{}, ErrInvalidMove
	} else if len(f.cards) == 0 && cards[0].Type != f.cardType {
		return Foundation{}, ErrInvalidMove
	}

	f.cards = append([]Card{c}, f.cards...)
	return f, nil
}
