package resource

import (
	"io/fs"

	"github.com/mlange-42/tiny-world/cmd/util"
)

type Resource uint

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
