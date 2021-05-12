package pre

import "github.com/Walker-PI/edgex-gateway/gateway/agw_context"

// AuthFilter 鉴权 jwt (可扩展)
type AuthFilter struct {
}

func NewAuthFilter() *AuthFilter {
	return &AuthFilter{}
}

func (f *AuthFilter) FilterType() string {
	return "pre"
}

func (f *AuthFilter) FilterOrder() int {
	return -9
}

func (f *AuthFilter) Run(ctx *agw_context.AGWContext) error {
	// 根据鉴权类型确定使用什么鉴权
	// TODO:

	return nil
}
