package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type LsHandler struct {
	ApiDir *string
}

type resources struct {
	Scripts []string `json:"api"`
	Pages   []string `json:"html"`
}

func (h *LsHandler) ListResources(c *gin.Context) {
	scripts := make([]string, 0)
	pages := make([]string, 0)

	hostname := c.Request.Host
	hostname = strings.Replace(hostname, "/", "", -1)
	baseUrl := fmt.Sprintf("http://%v", hostname)

	// List scripts
	err := filepath.Walk(*h.ApiDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "sh") {
			scriptUrl := fmt.Sprintf("%v/%v", baseUrl, strings.Replace(path, ".sh", "", -1))
			scripts = append(scripts, scriptUrl)
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	staticDir := "_static"
	htmlDir := fmt.Sprintf("%s/%s", *h.ApiDir, staticDir)

	// List static files
	err = filepath.Walk(htmlDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			staticPath := strings.Replace(path, "api/"+staticDir, "s", -1)
			url := fmt.Sprintf("%v/%v", baseUrl, staticPath)
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
