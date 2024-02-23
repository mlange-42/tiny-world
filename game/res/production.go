package res

import "github.com/mlange-42/tiny-world/game/resource"

type Production struct {
	Res [resource.EndResources]int
}

func (p *Production) Reset() {
	for i := resource.Resource(0); i < resource.EndResources; i++ {
		p.Res[i] = 0
	}
}

type Stock struct {
	Res [resource.EndResources]int
}
