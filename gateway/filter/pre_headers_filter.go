package filter

import (
	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
}

type PreHeadersFilter struct {
	baseFilter
}

func newPreHeadersFilter() Filter {
	return &PreHeadersFilter{}
}

func (f *PreHeadersFilter) Name() FilterName {
	return PreHeadersFilter_
}

func (f *PreHeadersFilter) Type() FilterType {
	return PreFilter
}

func (f *PreHeadersFilter) Priority() int {
	return priority[PreHeadersFilter_]
}

func (f *PreHeadersFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	realIP := ctx.Get("Real-IP").(string)
	ctx.ForwardRequest.Header.Set("X-Forwarded-For", realIP)
	ctx.ForwardRequest.Header.Set("X-Real-IP", realIP)
	for _, header := range hopHeaders {
		ctx.ForwardRequest.Header.Del(header)
	}
	return f.baseFilter.Run(ctx)
}
