package filter

import "github.com/Walker-PI/edgex-gateway/gateway/agw_context"

type Filter interface {
	FilterType() string
	FilterOrder() int
	Run(*agw_context.AGWContext) error
}
