package filter

import (
	"net/http"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

type Filter interface {
	Init(configFile string) error
	Type() FilterType
	Run(*agw_context.AGWContext) (code int, err error)
}

type baseFilter struct{}

func (f baseFilter) Init(configFile string) error {
	return nil
}
func (f baseFilter) Run(ctx *agw_context.AGWContext) (code int, err error) {
	return http.StatusOK, nil
}
