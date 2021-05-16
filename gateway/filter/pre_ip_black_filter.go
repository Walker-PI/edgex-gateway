package filter

import (
	"bytes"
	"net"
	"net/http"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

type IPBlackFilter struct {
	baseFilter
}

func newIPBlackFilter() Filter {
	return &IPBlackFilter{}
}

func (f *IPBlackFilter) Name() FilterName {
	return PreIPBlackFilter
}

func (f *IPBlackFilter) Type() FilterType {
	return PreFilter
}

func (f *IPBlackFilter) Priority() int {
	return priority[PreIPBlackFilter]
}

func (f *IPBlackFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {

	if len(ctx.RouteDetail.IPBlackList) == 0 {
		return f.baseFilter.Run(ctx)
	}

	realIP, ok := ctx.Get("Real-IP").(string)
	if !ok {
		return http.StatusForbidden, nil
	}

	netIP := net.ParseIP(realIP)
	if netIP == nil {
		return http.StatusForbidden, nil
	}

	for _, blackIP := range ctx.RouteDetail.IPBlackList {
		if bytes.Compare(blackIP, netIP) == 0 {
			return http.StatusForbidden, nil
		}
	}
	return f.baseFilter.Run(ctx)
}
