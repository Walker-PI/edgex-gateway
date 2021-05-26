package consul

import (
	"errors"
	"fmt"

	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
	"github.com/Walker-PI/iot-gateway/gateway/lb"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	consulapi "github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
)

var consulClient *consulapi.Client

func initClient() {
	var err error
	config := consulapi.DefaultConfig()
	config.Address = conf.ConsulConf.ConsulAddress
	consulClient, err = consulapi.NewClient(config)
	if err != nil {
		logger.Error("[EnableDiscovery] new consul client failed: config=%+v, err=%v", config, err)
		panic(err)
	}
}

type Consul struct{}

func NewConsul() *Consul {
	return &Consul{}
}

func (c *Consul) Register() {
	initClient()
	registration := &consulapi.AgentServiceRegistration{
		ID:      conf.ConsulConf.ServiceName + "-" + uuid.NewV4().String(),
		Name:    conf.ConsulConf.ServiceName,
		Port:    conf.Server.Port,
		Address: conf.ConsulConf.ServiceHost,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/ping", conf.ConsulConf.ServiceHost, conf.Server.Port),
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

func (c *Consul) GetInstance(ctx *agw_context.AGWContext, serviceName string, lbType string) (*model.Instance, error) {
	if serviceName == "" {
		return nil, errors.New("service_name is empty")
	}
	serviceList, _, err := consulClient.Catalog().Service(serviceName, "", nil)
	if err != nil {
		return nil, err
	}
	instanceList := make([]*model.Instance, 0)
	for _, service := range serviceList {
		instanceList = append(instanceList, &model.Instance{
			ServiceID:      service.ID,
			ServiceName:    service.ServiceName,
			ServiceAddress: service.Address,
			ServicePort:    service.ServicePort,
		})
	}
	// 负载均衡
	balance, err := lb.NewLoadBalance(ctx, lbType)
	if err != nil {
		return nil, err
	}

	return balance.Select(instanceList), nil
}
