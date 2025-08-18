package server

import (
	"soa-socialnetwork/services/gateway/internal/httperr"

	"github.com/gin-gonic/gin"
)

type httpContext struct {
	*gin.Context
}

func (ctx *httpContext) submitError(err httperr.Err) {
	ctx.Status(err.StatusCode)
	ctx.Error(err.Err)
}
