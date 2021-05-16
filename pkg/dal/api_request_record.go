package dal

import (
	"time"

	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"gorm.io/gorm"
)

// APIRequestRecord ...
type APIRequestRecord struct {
	ID          int64     `gorm:"column:id" json:"id"`
	Path        string    `gorm:"column:path" json:"path"`
	Method      string    `gorm:"column:method" json:"method"`
	StatusCode  int       `gorm:"column:status_code" json:"status_code"`
	CostTime    int64     `gorm:"column:cost_time" json:"cost_time"`
	StartTime   time.Time `gorm:"column:start_time" json:"start_time"`
	EndTime     time.Time `gorm:"column:end_time" json:"end_time"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	Result      string    `gorm:"column:result" json:"result"`
}

// AddRecord ...
func AddRecord(db *gorm.DB, record *APIRequestRecord) (err error) {
	if record == nil {
		return nil
	}
	dbRes := db.Debug().Model(&APIRequestRecord{}).Create(record)
	if dbRes.Error != nil {
		err = dbRes.Error
		logger.Error("[AddRecord] insert api record failed: record=%+v, err=%v", err)
		return
	}
	return
}
