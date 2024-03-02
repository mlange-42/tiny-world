package terr

import (
	"fmt"
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
	"github.com/mlange-42/tiny-world/game/resource"
)

type TerrainBit uint16

const (
	IsTerrain TerrainBit = iota
	IsPath
	IsBridge
	IsBuilding
	IsWarehouse
	CanBuild
	CanBuy
)

type TerrainBits uint16

func NewTerrainBits(bits ...TerrainBit) TerrainBits {
	d := TerrainBits(0)
	for _, dir := range bits {
		d |= (1 << dir)
	}
	return d
}

func (d *TerrainBits) Set(dir TerrainBit) {
	*d |= (1 << dir)
}

// Contains checks whether all the argument's bits are contained in this Subscription.
func (d TerrainBits) Contains(dir TerrainBit) bool {
	bits := TerrainBits(1 << dir)
	return (bits & d) == bits
}

type Terrain uint8

var Air Terrain
var Buildable Terrain
var Default Terrain
var Warehouse Terrain

type Terrains uint32

var Buildings Terrains
var Paths Terrains

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
		if t.Consumption.Amount > 0 {
			var ok bool
			consResource, ok = resource.ResourceID(t.Consumption.Resource)
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
		var productionTerrain Terrains
		if len(t.Production.ProductionTerrain) > 0 {
			productionTerrain = toTerrains(idLookup, t.Production.ProductionTerrain...)
		}

		storage := make([]int, len(resource.Properties))
		for _, entry := range t.Storage {
			res, ok := resource.ResourceID(entry.Resource)
			if !ok {
				panic(fmt.Sprintf("unknown resource %s", entry.Resource))
			}
			storage[res] = entry.Amount
		}

		bits := TerrainBits(0)
		if t.IsTerrain {
			bits.Set(IsTerrain)
		}
		if t.IsPath || t.IsBridge {
			bits.Set(IsPath)
		}
		if t.IsBridge {
			bits.Set(IsBridge)
		}
		if t.IsBuilding {
			bits.Set(IsBuilding)
		}
		if t.IsWarehouse {
			bits.Set(IsWarehouse)
		}
		if t.CanBuild {
			bits.Set(CanBuild)
		}
		if t.CanBuy {
			bits.Set(CanBuy)
		}

		p := TerrainProps{
			Name:         t.Name,
			TerrainBits:  bits,
			BuildOn:      toTerrains(idLookup, t.BuildOn...),
			TerrainBelow: terrBelow,
			ConnectsTo:   toTerrains(idLookup, t.ConnectsTo...),
			BuildRadius:  t.BuildRadius,
			Description:  t.Description,
			BuildCost:    cost,
			Storage:      storage,
			Production: Production{
				Resource:          prodRes,
				MaxProduction:     t.Production.MaxProduction,
				RequiredTerrain:   requiredTerrain,
				ProductionTerrain: productionTerrain,
				HaulCapacity:      t.Production.HaulCapacity,
			},
			Consumption: Consumption{
				Resource: consResource,
				Amount:   t.Consumption.Amount,
			},
		}

		if p.TerrainBits.Contains(IsBuilding) {
			Buildings.Set(Terrain(i))
		}
		if p.TerrainBits.Contains(IsPath) {
			Paths.Set(Terrain(i))
		}
		if p.TerrainBits.Contains(IsWarehouse) {
			Warehouse = Terrain(i)
		}

		props = append(props, p)
	}

	if Warehouse == Air {
		panic("no warehouse building defined")
	}

	Air = toTerrain(idLookup, propsHelper.ZeroTerrain)
	Buildable = toTerrain(idLookup, propsHelper.Buildable)
	Default = toTerrain(idLookup, propsHelper.Default)

	Properties = props
}

func TerrainID(name string) (Terrain, bool) {
	t, ok := idLookup[name]
	return t, ok
}

type TerrainProps struct {
	Name         string
	BuildOn      Terrains
	ConnectsTo   Terrains
	TerrainBits  TerrainBits
	TerrainBelow Terrain
	BuildRadius  uint8
	Description  string
	BuildCost    []ResourceAmount
	Storage      []int
	Production   Production
	Consumption  Consumption
}

type terrainPropsJs struct {
	Name         string             `json:"name"`
	IsTerrain    bool               `json:"is_terrain"`
	IsPath       bool               `json:"is_path"`
	IsBridge     bool               `json:"is_bridge"`
	IsBuilding   bool               `json:"is_building"`
	IsWarehouse  bool               `json:"is_warehouse"`
	BuildRadius  uint8              `json:"build_radius"`
	BuildOn      []string           `json:"build_on,omitempty"`
	TerrainBelow string             `json:"terrain_below"`
	ConnectsTo   []string           `json:"connects_to,omitempty"`
	CanBuild     bool               `json:"can_build"`
	CanBuy       bool               `json:"can_buy"`
	Production   productionJs       `json:"production"`
	Consumption  consumptionJs      `json:"consumption"`
	BuildCost    []resourceAmountJs `json:"build_cost,omitempty"`
	Storage      []resourceAmountJs `json:"storage,omitempty"`
	Description  string             `json:"description"`
}

type Production struct {
	Resource          resource.Resource
	MaxProduction     int
	RequiredTerrain   Terrain
	ProductionTerrain Terrains
	HaulCapacity      int
}

type productionJs struct {
	Resource          string   `json:"resource"`
	MaxProduction     int      `json:"max_production"`
	RequiredTerrain   string   `json:"required_terrain"`
	ProductionTerrain []string `json:"production_terrain"`
	HaulCapacity      int      `json:"haul_capacity"`
}

type Consumption struct {
	Resource resource.Resource
	Amount   int
}

type consumptionJs struct {
	Resource string `json:"resource"`
	Amount   int    `json:"amount"`
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
	ZeroTerrain string           `json:"zero_terrain"`
	Buildable   string           `json:"buildable"`
	Default     string           `json:"default"`
	Terrains    []terrainPropsJs `json:"terrains"`
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
