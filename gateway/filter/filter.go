package filter

import (
	"net/http"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
)

type FilterType string

const (
	PreFilter   = "PRE"
	PostFilter  = "POST"
	RouteFilter = "ROUTING"
)

type Filter interface {
	Name() FilterName
	Type() FilterType
	Priority() int
	Run(*agw_context.AGWContext) (code int, err error)
}

type baseFilter struct{}

func (f baseFilter) Run(ctx *agw_context.AGWContext) (code int, err error) {
	return http.StatusOK, nil
}
