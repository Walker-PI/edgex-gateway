package agw_context

import (
	"net/http"
	"sync"

	"github.com/Walker-PI/edgex-gateway/gateway/router"
)

type AGWContext struct {
	mutex          *sync.RWMutex
	data           map[string]interface{}
	ResponseWriter http.ResponseWriter
	OriginRequest  *http.Request
	ForwardRequest *http.Request
	Response       *http.Response

	RouteDetail *router.RouterInfo
	ParamsMap   map[string]string
}

func NewAGWContext(writer http.ResponseWriter, req *http.Request) *AGWContext {
	return &AGWContext{
		ResponseWriter: writer,
		OriginRequest:  req,
		ForwardRequest: copyRequest(req),
		data:           make(map[string]interface{}),
		mutex:          &sync.RWMutex{},
	}
}

func copyRequest(req *http.Request) *http.Request {
	return req.Clone(req.Context())
}

func (ctx *AGWContext) Get(key string) interface{} {
	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()
	value, _ := ctx.data[key]
	return value
}

func (ctx *AGWContext) Set(key string, value interface{}) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.data[key] = value
}
