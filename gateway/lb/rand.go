package lb

import (
	"math/rand"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	consulapi "github.com/hashicorp/consul/api"
)

type RandBalance struct {
	Ctx *agw_context.AGWContext
}

func newRandBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &RandBalance{Ctx: ctx}
}

func (rb *RandBalance) Select(serviceList []*consulapi.CatalogService) *consulapi.CatalogService {
	length := len(serviceList)
	if length == 0 {
		return nil
	}
	return serviceList[rand.Intn(length)]
}
