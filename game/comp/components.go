package comp

import (
	"image"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Tile struct {
	image.Point
}

type Terrain struct {
	terr.Terrain
}

type UpdateTick struct {
	Tick int64
}

type Path struct {
	Haulers []HaulerEntry
}

type HaulerEntry struct {
	Entity ecs.Entity
	YPos   float64
}

type Production struct {
	Type            resource.Resource
	Amount          int
	FoodConsumption int
	Countdown       int
	Stock           int
	IsHauling       bool
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
	Hauls        resource.Resource
	Home         ecs.Entity
	Path         []Tile
	Index        int
	PathFraction uint8
}

type HaulerSprite struct {
	SpriteIndex int
}

type RandomSprite struct {
	Rand uint16
}
