package discovery

import (
	"fmt"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	consulapi "github.com/hashicorp/consul/api"
)

func EnableDiscovery(_type string) {

	config := consulapi.DefaultConfig()
	config.Address = conf.ConsulConf.ConsulAddress
	consul, err := consulapi.NewClient(config)
	if err != nil {
		logger.Error("[EnableDiscovery] new consul client failed: config=%+v, err=%v", config, err)
		return
	}
	registration := &consulapi.AgentServiceRegistration{
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
	err = consul.Agent().ServiceRegister(registration)
	if err != nil {
		logger.Error("[EnableDiscovery] consul register failed: err=%v", err)
		return
	}
}

func ConsulDiscovery(serviceName string) {
	// // 创建连接consul服务配置
	// config := consulapi.DefaultConfig()
	// config.Address = conf.ConsulConf.ConsulAddress
	// client, err := consulapi.NewClient(config)
	// if err != nil {
	// 	logger.Error("[EnableDiscovery] new consul client failed: config=%+v, err=%v", config, err)
	// 	return
	// }

	// serviceList, queryMeta, err := client.Catalog().Service(serviceName, "", nil)
	// if err != nil {
	// 	logger.Error("[ConsulDiscovery] ")
	// }
	//  queryMeta.
	// // 获取指定service
	// service, _, err := client.Agent().Service("337", nil)
	// if err == nil {
	// 	fmt.Println(service.Address)
	// 	fmt.Println(service.Port)
	// }

}
