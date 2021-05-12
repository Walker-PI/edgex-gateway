package filter

import "github.com/Walker-PI/edgex-gateway/gateway/agw_context"

// DecorationFilter  路由配置信息 & Header之类 X-Forwarded-Host
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
