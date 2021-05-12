package discovery

import (
	"strconv"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	consulapi "github.com/hashicorp/consul/api"
)

func EnableDiscovery() {
	defaultConfig := consulapi.DefaultConfig()
	defaultConfig.Address = conf.ConsulConf.ConsulAddress + ":" + strconv.Itoa(conf.ConsulConf.ConsulPort)
	consul, err := consulapi.NewClient(defaultConfig)
	if err != nil {
		logger.Error("[EnableDiscovery] new consul client failed: config=%+v, err=%v", defaultConfig, err)
		panic(err)
	}

	// // Register the Service
	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		Name:    conf.ConsulConf.ServiceName,
		Address: conf.ConsulConf.ServiceAddress,
		Port:    conf.ConsulConf.ServicePort,
	})
	if err != nil {
		logger.Error("[EnableDiscovery] consul register failed: err=%v", err)
		panic(err)
	}

	// Register the Health Check
	err = consul.Agent().CheckRegister(&consulapi.AgentCheckRegistration{
		Name:      "Health Check",
		Notes:     "Check the health of the API",
		ServiceID: conf.ConsulConf.ServiceName,
		AgentServiceCheck: consulapi.AgentServiceCheck{
			HTTP:     conf.ConsulConf.CheckAddress,
			Interval: conf.ConsulConf.CheckInterval,
		},
	})
	if err != nil {
		logger.Error("[EnableDiscovery] consul check register failed: err=%v", err)
		panic(err)
	}
}
