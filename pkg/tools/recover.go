package tools

import (
	"runtime/debug"

	"github.com/Walker-PI/iot-gateway/pkg/logger"
)

// RecoverPanic ...
func RecoverPanic() {
	if x := recover(); x != nil {
		logger.Error("runtime panic: %v\n%v", x, string(debug.Stack()))
	}
}
