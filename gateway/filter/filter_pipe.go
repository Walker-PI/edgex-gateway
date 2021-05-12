package filter

import (
	"sync"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	// "github.com/Walker-PI/edgex-gateway/pkg/logger"
)

type FilterPipe struct {
	AGWContext *agw_context.AGWContext
	filters    []Filter
	mutex      sync.RWMutex
}

func NewFilterPipe(c *agw_context.AGWContext) *FilterPipe {
	return &FilterPipe{
		AGWContext: c,
	}
}

func (p *FilterPipe) RegisterFilter() {

	p.addFilter(
		NewDecorationFilter(),
		NewAuthFilter(),
	)

}

func (p *FilterPipe) addFilter(filters ...Filter) {

	if p.filters == nil {
		p.filters = make([]Filter, 0)
	}
	p.filters = append(p.filters, filters...)
}

func (p *FilterPipe) DoPreFilters() {

}

func (p *FilterPipe) DoRoutingFilters() {

}

func (p *FilterPipe) DoPostFilters() {

}

func (p *FilterPipe) DoErrorFilters() {

}

func (p *FilterPipe) RunFilters() {
	// for _, step := range p.filters {
	// 	if !step.ParseParams() {
	// 		logger.Warn("%s invalid params", step.ToString())
	// 	} else {
	// 		if err := step.Process(); err != nil {
	// 			logger.Error("%s error: %v", step.ToString(), err)
	// 		} else {
	// 			logger.Info("%s finished", step.ToString())
	// 		}
	// 	}
	// }
}
