package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := "*"
		c.Writer.Header().Set("Access-Control-Allow-Origin", domain)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	// Default routes
	router.GET("/", index)
	router.GET("/favicon.ico", favicon)

	// Basic authentication
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"zuperadmin": *adminPassword,
	}))

	// Version (commit and date)
	authorized.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"git_commit": gitCommit,
			"build_date": buildDate,
		})
	})

	// List scripts
	authorized.GET("/ls", func(c *gin.Context) {
		listResources(c)
	})

	// Scripts API
	authorized.GET("/api/*path", execScript)

	// Static HTML
	htmlDir := fmt.Sprintf("%s/_static", *scriptsDir)
	authorized.Static("/s/", htmlDir)

	return router
}

/** Base routes */

func index(c *gin.Context) {
	c.JSON(200, gin.H{
		"ok":     true,
		"status": 200,
		"name":   "go-apish",
	})
}

func favicon(c *gin.Context) {
	c.JSON(200, nil)
}
