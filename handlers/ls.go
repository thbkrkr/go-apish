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
	Static  []string `json:"static"`
}

func (h *LsHandler) ListResources(c *gin.Context) {
	scripts := make([]string, 0)
	pages := make([]string, 0)
	static := make([]string, 0)

	hostname := strings.Replace(c.Request.Host, "/", "", -1)

	// List scripts
	err := filepath.Walk(*h.ApiDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "sh") && !strings.Contains(path, "_static") {
			url := fileToUrl(hostname, "api", path, *h.ApiDir)
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

	staticDir := "_static"
	htmlDir := fmt.Sprintf("%s/%s", *h.ApiDir, staticDir)

	// List html files
	err = filepath.Walk(htmlDir, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, "html") {
			url := fileToUrl(hostname, "s", path, *h.ApiDir+"/_static")
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

	// List static files
	err = filepath.Walk(htmlDir, func(path string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() && !strings.HasSuffix(path, "html") {
			url := fileToUrl(hostname, "s", path, *h.ApiDir+"/_static")
			static = append(static, url)
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
		Static:  static,
	})
}

func fileToUrl(hostname string, prefix string, path string, apiDir string) string {
	// Remove ./ from apiDir
	apiDir = strings.Replace(apiDir, "./", "", -1)
	// Replace $apiDir by prefix
	filePath := strings.Replace(path, apiDir, prefix, -1)
	baseUrl := fmt.Sprintf("http://%v", hostname)

	if strings.Contains(path, "_static") {
		return fmt.Sprintf("%v/%v", baseUrl, filePath)
	} else {
		return fmt.Sprintf("%v/%v", baseUrl, strings.Replace(filePath, ".sh", "", -1))
	}

}
