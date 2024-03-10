//go:build !js

package game

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func run(g *Game, name string, load bool) {
	cobra.MousetrapHelpText = ""
	if err := command(g, name, load).Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func command(g *Game, name string, load bool) *cobra.Command {
	var tileSet string
	root := &cobra.Command{
		Use:           "tiny-world",
		Short:         "A tiny, slow-paced world and colony building game.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := runGame(g, load, name, tileSet)
			if err != nil {
				panic(err)
			}
		},
	}
	root.Flags().StringVarP(&tileSet, "tileset", "t", "paper", "Tileset to use.")

	return root
}

type canvasHelper struct{}

func newCanvasHelper() *canvasHelper {
	return &canvasHelper{}
}

func (c *canvasHelper) isMouseInside(width, height int) bool {
	x, y := ebiten.CursorPosition()
	return x >= 0 && y >= 0 && x < width && y < height
}
