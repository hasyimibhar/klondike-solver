package klondike

type moveType int

const (
	moveTypeDrawFromStock moveType = iota
	moveTypeMoveCard
)

type GameMove struct {
	t     moveType
	from  Location
	to    Location
	count int
}

func Draw() GameMove {
	return GameMove{t: moveTypeDrawFromStock}
}

func Move() moveBuilderFrom {
	return moveBuilderFrom{}
}

type moveBuilderFrom struct{}

type moveBuilderTo struct {
	from  Location
	count int
}

func (b moveBuilderFrom) FromPile(i int, n int) moveBuilderTo {
	return moveBuilderTo{GetPile(i), n}
}

func (b moveBuilderFrom) FromFoundation(cardType CardType) moveBuilderTo {
	return moveBuilderTo{GetFoundation(cardType), 1}
}

func (b moveBuilderFrom) FromStock() moveBuilderTo {
	return moveBuilderTo{LocationStock, 1}
}

func (b moveBuilderTo) ToPile(i int) GameMove {
	return GameMove{
		t:     moveTypeMoveCard,
		from:  b.from,
		to:    GetPile(i),
		count: b.count,
	}
}

func (b moveBuilderTo) ToFoundation(cardType CardType) GameMove {
	return GameMove{
		t:     moveTypeMoveCard,
		from:  b.from,
		to:    GetFoundation(cardType),
		count: b.count,
	}
}
