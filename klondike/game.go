package klondike

import (
	"crypto/sha256"
	"errors"
	"math/rand"
)

type GameState struct {
	Stock       Stock
	Piles       [7]Pile
	Foundations [4]Foundation
}

func (s GameState) Hash() [32]byte {
	b := s.Stock.bytes()
	for _, p := range s.Piles {
		b = append(b, p.bytes()...)
	}
	for _, f := range s.Foundations {
		b = append(b, f.bytes()...)
	}

	return sha256.Sum256(b)
}

type Game struct {
	draws int
	state GameState
}

type Location int

const (
	LocationStock Location = iota
	LocationPile1
	LocationPile2
	LocationPile3
	LocationPile4
	LocationPile5
	LocationPile6
	LocationPile7
	LocationFoundationHeart
	LocationFoundationSpade
	LocationFoundationDiamond
	LocationFoundationClub
)

func GetPile(i int) Location {
	return LocationPile1 + Location(i)
}

func GetFoundation(cardType CardType) Location {
	return LocationFoundationHeart + Location(int(cardType)-1)
}

func NewGameWithSeed(seed int64, draws int) Game {
	return NewGame(rand.NewSource(seed), draws)
}

func NewGame(randSrc rand.Source, draws int) Game {
	if randSrc == nil {
		panic(errors.New("nil random source"))
	}
	if draws < 1 {
		panic(errors.New("invalid draws count"))
	}

	deck := make([]Card, 52)
	for i := 0; i < 4; i++ {
		for j := 0; j < 13; j++ {
			deck[(i*13)+j] = Card{
				Type:    CardType(i + 1),
				Number:  j + 1,
				Flipped: false,
			}
		}
	}

	// Shuffle the deck
	r := rand.New(randSrc)
	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

	game := Game{draws: draws}

	var d int
	for i := 0; i < 7; i++ {
		game.state.Piles[i].flippedCount = 1
		game.state.Piles[i].cards = make([]Card, i+1)
		copy(game.state.Piles[i].cards, deck[d:d+i+1])
		game.state.Piles[i].cards[0].Flipped = true
		d += i + 1
	}

	game.state.Stock.passesCount = 0
	game.state.Stock.drawn = []Card{}
	game.state.Stock.cards = make([]Card, 24)
	copy(game.state.Stock.cards, deck[d:])

	for i := 0; i < 4; i++ {
		game.state.Foundations[i].cardType = CardType(i + 1)
		game.state.Foundations[i].cards = []Card{}
	}

	return game
}

func (g Game) State() GameState {
	return g.state
}

func (g Game) ApplyMove(m GameMove) (Game, error) {
	var err error

	if m.t == moveTypeDrawFromStock {
		g.state.Stock, err = g.state.Stock.draw(g.draws)
		return g, err
	}

	var cards []Card

	switch m.from {
	case LocationStock:
		g.state.Stock, cards, err = g.state.Stock.pop(m.count)

	case LocationFoundationClub, LocationFoundationSpade, LocationFoundationHeart, LocationFoundationDiamond:
		f := g.state.Foundations[m.from-8]
		g.state.Foundations[m.from-8], cards, err = f.pop(m.count)

	case LocationPile1, LocationPile2, LocationPile3, LocationPile4, LocationPile5, LocationPile6, LocationPile7:
		p := g.state.Piles[m.from-1]
		g.state.Piles[m.from-1], cards, err = p.pop(m.count)
	}

	if err != nil {
		return g, err
	}

	popped := make([]Card, len(cards))
	copy(popped, cards)

	switch m.to {
	case LocationStock:
		g.state.Stock, err = g.state.Stock.place(popped)

	case LocationFoundationClub, LocationFoundationSpade, LocationFoundationHeart, LocationFoundationDiamond:
		f := g.state.Foundations[m.to-8]
		g.state.Foundations[m.to-8], err = f.place(popped)

	case LocationPile1, LocationPile2, LocationPile3, LocationPile4, LocationPile5, LocationPile6, LocationPile7:
		p := g.state.Piles[m.to-1]
		g.state.Piles[m.to-1], err = p.place(popped)
	}

	if err != nil {
		return g, err
	}

	return g, nil
}

// Solved returns true if the game is solved.
// A game is considered solved if all cards in all piles have been flipped,
// and the stock pile is empty.
func (g Game) Solved() bool {
	for i := 0; i < 7; i++ {
		if g.state.Piles[i].FlippedCount() != g.state.Piles[i].Len() {
			return false
		}
	}

	return g.state.Stock.Len() == 0
}
