package res

import (
	"fmt"
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Rules struct {
	WorldSize           int
	InitialBuildRadius  int
	RandomTerrainsCount int
	RandomTerrains      []terr.Terrain

	InitialResources []int
}

func NewRules(f fs.FS, file string) Rules {
	rulesHelper := rulesJs{}
	err := util.FromJsonFs(f, file, &rulesHelper)
	if err != nil {
		panic(err)
	}

	storage := make([]int, len(resource.Properties))
	for _, entry := range rulesHelper.InitialResources {
		res, ok := resource.ResourceID(entry.Resource)
		if !ok {
			panic(fmt.Sprintf("unknown resource %s", entry.Resource))
		}
		storage[res] = entry.Amount
	}

	randTerr := []terr.Terrain{}
	for _, str := range rulesHelper.RandomTerrains {
		randTerr = append(randTerr, toTerrain(str))
	}

	return Rules{
		WorldSize:           rulesHelper.WorldSize,
		InitialBuildRadius:  rulesHelper.InitialBuildRadius,
		InitialResources:    storage,
		RandomTerrainsCount: rulesHelper.RandomTerrainsCount,
		RandomTerrains:      randTerr,
	}
}

type rulesJs struct {
	WorldSize           int `json:"world_size"`
	InitialBuildRadius  int `json:"initial_build_radius"`
	RandomTerrainsCount int `json:"random_terrains_count"`

	RandomTerrains   []string           `json:"random_terrains"`
	InitialResources []resourceAmountJs `json:"initial_resources"`
}

type resourceAmountJs struct {
	Resource string `json:"resource"`
	Amount   int    `json:"amount"`
}

func toTerrain(t string) terr.Terrain {
	id, ok := terr.TerrainID(t)
	if !ok {
		panic(fmt.Sprintf("unknown terrain %s", t))
	}
	return id
}
