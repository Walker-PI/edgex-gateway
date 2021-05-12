package pre

import "github.com/Walker-PI/edgex-gateway/gateway/agw_context"

type DecorationFilter struct {
}

func NewDecorationFilter() *DecorationFilter {
	return &DecorationFilter{}
}

func (f *DecorationFilter) FilterType() string {
	return "pre"
}

func (f *DecorationFilter) FilterOrder() int {
	return -10
}

func (f *DecorationFilter) Run(ctx *agw_context.AGWContext) error {
	// 查找路由
	// TODO:
	return nil
}
