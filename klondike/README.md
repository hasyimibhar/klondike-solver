# Klondike Game Engine

## Usage

```go
package main

import (
	"github.com/hasyimibhar/klondiker-solver/klondike"
)

func must(g klondike.Game, err error) klondike.Game {
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
	game = must(game.ApplyMove(klondike.Draw()))
	history = append(history, game)

 	// Move 1 card from pile 7 to 1
	game = must(game.ApplyMove(klondike.Move().FromPile(6, 1).ToPile(0)))
	history = append(history, game)

	// Move 1 card from stock to heart foundation
	game = must(game.ApplyMove(klondike.Move().FromStock().ToFoundation(klondike.Heart)))
	history = append(history, game)

	// Alternatively, store the list of moves and apply them
	// from the starting state

	anotherGame := klondike.NewGameWithSeed(42, 1)

	moves := []klondike.GameMove{
		klondike.Draw(),
		klondike.Move().FromPile(6, 1).ToPile(0),
		klondike.Move().FromStock().ToFoundation(klondike.Heart),
	}

	for _, m := range moves {
		anotherGame = must(anotherGame.ApplyMove(m))
	}
}
```
