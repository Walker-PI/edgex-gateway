package gateway

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery"
	"github.com/Walker-PI/iot-gateway/gateway/dispatch"
	"github.com/Walker-PI/iot-gateway/gateway/router"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func Start() {
	var confFilePath string
	var discoveryType string
	flag.StringVar(&confFilePath, "conf", "conf/app.ini", "Specify configuration file path")
	flag.StringVar(&discoveryType, "discovery_type", "eureka", "Specify discovery type")
	flag.Parse()

	conf.LoadConfig(confFilePath)
	logger.InitLogs()
	storage.InitStorage()
	router.InitRouter()

	discovery.EnableDiscovery(discoveryType)
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/test_service", Test)
	http.HandleFunc("/", dispatch.Handle)
	fmt.Printf("[Edgex-gateway] Listening and serving HTTP on :%d\n", conf.Server.Port)
	_ = http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), nil)
}

func Test(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service_name")
	ctx := agw_context.NewAGWContext(w, r, nil, nil, time.Now())
	service, err := discovery.GetServiceInstance(ctx, serviceName, "RAND")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
	} else {
		dataBytes, _ := json.Marshal(service)
		_, _ = w.Write(dataBytes)
	}
}
