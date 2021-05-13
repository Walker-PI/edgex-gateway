package proxy

import (
	"net/http"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := agw_context.NewAGWContext(w, r)
	proxy := NewProxy(ctx)
	proxy.buildDirector()
	proxy.buildModifyResponse()
	proxy.buildErrorHandler()
	proxy.ServeHTTP()
}
