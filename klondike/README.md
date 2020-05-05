# Klondike Game Engine

## Usage

```go
package main

import (
	"github.com/hasyimibhar/klondiker-solver/klondike"
)

func main() {
	game := klondike.NewGameWithSeed(42, 1)

	// For each move, the entire game state is copied.
	// This allows user to keep track of the state change,
	// and also perform undo/redo
	game, _ = game.DrawFromStock() // Draw from stock
	game, _ = game.MoveCard().FromPile(6).ToPile(0).Count(1) // Move 1 card from pile 1 to 7
}
```
