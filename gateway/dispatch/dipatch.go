package dispatch

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/proxy"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/metric"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
	"github.com/afex/hystrix-go/hystrix"
)

func (d *Dispatcher) Dispatch(ctx *agw_context.AGWContext) {

	logger.Info("[Dispatch-Request] url=%v, body=%+v", tools.GetMarshalStr(ctx.ForwardRequest.URL),
		tools.GetBodyStr(ctx.ForwardRequest))

	defer func() {
		if ctx.Response == nil {
			ErrorHandle(ctx, http.StatusBadGateway)
		}
		// Step6 请求记录上报
		metric.AsyncStatusEmit(ctx)
		logger.Info("[Dispatch-Response] cost=%dms, status=%v, resp=%+v", time.Since(ctx.StartTime).Milliseconds(),
			ctx.Response.StatusCode, ctx.Response)
	}()

	// Step1. 执行pre 过滤器
	statusCode, err := d.DoPreFilters(ctx)
	logger.Info("[Dispatch] DoPreFilters finished: statusCode=%v, err=%v", statusCode, err)
	if err != nil || statusCode != http.StatusOK {
		logger.Error("[Dispatch-DoPreFilters] failed: err=%v", err)
		ErrorHandle(ctx, statusCode)
		return
	}

	// Step2. 执行routing 过滤器
	statusCode, err = d.DoRoutingFilter(ctx)
	logger.Info("[Dispatch] DoRoutingFilter finished: statusCode=%v, err=%v", statusCode, err)
	if err != nil || statusCode != http.StatusOK {
		logger.Error("[Dispatch-DoRoutingFilter] failed: err=%v", err)
		ErrorHandle(ctx, statusCode)
		return
	}

	// Step3. 反向代理构建
	proxy := proxy.NewProxy()
	proxy.BuildProxy(ctx)

	// Step4. 执行post过滤器
	proxy.SetModifyResponse(func(resp *http.Response) error {
		ctx.Response = resp
		_, err := d.DoPostFilters(ctx)
		if err != nil {
			logger.Info("[Dispatch-DoPostFilters] failed: err=%v", err)
			return err
		}
		return nil
	})

	// Step5. 执行error处理
	proxy.SetErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			ErrorHandle(ctx, http.StatusBadGateway)
		}
	})

	// Step6 熔断&反向代理
	timeout := 3 * time.Second
	if ctx.RouteInfo.Target.Timeout > 0 {
		timeout = ctx.RouteInfo.Target.Timeout * time.Millisecond
	}
	command := fmt.Sprintf("%s-%s", ctx.ForwardRequest.Method, ctx.RouteInfo.Pattern)
	hystrix.ConfigureCommand(command, hystrix.CommandConfig{
		Timeout:                int(timeout / time.Millisecond), // 超时时间
		MaxConcurrentRequests:  100,                             // 最大并发量
		SleepWindow:            int(time.Second * 5),            // 熔断之后，等待尝试时间
		RequestVolumeThreshold: 30,                              // 10秒内请求数量，达到这个值开始判断是否熔断
		ErrorPercentThreshold:  50,                              // 错误百分比
	})
	err = hystrix.Do(command, func() error {
		err = proxy.ReverseProxy(ctx)
		return err
	}, nil)
	if err != nil {
		logger.Error("[Dispatch] ReverseProxy failed: err=%v", err)
		return
	}
}
