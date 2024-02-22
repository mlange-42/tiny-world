package terr

type Terrain uint8

const (
	Air Terrain = iota
	Buildable
	Grass
	Path
	EndTerrain
)

var Names = [EndTerrain]string{
	"air",
	"buildable",
	"grass",
	"path",
}
