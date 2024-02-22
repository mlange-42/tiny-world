package terr

type Terrain uint8

const (
	Air Terrain = iota
	Grass
	Path
	EndTerrain
)

var Names = [EndTerrain]string{
	"air",
	"grass",
	"path",
}
