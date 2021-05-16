package filter

import (
	"bytes"
	"net"
	"net/http"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
)

type IPWhiteFilter struct {
	baseFilter
}

func newIPWhiteFilter() Filter {
	return &IPWhiteFilter{}
}

func (f *IPWhiteFilter) Name() FilterName {
	return PreIPWhiteFilter
}

func (f *IPWhiteFilter) Type() FilterType {
	return PreFilter
}

func (f *IPWhiteFilter) Priority() int {
	return priority[PreIPWhiteFilter]
}

func (f *IPWhiteFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	if len(ctx.RouteDetail.IPWhiteList) == 0 {
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

	for _, whiteIP := range ctx.RouteDetail.IPWhiteList {
		if bytes.Compare(whiteIP, netIP) != 0 {
			return f.baseFilter.Run(ctx)
		}
	}
	return http.StatusForbidden, nil
}
