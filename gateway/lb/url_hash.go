package lb

import (
	"hash/fnv"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
)

type URLHashBalance struct {
	Ctx *agw_context.AGWContext
}

func newURLHashBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &URLHashBalance{Ctx: ctx}
}

func (urlHash *URLHashBalance) Select(serviceList []*model.Instance) *model.Instance {
	length := len(serviceList)
	if length == 0 {
		return nil
	}
	path := urlHash.Ctx.OriginRequest.URL.Path
	h := fnv.New32a()
	_, _ = h.Write([]byte(path))
	return serviceList[h.Sum32()%uint32(length)]
}
