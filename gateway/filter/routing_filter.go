package filter

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery"
	"github.com/Walker-PI/iot-gateway/gateway/router"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
)

type RoutingFilter struct {
	baseFilter
}

func newRoutingFilter() Filter {
	return &RoutingFilter{}
}

func (f *RoutingFilter) Name() FilterName {
	return RoutingFilterName
}

func (f *RoutingFilter) Type() FilterType {
	return RouteFilter
}

func (f *RoutingFilter) Priority() int {
	return priority[RoutingFilterName]
}

func (f *RoutingFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {

	targetObj := ctx.RouteInfo.Target
	targetURL := copyURL(ctx.ForwardRequest.URL)

	switch targetObj.Discovery {
	case router.DefaultDiscovery:
		defTargetURL := targetObj.URL
		targetURL.Host = defTargetURL.Host
		targetURL.Scheme = defTargetURL.Scheme
		targetURL.Path = combinePath(defTargetURL.Path, strings.TrimPrefix(targetURL.Path, ctx.RouteInfo.Pattern))
		targetURL.RawPath = targetURL.Path
	case router.DiscoveryEureka, router.DiscoveryConsul:
		service, err := discovery.GetServiceInstance(ctx, targetObj.ServiceName, targetObj.LoadBalance)
		if err != nil || service == nil {
			return http.StatusNotFound, err
		}
		targetURL.Host = service.ServiceAddress + ":" + strconv.Itoa(service.ServicePort)
		targetURL.Scheme = "http"

		// 指定路径转换
		if targetObj.DiscoveryPath != "" {
			path, err := getTargetPath(targetObj.DiscoveryPath, ctx)
			if err != nil {
				logger.Error("[RoutingFilter] requestPath or discoveryPath is invalid: requestPath=%v, discoveryPath=%v",
					targetURL.Path, targetObj.DiscoveryPath)
				return http.StatusBadGateway, err
			}
			targetURL.Path = path
		}
	}
	ctx.TargetURL = targetURL
	return f.baseFilter.Run(ctx)
}

func combinePath(pathA string, pathB string) string {
	pathA = strings.TrimSuffix(pathA, "/")
	pathB = strings.TrimPrefix(pathB, "/")
	return pathA + "/" + pathB
}

func getTargetPath(pathFmt string, ctx *agw_context.AGWContext) (path string, err error) {
	path = "/"
	parts := strings.Split(pathFmt, "/")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			v, exsit := ctx.GetPathParam(part[1:])
			if !exsit || v == "" {
				err = fmt.Errorf("Params is invalid: pathFmt=%v, key=%v", pathFmt, part[1:])
				return
			}
			path = path + "/" + v
			continue
		}
		path = path + "/" + part
	}
	return
}

// 深度拷贝
func copyURL(v *url.URL) *url.URL {
	if v == nil {
		return nil
	}
	u := &url.URL{
		Scheme:      v.Scheme,
		Opaque:      v.Opaque,
		Host:        v.Host,
		Path:        v.Path,
		RawPath:     v.RawPath,
		ForceQuery:  v.ForceQuery,
		RawQuery:    v.RawQuery,
		Fragment:    v.Fragment,
		RawFragment: v.RawFragment,
	}

	if v.User != nil {
		if password, ok := v.User.Password(); ok {
			u.User = url.UserPassword(v.User.Username(), password)
		} else {
			u.User = url.User(v.User.Username())
		}
	}
	return u
}
