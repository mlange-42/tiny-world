package terr

import (
	"fmt"
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/resource"
)

type Terrain uint8

var Air Terrain
var Buildable Terrain
var Default Terrain
var Warehouse Terrain

type Terrains uint32

var Buildings Terrains

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

var Properties []TerrainProps

var idLookup map[string]Terrain

var RandomTerrain []Terrain

func Prepare(f fs.FS, file string) {
	propsHelper := props{}
	err := util.FromJsonFs(f, file, &propsHelper)
	if err != nil {
		panic(err)
	}

	idLookup = map[string]Terrain{}
	for i, t := range propsHelper.Terrains {
		idLookup[t.Name] = Terrain(i)
	}

	props := []TerrainProps{}
	for i, t := range propsHelper.Terrains {
		cost := []ResourceAmount{}
		for _, cst := range t.BuildCost {
			id, ok := resource.ResourceID(cst.Resource)
			if !ok {
				panic(fmt.Sprintf("unknown resource %s", cst.Resource))
			}
			cost = append(cost, ResourceAmount{
				Resource: id,
				Amount:   cst.Amount,
			})
		}
		var prodRes resource.Resource
		var consResource resource.Resource
		if t.Production.MaxProduction > 0 {
			var ok bool
			prodRes, ok = resource.ResourceID(t.Production.Resource)
			if !ok {
				panic(fmt.Sprintf("unknown resource %s", t.Production.Resource))
			}
		}
		if t.Production.ConsumesAmount > 0 {
			var ok bool
			consResource, ok = resource.ResourceID(t.Production.ConsumesResource)
			if !ok {
				panic(fmt.Sprintf("unknown resource %s", t.Production.Resource))
			}
		}

		var terrBelow Terrain
		if t.TerrainBelow != "" {
			terrBelow = toTerrain(idLookup, t.TerrainBelow)
		}
		var requiredTerrain Terrain
		if t.Production.RequiredTerrain != "" {
			requiredTerrain = toTerrain(idLookup, t.Production.RequiredTerrain)
		}
		var productionTerrain Terrain
		if t.Production.ProductionTerrain != "" {
			productionTerrain = toTerrain(idLookup, t.Production.ProductionTerrain)
		}

		storage := make([]int, len(resource.Properties))
		for _, entry := range t.Storage {
			res, ok := resource.ResourceID(entry.Resource)
			if !ok {
				panic(fmt.Sprintf("unknown resource %s", entry.Resource))
			}
			storage[res] = entry.Amount
		}

		p := TerrainProps{
			Name:             t.Name,
			IsTerrain:        t.IsTerrain,
			IsPath:           t.IsPath,
			IsBuilding:       t.IsBuilding,
			IsWarehouse:      t.IsWarehouse,
			BuildOn:          toTerrains(idLookup, t.BuildOn...),
			TerrainBelow:     terrBelow,
			SelfConnectBelow: t.SelfConnectBelow,
			ConnectsTo:       toTerrains(idLookup, t.ConnectsTo...),
			CanBuild:         t.CanBuild,
			CanBuy:           t.CanBuy,
			Description:      t.Description,
			BuildCost:        cost,
			Storage:          storage,
			Production: Production{
				Resource:          prodRes,
				MaxProduction:     t.Production.MaxProduction,
				ConsumesResource:  consResource,
				ConsumesAmount:    t.Production.ConsumesAmount,
				RequiredTerrain:   requiredTerrain,
				ProductionTerrain: productionTerrain,
			},
		}

		if t.IsBuilding {
			Buildings.Set(Terrain(i))
		}
		if t.IsWarehouse {
			Warehouse = Terrain(i)
		}

		props = append(props, p)
	}

	randTerr := []Terrain{}

	for _, str := range propsHelper.RandomTerrains {
		randTerr = append(randTerr, toTerrain(idLookup, str))
	}

	if Warehouse == Air {
		panic("no warehouse building defined")
	}

	Air = toTerrain(idLookup, propsHelper.ZeroTerrain)
	Buildable = toTerrain(idLookup, propsHelper.Buildable)
	Default = toTerrain(idLookup, propsHelper.Default)

	RandomTerrain = randTerr
	Properties = props
}

func TerrainID(name string) (Terrain, bool) {
	t, ok := idLookup[name]
	return t, ok
}

type TerrainProps struct {
	Name             string
	IsTerrain        bool
	IsPath           bool
	IsBuilding       bool
	IsWarehouse      bool
	BuildOn          Terrains
	TerrainBelow     Terrain
	SelfConnectBelow bool
	ConnectsTo       Terrains
	CanBuild         bool
	CanBuy           bool
	Production       Production
	BuildCost        []ResourceAmount
	Storage          []int
	Description      string
}

type terrainPropsJs struct {
	Name             string             `json:"name"`
	IsTerrain        bool               `json:"is_terrain"`
	IsPath           bool               `json:"is_path"`
	IsBuilding       bool               `json:"is_building"`
	IsWarehouse      bool               `json:"is_warehouse"`
	BuildOn          []string           `json:"build_on,omitempty"`
	TerrainBelow     string             `json:"terrain_below"`
	SelfConnectBelow bool               `json:"self_connect_below"`
	ConnectsTo       []string           `json:"connects_to,omitempty"`
	CanBuild         bool               `json:"can_build"`
	CanBuy           bool               `json:"can_buy"`
	Production       productionJs       `json:"production"`
	BuildCost        []resourceAmountJs `json:"build_cost,omitempty"`
	Storage          []resourceAmountJs `json:"storage,omitempty"`
	Description      string             `json:"description,omitempty"`
}

type Production struct {
	Resource          resource.Resource
	MaxProduction     int
	ConsumesResource  resource.Resource
	ConsumesAmount    int
	RequiredTerrain   Terrain
	ProductionTerrain Terrain
}

type productionJs struct {
	Resource          string `json:"resource"`
	MaxProduction     int    `json:"max_production"`
	ConsumesResource  string `json:"consumes_resource"`
	ConsumesAmount    int    `json:"consumes_amount"`
	RequiredTerrain   string `json:"required_terrain"`
	ProductionTerrain string `json:"production_terrain"`
}

type ResourceAmount struct {
	Resource resource.Resource
	Amount   int
}

type resourceAmountJs struct {
	Resource string `json:"resource"`
	Amount   int    `json:"amount"`
}

type props struct {
	ZeroTerrain    string           `json:"zero_terrain"`
	Buildable      string           `json:"buildable"`
	Default        string           `json:"default"`
	RandomTerrains []string         `json:"random_terrains"`
	Terrains       []terrainPropsJs `json:"terrains"`
}

func toTerrains(lookup map[string]Terrain, terr ...string) Terrains {
	var ret Terrains
	for _, t := range terr {
		id, ok := lookup[t]
		if !ok {
			panic(fmt.Sprintf("unknown terrain %s", t))
		}
		ret.Set(id)
	}
	return ret
}

func toTerrain(lookup map[string]Terrain, t string) Terrain {
	id, ok := lookup[t]
	if !ok {
		panic(fmt.Sprintf("unknown terrain %s", t))
	}
	return id
}
