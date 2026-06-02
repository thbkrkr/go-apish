package middlewares

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
)

var AuthHeaderKey = "X-Auth"

func AuthMiddleware(apiKey string, accounts gin.Accounts) gin.HandlerFunc {
	basicAuth := gin.BasicAuthForRealm(accounts, "")

	return func(c *gin.Context) {
		// Try header auth, comparing in constant time to avoid leaking the
		// key through response-timing differences.
		got := c.Request.Header.Get(AuthHeaderKey)
		if subtle.ConstantTimeCompare([]byte(got), []byte(apiKey)) == 1 {
			return
		}
		// Fall back to basic auth
		basicAuth(c)
	}
}
