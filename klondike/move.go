package klondike

type moveType int

const (
	moveTypeDrawFromStock moveType = iota
	moveTypeMoveCard
)

type move struct {
	Type      moveType
	From      Location
	To        Location
	CardCount int
}

type moveCardBuilder struct {
	game Game
}

type moveCardBuilderFrom struct {
	game Game
	from Location
}

type moveCardBuilderTo struct {
	game Game
	from Location
	to   Location
}

func (b moveCardBuilder) FromPile(i int) moveCardBuilderFrom {
	return moveCardBuilderFrom{b.game, GetPile(i)}
}

func (b moveCardBuilder) FromFoundation(cardType CardType) moveCardBuilderFrom {
	return moveCardBuilderFrom{b.game, GetFoundation(cardType)}
}

func (b moveCardBuilder) FromStock() moveCardBuilderFrom {
	return moveCardBuilderFrom{b.game, LocationStock}
}

func (b moveCardBuilderFrom) ToPile(i int) moveCardBuilderTo {
	return moveCardBuilderTo{b.game, b.from, GetPile(i)}
}

func (b moveCardBuilderFrom) ToFoundation(cardType CardType) moveCardBuilderTo {
	return moveCardBuilderTo{b.game, b.from, GetFoundation(cardType)}
}

func (b moveCardBuilderTo) Count(n int) (Game, error) {
	return b.game.move(move{
		Type:      moveTypeMoveCard,
		From:      b.from,
		To:        b.to,
		CardCount: n,
	})
}
