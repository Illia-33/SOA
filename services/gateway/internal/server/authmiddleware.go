package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"soa-socialnetwork/internal/soajwt"

	"github.com/gin-gonic/gin"
)

func authMiddleware(verifier *soajwt.Verifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		authRegex := regexp.MustCompile(`^Bearer [\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}\.[\-A-Za-z0-9\+\/_]*={0,3}$`)
		log.Printf("suka: %s", auth)
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

		id := ctx.Param("id")
		if len(id) > 0 && token.Subject != id {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "not enough rights")
			return
		}
	}
}
