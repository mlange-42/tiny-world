package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
)

// Production, accumulated globally.
type Production struct {
	// Current production, indexed by [resource.Resource].
	Prod []int
	// Current consumption, indexed by [resource.Resource].
	Cons []int
}

// NewProduction creates a new Production resource.
func NewProduction() Production {
	return Production{
		Prod: make([]int, len(resource.Properties)),
		Cons: make([]int, len(resource.Properties)),
	}
}

// Reset all production and consumption zo zero.
func (p *Production) Reset() {
	for i := range resource.Properties {
		p.Prod[i] = 0
		p.Cons[i] = 0
	}
}
