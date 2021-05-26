package lb

import (
	"math/rand"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
)

type RandBalance struct {
	Ctx *agw_context.AGWContext
}

func newRandBalance(ctx *agw_context.AGWContext) LoadBalance {
	return &RandBalance{Ctx: ctx}
}

func (rb *RandBalance) Select(serviceList []*model.Instance) *model.Instance {
	length := len(serviceList)
	if length == 0 {
		return nil
	}
	return serviceList[rand.Intn(length)]
}
