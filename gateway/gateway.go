package gateway

import (
	"net/http"

	"github.com/Walker-PI/edgex-gateway/conf"
	"github.com/Walker-PI/edgex-gateway/gateway/discovery"
	"github.com/gorilla/mux"
)

func init() {

}

func Start() {
	// Step1. 服务注册
	discovery.EnableDiscovery()

	// 路由注册
	router := mux.NewRouter()
	router.HandleFunc("/", http.NotFound)
	router.HandleFunc("/ping", Ping)
	http.ListenAndServe(conf.Server.Port, router)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
