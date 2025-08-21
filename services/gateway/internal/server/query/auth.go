package query

import (
	"fmt"
	"net/http"
	"regexp"
	"soa-socialnetwork/internal/soajwt"

	"github.com/gin-gonic/gin"
)

func WithJwtAuth(verifier *soajwt.Verifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		authRegex := regexp.MustCompile(`^Bearer [\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}$`)
		if !authRegex.MatchString(auth) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "broken jwt token")
			return
		}

		rawToken := auth[7:] // skip 'Bearer '
		token, err := verifier.Verify(rawToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("bad jwt: %v", err))
			return
		}

		profileId := ctx.Param("profile_id")
		if len(profileId) > 0 && token.Subject != profileId {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "not enough rights")
			return
		}

		params := ExtractParams(ctx)
		params.RawJwtToken = rawToken
	}
}
