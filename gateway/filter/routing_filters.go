package filter

import "github.com/Walker-PI/edgex-gateway/gateway/agw_context"

// 负载均衡 & 熔断
type RibbonFilter struct {
}

func NewRibbonFilter() *RibbonFilter {
	return &RibbonFilter{}
}

func (f *RibbonFilter) FilterType() string {
	return "pre"
}

func (f *RibbonFilter) FilterOrder() int {
	return -9
}

func (f *RibbonFilter) Run(ctx *agw_context.AGWContext) error {
	// 负载均衡 & 熔断
	// TODO:

	return nil
}

// http_client
// 负载均衡 & 熔断
type HttpProxyFilter struct {
}

func NewHttpProxyFilter() *HttpProxyFilter {
	return &HttpProxyFilter{}
}

func (f *HttpProxyFilter) FilterType() string {
	return "pre"
}

func (f *HttpProxyFilter) FilterOrder() int {
	return -9
}

func (f *HttpProxyFilter) Run(ctx *agw_context.AGWContext) error {
	// 请求后端服务
	// TODO:

	return nil
}
