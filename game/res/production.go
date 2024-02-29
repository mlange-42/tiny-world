package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
)

type Production struct {
	Prod []int
	Cons []int
}

func NewProduction() Production {
	return Production{
		Prod: make([]int, len(resource.Properties)),
		Cons: make([]int, len(resource.Properties)),
	}
}

func (p *Production) Reset() {
	for i := range resource.Properties {
		p.Prod[i] = 0
		p.Cons[i] = 0
	}
}
