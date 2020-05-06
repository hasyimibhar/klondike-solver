## Usage

To play:

```sh
$ go run .
```

Commands:

- `u`: Undo a move
- `d`: Draw from stock
- `new [seed]`: Start a new game. If `seed` is not provided, it uses the current time as seed.
- `<from> [<to> <count>]`: Move `count` cards from `from` to `to`. Possible locations:

    - `p1` - `p7`: Pile 1 - 7
    - `s`: Stock
    - `fh`, `fd`, `fs`, `fc`: Foundations (second letter is the card shape)
