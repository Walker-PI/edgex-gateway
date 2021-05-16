package lb

import (
	"sync/atomic"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	consulapi "github.com/hashicorp/consul/api"
)

type RRBalance struct {
	Ctx    *agw_context.AGWContext
	Option *int64
}

func newRRBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &RRBalance{
		Ctx:    ctx,
		Option: new(int64),
	}
}

func (rr *RRBalance) Select(serviceList []*consulapi.CatalogService) *consulapi.CatalogService {
	length := int64(len(serviceList))
	if length == 0 {
		return nil
	}
	return serviceList[atomic.AddInt64(rr.Option, 1)%length]
}
