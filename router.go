package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	h "github.com/thbkrkr/go-apish/handlers"
	m "github.com/thbkrkr/go-apish/middlewares"
)

var basicAuthUser = "zuperadmin"

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(m.CORSMiddleware())

	// Default routes
	router.GET("/", index)
	router.GET("/favicon.ico", favicon)

	// Authentication
	authorized := router.Group("/")

	if *password != "" {
		authorized = router.Group("/", m.AuthMiddleware(
			*apiKey,
			gin.Accounts{
				basicAuthUser: *password,
			},
		))
	}

	// Version (commit and date)
	authorized.GET("/version", version)

	lsHandler := &h.LsHandler{ApiDir: apiDir}
	execHandler := &h.ExecHandler{ApiDir: apiDir}

	// List resources
	authorized.GET("/ls", func(c *gin.Context) {
		lsHandler.ListResources(c)
	})

	// API propulsed by shell scripts
	authorized.GET("/api/*path", execHandler.ExecScript)
	authorized.POST("/api/*path", execHandler.PostExecScript)

	authorized.POST("/docker", h.DockerRun)

	// Static files
	authorized.Static("/s/", *apiDir+"/_static")

	return router
}

/** Base routes */

func indexExists() bool {
	if _, err := os.Stat(*apiDir + "/_static/index.html"); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func index(c *gin.Context) {
	if indexExists() {
		c.Redirect(http.StatusMovedPermanently, "/s")
	} else {
		c.JSON(200, gin.H{
			"ok":     true,
			"status": 200,
			"name":   "go-apish",
		})
	}
}

func favicon(c *gin.Context) {
	c.JSON(200, nil)
}

func version(c *gin.Context) {
	c.JSON(200, gin.H{
		"git_commit": gitCommit,
		"build_date": buildDate,
	})
}
