package filter

import (
	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

type PrepareFilter struct {
	baseFilter
}

func newPrepareFilter() Filter {
	return &PrepareFilter{}
}

func (f *PrepareFilter) Name() FilterName {
	return PrePrepareFilter
}

func (f *PrepareFilter) Type() FilterType {
	return PreFilter
}

func (f *PrepareFilter) Priority() int {
	return priority[PrePrepareFilter]
}

func (f *PrepareFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	realIP := tools.RealIP(ctx.ForwardRequest)
	ctx.Set("Real-IP", realIP)

	ctx.ForwardRequest.Header.Set("X-Forwarded-For", realIP)
	ctx.ForwardRequest.Header.Set("X-Real-IP", realIP)

	return f.baseFilter.Run(ctx)
}
