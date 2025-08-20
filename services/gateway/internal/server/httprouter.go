package server

import (
	"net/http"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/jsonextractor"

	"github.com/gin-gonic/gin"
)

type httpRouter struct {
	*gin.Engine
}

type requestPerformer[TRequest any, TResponse any] func(*TRequest) (TResponse, httperr.Err)
type requestPerformerWithID[TRequest any, TResponse any] func(string, *TRequest) (TResponse, httperr.Err)

type emptyResponse struct{}

func createHandler[TRequest any, TResponse any](doRequest requestPerformer[TRequest, TResponse]) func(*gin.Context) {
	return createHandlerWithID(func(id string, r *TRequest) (TResponse, httperr.Err) {
		return doRequest(r)
	})
}

func createHandlerWithID[TRequest any, TResponse any](doRequest requestPerformerWithID[TRequest, TResponse]) func(*gin.Context) {
	ext := jsonextractor.New()
	return func(ctx *gin.Context) {
		httpCtx := httpContext{ctx}
		profileID := httpCtx.Param("id")

		var request TRequest
		err := ext.Extract(&request, ctx)
		if !err.IsOk() {
			httpCtx.submitError(err)
			return
		}

		response, err := doRequest(profileID, &request)
		if !err.IsOk() {
			httpCtx.submitError(err)
			return
		}

		httpCtx.JSON(http.StatusOK, response)
	}
}

func createRouter(serviceCtx *GatewayService) httpRouter {
	router := gin.Default()

	router.POST("/api/v1/profile", createHandler(
		func(r *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
			return serviceCtx.RegisterProfile(r)
		},
	))

	router.GET("/api/v1/profile/:id", createHandlerWithID(
		func(id string, r *jsonextractor.EmptyRequest) (api.GetProfileResponse, httperr.Err) {
			return serviceCtx.GetProfileInfo(id)
		},
	))

	authProfileGroup := router.Group("/api/v1/profile/:id")
	authProfileGroup.Use(authMiddleware(&serviceCtx.JwtVerifier))
	{
		authProfileGroup.PUT("", createHandlerWithID(
			func(id string, r *api.EditProfileRequest) (emptyResponse, httperr.Err) {
				return emptyResponse{}, serviceCtx.EditProfileInfo(id, r)
			},
		))

		authProfileGroup.DELETE("", createHandlerWithID(
			func(id string, r *jsonextractor.EmptyRequest) (emptyResponse, httperr.Err) {
				return emptyResponse{}, serviceCtx.DeleteProfile(id)
			},
		))
	}

	router.POST("/api/v1/auth", createHandler(
		func(r *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
			return serviceCtx.Authenticate(r)
		},
	))

	return httpRouter{router}
}
