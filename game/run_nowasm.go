//go:build !js

package game

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func run(g *Game) {
	if err := command(g).Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func command(g *Game) *cobra.Command {
	var tileSet, saveFile string
	root := &cobra.Command{
		Use:           "tiny-world",
		Short:         "A tiny, slow-paced world and colony building game.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			runGame(g, saveFile, tileSet)
		},
	}
	root.Flags().StringVarP(&tileSet, "tileset", "t", "paper", "Tileset to use.")
	root.Flags().StringVarP(&saveFile, "savefile", "s", "", "Savefile to load.")

	return root
}
