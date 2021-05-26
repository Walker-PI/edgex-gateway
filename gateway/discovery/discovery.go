package discovery

import (
	"strings"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/consul"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/eureka"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
)

type Discovery interface {
	Register()
	GetInstance(ctx *agw_context.AGWContext, serviceName string, lbType string) (*model.Instance, error)
}

var defaultDiscovery Discovery

func EnableDiscovery(discoveryType string) {
	discoveryType = strings.ToUpper(discoveryType)
	switch discoveryType {
	case "CONSUL":
		defaultDiscovery = consul.NewConsul()
		defaultDiscovery.Register()
	case "EUREKA":
		defaultDiscovery = eureka.NewEureka()
		defaultDiscovery.Register()
	default:
		panic("discovery_type is invalid")
	}
}

func GetServiceInstance(ctx *agw_context.AGWContext, serviceName string, lbType string) (*model.Instance, error) {
	return defaultDiscovery.GetInstance(ctx, serviceName, lbType)
}
