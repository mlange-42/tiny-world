package res

type Rules struct {
	WorldSize      int
	RandomTerrains int

	InitialResources  []int
	StockPerWarehouse []int

	StockPerBuilding int
	HaulerCapacity   int
}
