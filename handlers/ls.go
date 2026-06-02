package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type LsHandler struct {
	ApiDir string
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

	hostname := strings.ReplaceAll(c.Request.Host, "/", "")
	baseURL := requestScheme(c) + "://" + hostname
	staticDir := h.ApiDir + "/_static"

	// List API scripts (every .sh outside _static), propagating walk errors.
	err := filepath.Walk(h.ApiDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "sh") && !strings.Contains(path, "_static") {
			scripts = append(scripts, fileToUrl(baseURL, "api", path, h.ApiDir))
		}
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// List static resources in a single pass, splitting HTML pages from other
	// files. The _static directory is optional, so skip it when absent.
	if _, statErr := os.Stat(staticDir); statErr == nil {
		err = filepath.Walk(staticDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			url := fileToUrl(baseURL, "s", path, staticDir)
			if strings.HasSuffix(path, "html") {
				pages = append(pages, url)
			} else {
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
	}

	c.JSON(200, resources{
		Scripts: scripts,
		Pages:   pages,
		Static:  static,
	})
}

// requestScheme reports whether the request arrived over https, honoring a
// reverse proxy's X-Forwarded-Proto header.
func requestScheme(c *gin.Context) string {
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}
	return "http"
}

func fileToUrl(baseURL string, prefix string, path string, apiDir string) string {
	// Remove ./ from apiDir
	apiDir = strings.ReplaceAll(apiDir, "./", "")
	// Replace $apiDir by prefix
	filePath := strings.ReplaceAll(path, apiDir, prefix)

	if strings.Contains(path, "_static") {
		return fmt.Sprintf("%v/%v", baseURL, filePath)
	}
	return fmt.Sprintf("%v/%v", baseURL, strings.ReplaceAll(filePath, ".sh", ""))
}
