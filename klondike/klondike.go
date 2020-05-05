package klondike

import (
	"crypto/sha256"
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

type moveType int

const (
	moveTypeDrawFromStock moveType = iota
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

	case TableauLocationFoundationClub, TableauLocationFoundationSpade, TableauLocationFoundationHeart, TableauLocationFoundationDiamond:
		f := g.state.Foundations[m.From-8]
		g.state.Foundations[m.From-8], cards, err = f.pop(m.CardCount)

	case TableauLocationPile1, TableauLocationPile2, TableauLocationPile3, TableauLocationPile4, TableauLocationPile5, TableauLocationPile6, TableauLocationPile7:
		p := g.state.Piles[m.From-1]
		g.state.Piles[m.From-1], cards, err = p.pop(m.CardCount)
	}

	if err != nil {
		return g, err
	}

	popped := make([]Card, len(cards))
	copy(popped, cards)

	switch m.To {
	case TableauLocationStock:
		g.state.Stock, err = g.state.Stock.place(popped)

	case TableauLocationFoundationClub, TableauLocationFoundationSpade, TableauLocationFoundationHeart, TableauLocationFoundationDiamond:
		f := g.state.Foundations[m.To-8]
		g.state.Foundations[m.To-8], err = f.place(popped)

	case TableauLocationPile1, TableauLocationPile2, TableauLocationPile3, TableauLocationPile4, TableauLocationPile5, TableauLocationPile6, TableauLocationPile7:
		p := g.state.Piles[m.To-1]
		g.state.Piles[m.To-1], err = p.place(popped)
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
