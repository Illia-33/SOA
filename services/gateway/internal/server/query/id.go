package query

import "github.com/gin-gonic/gin"

func WithProfileID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params := ExtractParams(ctx)
		profileId := ctx.Param("profile_id")
		params.ProfileId = profileId
	}
}
