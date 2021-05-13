package gateway

import (
	"flag"
	"net/http"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/gateway/discovery"
	"github.com/Walker-PI/edgex-gateway/gateway/proxy"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	"github.com/Walker-PI/edgex-gateway/pkg/storage"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func Start() {
	var (
		confFilePath  string
		discoveryType string
	)
	flag.StringVar(&confFilePath, "conf", "conf/app.ini", "Specify configuration file path")
	flag.StringVar(&discoveryType, "discovery_type", "consul", "Specify discovery type")
	flag.Parse()

	conf.LoadConfig(confFilePath)
	logger.InitLogs()
	storage.InitStorage()

	discovery.EnableDiscovery(discoveryType)
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/", proxy.ProxyHandler)
	http.ListenAndServe(conf.Server.ListenAddress, nil)
}
