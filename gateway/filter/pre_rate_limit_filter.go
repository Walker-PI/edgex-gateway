package filter

import (
	"net/http"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
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
	if ctx.RouteInfo.Limiter == nil {
		return f.baseFilter.Run(ctx)
	}
	beforeAvailable := ctx.RouteInfo.Limiter.Available()
	if !ctx.RouteInfo.Limiter.Do(1) {
		logger.Warn("[RateLimitFilter-Run] rate limit: before_available=%v", beforeAvailable)
		return http.StatusServiceUnavailable, nil
	}
	logger.Info("[RateLimitFilter-Run] before_available=%v, after_available=%v", beforeAvailable,
		ctx.RouteInfo.Limiter.Available())
	return f.baseFilter.Run(ctx)
}
