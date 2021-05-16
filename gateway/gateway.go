package gateway

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/gateway/discovery"
	"github.com/Walker-PI/edgex-gateway/gateway/dispatch"
	"github.com/Walker-PI/edgex-gateway/gateway/router"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	"github.com/Walker-PI/edgex-gateway/pkg/storage"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func Start() {
	var confFilePath string
	// var discoveryType string
	flag.StringVar(&confFilePath, "conf", "conf/app.ini", "Specify configuration file path")
	// flag.StringVar(&discoveryType, "discovery_type", "consul", "Specify discovery type")
	flag.Parse()

	conf.LoadConfig(confFilePath)
	logger.InitLogs()
	storage.InitStorage()
	router.InitRouter()

	discovery.EnableDiscovery()
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/", dispatch.Dsipatch)
	fmt.Printf("[Edgex-gateway] Listening and serving HTTP on :%d\n", conf.Server.Port)
	http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), nil)
}
