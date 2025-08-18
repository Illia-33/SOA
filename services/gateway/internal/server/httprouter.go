package server

import (
	"net/http"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/httperr"

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
	ext := jsonExtractor{}
	return func(ctx *gin.Context) {
		httpCtx := httpContext{ctx}
		profileID := httpCtx.Param("id")

		var request TRequest
		err := ext.extract(request, httpCtx)
		if !err.IsOK() {
			httpCtx.submitError(err)
			return
		}

		response, err := doRequest(profileID, &request)
		if !err.IsOK() {
			httpCtx.submitError(err)
			return
		}

		httpCtx.JSON(http.StatusOK, response)
	}
}

func createRouter(serviceCtx *gatewayService) httpRouter {
	router := gin.Default()

	router.POST("/api/v1/profile", createHandler(
		func(r *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
			return serviceCtx.RegisterProfile(r)
		},
	))
	router.GET("/api/v1/profile/:id", createHandlerWithID(
		func(id string, r *emptyRequest) (api.GetProfileResponse, httperr.Err) {
			return serviceCtx.GetProfileInfo(id)
		},
	))
	router.PUT("/api/v1/profile/:id", createHandlerWithID(
		func(id string, r *api.EditProfileRequest) (emptyResponse, httperr.Err) {
			return emptyResponse{}, serviceCtx.EditProfileInfo(id, r)
		},
	))
	router.DELETE("/api/v1/profile/:id", createHandlerWithID(
		func(id string, r *emptyRequest) (emptyResponse, httperr.Err) {
			return emptyResponse{}, serviceCtx.DeleteProfile(id)
		},
	))

	return httpRouter{router}
}
