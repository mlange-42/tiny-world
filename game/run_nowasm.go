//go:build !js

package game

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Run(data embed.FS) {
	gameData = data
	if err := command().Execute(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func command() *cobra.Command {
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
			run(saveFile, tileSet)
		},
	}
	root.Flags().StringVarP(&tileSet, "tileset", "t", "paper", "Tileset to use.")
	root.Flags().StringVarP(&saveFile, "savefile", "s", "", "Savefile to load.")

	return root
}
