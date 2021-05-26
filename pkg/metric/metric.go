package metric

import (
	"time"

	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/pkg/dal"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

// AsyncStatusEmit 异步Status记录
func AsyncStatusEmit(ctx *agw_context.AGWContext) {
	go func() {
		defer func() {
			tools.RecoverPanic()
		}()
		endTime := time.Now()
		record := &dal.RequestRecord{
			Source:     conf.Server.Source,
			GroupName:  ctx.RouteInfo.GroupName,
			Path:       ctx.RouteInfo.Pattern,
			Method:     ctx.ForwardRequest.Method,
			StatusCode: ctx.Response.StatusCode,
			CostTime:   endTime.Sub(ctx.StartTime).Milliseconds(),
		}
		if err := dal.AddRequestRecord(storage.MysqlClient, record); err != nil {
			logger.Error("[AsyncStatusEmit] add record to DB failed: record=%+v, err=%v", *record, err)
			return
		}
	}()
}
