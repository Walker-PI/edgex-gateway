package router

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/ratelimit"
)

type Router struct {
	root         map[string]*node
	routeInfoMap map[string]*RouteInfo
}

type node struct {
	part     string
	children map[string]*node
	isParam  bool
}

type RouteInfo struct {
	Pattern     string
	Target      *target  // Target
	Auth        string   // 鉴权类型
	IPWhiteList []net.IP // IP白名单
	IPBlackList []net.IP // IP黑名单
	Limiter     *ratelimit.RateLimiter
	GroupName   string
}

type target struct {
	Discovery     string
	DiscoveryPath string
	ServiceName   string
	URL           *url.URL
	Timeout       time.Duration // 超时时间
	LoadBalance   string        // 负载均衡
}

func (r *Router) addRoute(method string, info *RouteInfo) error {
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
	curNode := r.root[method]
	key := method + "-" + "/" + strings.Join(parts, "/")

	if _, exsit := r.routeInfoMap[key]; exsit {
		logger.Warn("[Router-addRoute] http: multiple registrations for %s", info.Pattern)
		return errors.New("http: multiple registrations for " + info.Pattern)
	}

	for _, part := range parts {
		if curNode.children[part] == nil {
			curNode.children[part] = &node{
				part:     part,
				children: make(map[string]*node),
				isParam:  part[0] == ':',
			}
		}
		curNode = curNode.children[part]
	}
	r.routeInfoMap[key] = info
	return nil
}

func (r *Router) Match(method string, path string) (*RouteInfo, map[string]string) {

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

	routeInfo, exsit := r.routeInfoMap[key]
	if !exsit {
		return nil, nil
	}
	return routeInfo, params
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
