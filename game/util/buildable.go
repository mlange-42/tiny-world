package util

import (
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
)

func IsBuildable(x, y int, query generic.Query2[comp.Tile, comp.BuildRadius]) bool {
	for query.Next() {
		tile, rad := query.Get()

		dx, dy := tile.X-x, tile.Y-y
		r2 := int(rad.Radius) * int(rad.Radius)
		if dx*dx+dy*dy <= r2 {
			query.Close()
			return true
		}
	}
	return false
}
