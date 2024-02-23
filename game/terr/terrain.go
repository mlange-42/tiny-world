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
	Name       string
	IsTerrain  bool
	BuildOn    Terrains
	CanBuild   bool
	Production Production
}

type Production struct {
	Produces          resource.Resource
	ConsumesFood      int
	RequiredTerrain   Terrain
	RequiredLandUse   Terrain
	ProductionTerrain Terrain
	ProductionLandUse Terrain
}

var Properties = [EndTerrain]TerrainProps{
	{Name: "air", IsTerrain: true,
		BuildOn:    NewTerrains(),
		CanBuild:   false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "buildable", IsTerrain: true,
		BuildOn:    NewTerrains(),
		CanBuild:   false,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "grass", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable, Water, Desert),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "water", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable, Grass, Desert),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "desert", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable, Grass, Water),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "path", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert, Water),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "field", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "tree", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "rock", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.EndResources},
	},
	{Name: "farm", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.Food, RequiredLandUse: Path, ProductionLandUse: Field, ConsumesFood: 1},
	},
	{Name: "fisherman", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.Food, RequiredLandUse: Path, ProductionTerrain: Water, ConsumesFood: 1},
	},
	{Name: "lumberjack", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.Wood, RequiredLandUse: Path, ProductionLandUse: Tree, ConsumesFood: 5},
	},
	{Name: "mason", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		Production: Production{Produces: resource.Stones, RequiredLandUse: Path, ProductionLandUse: Rock, ConsumesFood: 5},
	},
}
