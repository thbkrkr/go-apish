package handlers

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ExecHandler struct {
	ApiDir *string
}

// scriptPath resolves the wildcard request path to an absolute `.sh` path and
// guarantees it stays within ApiDir, preventing path-traversal escapes such as
// `/api/../../../tmp/evil`. It returns false when the path escapes ApiDir.
func (h *ExecHandler) scriptPath(reqPath string) (string, bool) {
	base, err := filepath.Abs(*h.ApiDir)
	if err != nil {
		return "", false
	}
	// filepath.Join cleans the result, collapsing any `..` segments.
	script := filepath.Join(base, reqPath+".sh")
	if script != base && !strings.HasPrefix(script, base+string(os.PathSeparator)) {
		return "", false
	}
	return script, true
}

// resolve validates the request path and ensures the target script exists,
// writing a 404 and returning ok=false when it cannot be served.
func (h *ExecHandler) resolve(c *gin.Context) (string, bool) {
	path := c.Param("path")
	script, ok := h.scriptPath(path)
	if !ok {
		c.JSON(404, gin.H{"error": "Resource not found"})
		logrus.Errorf("invalid resource path: %s", path)
		return "", false
	}
	if _, err := os.Stat(script); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Resource not found"})
		logrus.Errorf("resource not found: %s", script)
		return "", false
	}
	return script, true
}

// run executes cmd, which is expected to print JSON to stdout, and writes the
// parsed JSON (or an appropriate error) to the response.
func (h *ExecHandler) run(c *gin.Context, script string, cmd *exec.Cmd) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		logrus.Errorf("executing `%s`: %s: %s", script, err, stderr.String())
		return
	}

	var payload any
	if err := json.Unmarshal(stdout.Bytes(), &payload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		logrus.Errorf("invalid JSON for `%s`: %s", script, stdout.Bytes())
		return
	}

	c.JSON(200, payload)
}

// ExecScript runs the script for a GET request, optionally passing the `q`
// query parameter as the script's first argument.
func (h *ExecHandler) ExecScript(c *gin.Context) {
	script, ok := h.resolve(c)
	if !ok {
		return
	}

	cmd := exec.Command(script)
	if q, isParam := c.Request.URL.Query()["q"]; isParam {
		cmd = exec.Command(script, q[0])
	}
	h.run(c, script, cmd)
}

// PostExecScript runs the script for a POST request, piping the request body
// to the script's stdin.
func (h *ExecHandler) PostExecScript(c *gin.Context) {
	script, ok := h.resolve(c)
	if !ok {
		return
	}

	cmd := exec.Command(script)
	cmd.Stdin = c.Request.Body
	h.run(c, script, cmd)
}
