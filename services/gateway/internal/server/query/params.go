package query

import "github.com/gin-gonic/gin"

type AuthTokenKind int

const (
	AUTH_TOKEN_EMPTY AuthTokenKind = iota
	AUTH_TOKEN_JWT
	AUTH_TOKEN_SOA
)

type Params struct {
	ProfileId string
	PostId    int32
	AuthToken string
	AuthKind  AuthTokenKind
}

const QUERY_PARAMS_KEY = "SOAQUERYPARAMS"

func ExtractParams(ctx *gin.Context) *Params {
	p, exists := ctx.Get(QUERY_PARAMS_KEY)
	if !exists {
		newParams := &Params{}
		ctx.Set(QUERY_PARAMS_KEY, newParams)
		return newParams
	}

	return p.(*Params)
}
