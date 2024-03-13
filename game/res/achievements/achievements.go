package achievements

import (
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/comp"
	"github.com/mlange-42/tiny-world/game/res"
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/save"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Achievement struct {
	Name        string
	Description string
	Conditions  []Condition
	Completed   bool
}

type Condition struct {
	Type   string
	IDs    uint32
	Number int
}

type Achievements struct {
	Achievements []Achievement
	Completed    []string

	world *ecs.World

	terrainFilter generic.Filter1[comp.Terrain]
	stock         *res.Stock
	production    *res.Production

	checks map[string]func(uint32, int) bool
}

func New(world *ecs.World, f fs.FS, file string, playerFile string) *Achievements {
	var stock *res.Stock
	var prod *res.Production

	sRes := generic.NewResource[res.Stock](world)
	pRes := generic.NewResource[res.Production](world)
	if sRes.Has() {
		stock = sRes.Get()
	}
	if pRes.Has() {
		prod = pRes.Get()
	}

	a := Achievements{
		world:         world,
		terrainFilter: *generic.NewFilter1[comp.Terrain](),
		stock:         stock,
		production:    prod,
	}

	a.checks = map[string]func(uint32, int) bool{
		"terrain":          a.checkTerrain,
		"stock":            a.checkStock,
		"production":       a.checkProduction,
		"consumption":      a.checkConsumption,
		"total_production": a.checkTotalProduction,
		"net_production":   a.checkNetProduction,
	}
	parse := map[string]func(...string) uint32{
		"terrain":          a.parseTerrains,
		"stock":            a.parseResources,
		"production":       a.parseResources,
		"consumption":      a.parseResources,
		"total_production": a.parseResources,
		"net_production":   a.parseResources,
	}

	ach := []achievementJs{}
	err := util.FromJsonFs(f, file, &ach)
	if err != nil {
		panic(err)
	}

	for _, achieve := range ach {
		if strings.Contains(achieve.Name, " ") {
			log.Fatalf("disallowed spaces in achievement name '%s'", achieve.Name)
		}

		conditions := []Condition{}

		for _, c := range achieve.Conditions {
			conditions = append(conditions,
				Condition{
					Type:   c.Type,
					IDs:    parse[c.Type](c.IDs...),
					Number: c.Number,
				},
			)
		}

		a.Achievements = append(a.Achievements,
			Achievement{
				Name:        achieve.Name,
				Description: achieve.Description,
				Conditions:  conditions,
			},
		)
	}

	err = save.LoadAchievements(playerFile, &a.Completed)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			log.Fatal("error parsing achievement: ", err)
		}
	}

	for i := range a.Achievements {
		if slices.Contains(a.Completed, a.Achievements[i].Name) {
			a.Achievements[i].Completed = true
		}
	}

	return &a
}

func (a *Achievements) Check(ach *Achievement) {
	if ach.Completed {
		return
	}

	for i := range ach.Conditions {
		cond := &ach.Conditions[i]
		if !a.checks[cond.Type](cond.IDs, cond.Number) {
			return
		}
	}

	ach.Completed = true
}

func (a *Achievements) checkTerrain(ids uint32, num int) bool {
	cnt := 0
	query := a.terrainFilter.Query(a.world)
	for query.Next() {
		t := query.Get()
		bits := uint32(1) << t.Terrain
		if (bits & ids) == bits {
			cnt++
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) checkProduction(ids uint32, num int) bool {
	cnt := 0
	for i := range resource.Properties {
		bits := uint32(1) << i
		if (bits & ids) == bits {
			cnt += a.production.Prod[i]
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) checkConsumption(ids uint32, num int) bool {
	cnt := 0
	for i := range resource.Properties {
		bits := uint32(1) << i
		if (bits & ids) == bits {
			cnt += a.production.Cons[i]
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) checkStock(ids uint32, num int) bool {
	cnt := 0
	for i := range resource.Properties {
		bits := uint32(1) << i
		if (bits & ids) == bits {
			cnt += a.stock.Res[i]
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) checkNetProduction(ids uint32, num int) bool {
	cnt := 0
	for i := range resource.Properties {
		bits := uint32(1) << i
		if (bits & ids) == bits {
			cnt += a.production.Prod[i] - a.production.Cons[i]
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) checkTotalProduction(ids uint32, num int) bool {
	cnt := 0
	for i := range resource.Properties {
		bits := uint32(1) << i
		if (bits & ids) == bits {
			cnt += a.stock.Total[i]
			if cnt >= num {
				return true
			}
		}
	}
	return false
}

func (a *Achievements) parseTerrains(ids ...string) uint32 {
	return uint32(terr.ToTerrains(ids...))
}

func (a *Achievements) parseResources(ids ...string) uint32 {
	return uint32(resource.ToResources(ids...))
}

type achievementJs struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Conditions  []conditionJs `json:"conditions"`
}

type conditionJs struct {
	Type   string   `json:"type"`
	IDs    []string `json:"ids"`
	Number int      `json:"number"`
}
