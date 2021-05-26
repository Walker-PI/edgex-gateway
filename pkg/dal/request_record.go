package dal

import (
	"time"

	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"gorm.io/gorm"
)

type RequestRecord struct {
	ID          int64     `gorm:"column:id" json:"id" form:"id"`
	Source      string    `gorm:"column:source" json:"source" form:"source"`
	GroupName   string    `gorm:"column:group_name" json:"group_name" form:"group_name"`
	Path        string    `gorm:"column:path" json:"path" form:"path"`
	Method      string    `gorm:"column:method" json:"method" form:"method"`
	StatusCode  int       `gorm:"column:status_code" json:"status_code" form:"status_code"`
	CostTime    int64     `gorm:"column:cost_time" json:"cost_time" form:"cost_time"`
	CreatedTime time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time" form:"created_time"`
}

func AddRequestRecord(db *gorm.DB, record *RequestRecord) error {
	if record == nil {
		return nil
	}
	dbRes := db.Debug().Model(&RequestRecord{}).Create(record)
	if dbRes.Error != nil {
		logger.Error("[AddRequestRecord] write a record failed: err=%v", dbRes.Error)
		return dbRes.Error
	}
	return nil

}
