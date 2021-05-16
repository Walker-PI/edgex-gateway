package metric

import (
	"net/http"
	"time"

	"github.com/Walker-PI/edgex-gateway/pkg/dal"
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	"github.com/Walker-PI/edgex-gateway/pkg/storage"
	"github.com/Walker-PI/edgex-gateway/pkg/tools"
)

const (
	resultFail      = "fail"
	resultSucceed   = "succeed"
	resultRateLimit = "rateLimit"
	resultReject    = "reject"
)

// 异步Status记录
func AsyncStatusEmit(startTime time.Time, req *http.Request, resp *http.Response) {
	go func() {
		defer func() {
			tools.RecoverPanic()
		}()
		endTime := time.Now()
		result := resultSucceed
		// TODO: 通过status_code判断 result
		record := &dal.APIRequestRecord{
			Path:        req.URL.Path,
			Method:      req.Method,
			CostTime:    endTime.Sub(startTime).Milliseconds(),
			StatusCode:  resp.StatusCode,
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedTime: time.Now(),
			Result:      result,
		}
		if err := dal.AddRecord(storage.MysqlClient, record); err != nil {
			logger.Error("[AsyncStatusEmit] add record to DB failed: record=%+v, err=%v", *record, err)
			return
		}
	}()
}
