package agw_context

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Walker-PI/iot-gateway/gateway/router"
)

type AGWContext struct {
	mutex *sync.RWMutex
	Keys  map[string]interface{}

	ResponseWriter http.ResponseWriter
	OriginRequest  *http.Request
	ForwardRequest *http.Request
	Response       *http.Response

	RouteInfo  *router.RouteInfo
	PathParams map[string]string

	TargetURL *url.URL

	StartTime time.Time
}

func NewAGWContext(writer http.ResponseWriter, req *http.Request, routeInfo *router.RouteInfo, pathParams map[string]string, startTime time.Time) *AGWContext {
	return &AGWContext{
		ResponseWriter: writer,
		OriginRequest:  req,
		ForwardRequest: copyRequest(req),
		Keys:           make(map[string]interface{}),
		mutex:          &sync.RWMutex{},
		RouteInfo:      routeInfo,
		PathParams:     pathParams,
		StartTime:      startTime,
	}
}

func copyRequest(req *http.Request) *http.Request {
	return req.Clone(req.Context())
}

func (ctx *AGWContext) Get(key string) (value interface{}, exists bool) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	value, exists = ctx.Keys[key]
	return
}

func (ctx *AGWContext) GetString(key string) (s string) {
	if val, ok := ctx.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (ctx *AGWContext) Set(key string, value interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.Keys[key] = value
}

func (ctx *AGWContext) GetPathParam(key string) (value string, exists bool) {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	value, exists = ctx.PathParams[key]
	return
}

func (ctx *AGWContext) Deadline() (deadline time.Time, ok bool) {
	return ctx.ForwardRequest.Context().Deadline()
}

func (ctx *AGWContext) Done() <-chan struct{} {
	return ctx.ForwardRequest.Context().Done()
}

func (ctx *AGWContext) Err() error {
	return ctx.ForwardRequest.Context().Err()
}

func (ctx *AGWContext) Value(key interface{}) interface{} {
	if key == 0 {
		return ctx.ForwardRequest
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := ctx.Get(keyAsString)
		return val
	}
	return nil
}
