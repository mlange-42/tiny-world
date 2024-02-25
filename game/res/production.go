package res

import (
	"github.com/mlange-42/tiny-world/game/resource"
	"github.com/mlange-42/tiny-world/game/terr"
)

type Production struct {
	Prod [resource.EndResources]int
	Cons [resource.EndResources]int
}

func (p *Production) Reset() {
	for i := resource.Resource(0); i < resource.EndResources; i++ {
		p.Prod[i] = 0
		p.Cons[i] = 0
	}
}

type Stock struct {
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
