package middlewares

import (
	"crypto/subtle"

	"github.com/gin-gonic/gin"
)

var AuthHeaderKey = "X-Auth"

func AuthMiddleware(apiKey string, accounts gin.Accounts) gin.HandlerFunc {
	basicAuth := gin.BasicAuthForRealm(accounts, "")

	return func(c *gin.Context) {
		// Try header auth when a key is configured, comparing in constant time
		// to avoid leaking the key through response-timing differences. An
		// empty key disables header auth (so it can't match a missing header).
		if apiKey != "" {
			got := c.Request.Header.Get(AuthHeaderKey)
			if subtle.ConstantTimeCompare([]byte(got), []byte(apiKey)) == 1 {
				return
			}
		}
		// Fall back to basic auth
		basicAuth(c)
	}
}
