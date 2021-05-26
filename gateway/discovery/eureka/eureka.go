package eureka

import (
	"strconv"

	"github.com/SimonWang00/goeureka"
	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
	"github.com/Walker-PI/iot-gateway/gateway/lb"
)

type Eureka struct{}

func NewEureka() *Eureka {
	return &Eureka{}
}

func (e *Eureka) Register() {
	goeureka.RegisterClient(conf.EurekaConf.EurekaURL, conf.EurekaConf.LocalIP, conf.EurekaConf.ServiceName,
		strconv.Itoa(conf.Server.Port), "43", nil)
}

func (e *Eureka) GetInstance(ctx *agw_context.AGWContext, serviceName string, lbType string) (*model.Instance, error) {
	serviceList, err := goeureka.GetServiceInstances(serviceName)
	if err != nil {
		return nil, err
	}
	instanceList := make([]*model.Instance, 0)
	for _, service := range serviceList {
		instanceList = append(instanceList, &model.Instance{
			ServiceID:      service.InstanceId,
			ServiceName:    service.App,
			ServiceAddress: service.IpAddr,
			ServicePort:    service.Port.Port,
		})
	}
	// 负载均衡
	balance, err := lb.NewLoadBalance(ctx, lbType)
	if err != nil {
		return nil, err
	}
	return balance.Select(instanceList), nil
}
