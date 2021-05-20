package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/go-redis/redis/v8"
)

// Redis key
const (
	AllAPIConfigID  = "all-api-config-id"
	APIConfigKeyFmt = "api-config:api_id:%d"
)

type APIGatewayConfig struct {
	ID                int64  `gorm:"column:id" json:"id"`
	Pattern           string `gorm:"column:pattern" json:"pattern"`
	Method            string `gorm:"column:method" json:"method"`
	APIName           string `gorm:"column:api_name" json:"api_name"`
	TargetMode        int32  `gorm:"column:target_mode" json:"target_mode"`
	TargetHost        string `gorm:"column:target_host" json:"target_host"`
	TargetScheme      string `gorm:"column:target_scheme" json:"target_scheme"`
	TargetPath        string `gorm:"column:target_path" json:"target_path"`
	TargetServiceName string `gorm:"column:target_service_name" json:"target_service_name"`
	TargetStripPrefix int32  `gorm:"column:target_strip_prefix" json:"target_strip_prefix"`
	TargetLb          string `gorm:"column:target_lb" json:"target_lb"`
	TargetTimeout     int64  `gorm:"column:target_timeout" json:"target_timeout"`
	MaxQPS            int32  `gorm:"column:max_qps" json:"max_qps"`
	Auth              string `gorm:"column:auth" json:"auth"`
	IPWhiteList       string `gorm:"column:ip_white_list" json:"ip_white_list"`
	IPBlackList       string `gorm:"column:ip_black_list" json:"ip_black_list"`
	CreatedTime       string `gorm:"column:created_time" json:"created_time"`
	ModifiedTime      string `gorm:"column:modified_time" json:"modified_time"`
	Status            int32  `gorm:"column:status" json:"status"`
	Description       string `gorm:"column:description" json:"description"`
}

func GetAllAPIConfig() (apiConfigList []*APIGatewayConfig, err error) {
	var (
		ctx        = context.Background()
		apiStrList []string
	)

	apiStrList, err = storage.RedisClient.SMembers(ctx, AllAPIConfigID).Result()
	if err != nil {
		logger.Error("[GetAllAPIConfig] get all api_id failed: err=%v", err)
		return
	}

	redisPipeline := storage.RedisClient.Pipeline()
	defer redisPipeline.Close()

	cmds := make([]*redis.StringCmd, len(apiStrList))

	for i, apiStr := range apiStrList {
		apiID, innErr := strconv.ParseInt(apiStr, 10, 64)
		if innErr != nil {
			continue
		}
		key := fmt.Sprintf(APIConfigKeyFmt, apiID)
		cmds[i] = redisPipeline.Get(ctx, key)
	}
	_, err = redisPipeline.Exec(ctx)
	if err != nil {
		logger.Error("[GetAllAPIConfig] redis pipeline exec failed: err=%v", err)
		return
	}

	apiConfigList = make([]*APIGatewayConfig, 0)

	for i, apiStr := range apiConfigList {
		if cmds[i] == nil {
			continue
		}
		data, innErr := cmds[i].Result()
		if innErr != nil || data == "" {
			logger.Error("[GetAllAPIConfig] api_id=%v, data=%v, err=%v", apiStr, data, innErr)
			continue
		}
		var apiConfig = &APIGatewayConfig{}
		innErr = json.Unmarshal([]byte(data), &apiConfig)
		if innErr != nil {
			logger.Error("[GetAllAPIConfig] json marshal failed: data=%v, err=%v", data, err)
			continue
		}
		apiConfigList = append(apiConfigList, apiConfig)
	}
	return
}
