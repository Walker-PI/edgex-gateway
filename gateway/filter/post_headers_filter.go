package filter

import (
	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

type PostHeadersFilter struct {
	baseFilter
}

func newPostHeadersFilter() Filter {
	return &PostHeadersFilter{}
}

func (f *PostHeadersFilter) Name() FilterName {
	return PostHeadersFilter_
}

func (f *PostHeadersFilter) Type() FilterType {
	return PostFilter
}

func (f *PostHeadersFilter) Priority() int {
	return priority[PostHeadersFilter_]
}

func (f *PostHeadersFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {

	for _, header := range hopHeaders {
		ctx.ForwardRequest.Header.Del(header)
	}
	return f.baseFilter.Run(ctx)
}
