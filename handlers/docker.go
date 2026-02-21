package handlers

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// DockerRun makes possible the execution of any docker run command
func DockerRun(c *gin.Context) {
	// Parse command in json body
	var form struct {
		Cmd string `json:"run"`
	}
	if err := c.BindJSON(&form); err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{"type": "error", "message": "Invalid docker run command"})
		return
	}

	// Invalid empty command
	if form.Cmd == "" {
		c.JSON(400, gin.H{"type": "error", "message": "Docker run command empty"})
		return
	}

	// Exec docker run
	args := append([]string{"run"}, strings.Split(form.Cmd, " ")...)
	output, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		message := err.Error() + ": " + strings.Replace(string(output), "\n", " ", -1)
		c.JSON(400, gin.H{"type": "error", "message": message})
		return
	}

	// Try to unmarshal the output to format json
	var obj interface{}
	err = json.Unmarshal(output, &obj)
	if err == nil {
		c.JSON(200, obj)
		return
	}

	c.String(200, string(output))
}
