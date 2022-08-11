package blog

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

	header := c.GetHeader("Authorization")
	if len(header) == 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Authorization header is incomplete")
		return
	}
	splittedheader := strings.Fields(header)
	if len(splittedheader) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Are you a hacker?")
		return
	}

	authorizationType := strings.ToLower(splittedheader[0])
	if authorizationType != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "You got the wrong agent")
		return

	}

	signedToken := splittedheader[1]

	ValidateToken(signedToken)

	c.Next()
}
