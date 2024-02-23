package comp

import (
	"image"
)

type Tile struct {
	image.Point
}

type UpdateTick struct {
	Tick int64
}

type Production struct {
	Amount int
}
