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
	Warehouse
	EndTerrain
)

type Terrains uint32

var Buildings Terrains = NewTerrains(
	Farm,
	Fisherman,
	Lumberjack,
	Mason,
	Warehouse,
)

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
	Name        string      `json:"name"`
	IsTerrain   bool        `json:"is_terrain"`
	IsPath      bool        `json:"is_path"`
	IsWarehouse bool        `json:"is_warehouse"`
	BuildOn     Terrains    `json:"build_on"`
	Below       Terrain     `json:"terrain_below"`
	ConnectsTo  Terrains    `json:"connects_to"`
	CanBuild    bool        `json:"can_build"`
	CanBuy      bool        `json:"can_buy"`
	Production  Production  `json:"production"`
	BuildCost   []BuildCost `json:"build_cost,omitempty"`
}

type Production struct {
	Produces          resource.Resource
	MaxProduction     int
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
	"A basic land tile",
	"A land tile with water. Can be used by fisherman",
	"A desert land tile",
	"A path. Required by all buildings",
	"Can be used by farms to produce food",
	"Can be used by lumberjacks to produce wood",
	"Can be used by masons to produce stones",
	"Produces 1 food/min per neighboring field",
	"Produces 1 food/min per neighboring water",
	"Produces 1 wood/min per neighboring tree",
	"Produces 1 stone/min per neighboring rock",
	"Stores resources",
}

var Properties = [EndTerrain]TerrainProps{
	{Name: "air", IsTerrain: true,
		CanBuild: false,
		CanBuy:   false,
	},
	{Name: "buildable", IsTerrain: true,
		CanBuild: false,
		CanBuy:   false,
	},
	{Name: "grass", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		ConnectsTo: NewTerrains(Grass, Water, Desert),
		CanBuild:   true,
		CanBuy:     false,
	},
	{Name: "water", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		ConnectsTo: NewTerrains(Water),
		Below:      Grass,
		CanBuild:   true,
		CanBuy:     false,
	},
	{Name: "desert", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		ConnectsTo: NewTerrains(Desert),
		Below:      Grass,
		CanBuild:   true,
		CanBuy:     false,
	},
	{Name: "path", IsTerrain: false,
		IsPath:     true,
		BuildOn:    NewTerrains(Grass, Desert, Water),
		ConnectsTo: NewTerrains(Path, Farm, Fisherman, Lumberjack, Mason, Warehouse),
		CanBuild:   true,
		CanBuy:     true,
		BuildCost: []BuildCost{
			{resource.Wood, 1},
		},
	},
	{Name: "field", IsTerrain: false,
		BuildOn:  NewTerrains(Grass),
		CanBuild: true,
		CanBuy:   true,
		BuildCost: []BuildCost{
			{resource.Wood, 1},
			{resource.Stones, 1},
		},
	},
	{Name: "tree", IsTerrain: false,
		BuildOn:  NewTerrains(Grass),
		CanBuild: true,
		CanBuy:   false,
	},
	{Name: "rock", IsTerrain: false,
		BuildOn:  NewTerrains(Grass, Desert),
		CanBuild: true,
		CanBuy:   false,
	},
	{Name: "farm", IsTerrain: false,
		BuildOn:  NewTerrains(Grass),
		CanBuild: true,
		CanBuy:   true,
		Production: Production{
			Produces: resource.Food, MaxProduction: 7,
			RequiredLandUse: Path, ProductionLandUse: Field, ConsumesFood: 1},
		BuildCost: []BuildCost{
			{resource.Wood, 5},
			{resource.Stones, 2},
		},
	},
	{Name: "fisherman", IsTerrain: false,
		BuildOn:  NewTerrains(Grass, Desert),
		CanBuild: true,
		CanBuy:   true,
		Production: Production{
			Produces: resource.Food, MaxProduction: 5,
			RequiredLandUse: Path, ProductionTerrain: Water, ConsumesFood: 1},
		BuildCost: []BuildCost{
			{resource.Wood, 3},
			{resource.Stones, 0},
		},
	},
	{Name: "lumberjack", IsTerrain: false,
		BuildOn:  NewTerrains(Grass),
		CanBuild: true,
		CanBuy:   true,
		Production: Production{
			Produces: resource.Wood, MaxProduction: 7,
			RequiredLandUse: Path, ProductionLandUse: Tree, ConsumesFood: 5},
		BuildCost: []BuildCost{
			{resource.Wood, 2},
			{resource.Stones, 3},
		},
	},
	{Name: "mason", IsTerrain: false,
		BuildOn:  NewTerrains(Grass, Desert),
		CanBuild: true,
		CanBuy:   true,
		Production: Production{
			Produces: resource.Stones, MaxProduction: 3,
			RequiredLandUse: Path, ProductionLandUse: Rock, ConsumesFood: 5},
		BuildCost: []BuildCost{
			{resource.Wood, 10},
		},
	},
	{Name: "warehouse", IsTerrain: false,
		IsWarehouse: true,
		BuildOn:     NewTerrains(Grass, Desert),
		CanBuild:    true,
		CanBuy:      true,
		Production:  Production{Produces: resource.EndResources},
		BuildCost: []BuildCost{
			{resource.Wood, 25},
			{resource.Stones, 25},
		},
	},
}

var idLookup map[string]Terrain

var RandomTerrain = []Terrain{
	Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass,
	Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass, Grass,
	Water, Water, Water, Water, Water, Water,
	Desert,
	Tree, Tree, Tree, Tree,
	Rock,
}

func init() {
	idLookup = map[string]Terrain{}

	for i := Terrain(0); i < EndTerrain; i++ {
		idLookup[Properties[i].Name] = i
	}
}

func TerrainID(name string) (Terrain, bool) {
	t, ok := idLookup[name]
	return t, ok
}
