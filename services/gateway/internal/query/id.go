package query

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func WithProfileId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := ExtractParams(ctx)
		profileId := ctx.Param("profile_id")
		params.ProfileId = profileId
	}
}

func WithPostId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := ExtractParams(ctx)
		postIdStr := ctx.Param("post_id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("bad post id"))
			return
		}
		params.PostId = int32(postId)
	}
}
