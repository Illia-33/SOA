package server

import "github.com/gin-gonic/gin"

type httpContext struct {
	*gin.Context
}

func (ctx *httpContext) submitError(err httpError) {
	ctx.Status(err.StatusCode)
	ctx.Error(err.Err)
}
