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
	args := append([]string{"run"}, splitArgs(form.Cmd)...)
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

// splitArgs splits a command string into arguments, honoring single and double
// quotes so that quoted arguments containing spaces stay intact and consecutive
// spaces don't produce empty arguments.
func splitArgs(s string) []string {
	args := make([]string, 0)
	var cur strings.Builder
	var quote rune
	inWord := false

	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
			inWord = true
		case r == '\'' || r == '"':
			quote = r
			inWord = true
		case r == ' ' || r == '\t':
			if inWord {
				args = append(args, cur.String())
				cur.Reset()
				inWord = false
			}
		default:
			cur.WriteRune(r)
			inWord = true
		}
	}
	if inWord {
		args = append(args, cur.String())
	}
	return args
}
