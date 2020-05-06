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

func Move() MoveBuilder {
	return MoveBuilder{}
}

type MoveBuilder struct {
	from  Location
	count int
}

func (b MoveBuilder) FromPile(i int, n int) MoveBuilder {
	return MoveBuilder{GetPile(i), n}
}

func (b MoveBuilder) FromFoundation(cardType CardType) MoveBuilder {
	return MoveBuilder{GetFoundation(cardType), 1}
}

func (b MoveBuilder) FromStock() MoveBuilder {
	return MoveBuilder{LocationStock, 1}
}

func (b MoveBuilder) ToPile(i int) GameMove {
	return GameMove{
		t:     moveTypeMoveCard,
		from:  b.from,
		to:    GetPile(i),
		count: b.count,
	}
}

func (b MoveBuilder) ToFoundation(cardType CardType) GameMove {
	return GameMove{
		t:     moveTypeMoveCard,
		from:  b.from,
		to:    GetFoundation(cardType),
		count: b.count,
	}
}
