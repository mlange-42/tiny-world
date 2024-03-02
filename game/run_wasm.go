//go:build js

package game

func run(g *Game, name string, load bool) {
	if err := runGame(g, load, name, "paper"); err != nil {
		panic(err)
	}
}
