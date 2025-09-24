package server

import (
	"net/http"
	"soa-socialnetwork/services/gateway/api"
	"soa-socialnetwork/services/gateway/internal/server/httperr"
	"soa-socialnetwork/services/gateway/internal/server/query"

	"github.com/gin-gonic/gin"
)

type httpRouter struct {
	*gin.Engine
}

type empty struct{}

type requestPerformer[TRequest any, TResponse any] func(*query.Params, *TRequest) (TResponse, httperr.Err)

func createHandler[TRequest any, TResponse any](doRequest requestPerformer[TRequest, TResponse]) func(*gin.Context) {
	return func(ctx *gin.Context) {
		params := query.ExtractParams(ctx)
		var request TRequest
		if err := ctx.BindJSON(&request); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
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
	withAuth := query.WithAuth(&service.JwtVerifier)
	withProfileId := query.WithProfileId()
	withPostId := query.WithPostId()

	{
		profileGroup := restApi.Group("/profile")
		profileGroup.POST("", createHandler(
			func(qp *query.Params, r *api.RegisterProfileRequest) (api.RegisterProfileResponse, httperr.Err) {
				return service.RegisterProfile(qp, r)
			},
		))

		profileIdGroup := restApi.Group("/profile/:profile_id")
		profileIdGroup.Use(withProfileId)
		profileIdGroup.GET("", createHandler(
			func(qp *query.Params, r *empty) (api.GetProfileResponse, httperr.Err) {
				return service.GetProfileInfo(qp)
			},
		))

		profileIdGroup.PUT("", withAuth, createHandler(
			func(qp *query.Params, r *api.EditProfileRequest) (empty, httperr.Err) {
				return empty{}, service.EditProfileInfo(qp, r)
			},
		))

		profileIdGroup.DELETE("", withAuth, createHandler(
			func(qp *query.Params, r *empty) (empty, httperr.Err) {
				return empty{}, service.DeleteProfile(qp)
			},
		))
	}

	{
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
	}

	{
		restApi.GET("/profile/:profile_id/page/settings", withProfileId, createHandler(
			func(qp *query.Params, r *empty) (api.GetPageSettingsResponse, httperr.Err) {
				return service.GetPageSettings(qp)
			},
		))
		restApi.PUT("/profile/:profile_id/page/settings", withProfileId, withAuth, createHandler(
			func(qp *query.Params, r *api.EditPageSettingsRequest) (empty, httperr.Err) {
				return empty{}, service.EditPageSettings(qp, r)
			},
		))
	}

	{
		restApi.GET("/profile/:profile_id/page/posts", withProfileId, createHandler(
			func(qp *query.Params, r *api.GetPostsRequest) (api.GetPostsResponse, httperr.Err) {
				return service.GetPosts(qp, r)
			},
		))

		restApi.POST("/profile/:profile_id/page/posts", withProfileId, withAuth, createHandler(
			func(qp *query.Params, r *api.NewPostRequest) (api.NewPostResponse, httperr.Err) {
				return service.NewPost(qp, r)
			},
		))
	}

	{
		postGroup := restApi.Group("/post/:post_id")
		postGroup.Use(withPostId)
		postGroup.GET("", createHandler(
			func(qp *query.Params, r *empty) (api.Post, httperr.Err) {
				return service.GetPost(qp)
			},
		))
		postGroup.PUT("", withAuth, createHandler(
			func(qp *query.Params, r *api.EditPostRequest) (empty, httperr.Err) {
				return empty{}, service.EditPost(qp, r)
			},
		))
		postGroup.DELETE("", withAuth, createHandler(
			func(qp *query.Params, r *empty) (empty, httperr.Err) {
				return empty{}, service.DeletePost(qp)
			},
		))
	}

	{
		restApi.POST("/post/:post_id/comments", withPostId, withAuth, createHandler(
			func(qp *query.Params, r *api.NewCommentRequest) (api.NewCommentResponse, httperr.Err) {
				return service.NewComment(qp, r)
			},
		))
		restApi.GET("/post/:post_id/comments", withPostId, createHandler(
			func(qp *query.Params, r *api.GetCommentsRequest) (api.GetCommentsResponse, httperr.Err) {
				return service.GetComments(qp, r)
			},
		))

		restApi.POST("/post/:post_id/views", withPostId, withAuth, createHandler(
			func(qp *query.Params, r *empty) (empty, httperr.Err) {
				return empty{}, service.NewView(qp)
			},
		))
		restApi.POST("/post/:post_id/likes", withPostId, withAuth, createHandler(
			func(qp *query.Params, r *empty) (empty, httperr.Err) {
				return empty{}, service.NewLike(qp)
			},
		))
	}

	{
		restApi.GET("/post/:post_id/metric", withPostId, createHandler(
			func(qp *query.Params, r *api.GetPostMetricRequest) (api.GetPostMetricResponse, httperr.Err) {
				return service.GetPostMetric(qp, r)
			},
		))

		restApi.GET("/post/:post_id/metric_dynamics", withPostId, createHandler(
			func(qp *query.Params, r *api.GetPostMetricDynamicsRequest) (api.GetPostMetricDynamicsResponse, httperr.Err) {
				return service.GetPostMetricDynamics(qp, r)
			},
		))
	}

	{
		restApi.GET("/top10/posts", createHandler(
			func(qp *query.Params, r *api.GetTop10PostsRequest) (api.GetTop10PostsResponse, httperr.Err) {
				return service.GetTop10Posts(qp, r)
			},
		))

		restApi.GET("/top10/users", createHandler(
			func(qp *query.Params, r *api.GetTop10UsersRequest) (api.GetTop10UsersResponse, httperr.Err) {
				return service.GetTop10Users(qp, r)
			},
		))
	}

	return httpRouter{router}
}
