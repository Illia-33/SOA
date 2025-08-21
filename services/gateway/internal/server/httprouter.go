package server

import (
	"net/http"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/jsonextractor"
	"soa-socialnetwork/services/gateway/internal/server/query"

	"github.com/gin-gonic/gin"
)

type httpRouter struct {
	*gin.Engine
}

type requestPerformer[TRequest any, TResponse any] func(*query.Params, *TRequest) (TResponse, httperr.Err)

type emptyResponse struct{}

func createHandler[TRequest any, TResponse any](doRequest requestPerformer[TRequest, TResponse]) func(*gin.Context) {
	extractor := jsonextractor.New()
	return func(ctx *gin.Context) {
		params := query.ExtractParams(ctx)
		var request TRequest
		err := extractor.Extract(&request, ctx)

		if !err.IsOk() {
			ctx.AbortWithError(err.StatusCode, err.Err)
			return
		}

		response, err := doRequest(params, &request)
		if !err.IsOk() {
			ctx.AbortWithError(err.StatusCode, err.Err)
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func createRouter(serviceCtx *GatewayService) httpRouter {
	router := gin.Default()
	restApi := router.Group("/api/v1")

	profileGroup := restApi.Group("/profile")
	{
		profileGroup.POST("", createHandler(
			func(qp *query.Params, r *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
				return serviceCtx.RegisterProfile(qp, r)
			},
		))

		idGroup := profileGroup.Group("/:profile_id")
		idGroup.Use(query.WithProfileID())
		idGroup.GET("", createHandler(
			func(qp *query.Params, r *jsonextractor.EmptyRequest) (api.GetProfileResponse, httperr.Err) {
				return serviceCtx.GetProfileInfo(qp)
			},
		))

		idAuthGroup := idGroup.Group("")
		idAuthGroup.Use(query.WithJwtAuth(&serviceCtx.JwtVerifier))
		idAuthGroup.PUT("", createHandler(
			func(qp *query.Params, r *api.EditProfileRequest) (emptyResponse, httperr.Err) {
				return emptyResponse{}, serviceCtx.EditProfileInfo(qp, r)
			},
		))

		idAuthGroup.DELETE("", createHandler(
			func(qp *query.Params, r *jsonextractor.EmptyRequest) (emptyResponse, httperr.Err) {
				return emptyResponse{}, serviceCtx.DeleteProfile(qp)
			},
		))
	}

	restApi.POST("/auth", createHandler(
		func(qp *query.Params, r *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
			return serviceCtx.Authenticate(qp, r)
		},
	))

	return httpRouter{router}
}
