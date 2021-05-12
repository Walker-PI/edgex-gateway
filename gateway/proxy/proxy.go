package proxy

import (
	"net/http"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
)

type ProxyHandler struct {
}

func (p *ProxyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := agw_context.NewAGWContext(w, r)

	// Pre获取路由信息 装饰

	p.proxyHTTP(ctx)
}

// 反向代理http协议
func (p *ProxyHandler) proxyHTTP(ctx *agw_context.AGWContext) {
	// // 匹配路由
	// path := ctx.Request.URL.Path
	// remoteUrl, route, err := router.Match(path)
	// if err != nil {
	// 	logger.Error("route is: ", route, ", error is:", err)
	// 	return response.NewError(response.ProxyUrlNotFound, err.Error())
	// }
	// ctx.RemoteURL = remoteUrl

	// // 创建代理对象
	// proxy := &httputil.ReverseProxy{
	// 	Director: func(req *http.Request) {
	// 		// FIXME: 放在pre filter-decoration里处理
	// 		req.Host = remoteUrl.Host
	// 		req.URL.Scheme = remoteUrl.Scheme
	// 		req.URL.Host = remoteUrl.Host
	// 		req.URL.Path = remoteUrl.Path

	// 		if ctx.Request.URL.RawQuery == "" || req.URL.RawQuery == "" {
	// 			req.URL.RawQuery = ctx.Request.URL.RawQuery + req.URL.RawQuery
	// 		} else {
	// 			req.URL.RawQuery = ctx.Request.URL.RawQuery + "&" + req.URL.RawQuery
	// 		}
	// 		ctx.Request = req
	// 		logger.Debug("director: ", remoteUrl.String())
	// 	},
	// 	ModifyResponse: func(resp *http.Response) error {
	// 		logger.Debug("modify response:", remoteUrl.String())
	// 		ctx.Response = resp
	// 		err := filter.BeforeResponseFilter(ctx)
	// 		if err != nil {
	// 			logger.Info("BeforeResponse Stop")
	// 		}
	// 		return nil
	// 	},
	// 	ErrorHandler: r.ErrorHandler,
	// 	Transport:    defaultGatewayTransport.GetTransport(route),
	// }

	// proxy.ServeHTTP(ctx.ResponseWriter, ctx.Request)

	// logger.Info("请求完成, 请求地址：", path, "，目标地址：", remoteUrl)
}

// TODO：转换gRpc协议
// func (p *ProxyHandler) proxyGRPC() {

// }
