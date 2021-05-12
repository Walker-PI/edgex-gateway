package main

import (
	"flag"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/gateway"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	"github.com/Walker-PI/edgex-gateway/pkg/storage"
)

func main() {

	var confFilePath string
	flag.StringVar(&confFilePath, "conf", "conf/app.ini", "Specify local configuration file path")
	flag.Parse()

	// 基本配置初始化
	conf.LoadConfig(confFilePath)
	logger.InitLogs()
	storage.InitStorage()

	gateway.Start()
}
