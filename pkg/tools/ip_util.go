package tools

import (
	"net"
	"net/http"
	"strings"
)

var localNetworks []*net.IPNet

func init() {
	localNetworks = append(localNetworks,
		&net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		&net.IPNet{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)},
		&net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 17, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 18, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 19, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 20, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 21, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 22, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 23, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 24, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 25, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 26, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 27, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 28, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 29, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 30, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(172, 31, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
	)
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// ClientPublicIP 尽最大努力实现获取客户端公网 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}

	return ""
}

// HasLocalIPddr 检测 IP 地址字符串是否是内网地址
func HasLocalIPddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	for _, network := range localNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return ip.IsLoopback()
}

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
