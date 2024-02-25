package res

import "github.com/mlange-42/tiny-world/game/resource"

type Rules struct {
	AllowStroke         bool
	AllowReplaceTerrain bool
	AllowRemoveNatural  bool
	AllowRemoveBuilt    bool

	WorldSize      int
	RandomTerrains int

	InitialResources  [resource.EndResources]int
	StockPerWarehouse [resource.EndResources]int
}
