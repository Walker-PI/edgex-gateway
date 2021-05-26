package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Walker-PI/iot-gateway/conf"
	"github.com/Walker-PI/iot-gateway/pkg/logger"
	"github.com/Walker-PI/iot-gateway/pkg/storage"
	"github.com/go-redis/redis/v8"
)

// Redis key
const (
	AllRouteConfigIDFmt = "all-route-config-id:source:%s"
	RouteConfigKeyFmt   = "route-config:api_id:%d:source:%s"
)

type RouteConfig struct {
	ID                   int64     `gorm:"column:id" json:"id"`
	GroupID              int64     `gorm:"column:group_id" json:"group_id"`
	GroupName            string    `gorm:"column:group_name" json:"group_name"`
	Source               string    `gorm:"column:source" json:"source"`
	Pattern              string    `gorm:"column:pattern" json:"pattern"`
	Methods              string    `gorm:"column:methods" json:"methods"`
	RateLimit            int32     `gorm:"column:rate_limit" json:"rate_limit"`
	AuthType             string    `gorm:"column:auth_type" json:"auth_type"`
	IPWhiteList          string    `gorm:"column:ip_white_list" json:"ip_white_list"`
	IPBlackList          string    `gorm:"column:ip_black_list" json:"ip_black_list"`
	TargetURL            string    `gorm:"column:target_url" json:"target_url"`
	TargetTimeout        int32     `gorm:"column:target_timeout" json:"target_timeout"`
	Discovery            string    `gorm:"column:discovery" json:"discovery"`
	DiscoveryPath        string    `gorm:"column:discovery_path" json:"discovery_path"`
	DiscoveryServiceName string    `gorm:"column:discovery_service_name" json:"discovery_service_name"`
	DiscoveryLoadBalance string    `gorm:"column:discovery_load_balance" json:"discovery_load_balance"`
	Deleted              int32     `gorm:"column:deleted" json:"deleted"`
	CreatedTime          time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	ModifiedTime         time.Time `gorm:"column:modified_time;default:CURRENT_TIMESTAMP" json:"modified_time"`
}

func GetAllRouteConfig() (routeConfigList []*RouteConfig, err error) {
	var (
		ctx          = context.Background()
		routeStrList []string
	)

	key := fmt.Sprintf(AllRouteConfigIDFmt, conf.Server.Source)
	routeStrList, err = storage.RedisClient.SMembers(ctx, key).Result()
	if err != nil {
		logger.Error("[GetAllRouteConfig] get all route_id failed: err=%v", err)
		return
	}

	redisPipeline := storage.RedisClient.Pipeline()
	defer redisPipeline.Close()

	cmds := make([]*redis.StringCmd, len(routeStrList))

	for i, routeStr := range routeStrList {
		routeID, innErr := strconv.ParseInt(routeStr, 10, 64)
		if innErr != nil {
			continue
		}
		key := fmt.Sprintf(RouteConfigKeyFmt, routeID, conf.Server.Source)
		cmds[i] = redisPipeline.Get(ctx, key)
	}
	_, err = redisPipeline.Exec(ctx)
	if err != nil {
		logger.Error("[GetAllRouteConfig] redis pipeline exec failed: err=%v", err)
		return
	}

	routeConfigList = make([]*RouteConfig, 0)

	for i, routeStr := range routeStrList {
		if cmds[i] == nil {
			continue
		}
		data, innErr := cmds[i].Result()
		if innErr != nil || data == "" {
			logger.Error("[GetAllRouteConfig] route_id=%v, data=%v, err=%v", routeStr, data, innErr)
			continue
		}
		var routeConfig = &RouteConfig{}
		innErr = json.Unmarshal([]byte(data), &routeConfig)
		if innErr != nil {
			logger.Error("[GetAllRouteConfig] json marshal failed: data=%v, err=%v", data, err)
			continue
		}
		routeConfigList = append(routeConfigList, routeConfig)
	}
	return
}
