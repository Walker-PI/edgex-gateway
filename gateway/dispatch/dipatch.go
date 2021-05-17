package dispatch

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Walker-PI/iot-gateway/gateway/discovery"
	"github.com/Walker-PI/iot-gateway/gateway/proxy"
	"github.com/Walker-PI/iot-gateway/gateway/router"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/metric"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

func Dsipatch(w http.ResponseWriter, r *http.Request) {

	logger.Info("[Dispatch-Request] url=%v, body=%+v", tools.GetMarshalStr(r.URL), tools.GetBodyStr(r))

	proxy := proxy.NewProxy(w, r)

	defer func() {
		if proxy.Ctx.Response == nil {
			proxy.ErrorHandle(proxy.Ctx.ResponseWriter, http.StatusBadGateway)
		}
		// 异步监控上报
		metric.AsyncStatusEmit(proxy.StartTime, proxy.Ctx.OriginRequest, proxy.Ctx.Response)
		logger.Info("[Dispatch-Response] cost=%dms, status=%v, resp=%+v", time.Now().Sub(proxy.StartTime).Milliseconds(),
			proxy.Ctx.Response.StatusCode, proxy.Ctx.Response)
	}()

	// Step1 匹配路由
	match := proxy.MatchRoute()
	if !match {
		proxy.ErrorHandle(proxy.Ctx.ResponseWriter, http.StatusNotFound)
		return
	}

	// Step2 doPreFilters
	statusCode, err := proxy.DoPreFilters()
	if err != nil {
		logger.Error("[Dispatch-DoPreFilters] failed: err=%v", err)
		proxy.ErrorHandle(proxy.Ctx.ResponseWriter, statusCode)
		return
	}
	if statusCode != http.StatusOK {
		logger.Warn("[Dispatch-DoPreFilters] emerge execption: statusCode=%v", statusCode)
		proxy.ErrorHandle(proxy.Ctx.ResponseWriter, statusCode)
		return
	}

	// Step3 服务发现
	var (
		targetMode int32
		targetHost string
		targetPath string
		target     = proxy.Ctx.RouteDetail.Target
	)
	switch target.Mode {
	case router.ConsulTargetMode:
		// 服务发现
		service, err := discovery.GetInstance(proxy.Ctx, target.ServiceName, target.LoadBalance)
		if err != nil || service == nil {
			logger.Error("[Dispatch] discovery.GetInstance falild: serviceName=%s, service=%+v, err=%v",
				target.ServiceName, service, err)
			proxy.ErrorHandle(proxy.Ctx.ResponseWriter, http.StatusNotFound)
			return
		}
		targetMode = router.ConsulTargetMode
		targetHost = service.Address + ":" + strconv.Itoa(service.ServicePort)
		targetPath = proxy.Ctx.OriginRequest.URL.Path
		if target.StripPrefix {
			targetPath = strings.TrimPrefix(proxy.Ctx.OriginRequest.URL.Path, proxy.Ctx.RouteDetail.Pattern)
		}
	default:
		targetMode = router.DefaultTargetMode
		targetHost = target.URL.Host
		targetPath = strings.TrimSuffix(target.URL.Path, "/") + "/" +
			strings.TrimPrefix(strings.TrimPrefix(proxy.Ctx.OriginRequest.URL.Path, proxy.Ctx.RouteDetail.Pattern), "/")
	}

	// Step4 build Director
	director := func(req *http.Request) {
		req.Host = targetHost
		req.URL.Scheme = "http"
		if targetMode == router.DefaultTargetMode {
			req.URL.Scheme = target.URL.Scheme
		}
		req.URL.Host = targetHost
		req.URL.Path = targetPath
		if proxy.Ctx.ForwardRequest.URL.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = proxy.Ctx.ForwardRequest.URL.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = proxy.Ctx.ForwardRequest.URL.RawQuery + "&" + req.URL.RawQuery
		}
	}

	// Step5 build ModifyResponse & doPostFilters
	modifyResponse := func(resp *http.Response) error {
		proxy.Ctx.Response = resp
		_, err := proxy.DoPostFilters()
		if err != nil {
			logger.Info("[Dispatch-DoPostFilters] failed: err=%v", err)
			return err
		}
		return nil
	}

	// Step6 build ErrorHandler
	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		proxy.ErrorHandle(proxy.Ctx.ResponseWriter, http.StatusBadGateway)
	}

	// Step7 doProxy
	proxy.DoProxy(director, modifyResponse, errorHandler)
}
