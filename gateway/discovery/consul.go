package discovery

import (
	"errors"
	"fmt"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	"github.com/Walker-PI/edgex-gateway/gateway/lb"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	consulapi "github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
)

var consulClient *consulapi.Client

func initCousulClient() {
	var err error
	config := consulapi.DefaultConfig()
	config.Address = conf.ConsulConf.ConsulAddress
	consulClient, err = consulapi.NewClient(config)
	if err != nil {
		logger.Error("[EnableDiscovery] new consul client failed: config=%+v, err=%v", config, err)
		panic(err)
	}
}

func EnableDiscovery() {
	initCousulClient()
	registration := &consulapi.AgentServiceRegistration{
		ID:      conf.ConsulConf.ServiceName + "-" + uuid.NewV4().String(),
		Name:    conf.ConsulConf.ServiceName,
		Port:    conf.Server.Port,
		Address: conf.Server.Host,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/ping", conf.Server.Host, conf.Server.Port),
			Timeout:                        conf.ConsulConf.CheckTimeout,
			Interval:                       conf.ConsulConf.CheckInterval,
			DeregisterCriticalServiceAfter: "30s", // 故障检查失败30s后 consul自动将注册服务删除
		},
	}
	// Register the Service
	err := consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		logger.Error("[EnableDiscovery] consul register failed: err=%v", err)
		panic(err)
	}
}

func GetInstance(ctx *agw_context.AGWContext, serviceName string, lbType string) (*consulapi.CatalogService, error) {
	if serviceName == "" {
		return nil, errors.New("service_name is empty")
	}
	serviceList, _, err := consulClient.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return nil, err
	}

	// 负载均衡
	balance, err := lb.NewLoadBalance(ctx, lbType)
	if err != nil {
		return nil, err
	}
	return balance.Select(serviceList), nil
}
