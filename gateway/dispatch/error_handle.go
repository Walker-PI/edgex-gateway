package dispatch

import (
	"fmt"
	"net/http"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
)

func ErrorHandle(ctx *agw_context.AGWContext, statusCode int) {
	status := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	ctx.Response = &http.Response{
		Status:     status,
		StatusCode: statusCode,
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.ResponseWriter.WriteHeader(statusCode)
	fmt.Fprintln(ctx.ResponseWriter, status)
}
