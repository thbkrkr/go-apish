package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	gitCommit = "undefined"
	buildDate = "undefined"

	port         = flag.Int("port", 4242, "HTTP port to listen")
	password     = flag.String("password", "", "Admin password for basic auth")
	apiKey       = flag.String("apiKey", "42", "API key for header auth")
	apiDir       = flag.String("apiDir", "./api", "API directory (sh scripts and html pages)")
	enableDocker = flag.Bool("enableDocker", false, "Enable the /docker endpoint (grants full host access via docker run)")
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	logrus.Infof("Running with %d CPUs", nuCPU)
}

func StartGin() {
	start := time.Now()
	gin.SetMode(gin.ReleaseMode)
	router := Router()
	sport := fmt.Sprintf(":%d", *port)

	s := &http.Server{
		Addr:           sport,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logrus.Infof("API ready in %v, listening on %s", time.Since(start), sport)

	// ListenAndServe only returns on error; log it and exit instead of
	// silently busy-looping a restart.
	if err := s.ListenAndServe(); err != nil {
		logrus.Fatalf("server stopped: %v", err)
	}
}

func main() {
	flag.Parse()
	ConfigRuntime()
	StartGin()
}
