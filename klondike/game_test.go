package klondike

import (
	"math/rand"
	"testing"
)

func TestKlondike_NewGame(t *testing.T) {
	randSrc := rand.NewSource(42)

	for i := 0; i < 10; i++ {
		game := NewGame(randSrc, 1)
		state := game.State()

		for j := 0; j < 7; j++ {
			pile := state.Piles[j]
			if pile.Len() != j+1 {
				t.Fatal("pile.Len() is wrong")
			}
			if pile.FlippedCount() != 1 {
				t.Fatal("pile.FlippedCount() is wrong")
			}
			if !pile.Card(0).Flipped {
				t.Fatal("top card in pile should be flipped")
			}
		}

		for j := 0; j < 4; j++ {
			foundation := state.Foundations[j]
			if foundation.Len() != 0 {
				t.Fatal("foundation should be empty")
			}
		}

		if state.Stock.Len() != 24 {
			t.Fatal("stock.FlippedCount() is wrong")
		}

		if game.Solved() {
			t.Fatal("game should not be solved")
		}
	}
}

func TestKlondike_NewGameWithSeed(t *testing.T) {
	g1 := NewGameWithSeed(69, 1)
	g2 := NewGameWithSeed(69, 1)
	g3 := NewGameWithSeed(1337, 1)

	if g1.State().Hash() != g2.State().Hash() {
		t.Fatal("two games from the same seed should be equal")
	}

	if g1.State().Hash() == g3.State().Hash() {
		t.Fatal("two games from different seeds should not be equal")
	}
}
