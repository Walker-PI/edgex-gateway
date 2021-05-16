package filter

import (
	"net/http"
	"strings"

	"github.com/Walker-PI/iot-gateway/gateway/agw_context"
	"github.com/Walker-PI/iot-gateway/pkg/tools"
	"github.com/dgrijalva/jwt-go"
)

const (
	KeyLess = "KEYLESS"
	AuthJWT = "JWT"
)

type AuthFilter struct {
	baseFilter
}

func newAuthFilter() Filter {
	return &AuthFilter{}
}

func (f *AuthFilter) Name() FilterName {
	return PreAuthFilter
}

func (f *AuthFilter) Priority() int {
	return priority[PreAuthFilter]
}

func (f *AuthFilter) Type() FilterType {
	return PreFilter
}

func (f *AuthFilter) Run(ctx *agw_context.AGWContext) (Code int, err error) {
	auth := ctx.RouteDetail.Auth
	auth = strings.ToUpper(auth)
	switch auth {
	case AuthJWT:
		tokenStr := ctx.ForwardRequest.Header.Get("Authorization")
		token, err := jwt.ParseWithClaims(tokenStr, &tools.LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(tools.SecretKey), nil
		})
		if err != nil {
			return http.StatusUnauthorized, err
		}
		if _, ok := token.Claims.(*tools.LoginClaims); !ok || !token.Valid {
			return http.StatusUnauthorized, nil
		}
		return f.baseFilter.Run(ctx)
	default:
		return f.baseFilter.Run(ctx)
	}
}
