# Klondike Game Engine

## Usage

```go
package main

import (
	"github.com/hasyimibhar/klondiker-solver/klondike"
)

func mustApplyMove(g klondike.Game, err error) klondike.Game {
	if err != nil {
		panic(err)
	}
	return g
}

func main() {
	game := klondike.NewGameWithSeed(42, 1)

	// For each move, the entire game state is copied.
	// This allows user to keep track of the state change,
	// and also perform undo/redo

	history := []klondike.Game{game}

 	// Draw from stock
	game = mustApplyMove(game.ApplyMove(klondike.DrawFromStock()))
	history = append(history, game)

 	// Move 1 card from pile 7 to 1
	game = mustApplyMove(game.ApplyMove(klondike.MoveCard().FromPile(6).ToPile(0).Count(1)))
	history = append(history, game)

	// Move 1 card from stock to heart foundation
	game = mustApplyMove(game.ApplyMove(klondike.MoveCard().FromStock().ToFoundation(klondike.CardTypeHeart).Count(1)))
	history = append(history, game)

	// Alternatively, store the list of moves and apply them
	// from the starting state

	anotherGame := klondike.NewGameWithSeed(42, 1)

	moves := []klondike.Move{
		klondike.DrawFromStock(),
		klondike.MoveCard().FromPile(6).ToPile(0).Count(1),
		klondike.MoveCard().FromStock().ToFoundation(klondike.CardTypeHeart).Count(1),
	}

	for _, m := range moves {
		anotherGame = mustApplyMove(anotherGame.ApplyMove(m))
	}
}
```
