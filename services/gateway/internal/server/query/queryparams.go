package query

import "github.com/gin-gonic/gin"

type Params struct {
	ProfileId   string
	RawJwtToken string
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
