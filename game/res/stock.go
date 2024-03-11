package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

// Stock resource, holding global stock information.
type Stock struct {
	// Total storage capacity, indexed by [resource.Resource].
	Cap []int
	// Total storage, indexed by [resource.Resource].
	Res []int

	// Total population.
	Population int
	// Total population limit.
	MaxPopulation int
}

// NewStock creates a new Stock resource with the given initial resources.
func NewStock(initial []int) Stock {
	if len(initial) != len(resource.Properties) {
		panic("initial resources don't match number of actual resources")
	}
	return Stock{
		Cap: make([]int, len(resource.Properties)),
		Res: initial,
	}
}

// CanPay checks whether there are sufficient resources in the stock to pay the given amounts.
func (s *Stock) CanPay(cost []terr.ResourceAmount) bool {
	for _, c := range cost {
		if s.Res[c.Resource] < int(c.Amount) {
			return false
		}
	}
	return true
}

// Pay the given amounts be subtracting them from the stock.
func (s *Stock) Pay(cost []terr.ResourceAmount) {
	for _, c := range cost {
		s.Res[c.Resource] -= int(c.Amount)
	}
}
