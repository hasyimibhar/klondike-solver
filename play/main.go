package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hasyimibhar/klondiker-solver/klondike"
	"github.com/logrusorgru/aurora"
)

var au aurora.Aurora
var CardTypeSymbol map[klondike.CardType]string

func init() {
	au = aurora.NewAurora(true)

	CardTypeSymbol = map[klondike.CardType]string{
		klondike.Spade:   `♠`,
		klondike.Club:    `♣`,
		klondike.Heart:   `♥`,
		klondike.Diamond: `♦`,
	}
}

func PrintCard(c klondike.Card) {
	if !c.Flipped {
		fmt.Print("(? ?)")
		return
	}

	var n string
	if c.Number == 1 {
		n = " A"
	} else if c.Number == 11 {
		n = " J"
	} else if c.Number == 12 {
		n = " Q"
	} else if c.Number == 13 {
		n = " K"
	} else {
		n = strconv.FormatInt(int64(c.Number), 10)
		if c.Number < 10 {
			n = " " + n
		}
	}

	s := fmt.Sprintf("(%s%s)", CardTypeSymbol[c.Type], n)

	if c.Type == klondike.Diamond || c.Type == klondike.Heart {
		fmt.Printf("%s", au.Red(s).BgWhite())
	} else {
		fmt.Printf("%s", au.Black(s).BgWhite())
	}
}

func PrintPile(p klondike.Pile) {
	for i := p.Len() - 1; i >= 0; i-- {
		PrintCard(p.Card(i))
	}
}

func PrintStock(s klondike.Stock) {
	if s.Len() == 0 && len(s.Drawn()) == 0 {
		fmt.Printf("[X]")
		return
	}
	if len(s.Drawn()) == 0 {
		fmt.Printf("[S]")
		return
	}

	fmt.Printf("[S] ")
	for _, c := range s.Drawn() {
		PrintCard(c)
	}
}

func PrintFoundation(f klondike.Foundation) {
	if f.Len() == 0 {
		s := fmt.Sprintf("(%s  )", CardTypeSymbol[f.CardType()])
		if f.CardType() == klondike.Diamond || f.CardType() == klondike.Heart {
			fmt.Printf("%s", au.Red(s).BgWhite())
		} else {
			fmt.Printf("%s", au.Black(s).BgWhite())
		}

		return
	}

	PrintCard(f.Card())
}

func PrintState(s klondike.GameState) {
	fmt.Print("Foundations: ")
	for i := 0; i < 4; i++ {
		PrintFoundation(s.Foundations[i])
	}

	fmt.Println()
	fmt.Println()

	PrintStock(s.Stock)
	fmt.Println()
	fmt.Println()
	fmt.Println("Piles:")
	for i := 0; i < 7; i++ {
		fmt.Printf("[%d] ", i+1)
		PrintPile(s.Piles[i])
		fmt.Println()
	}
}

func ParseMove(game klondike.Game, cmd string) (klondike.GameMove, error) {
	if cmd == "d" || cmd == "draw" {
		return klondike.Draw(), nil
	}

	tokens := strings.Split(cmd, " ")

	if len(tokens) < 1 {
		return klondike.GameMove{}, errors.New("usage: from [to count]")
	}

	moveBuilder := klondike.Move()
	var card klondike.Card

	switch tokens[0] {
	case "s":
		moveBuilder = moveBuilder.FromStock()
		card = game.State().Stock.Card(0)

	case "fh":
		moveBuilder = moveBuilder.FromFoundation(klondike.Heart)
		card = game.State().Foundations[klondike.Heart-1].Card()

	case "fd":
		moveBuilder = moveBuilder.FromFoundation(klondike.Diamond)
		card = game.State().Foundations[klondike.Diamond-1].Card()

	case "fs":
		moveBuilder = moveBuilder.FromFoundation(klondike.Spade)
		card = game.State().Foundations[klondike.Spade-1].Card()

	case "fc":
		moveBuilder = moveBuilder.FromFoundation(klondike.Club)
		card = game.State().Foundations[klondike.Club-1].Card()

	case "p1", "p2", "p3", "p4", "p5", "p6", "p7":
		n := 1
		if len(tokens) == 3 {
			var err error
			n, err = strconv.Atoi(tokens[2])
			if err != nil {
				return klondike.GameMove{}, errors.New("invalid count")
			}
		}

		i, _ := strconv.Atoi(tokens[0][1:])
		moveBuilder = moveBuilder.FromPile(i-1, n)
		card = game.State().Piles[i-1].Card(0)

	default:
		return klondike.GameMove{}, errors.New("invalid from")
	}

	var move klondike.GameMove

	// If no destination is provided, set it to foundation
	if len(tokens) == 1 {
		tokens = append(tokens, "f")
	}

	switch tokens[1] {
	case "f":
		move = moveBuilder.ToFoundation(card.Type)

	case "p1", "p2", "p3", "p4", "p5", "p6", "p7":
		i, _ := strconv.Atoi(tokens[1][1:])
		move = moveBuilder.ToPile(i - 1)

	default:
		return klondike.GameMove{}, errors.New("invalid to")
	}

	return move, nil
}

func main() {
	game := klondike.NewGameWithSeed(time.Now().UnixNano(), 1)
	history := []klondike.Game{}
	reader := bufio.NewReader(os.Stdin)

	PrintState(game.State())

	for {
		fmt.Print("> ")
		txt, _ := reader.ReadString('\n')
		txt = strings.Replace(txt, "\n", "", -1)

		if txt == "" {
			continue
		}

		tokens := strings.Split(txt, " ")

		switch tokens[0] {
		case "new":
			var seed int64
			if len(tokens) == 1 {
				seed = time.Now().UnixNano()
			} else {
				var err error
				seed, err = strconv.ParseInt(tokens[1], 10, 64)
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
			}

			game = klondike.NewGameWithSeed(seed, 1)
			PrintState(game.State())

		case "u", "undo":
			if len(history) > 0 {
				game = history[len(history)-1]
				history = history[:len(history)-1]
			}

			PrintState(game.State())

		default:
			move, err := ParseMove(game, txt)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			next, err := game.ApplyMove(move)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			history = append(history, game)
			game = next

			PrintState(game.State())
		}
	}
}
