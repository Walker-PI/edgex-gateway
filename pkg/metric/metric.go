package metric

import (
	"net/http"
	"time"

	"github.com/Walker-PI/iot-gateway/pkg/dal"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
)

// const (
// 	resultFail      = "fail"
// 	resultSucceed   = "succeed"
// 	resultRateLimit = "rateLimit"
// 	resultReject    = "reject"
// )

// AsyncStatusEmit 异步Status记录
func AsyncStatusEmit(startTime time.Time, req *http.Request, resp *http.Response) {
	go func() {
		defer func() {
			tools.RecoverPanic()
		}()
		endTime := time.Now()
		record := &dal.APIRequestRecord{
			Path:        req.URL.Path,
			Method:      req.Method,
			CostTime:    endTime.Sub(startTime).Milliseconds(),
			StatusCode:  resp.StatusCode,
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedTime: time.Now(),
		}
		if err := dal.AddRecord(storage.MysqlClient, record); err != nil {
			logger.Error("[AsyncStatusEmit] add record to DB failed: record=%+v, err=%v", *record, err)
			return
		}
	}()
}
