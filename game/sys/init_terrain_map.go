package sys

import (
	"fmt"
	"image"
	"io/fs"
	"log"
	"math"
	"math/rand"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/terr"
)

// InitTerrainMap system.
type InitTerrainMap struct {
	FS        fs.FS
	MapFolder string
	MapFile   string
}

// Initialize the system
func (s *InitTerrainMap) Initialize(world *ecs.World) {
	rules := ecs.GetResource[res.Rules](world)
	terrain := ecs.GetResource[res.Terrain](world)
	bounds := ecs.GetResource[res.WorldBounds](world)
	fac := ecs.GetResource[res.EntityFactory](world)

	mapData, err := save.LoadMap(s.FS, s.MapFolder, s.MapFile)
	if err != nil {
		log.Fatal("error reading map file", err.Error())
	}

	xOff, yOff := (terrain.Width()+1-len(mapData[0]))/2, (terrain.Height()+1-len(mapData))/2

	x, y := terrain.Width()/2, terrain.Height()/2
	bounds.Min = image.Pt(x-1, y-1)
	bounds.Max = image.Pt(x+1, y+1)

	for y := 0; y < len(mapData); y++ {
		line := mapData[y]
		yy := y + yOff
		for x := 0; x < len(line); x++ {
			rn := line[x]
			ter, ok := terr.SymbolToTerrain[rn]
			if !ok {
				panic(fmt.Sprintf("unknown map symbol '%s'", string(rn)))
			}
			xx := x + xOff
			if ter.Terrain != terr.Air {
				fac.Set(world, xx, yy, ter.Terrain, uint16(rand.Int31n(math.MaxUint16)))
			}
			if ter.LandUse != terr.Air {
				fac.Set(world, xx, yy, ter.LandUse, uint16(rand.Int31n(math.MaxUint16)))
			}
		}
	}

	fac.SetBuildable(x, y, rules.InitialBuildRadius, true)

	radFilter := generic.NewFilter2[comp.Tile, comp.BuildRadius]()
	radQuery := radFilter.Query(world)
	for radQuery.Next() {
		tile, rad := radQuery.Get()
		fac.SetBuildable(tile.X, tile.Y, int(rad.Radius), true)
	}
}

// Update the system
func (s *InitTerrainMap) Update(world *ecs.World) {}

// Finalize the system
func (s *InitTerrainMap) Finalize(world *ecs.World) {}
