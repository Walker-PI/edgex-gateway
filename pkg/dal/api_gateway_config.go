package dal

import (
	"github.com/Walker-PI/edgex-gateway/pkg/logger"
	"github.com/Walker-PI/edgex-gateway/pkg/storage"
)

// APIGatewayConfig ...
type APIGatewayConfig struct {
	ID                int64  `gorm:"column:id" json:"id"`
	Pattern           string `gorm:"column:pattern" json:"pattern"`
	Method            string `gorm:"column:method" json:"method"`
	ApiName           string `gorm:"column:api_name" json:"api_name"`
	TargetMode        int32  `gorm:"column:target_mode" json:"target_mode"`
	TargetHost        string `gorm:"column:target_host" json:"target_host"`
	TargetScheme      string `gorm:"column:target_scheme" json:"target_scheme"`
	TargetPath        string `gorm:"column:target_path" json:"target_path"`
	TargetServiceName string `gorm:"column:target_service_name" json:"target_service_name"`
	TargetLb          string `gorm:"column:target_lb" json:"target_lb"`
	TargetTimeout     int64  `gorm:"column:target_timeout" json:"target_timeout"`
	MaxQps            int32  `gorm:"column:max_qps" json:"max_qps"`
	Auth              string `gorm:"column:auth" json:"auth"`
	IpWhiteList       string `gorm:"column:ip_white_list" json:"ip_white_list"`
	IpBlackList       string `gorm:"column:ip_black_list" json:"ip_black_list"`
	CreatedTime       string `gorm:"column:created_time" json:"created_time"`
	ModifiedTime      string `gorm:"column:modified_time" json:"modified_time"`
	Deleted           int32  `gorm:"column:deleted" json:"deleted"`
	Description       string `gorm:"column:description" json:"description"`
}

func GetAllAPIConfig() (apiConfigList []*APIGatewayConfig, err error) {
	apiConfigList = make([]*APIGatewayConfig, 0)
	dbRes := storage.MysqlClient.Debug().Model(&APIGatewayConfig{}).Where("deleted = 0").Find(&apiConfigList)
	if dbRes.Error != nil {
		logger.Error("[GetAllAPIConfig] get all apiConfig failed: err=%v", dbRes.Error)
		err = dbRes.Error
		return
	}
	return
}
