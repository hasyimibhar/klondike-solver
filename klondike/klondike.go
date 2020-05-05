package klondike

import (
	"errors"
	"math/rand"
)

type CardType int

var (
	ErrInvalidMove = errors.New("invalid move")
)

const (
	CardTypeUnknown CardType = iota
	CardTypeHeart
	CardTypeSpade
	CardTypeDiamond
	CardTypeClub
)

type Card struct {
	Type    CardType
	Number  int
	Flipped bool
}

func (c Card) stackablePile(cc Card) bool {
	// Must be different color
	if c.Type%2 != cc.Type%2 {
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

type Pile struct {
	cards        []Card
	flippedCount int
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
	if !p.cards[0].Flipped {
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

type Foundation struct {
	cards []Card
}

func (f Foundation) Len() int {
	return len(f.cards)
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
	if n != 0 {
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
	}

	f.cards = append([]Card{c}, f.cards...)
	return f, nil
}

type Stock struct {
	passesCount  int
	currentIndex int
	cards        []Card
}

func (s Stock) Len() int {
	return len(s.cards)
}

func (s Stock) draw(n int) (Stock, error) {
	if s.currentIndex == -1 {
		s.currentIndex = 0
		return s, nil
	}

	if s.currentIndex+n == len(s.cards) {
		s.currentIndex = -1
		s.passesCount++
		return s, nil
	}

	if s.currentIndex+n > len(s.cards) {
		n = len(s.cards) - s.currentIndex - 1
	}

	s.currentIndex += n
	return s, nil
}

func (s Stock) pop(n int) (Stock, []Card, error) {
	if len(s.cards) == 0 {
		panic(errors.New("invalid pop"))
	}
	if s.currentIndex == -1 {
		panic(errors.New("cannot pop from undrawn stock"))
	}
	if n != 0 {
		panic(errors.New("can only pop 1 card from stock"))
	}

	popped := s.cards[s.currentIndex]
	s.cards = append(s.cards[:s.currentIndex], s.cards[s.currentIndex+1:]...)
	return s, []Card{popped}, nil
}

func (s Stock) place(cards []Card) (Stock, error) {
	return Stock{}, ErrInvalidMove
}

type GameState struct {
	Stock       Stock
	Piles       [7]Pile
	Foundations [4]Foundation
}

type Game struct {
	draws int
	state GameState
}

type moveType int

const (
	moveTypeDrawFromStock moveType = 1
	moveTypeMoveCard
)

type TableauLocation int

const (
	TableauLocationStock TableauLocation = iota
	TableauLocationPile1
	TableauLocationPile2
	TableauLocationPile3
	TableauLocationPile4
	TableauLocationPile5
	TableauLocationPile6
	TableauLocationPile7
	TableauLocationFoundationHeart
	TableauLocationFoundationSpade
	TableauLocationFoundationDiamond
	TableauLocationFoundationClub
)

type move struct {
	Type      moveType
	From      TableauLocation
	CardCount int
	To        TableauLocation
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
			deck[(i*4)+j] = Card{
				Type:    CardType(i + 1),
				Number:  j,
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
		game.state.Piles[i].cards = deck[d : d+i+1]
		game.state.Piles[i].cards[0].Flipped = true
		d += i + 1
	}

	game.state.Stock.passesCount = 0
	game.state.Stock.currentIndex = -1
	game.state.Stock.cards = deck[d:]

	for i := 0; i < 4; i++ {
		game.state.Foundations[i].cards = []Card{}
	}

	return game
}

func (g Game) State() GameState {
	return g.state
}

func (g Game) move(m move) (Game, error) {
	var err error

	if m.Type == moveTypeDrawFromStock {
		g.state.Stock, err = g.state.Stock.draw(g.draws)
		return g, err
	}

	var cards []Card

	switch m.From {
	case TableauLocationStock:
		g.state.Stock, cards, err = g.state.Stock.pop(m.CardCount)

	case TableauLocationFoundationClub:
	case TableauLocationFoundationSpade:
	case TableauLocationFoundationHeart:
	case TableauLocationFoundationDiamond:
		f := g.state.Foundations[m.From-8]
		g.state.Foundations[m.From-8], cards, err = f.pop(m.CardCount)

	case TableauLocationPile1:
	case TableauLocationPile2:
	case TableauLocationPile3:
	case TableauLocationPile4:
	case TableauLocationPile5:
	case TableauLocationPile6:
	case TableauLocationPile7:
		p := g.state.Piles[m.From-1]
		g.state.Piles[m.From-1], cards, err = p.pop(m.CardCount)
	}

	if err != nil {
		return g, err
	}

	switch m.To {
	case TableauLocationStock:
		g.state.Stock, err = g.state.Stock.place(cards)

	case TableauLocationFoundationClub:
	case TableauLocationFoundationSpade:
	case TableauLocationFoundationHeart:
	case TableauLocationFoundationDiamond:
		f := g.state.Foundations[m.To-8]
		g.state.Foundations[m.To-8], err = f.place(cards)

	case TableauLocationPile1:
	case TableauLocationPile2:
	case TableauLocationPile3:
	case TableauLocationPile4:
	case TableauLocationPile5:
	case TableauLocationPile6:
	case TableauLocationPile7:
		p := g.state.Piles[m.To-1]
		g.state.Piles[m.To-1], err = p.place(cards)
	}

	if err != nil {
		return g, err
	}

	return g, nil
}

func (g Game) DrawFromStock() (Game, error) {
	return g.move(move{Type: moveTypeDrawFromStock})
}

func (g Game) MoveCard(from, to TableauLocation, count int) (Game, error) {
	if from == to {
		return g, nil
	}

	return g.move(move{
		Type:      moveTypeMoveCard,
		From:      from,
		To:        to,
		CardCount: count,
	})
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
