package terr

import "github.com/hajimehoshi/ebiten/v2"

type Terrain uint8

const (
	Air Terrain = iota
	Buildable
	Grass
	Water
	Desert
	Field
	Path
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
	ShortKey   ebiten.Key
	Production Production
}

type Production struct {
	Produces          bool
	RequiredTerrain   Terrain
	RequiredLandUse   Terrain
	ProductionTerrain Terrain
	ProductionLandUse Terrain
}

var Properties = [EndTerrain]TerrainProps{
	{Name: "air", IsTerrain: true,
		BuildOn:    NewTerrains(),
		CanBuild:   false,
		ShortKey:   ebiten.KeyEscape,
		Production: Production{Produces: false},
	},
	{Name: "buildable", IsTerrain: true,
		BuildOn:    NewTerrains(),
		CanBuild:   false,
		ShortKey:   ebiten.KeyEscape,
		Production: Production{Produces: false},
	},
	{Name: "grass", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		CanBuild:   true,
		ShortKey:   ebiten.Key1,
		Production: Production{Produces: false},
	},
	{Name: "water", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		CanBuild:   true,
		ShortKey:   ebiten.Key2,
		Production: Production{Produces: false},
	},
	{Name: "desert", IsTerrain: true,
		BuildOn:    NewTerrains(Buildable),
		CanBuild:   true,
		ShortKey:   ebiten.Key3,
		Production: Production{Produces: false},
	},
	{Name: "field", IsTerrain: false,
		BuildOn:    NewTerrains(Grass),
		CanBuild:   true,
		ShortKey:   ebiten.Key4,
		Production: Production{Produces: false},
	},
	{Name: "path", IsTerrain: false,
		BuildOn:    NewTerrains(Grass, Desert),
		CanBuild:   true,
		ShortKey:   ebiten.Key5,
		Production: Production{Produces: false},
	},
}
