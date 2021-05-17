package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Walker-PI/iot-gateway/pkg/dal"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/ratelimit"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

var defaultRouter *Router

const (
	DefaultTargetMode int32 = 1
	ConsulTargetMode  int32 = 2
)

// UpdateGatewayRoute ... Redis Sub channel
const UpdateGatewayRoute = "update-gateway-route"

type Router struct {
	root          map[string]*node
	routerInfoMap map[string]*RouterInfo
}

type RouterInfo struct {
	Pattern     string
	Target      *target  // Target
	Auth        string   // 鉴权类型
	IPWhiteList []net.IP // IP白名单
	IPBlackList []net.IP // IP黑名单
	Limiter     *ratelimit.RateLimiter
}

type target struct {
	Mode        int32
	ServiceName string
	URL         *url.URL
	Timeout     time.Duration // 超时时间
	LoadBalance string        // 负载均衡
}

type node struct {
	part     string
	children map[string]*node
	isParam  bool
}

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
		pubSub := storage.RedisClient.Subscribe(ctx, UpdateGatewayRoute)
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

	fmt.Println("[Edgex-gateway] Router initialized!")
}

func DefaultRouter() *Router {
	return defaultRouter
}

func newRouter() (router *Router, err error) {
	apiConfigList, err := dal.GetAllAPIConfig()
	if err != nil {
		return
	}

	router = &Router{
		root:          make(map[string]*node),
		routerInfoMap: make(map[string]*RouterInfo),
	}

	for _, apiConfig := range apiConfigList {
		if apiConfig == nil {
			continue
		}
		var routerInfo *RouterInfo
		routerInfo, err = packRouterInfo(apiConfig)
		if err != nil || routerInfo == nil {
			logger.Error("[Router-newRouter] pack routerInfo failed: err=%v", err)
			return
		}
		err = router.addRoute(apiConfig.Method, routerInfo)
		if err != nil {
			logger.Error("[Router-newRouter] addRoute failed: err=%v", err)
			return
		}
	}
	return
}

func packRouterInfo(apiConfig *dal.APIGatewayConfig) (*RouterInfo, error) {
	var err error
	if apiConfig == nil {
		err = errors.New("APIConfig is nil!")
		return nil, err
	}
	routerInfo := &RouterInfo{
		Pattern:     apiConfig.Pattern,
		Auth:        apiConfig.Auth,
		IPWhiteList: make([]net.IP, 0),
		IPBlackList: make([]net.IP, 0),
	}
	if apiConfig.MaxQPS > 0 {
		routerInfo.Limiter = ratelimit.NewRateLimiter(int64(apiConfig.MaxQPS))
	}

	if apiConfig.IPBlackList != "" {
		ipStrs := strings.Split(apiConfig.IPBlackList, ",")
		for _, ipStr := range ipStrs {
			netIP := net.ParseIP(ipStr)
			if netIP == nil {
				logger.Warn("[Router-IPBlackList] ip is invalid: ip=%s", ipStr)
				continue
			}
			routerInfo.IPBlackList = append(routerInfo.IPBlackList, netIP)
		}
	}

	if apiConfig.IPWhiteList != "" {
		ipStrs := strings.Split(apiConfig.IPWhiteList, ",")
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
	timeout := 5 * time.Second
	if apiConfig.TargetTimeout > 0 {
		timeout = time.Duration(apiConfig.TargetTimeout) * time.Millisecond
	}

	switch apiConfig.TargetMode {
	case DefaultTargetMode:
		routerInfo.Target = &target{
			Mode:    DefaultTargetMode,
			Timeout: timeout,
			URL: &url.URL{
				Host:   apiConfig.TargetHost,
				Path:   apiConfig.TargetPath,
				Scheme: apiConfig.TargetScheme,
			},
		}
	case ConsulTargetMode:
		routerInfo.Target = &target{
			Mode:        apiConfig.TargetMode,
			ServiceName: apiConfig.TargetServiceName,
			LoadBalance: apiConfig.TargetLb,
		}
	default:
		err = errors.New("Unknown target mode!")
	}
	return routerInfo, err
}

func (r *Router) addRoute(method string, info *RouterInfo) error {
	if info == nil {
		return errors.New("RouterInfo is empty!")
	}
	parts := parsePath(info.Pattern)
	if len(parts) == 0 {
		logger.Warn("[Router-addRoute] http: invalid pattern, pattern=%v", info.Pattern)
		return errors.New("http: invalid partten, pattern = " + info.Pattern)
	}
	if _, ok := r.root[method]; !ok {
		r.root[method] = &node{children: make(map[string]*node)}
	}
	root := r.root[method]
	key := method + "-" + "/" + strings.Join(parts, "/")

	if _, exsit := r.routerInfoMap[key]; exsit {
		logger.Warn("[Router-addRoute] http: multiple registrations for %s", info.Pattern)
		return errors.New("http: multiple registrations for " + info.Pattern)
	}

	for _, part := range parts {
		if root.children[part] == nil {
			root.children[part] = &node{
				part:     part,
				children: make(map[string]*node),
				isParam:  part[0] == ':',
			}
		}
		root = root.children[part]
	}
	r.routerInfoMap[key] = info
	return nil
}

func (r *Router) Match(method string, path string) (*RouterInfo, map[string]string) {

	logger.Info("[Router-Match] method=%v, path=%+v", method, path)

	curNode, exsit := r.root[method]
	if !exsit {
		return nil, nil
	}

	key := method + "-"
	searchParts := parsePath(path)
	params := make(map[string]string)

	for _, part := range searchParts {
		var nextNode *node
		for _, child := range curNode.children {
			if child.part == part || child.isParam {
				nextNode = child
				key = key + "/" + child.part
				if child.part[0] == ':' {
					params[child.part[1:]] = part
				}
				break
			}
		}
		if nextNode == nil {
			break
		}
		curNode = nextNode
	}

	routerInfo, exsit := r.routerInfoMap[key]
	if !exsit {
		return nil, nil
	}
	return routerInfo, params
}

func parsePath(partten string) []string {
	partList := strings.Split(partten, "/")
	parts := make([]string, 0)
	for _, v := range partList {
		if v == "" {
			continue
		}
		parts = append(parts, v)
	}
	return parts
}
