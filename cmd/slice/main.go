package main

import (
	"fmt"

	"github.com/mlange-42/tiny-world/cmd/util"
)

const (
	basePath = "artwork"
	tileSet  = "sprites"
)

func main() {
	util.Walk(basePath, tileSet, func(sheet util.TileSheet, dir util.Directory) {
		fmt.Println(dir)
	})
}
