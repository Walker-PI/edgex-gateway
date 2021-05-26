package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/logic"
	"github.com/Walker-PI/iot-gateway/pkg/ratelimit"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

var defaultRouter *Router

// Discovery
const (
	DefaultDiscovery = ""
	DiscoveryEureka  = "EUREKA"
	DiscoveryConsul  = "CONSUL"
)

// UpdateGatewayRoute ... Redis Sub channel
const UpdateGatewayRouteFmt = "%s-update-gateway-route"

func InitRouter() {
	var err error
	defaultRouter, err = newRouter()
	if err != nil {
		panic(err)
	}

	// Subscribe
	go func() {
		defer func() {
			tools.RecoverPanic()
		}()
		ctx := context.Background()
		channel := fmt.Sprintf(UpdateGatewayRouteFmt, conf.Server.Source)
		pubSub := storage.RedisClient.Subscribe(ctx, channel)
		if _, err := pubSub.Receive(ctx); err != nil {
			logger.Error("[Update-Router] Receive failed: err=%v", err)
			return
		}
		messageChannel := pubSub.Channel()
		for msg := range messageChannel {
			logger.Info("[Update-Router] Subscribe info: channel=%v, payload=%v", msg.Channel, msg.Payload)
			router, err := newRouter()
			if err != nil {
				logger.Error("[Update-Router] failed: err=%v", err)
				continue
			}
			defaultRouter = router
			logger.Info("[Update-Router] succeed!")
		}
	}()

	fmt.Println("[API-Gateway] Router initialized!")
}

func DefaultRouter() *Router {
	return defaultRouter
}

func newRouter() (router *Router, err error) {
	routeConfigList, err := logic.GetAllRouteConfig()
	if err != nil {
		return
	}

	router = &Router{
		root:         make(map[string]*node),
		routeInfoMap: make(map[string]*RouteInfo),
	}

	for _, routeConfig := range routeConfigList {
		if routeConfig == nil {
			continue
		}
		var routerInfo *RouteInfo
		routerInfo, err = packRouterInfo(routeConfig)
		if err != nil || routerInfo == nil {
			logger.Error("[Router-newRouter] pack routerInfo failed: err=%v", err)
			return
		}

		methods := strings.Split(routeConfig.Methods, ",")

		for _, method := range methods {
			err = router.addRoute(method, routerInfo)
			if err != nil {
				logger.Error("[Router-newRouter] addRoute failed: err=%v", err)
				return
			}
		}
	}

	// FIXME: delete
	logger.Info("[Route-newRouter] update %d apis", len(routeConfigList))
	return
}

func packRouterInfo(routeConfig *logic.RouteConfig) (*RouteInfo, error) {
	var err error
	if routeConfig == nil {
		err = errors.New("RouteConfig is nil!")
		return nil, err
	}
	routerInfo := &RouteInfo{
		Pattern:     routeConfig.Pattern,
		Auth:        routeConfig.AuthType,
		IPWhiteList: make([]net.IP, 0),
		IPBlackList: make([]net.IP, 0),
		GroupName:   routeConfig.GroupName,
	}
	if routeConfig.RateLimit > 0 {
		routerInfo.Limiter = ratelimit.NewRateLimiter(int64(routeConfig.RateLimit))
	}

	if routeConfig.IPBlackList != "" {
		ipStrs := strings.Split(routeConfig.IPBlackList, ",")
		for _, ipStr := range ipStrs {
			netIP := net.ParseIP(ipStr)
			if netIP == nil {
				logger.Warn("[Router-IPBlackList] ip is invalid: ip=%s", ipStr)
				continue
			}
			routerInfo.IPBlackList = append(routerInfo.IPBlackList, netIP)
		}
	}

	if routeConfig.IPWhiteList != "" {
		ipStrs := strings.Split(routeConfig.IPWhiteList, ",")
		for _, ipStr := range ipStrs {
			netIP := net.ParseIP(ipStr)
			if netIP == nil {
				logger.Warn("[Router-IPWhiteList] ip is invalid: ip=%s", ipStr)
				continue
			}
			routerInfo.IPWhiteList = append(routerInfo.IPWhiteList, netIP)
		}
	}

	// 默认超时时间
	timeout := 3 * time.Second
	if routeConfig.TargetTimeout > 0 {
		timeout = time.Duration(routeConfig.TargetTimeout) * time.Millisecond
	}

	switch routeConfig.Discovery {
	case DefaultDiscovery:
		targetURL, innErr := url.Parse(routeConfig.TargetURL)
		if innErr != nil {
			err = innErr
			break
		}
		routerInfo.Target = &target{
			Discovery: DefaultDiscovery,
			Timeout:   timeout,
			URL:       targetURL,
		}
	case DiscoveryConsul, DiscoveryEureka:
		routerInfo.Target = &target{
			Discovery:     routeConfig.Discovery,
			ServiceName:   routeConfig.DiscoveryServiceName,
			LoadBalance:   routeConfig.DiscoveryLoadBalance,
			DiscoveryPath: routeConfig.DiscoveryPath,
			Timeout:       timeout,
		}
	default:
		err = errors.New("Unknown target mode!")
	}
	return routerInfo, err
}
