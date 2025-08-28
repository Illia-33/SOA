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

func newHttpRouter(service *GatewayService) httpRouter {
	router := gin.Default()
	restApi := router.Group("/api/v1")

	profileGroup := restApi.Group("/profile")
	{
		profileGroup.POST("", createHandler(
			func(qp *query.Params, r *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
				return service.RegisterProfile(qp, r)
			},
		))

		idGroup := profileGroup.Group("/:profile_id")
		idGroup.Use(query.WithProfileId())
		idGroup.GET("", createHandler(
			func(qp *query.Params, r *api.Empty) (api.GetProfileResponse, httperr.Err) {
				return service.GetProfileInfo(qp)
			},
		))

		{
			idAuthGroup := idGroup.Group("")
			idAuthGroup.Use(query.WithAuth(&service.jwtVerifier))
			idAuthGroup.PUT("", createHandler(
				func(qp *query.Params, r *api.EditProfileRequest) (api.Empty, httperr.Err) {
					return api.Empty{}, service.EditProfileInfo(qp, r)
				},
			))

			idAuthGroup.DELETE("", createHandler(
				func(qp *query.Params, r *api.Empty) (api.Empty, httperr.Err) {
					return api.Empty{}, service.DeleteProfile(qp)
				},
			))
		}

		pageGroup := idGroup.Group("/page")
		{
			pageGroup.GET("/settings", createHandler(
				func(qp *query.Params, r *api.Empty) (api.GetPageSettingsResponse, httperr.Err) {
					return service.GetPageSettings(qp)
				},
			))
		}

		{
			authPageGroup := pageGroup.Group("")
			authPageGroup.Use(query.WithAuth(&service.jwtVerifier))
			authPageGroup.PUT("/settings", createHandler(
				func(qp *query.Params, r *api.EditPageSettingsRequest) (api.Empty, httperr.Err) {
					return api.Empty{}, service.EditPageSettings(qp, r)
				},
			))

			authPageGroup.POST("", createHandler(
				func(qp *query.Params, r *api.NewPostRequest) (api.NewPostResponse, httperr.Err) {
					return service.NewPost(qp, r)
				},
			))

			{
				postGroup := authPageGroup.Group("/:post_id")
				postGroup.Use(query.WithPostId())
				postGroup.POST("/comments", createHandler(
					func(qp *query.Params, r *api.NewCommentRequest) (api.NewCommentResponse, httperr.Err) {
						return service.NewComment(qp, r)
					},
				))
			}
		}
	}

	restApi.POST("/auth", createHandler(
		func(qp *query.Params, r *api.AuthenticateRequest) (api.AuthenticateResponse, httperr.Err) {
			return service.Authenticate(qp, r)
		},
	))

	restApi.POST("/api_token", createHandler(
		func(qp *query.Params, r *api.CreateApiTokenRequest) (api.CreateApiTokenResponse, httperr.Err) {
			return service.CreateApiToken(qp, r)
		},
	))

	return httpRouter{router}
}
