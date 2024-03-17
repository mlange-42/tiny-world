package res

import "github.com/mlange-42/tiny-world/game/terr"

type RandomTerrains struct {
	Terrains       []terr.Terrain
	AllowRemove    []bool
	TotalAvailable int
	TotalPlaced    int
}
