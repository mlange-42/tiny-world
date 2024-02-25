package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Stock struct {
	Cap [resource.EndResources]int
	Res [resource.EndResources]int
}

func (s *Stock) CanPay(cost []terr.BuildCost) bool {
	for _, c := range cost {
		if s.Res[c.Type] < c.Amount {
			return false
		}
	}
	return true
}

func (s *Stock) Pay(cost []terr.BuildCost) {
	for _, c := range cost {
		s.Res[c.Type] -= c.Amount
	}
}
