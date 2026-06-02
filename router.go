package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	h "github.com/thbkrkr/go-apish/handlers"
	m "github.com/thbkrkr/go-apish/middlewares"
)

var basicAuthUser = "zuperadmin"

func Router() *gin.Engine {
	router := gin.Default()

	router.Use(m.CORSMiddleware())

	// Default routes (no auth: useful for health checks)
	router.GET("/", index)
	router.GET("/favicon.ico", favicon)
	router.GET("/version", version)

	// Authentication
	authorized := router.Group("/")

	if *password != "" {
		authorized = router.Group("/", m.AuthMiddleware(
			*apiKey,
			gin.Accounts{
				basicAuthUser: *password,
			},
		))
	} else {
		logrus.Warn("no -password set: authentication is DISABLED and all endpoints are publicly accessible")
	}

	lsHandler := &h.LsHandler{ApiDir: *apiDir}
	execHandler := &h.ExecHandler{ApiDir: *apiDir}

	// List resources
	authorized.GET("/ls", func(c *gin.Context) {
		lsHandler.ListResources(c)
	})

	// API propulsed by shell scripts
	authorized.GET("/api/*path", execHandler.ExecScript)
	authorized.POST("/api/*path", execHandler.PostExecScript)

	// Arbitrary `docker run` execution — off by default as it grants full
	// host access. Enable explicitly with -enableDocker.
	if *enableDocker {
		authorized.POST("/docker", h.DockerRun)
	}

	// Static files
	authorized.Static("/s/", *apiDir+"/_static")

	return router
}

/** Base routes */

func indexExists() bool {
	// Any stat error (not found, permission, ...) means we can't serve it,
	// so treat it as absent rather than redirecting to an unreadable file.
	_, err := os.Stat(*apiDir + "/_static/index.html")
	return err == nil
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
	c.Status(http.StatusNoContent)
}

func version(c *gin.Context) {
	c.JSON(200, gin.H{
		"git_commit": gitCommit,
		"build_date": buildDate,
	})
}
