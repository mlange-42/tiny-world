package comp

import (
	"image"

	"github.com/mlange-42/tiny-world/game/resource"
)

type Tile struct {
	image.Point
}

type UpdateTick struct {
	Tick int64
}

type Production struct {
	Type   resource.Resource
	Amount int
}
