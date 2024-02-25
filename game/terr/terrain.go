package terr

import (
	"github.com/mlange-42/tiny-world/game/resource"
)

type Terrain uint8

const (
	// Terrain
	Air Terrain = iota
	Buildable
	Grass
	Water
	Desert
	// Land use
	Path
	Field
	Tree
	Rock
	Farm
	Fisherman
	Lumberjack
	Mason
	EndTerrain
)

type Terrains uint16

func NewTerrains(dirs ...Terrain) Terrains {
	d := Terrains(0)
	for _, dir := range dirs {
		d |= (1 << dir)
	}
	return d
}

func (d *Terrains) Set(dir Terrain) {
	*d |= (1 << dir)
}

func (d *Terrains) Unset(dir Terrain) {
	*d &= ^(1 << dir)
}

// Contains checks whether all the argument's bits are contained in this Subscription.
func (d Terrains) Contains(dir Terrain) bool {
	bits := Terrains(1 << dir)
	return (bits & d) == bits
}

type TerrainProps struct {
	Name        string
	IsTerrain   bool
	BuildOn     Terrains
	BuildOnFree Terrains
	CanBuild    bool
	CanBuy      bool
	Production  Production
	BuildCost   []BuildCost
}

type Production struct {
	Produces          resource.Resource
	ConsumesFood      int
	RequiredTerrain   Terrain
	RequiredLandUse   Terrain
	ProductionTerrain Terrain
	ProductionLandUse Terrain
}

type BuildCost struct {
	Type   resource.Resource
	Amount int
}

var Descriptions = [EndTerrain]string{
	"Nothing",
	"Nothing",
	"A basic land tile.",
	"A land tile with water. Can be used by fisherman.",
	"A desert land tile.",
	"A path. Required by all buildings.",
	"Can be used by farms to produce food.",
	"Can be used by lumberjacks to produce wood.",
	"Can be used by masons to produce stones.",
	"Produces 1 food/min per neighboring field.",
	"Produces 1 food/min per neighboring water.",
	"Produces 1 wood/min per neighboring tree.",
	"Produces 1 stone/min per neighboring rock.",
}

var Properties = [EndTerrain]TerrainProps{
	{Name: "air", IsTerrain: true,
		CanBuild:   false,
		CanBuy:     false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "buildable", IsTerrain: true,
		CanBuild:   false,
		CanBuy:     false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "grass", IsTerrain: true,
		BuildOn:     NewTerrains(Buildable),
		BuildOnFree: NewTerrains(Buildable, Water, Desert),
		CanBuild:    true,
		CanBuy:      false,
		Production:  Production{Produces: resource.EndResources},
	},
	{Name: "water", IsTerrain: true,
		BuildOn:     NewTerrains(Buildable),
		BuildOnFree: NewTerrains(Buildable, Grass, Desert),
		CanBuild:    true,
		CanBuy:      false,
		Production:  Production{Produces: resource.EndResources},
	},
	{Name: "desert", IsTerrain: true,
		BuildOn:     NewTerrains(Buildable),
		BuildOnFree: NewTerrains(Buildable, Grass, Water),
		CanBuild:    true,
		CanBuy:      false,
		Production:  Production{Produces: resource.EndResources},
	},
	{Name: "path", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert, Water),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.EndResources},
		BuildCost: []BuildCost{
			{resource.Wood, 1},
		},
	},
	{Name: "field", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.EndResources},
		BuildCost: []BuildCost{
			{resource.Wood, 1},
			{resource.Stones, 1},
		},
	},
	{Name: "tree", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		CanBuy:     false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "rock", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert),
		CanBuild:   true,
		CanBuy:     false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "farm", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.Food, RequiredLandUse: Path, ProductionLandUse: Field, ConsumesFood: 1},
		BuildCost: []BuildCost{
			{resource.Wood, 5},
			{resource.Stones, 2},
		},
	},
	{Name: "fisherman", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.Food, RequiredLandUse: Path, ProductionTerrain: Water, ConsumesFood: 1},
		BuildCost: []BuildCost{
			{resource.Wood, 3},
			{resource.Stones, 0},
		},
	},
	{Name: "lumberjack", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.Wood, RequiredLandUse: Path, ProductionLandUse: Tree, ConsumesFood: 5},
		BuildCost: []BuildCost{
			{resource.Wood, 2},
			{resource.Stones, 3},
		},
	},
	{Name: "mason", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert),
		CanBuild:   true,
		CanBuy:     true,
		Production: Production{Produces: resource.Stones, RequiredLandUse: Path, ProductionLandUse: Rock, ConsumesFood: 5},
		BuildCost: []BuildCost{
			{resource.Wood, 5},
			{resource.Stones, 1},
		},
	},
}

var RandomTerrain = []Terrain{
	Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass,
	Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass,
	Water, Water, Water, Water,
	Desert,
	Tree, Tree, Tree, Tree, Tree, Tree,
	Rock,
}
