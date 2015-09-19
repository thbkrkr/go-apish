package middlewares

import "github.com/gin-gonic/gin"

var AuthHeaderKey = "X-Auth"

func AuthMiddleware(apiKey string, accounts gin.Accounts) gin.HandlerFunc {
	basicAuth := gin.BasicAuthForRealm(accounts, "")

	return func(c *gin.Context) {
		// Try header auth
		if c.Request.Header.Get(AuthHeaderKey) == apiKey {
			return
		} else {
			// Try basic auth
			basicAuth(c)
		}
	}
}
