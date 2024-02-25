package comp

import (
	"image"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/resource"
)

type Tile struct {
	image.Point
}

type UpdateTick struct {
	Tick int64
}

type Production struct {
	Type            resource.Resource
	Amount          int
	FoodConsumption int
	Countdown       int
	Paused          bool
}

type Consumption struct {
	Amount    int
	Countdown int
}

type Warehouse struct{}

type ProductionMarker struct {
	StartTick int64
	Resource  resource.Resource
}

type Hauler struct {
	Hauls resource.Resource
	Home  ecs.Entity
	Path  []Tile
}
