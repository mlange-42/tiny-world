package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Stock struct {
	Cap []int
	Res []int
}

func NewStock(initial []int) Stock {
	if len(initial) != len(resource.Properties) {
		panic("initial resources don't match number of actual resources")
	}
	return Stock{
		Cap: make([]int, len(resource.Properties)),
		Res: initial,
	}
}

func (s *Stock) CanPay(cost []terr.ResourceAmount) bool {
	for _, c := range cost {
		if s.Res[c.Resource] < c.Amount {
			return false
		}
	}
	return true
}

func (s *Stock) Pay(cost []terr.ResourceAmount) {
	for _, c := range cost {
		s.Res[c.Resource] -= c.Amount
	}
}
