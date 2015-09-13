package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type resources struct {
	Scripts []string `json:"api"`
	Pages   []string `json:"html"`
}

func listResources(c *gin.Context) {
	scripts := make([]string, 0)
	pages := make([]string, 0)

	hostname := c.Request.Host
	hostname = strings.Replace(hostname, "/", "", -1)
	baseURL := fmt.Sprintf("http://%v", hostname)

	// List scripts
	err := filepath.Walk(*scriptsDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "sh") {
			url := fmt.Sprintf("%v/%v", baseURL, strings.Replace(path, ".sh", "", -1))
			scripts = append(scripts, url)
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	htmlDir := fmt.Sprintf("%s/_static", *scriptsDir)

	// List HTML files
	err = filepath.Walk(htmlDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			url := fmt.Sprintf("%v/%v", baseURL, strings.Replace(path, "api/_static", "s", -1))
			pages = append(pages, url)
		}
		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, resources{
		Scripts: scripts,
		Pages:   pages,
	})
}

func execScript(c *gin.Context) {
	path := c.Param("path")

	extension := ".sh"
	script := fmt.Sprintf("%s%s%s", *scriptsDir, path, extension)

	stdout, err := exec.Command(script).Output()

	if err != nil {
		serr := err.Error()
		log.Printf("Error executing `%s`: %s", path, serr)

		c.JSON(500, gin.H{
			"error": serr,
		})
		return
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(stdout, &objmap)

	//log.Printf("{\"script\": \"%s\", \"out\": %s}", script, stdout)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, objmap)
}
