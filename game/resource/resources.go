package resource

import (
	"fmt"
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
)

type Resource uint8

type Resources uint32

func (d *Resources) Set(dir Resource) {
	*d |= (1 << dir)
}

func (d *Resources) Unset(dir Resource) {
	*d &= ^(1 << dir)
}

// Contains checks whether all the argument's bits are contained in this mask.
func (d Resources) Contains(dir Resource) bool {
	bits := Resources(1 << dir)
	return (bits & d) == bits
}

type ResourceProps struct {
	Name  string `json:"name"`
	Short string `json:"short"`
}

var Properties []ResourceProps

func ResourceID(name string) (Resource, bool) {
	t, ok := idLookup[name]
	return t, ok
}

var idLookup map[string]Resource

func Prepare(f fs.FS, file string) {
	props := []ResourceProps{}
	err := util.FromJsonFs(f, file, &props)
	if err != nil {
		panic(err)
	}

	idLookup = map[string]Resource{}
	for i, t := range props {
		idLookup[t.Name] = Resource(i)
	}
	Properties = props
}

func ToResources(res ...string) Resources {
	var ret Resources
	for _, r := range res {
		id, ok := idLookup[r]
		if !ok {
			panic(fmt.Sprintf("unknown resource '%s'", r))
		}
		ret.Set(id)
	}
	return ret
}

func ToResource(r string) Resource {
	id, ok := idLookup[r]
	if !ok {
		panic(fmt.Sprintf("unknown resource '%s'", r))
	}
	return id
}
