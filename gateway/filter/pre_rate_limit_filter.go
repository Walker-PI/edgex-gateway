package filter

import (
	"net/http"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

type RateLimitFilter struct {
	baseFilter
}

func newRateLimitFilter() Filter {
	return &RateLimitFilter{}
}

func (f *RateLimitFilter) Name() FilterName {
	return PreRateLimitFilter
}

func (f *RateLimitFilter) Type() FilterType {
	return PreFilter
}

func (f *RateLimitFilter) Priority() int {
	return priority[PreRateLimitFilter]
}

func (f *RateLimitFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	if ctx.RouteDetail.Limiter == nil || ctx.RouteDetail.Limiter.Do(1) {
		return f.baseFilter.Run(ctx)
	}
	return http.StatusServiceUnavailable, nil
}
