package lb

import (
	"hash/fnv"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

type IPHashBalance struct {
	Ctx *agw_context.AGWContext
}

func newIPHashBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &IPHashBalance{Ctx: ctx}
}

func (ipHash *IPHashBalance) Select(serviceList []*model.Instance) *model.Instance {
	length := len(serviceList)
	if length == 0 {
		return nil
	}
	realIP := tools.RealIP(ipHash.Ctx.OriginRequest)
	h := fnv.New32a()
	_, _ = h.Write([]byte(realIP))
	return serviceList[h.Sum32()%uint32(length)]
}
