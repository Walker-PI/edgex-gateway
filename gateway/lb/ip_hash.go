package lb

import (
	"hash/fnv"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	"github.com/Walker-PI/edgex-gateway/pkg/tools"
	consulapi "github.com/hashicorp/consul/api"
)

type IPHashBalance struct {
	Ctx *agw_context.AGWContext
}

func newIPHashBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &IPHashBalance{Ctx: ctx}
}

func (ipHash *IPHashBalance) Select(serviceList []*consulapi.CatalogService) *consulapi.CatalogService {
	length := len(serviceList)
	if length == 0 {
		return nil
	}
	realIP := tools.RealIP(ipHash.Ctx.OriginRequest)
	h := fnv.New32a()
	h.Write([]byte(realIP))
	return serviceList[h.Sum32()%uint32(length)]
}
