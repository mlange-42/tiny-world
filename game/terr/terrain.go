package terr

import "github.com/hajimehoshi/ebiten/v2"

type Terrain uint8

const (
	Air Terrain = iota
	Buildable
	Grass
	Path
	Cursor
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
	Name      string
	IsTerrain bool
	BuildOn   Terrains
	CanBuild  bool
	ShortKey  ebiten.Key
}

var Properties = [EndTerrain]TerrainProps{
	{"air", true, NewTerrains(), false, ebiten.KeyEscape},
	{"buildable", true, NewTerrains(), false, ebiten.KeyEscape},
	{"grass", true, NewTerrains(Buildable), true, ebiten.Key1},
	{"path", false, NewTerrains(Grass), true, ebiten.Key2},
	{"cursor", false, NewTerrains(), false, ebiten.KeyEscape},
}
