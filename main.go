package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	gitCommit = "undefined"
	buildDate = "undefined"

	port          = flag.Int("port", 4242, "HTTP port to listen")
	adminPassword = flag.String("adminPassword", "42", "Admin password")
	scriptsDir    = flag.String("scriptsDir", "./api", "API directory (sh scripts and html pages)")
)

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("[info] Running with %d CPUs\n", nuCPU)
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

	time.Sleep(time.Second)
	log.Printf("[info] Magic API started in %v on %s\n", time.Since(start), sport)

	for {
		s.ListenAndServe()
	}
}

func main() {
	ConfigRuntime()
	StartGin()
}
