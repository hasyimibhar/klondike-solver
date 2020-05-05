package klondike

import (
	"errors"
)

type Stock struct {
	passesCount int
	cards       []Card
	drawn       []Card
}

func (s Stock) bytes() []byte {
	b := cardPile(s.cards).bytes()
	b = append(b, cardPile(s.drawn).bytes()...)
	return append([]byte{byte(s.passesCount)}, b...)
}

func (s Stock) Drawn() []Card {
	return s.drawn
}

func (s Stock) Len() int {
	return len(s.cards)
}

func (s Stock) Card(idx int) Card {
	if idx < 0 || idx >= len(s.drawn) {
		panic(errors.New("invalid index"))
	}

	return s.drawn[idx]
}

func (s Stock) draw(n int) (Stock, error) {
	// Recycle
	if len(s.cards) == 0 {
		// Put the drawn cards back to the stock pile and reverse the order
		s.cards = s.drawn
		for i := len(s.cards)/2 - 1; i >= 0; i-- {
			j := len(s.cards) - 1 - i
			s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
		}

		s.drawn = []Card{}
		s.passesCount++
		return s, nil
	}

	actualN := n
	if actualN > len(s.cards) {
		actualN = len(s.cards)
	}

	drawn := make([]Card, actualN)
	copy(drawn, s.cards[:actualN])

	// Flip drawn cards
	for i := range drawn {
		drawn[i].Flipped = true
	}

	s.drawn = append(drawn, s.drawn...)
	s.cards = s.cards[n:]

	return s, nil
}

func (s Stock) pop(n int) (Stock, []Card, error) {
	if len(s.drawn) == 0 {
		panic(errors.New("invalid pop"))
	}
	if n != 1 {
		panic(errors.New("can only pop 1 card from stock"))
	}

	popped := s.drawn[:n]
	s.drawn = s.drawn[n:]
	return s, popped, nil
}

func (s Stock) place(cards []Card) (Stock, error) {
	return Stock{}, ErrInvalidMove
}
