package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	"github.com/Walker-PI/edgex-gateway/gateway/filter"
	// "github.com/Walker-PI/edgex-gateway/pkg/logger"
)

type Proxy struct {
	StartTime      time.Time
	EndTime        time.Time
	Ctx            *agw_context.AGWContext
	Filters        map[string][]*filter.Filter
	Director       func(req *http.Request)
	ModifyResponse func(*http.Response) error
	ErrorHandler   func(http.ResponseWriter, *http.Request, error)
}

func NewProxy(ctx *agw_context.AGWContext) *Proxy {
	proxy := &Proxy{
		StartTime: time.Now(),
		Ctx:       ctx,
	}
	proxy.initFilters()
	return proxy
}

func (p *Proxy) initFilters() {
	// TODO:
}

func (p *Proxy) buildDirector() {
	origin, _ := url.Parse("http://106.14.157.113:6789")
	p.Director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.Host = origin.Host
		req.URL.Scheme = origin.Scheme
		req.URL.Host = origin.Host
		req.URL.Path = p.Ctx.OriginRequest.URL.Path

		if p.Ctx.OriginRequest.URL.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = p.Ctx.OriginRequest.URL.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = p.Ctx.OriginRequest.URL.RawQuery + "&" + req.URL.RawQuery
		}
		p.Ctx.ForwardRequest = req
	}
}

func (p *Proxy) buildModifyResponse() {
	// TODO:
}

func (p *Proxy) buildErrorHandler() {
	// TODO:
}

func (p *Proxy) ServeHTTP() {
	proxy := &httputil.ReverseProxy{
		Director:       p.Director,
		ModifyResponse: p.ModifyResponse,
		ErrorHandler:   p.ErrorHandler,
	}
	proxy.ServeHTTP(p.Ctx.ResponseWriter, p.Ctx.ForwardRequest)
}
