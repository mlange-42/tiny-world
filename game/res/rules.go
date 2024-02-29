package res

import "github.com/mlange-42/tiny-world/game/resource"

type Rules struct {
	WorldSize      int
	RandomTerrains int

	InitialResources  [resource.EndResources]int
	StockPerWarehouse [resource.EndResources]int

	StockPerBuilding int
	HaulerCapacity   int
}
