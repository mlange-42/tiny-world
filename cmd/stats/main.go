package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/terr"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please specify a map file!")
	}
	file := os.Args[1]
	if !strings.HasSuffix(file, ".asc") {
		file += ".asc"
	}
	mapStr, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Error reading map '%s': %s", file, err.Error())
	}

	mapData, err := save.ParseMap(string(mapStr))
	if err != nil {
		log.Fatalf("Error parsing map '%s': %s", file, err.Error())
	}

	resource.Prepare(os.DirFS("."), "data/json/resources.json")
	terr.Prepare(os.DirFS("."), "data/json/terrain.json")

	frequencies := map[terr.Terrain]int{}
	total := 0

	for y := 0; y < len(mapData.Data); y++ {
		line := mapData.Data[y]
		for x := 0; x < len(line); x++ {
			rn := line[x]
			ter, ok := terr.SymbolToTerrain[rn]
			if !ok {
				panic(fmt.Sprintf("unknown map symbol '%s'", string(rn)))
			}
			tBits := terr.Properties[ter.Terrain].TerrainBits
			luBits := terr.Properties[ter.LandUse].TerrainBits
			if ter.Terrain != terr.Air && tBits.Contains(terr.CanBuild) && !tBits.Contains(terr.CanBuy) {
				if cnt, ok := frequencies[ter.Terrain]; ok {
					frequencies[ter.Terrain] = cnt + 1
				} else {
					frequencies[ter.Terrain] = 1
				}
				total++
			}
			if ter.LandUse != terr.Air && luBits.Contains(terr.CanBuild) && !luBits.Contains(terr.CanBuy) {
				if cnt, ok := frequencies[ter.LandUse]; ok {
					frequencies[ter.LandUse] = cnt + 1
				} else {
					frequencies[ter.LandUse] = 1
				}
				total++
			}
		}
	}

	minCount := math.MaxInt
	for _, cnt := range frequencies {
		if cnt < minCount {
			minCount = cnt
		}
	}

	for t, cnt := range frequencies {
		fmt.Printf("%10s (%s): %4d (%5.1f%%) %5.1fx\n",
			terr.Properties[t].Name, string(terr.Properties[t].Symbols[0]),
			cnt,
			float64(cnt*100)/float64(total), float64(cnt)/float64(minCount))
	}
}
