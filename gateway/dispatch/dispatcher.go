package dispatch

import (
	"net/http"
	"sort"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/filter"
)

type Dispatcher struct {
	PreFilters    []filter.Filter
	RoutingFilter filter.Filter
	PostFilters   []filter.Filter
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		PreFilters:  make([]filter.Filter, 0),
		PostFilters: make([]filter.Filter, 0),
	}
	dispatcher.initFilters()
	return dispatcher
}

func (d *Dispatcher) initFilters() {
	d.PreFilters = append(d.PreFilters,
		filter.NewFilter(filter.PrePrepareFilter),
		filter.NewFilter(filter.PreHeadersFilterBefore),
		filter.NewFilter(filter.PreIPWhiteFilter),
		filter.NewFilter(filter.PreIPBlackFilter),
		filter.NewFilter(filter.PreAuthFilter),
		filter.NewFilter(filter.PreRateLimitFilter),
	)
	sort.SliceStable(d.PreFilters, func(i, j int) bool {
		return d.PreFilters[i].Priority() < d.PreFilters[j].Priority()
	})

	d.RoutingFilter = filter.NewFilter(filter.RoutingFilterName)

	d.PostFilters = append(d.PostFilters,
		filter.NewFilter(filter.PostHeadersFilterAfter),
	)
	sort.SliceStable(d.PostFilters, func(i, j int) bool {
		return d.PostFilters[i].Priority() < d.PostFilters[j].Priority()
	})

}

func (d *Dispatcher) DoPreFilters(ctx *agw_context.AGWContext) (statusCode int, err error) {
	for _, v := range d.PreFilters {
		statusCode, err = v.Run(ctx)
		if err != nil || statusCode != http.StatusOK {
			return
		}
	}
	return http.StatusOK, nil
}

func (d *Dispatcher) DoRoutingFilter(ctx *agw_context.AGWContext) (statusCode int, err error) {
	statusCode, err = d.RoutingFilter.Run(ctx)
	if err != nil || statusCode != http.StatusOK {
		return
	}
	return http.StatusOK, nil
}

func (d *Dispatcher) DoPostFilters(ctx *agw_context.AGWContext) (statusCode int, err error) {
	for _, v := range d.PostFilters {
		statusCode, err = v.Run(ctx)
		if err != nil || statusCode != http.StatusOK {
			return
		}
	}
	return http.StatusOK, nil
}
