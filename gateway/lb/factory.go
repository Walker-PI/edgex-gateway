package lb

import (
	"errors"
	"strings"

	"github.com/Walker-PI/edgex-gateway/gateway/agw_context"
	consulapi "github.com/hashicorp/consul/api"
)

const (
	RandBalanceType       = "RAND"
	IPHashBalanceType     = "IP_HASH"
	URLHashBalanceType    = "URL_HASH"
	RoundRobinBalanceType = "ROUND_ROBIN"
)

type LoadBalance interface {
	Select([]*consulapi.CatalogService) *consulapi.CatalogService
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
