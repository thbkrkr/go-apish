package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type ExecHandler struct {
	ApiDir *string
}

func (h *ExecHandler) ExecScript(c *gin.Context) {
	var stdout []byte
	var err error

	// Build script name
	path := c.Param("path")
	script := fmt.Sprintf("%s%s%s", *h.ApiDir, path, ".sh")

	// Check script exists
	if _, err := os.Stat(script); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"error": "Resource not found",
		})
		fmt.Printf("[error] resource not found: %s", script)
		return
	}

	// Exec script with or without param
	q := c.Request.URL.Query()
	param, isParam := q["q"]
	if isParam {
		stdout, err = exec.Command(script, param[0]).Output()
	} else {
		stdout, err = exec.Command(script).Output()
	}

	if err != nil {
		serr := err.Error()
		c.JSON(500, gin.H{
			"error": serr,
		})
		log.Printf("[error] executing `%s`: %s", path, serr)
		return
	}

	// Try to unmarshal JSON
	var someJson interface{}
	err = json.Unmarshal(stdout, &someJson)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		log.Printf("[error] invalid JSON for `%s`: %s", script, stdout)
		return
	}

	c.JSON(200, someJson)
	//log.Printf("[info] executing `%s`: %s", script, stdout)
}

func (h *ExecHandler) PostExecScript(c *gin.Context) {
	var stdout []byte
	var err error

	// Build script name
	path := c.Param("path")
	script := fmt.Sprintf("%s%s%s", *h.ApiDir, path, ".sh")

	// Check script exists
	if _, err := os.Stat(script); os.IsNotExist(err) {
		c.JSON(404, gin.H{
			"error": "Resource not found",
		})
		fmt.Printf("[error] resource not found: %s", script)
		return
	}

	// Exec script with or without body

	c1 := exec.Command(script)

	body := c.Request.Body
	var buf bytes.Buffer

	c1.Stdin = body
	c1.Stdout = &buf
	_ = c1.Start()
	_ = c1.Wait()

	if err != nil {
		serr := err.Error()
		c.JSON(500, gin.H{
			"error": serr,
		})
		log.Printf("[error] executing `%s`: %s", path, serr)
		return
	}

	stdout = buf.Bytes()

	// Try to unmarshal JSON
	var someJson interface{}
	err = json.Unmarshal(stdout, &someJson)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid JSON",
		})
		log.Printf("[error] invalid JSON for `%s`: %s", script, stdout)
		return
	}

	c.JSON(200, someJson)
	//log.Printf("[info] executing `%s`: %s", script, stdout)
}
