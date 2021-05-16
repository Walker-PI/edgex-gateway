package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"sort"
	"time"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	"github.com/Walker-PI/edgex-gateway/gateway/filter"
	"github.com/Walker-PI/edgex-gateway/gateway/router"
)

type Proxy struct {
	StartTime time.Time
	Ctx       *agw_context.AGWContext
	Filters   []filter.Filter
}

func NewProxy(w http.ResponseWriter, r *http.Request) *Proxy {
	proxy := &Proxy{
		StartTime: time.Now(),
		Ctx:       agw_context.NewAGWContext(w, r),
	}
	proxy.initFilters()
	return proxy
}

func (p *Proxy) initFilters() {
	p.Filters = make([]filter.Filter, 0)
	p.Filters = append(p.Filters,
		filter.NewFilter(filter.PrePrepareFilter),
		filter.NewFilter(filter.PreHeadersFilter_),
		filter.NewFilter(filter.PreIPWhiteFilter),
		filter.NewFilter(filter.PreIPBlackFilter),
		filter.NewFilter(filter.PreAuthFilter),
		filter.NewFilter(filter.PreRateLimitFilter),
		filter.NewFilter(filter.PostHeadersFilter_),
	)
	sort.SliceStable(p.Filters, func(i, j int) bool {
		return p.Filters[i].Priority() < p.Filters[j].Priority()
	})
}

func (p *Proxy) DoPreFilters() (statusCode int, err error) {
	for _, v := range p.Filters {
		if v.Type() == filter.PreFilter {
			statusCode, err = v.Run(p.Ctx)
			if err != nil || statusCode != http.StatusOK {
				return
			}
		}
	}
	return http.StatusOK, nil
}

func (p *Proxy) DoPostFilters() (statusCode int, err error) {
	for _, v := range p.Filters {
		if v.Type() == filter.PostFilter {
			statusCode, err = v.Run(p.Ctx)
			if err != nil || statusCode != http.StatusOK {
				return
			}
		}
	}
	return http.StatusOK, nil
}

func (p *Proxy) DoProxy(director func(*http.Request), modifyResponse func(*http.Response) error,
	errorHandler func(http.ResponseWriter, *http.Request, error)) {

	// 超时时间
	timeout := 30 * time.Second
	if p.Ctx.RouteDetail.Target.Timeout > 0 {
		timeout = p.Ctx.RouteDetail.Target.Timeout * time.Millisecond
	}
	proxy := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
		ErrorHandler:   errorHandler,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	proxy.ServeHTTP(p.Ctx.ResponseWriter, p.Ctx.ForwardRequest)
}

func (p *Proxy) MatchRoute() bool {
	route := router.DefaultRouter()
	p.Ctx.RouteDetail, p.Ctx.ParamsMap = route.Match(p.Ctx.OriginRequest.Method, p.Ctx.OriginRequest.URL.Path)
	return p.Ctx.RouteDetail != nil
}

func (p *Proxy) ErrorHandle(w http.ResponseWriter, statusCode int) {
	status := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	p.Ctx.Response = &http.Response{
		Status:     status,
		StatusCode: statusCode,
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, status)
}
