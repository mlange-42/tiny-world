package comp

import (
	"image"

	"github.com/mlange-42/ark/ecs"
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
	Resource    resource.Resource
	Amount      uint8
	Stock       uint8
	Countdown   int
	IsHauling   bool
	HasRequired bool
}

type Consumption struct {
	Amount      []uint8
	Countdown   []int16
	IsSatisfied bool
}

type Population struct {
	Pop uint8
}

type PopulationSupport struct {
	Pop         uint8
	HasRequired bool
}

type Warehouse struct{}

type UnlocksTerrain struct{}

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

func (r *RandomSprite) GetRand() uint16 {
	if r == nil {
		return 0
	}
	return r.Rand
}

type BuildRadius struct {
	Radius uint8
}

type CardAnimation struct {
	image.Point
	Target     image.Point
	Terrain    terr.Terrain
	RandSprite uint16
	StartTick  int64
}
