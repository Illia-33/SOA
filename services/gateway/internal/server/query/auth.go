package query

import (
	"fmt"
	"net/http"
	"regexp"
	"soa-socialnetwork/services/accounts/pkg/soajwt"
	"soa-socialnetwork/services/accounts/pkg/soatoken"
	"strings"

	"github.com/gin-gonic/gin"
)

var jwt_regexp = regexp.MustCompile(`^Bearer [\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}$`)
var soatoken_regexp = regexp.MustCompile(`^SoaToken [\-A-Za-z0-9\+\/_]={0,3}`)

func WithAuth(verifier *soajwt.Verifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer") {
			if !jwt_regexp.MatchString(auth) {
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
			params.AuthToken = rawToken
			params.AuthKind = AUTH_TOKEN_JWT
		} else if strings.HasPrefix(auth, "SoaToken") {
			if !soatoken_regexp.MatchString(auth) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, "broken soa token")
				return
			}

			rawToken := auth[9:] // skip 'SoaToken '
			token, err := soatoken.Parse(rawToken)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "broken soa token")
				return
			}

			profileId := ctx.Param("profile_id")
			if len(profileId) > 0 && token.ProfileID.String() != profileId {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, "not enough rights")
				return
			}

			params := ExtractParams(ctx)
			params.AuthToken = rawToken
			params.AuthKind = AUTH_TOKEN_SOA

		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "auth required")
			return
		}
	}
}
