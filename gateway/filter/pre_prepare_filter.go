package filter

import (
	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	"github.com/Walker-PI/edgex-gateway/pkg/tools"
)

type PrepareFilter struct {
	baseFilter
}

func newPrepareFilter() Filter {
	return &PrepareFilter{}
}

func (f *PrepareFilter) Init(configFile string) error {
	return f.baseFilter.Init(configFile)
}

func (f *PrepareFilter) Type() FilterType {
	return PrePrepareFilter
}

func (f *PrepareFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	realIP := tools.ClientPublicIP(ctx.OriginRequest)
	if realIP == "" {
		realIP = tools.ClientIP(ctx.OriginRequest)
	}
	if realIP == "" {
		realIP = tools.RemoteIP(ctx.OriginRequest)
	}
	return f.baseFilter.Run(ctx)
}
