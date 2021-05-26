package lb

import (
	"errors"
	"strings"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/gateway/discovery/model"
)

const (
	RandBalanceType       = "RAND"
	IPHashBalanceType     = "IP_HASH"
	URLHashBalanceType    = "URL_HASH"
	RoundRobinBalanceType = "ROUND_ROBIN"
)

type LoadBalance interface {
	Select([]*model.Instance) *model.Instance
}

func NewLoadBalance(ctx *agw_context.AGWContext, lbType string) (LoadBalance, error) {
	lbType = strings.ToUpper(lbType)
	switch lbType {
	case RandBalanceType:
		return newRandBalance(ctx), nil
	case IPHashBalanceType:
		return newIPHashBalance(ctx), nil
	case URLHashBalanceType:
		return newURLHashBalance(ctx), nil
	case RoundRobinBalanceType:
		return newRRBalance(ctx), nil
	default:
		return nil, errors.New("Unknown load balance type!")
	}
}
