package dispatch

import (
	"net/http"
	"time"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/router"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// 路由匹配
	routeInfo, params := router.DefaultRouter().Match(r.Method, r.URL.Path)
	if routeInfo == nil {
		http.NotFound(w, r)
		return
	}
	// 创建请求处理上下文信息AGWContext
	ctx := agw_context.NewAGWContext(w, r, routeInfo, params, startTime)

	// 请求分发处理
	dispatcher := NewDispatcher()
	dispatcher.Dispatch(ctx)
}
