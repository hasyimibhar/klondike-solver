package klondike

type moveType int

const (
	moveTypeDrawFromStock moveType = iota
	moveTypeMoveCard
)

type Move struct {
	t     moveType
	from  Location
	to    Location
	count int
}

func DrawFromStock() Move {
	return Move{t: moveTypeDrawFromStock}
}

func MoveCard() moveCardBuilder {
	return moveCardBuilder{}
}

type moveCardBuilder struct {
}

type moveCardBuilderFrom struct {
	from Location
}

type moveCardBuilderTo struct {
	from Location
	to   Location
}

func (b moveCardBuilder) FromPile(i int) moveCardBuilderFrom {
	return moveCardBuilderFrom{GetPile(i)}
}

func (b moveCardBuilder) FromFoundation(cardType CardType) moveCardBuilderFrom {
	return moveCardBuilderFrom{GetFoundation(cardType)}
}

func (b moveCardBuilder) FromStock() moveCardBuilderFrom {
	return moveCardBuilderFrom{LocationStock}
}

func (b moveCardBuilderFrom) ToPile(i int) moveCardBuilderTo {
	return moveCardBuilderTo{b.from, GetPile(i)}
}

func (b moveCardBuilderFrom) ToFoundation(cardType CardType) moveCardBuilderTo {
	return moveCardBuilderTo{b.from, GetFoundation(cardType)}
}

func (b moveCardBuilderTo) Count(n int) Move {
	return Move{
		t:     moveTypeMoveCard,
		from:  b.from,
		to:    b.to,
		count: n,
	}
}
