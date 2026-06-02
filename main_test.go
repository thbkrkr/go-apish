package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	test "github.com/thbkrkr/go-apish/test"
)

var server *httptest.Server

var auth = &test.BasicAuth{Username: "zuperadmin", Password: "42"}

func init() {
	gin.SetMode(gin.TestMode)
	*apiDir = "example/api"
	*password = "42"
	*apiKey = "42"
	server = httptest.NewServer(Router())

	test.ServerURL = server.URL
}

func TestBase(t *testing.T) {
	test.PrefixURL = ""

	// / redirects to /s (index.html exists in the example), which is behind
	// auth, so the request must carry credentials to follow through to 200.
	status, _ := test.Get(t, "/", auth)
	assert.Equal(t, 200, status, "should get a 200")

	// /version is exposed without auth.
	status, _ = test.Get(t, "/version", nil)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/favicon.ico", nil)
	assert.Equal(t, 204, status, "should get a 204")

	status, _ = test.Get(t, "/blablabla", nil)
	assert.Equal(t, 404, status, "should get a 404")
}

func TestAuthentication(t *testing.T) {
	status, _ := test.Get(t, "/api/time/date", nil)
	assert.Equal(t, 401, status, "should get a 401")

	auth := &test.BasicAuth{Username: "zuperadmin", Password: "42"}
	status, _ = test.Get(t, "/api/time/date", auth)
	assert.Equal(t, 200, status, "should get a 200")

	apiKey := new(string)
	*apiKey = "42"
	status, _ = test.Get2(t, "/api/time/date", apiKey)
	assert.Equal(t, 200, status, "should get a 200")
}

func TestScripts(t *testing.T) {
	status, _ := test.Get(t, "/api/nothing", auth)
	assert.Equal(t, 404, status, "should get a 404")

	status, _ = test.Get(t, "/api/time/date", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/api/test/param?q=hello", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/api/time/date", auth)
	assert.Equal(t, 200, status, "should get a 200")
}

func TestPages(t *testing.T) {
	status, _ := test.Get(t, "/s/", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/s/date.html", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/s/css/styles.css", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/s/css/styles.css", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/s/js/script.js", auth)
	assert.Equal(t, 200, status, "should get a 200")
}

func TestPost(t *testing.T) {
	status, body := test.Post(t, "/api/test/post", `{"o": 42}`, auth)
	assert.Equal(t, 200, status, "should get a 200")
	assert.Contains(t, body, "jackpot")
}

func TestInvalidJSON(t *testing.T) {
	status, _ := test.Get(t, "/api/test/invalid-json", auth)
	assert.Equal(t, 400, status, "invalid JSON output should yield a 400")
}

func TestListResources(t *testing.T) {
	status, body := test.Get(t, "/ls", auth)
	assert.Equal(t, 200, status, "should get a 200")
	assert.Contains(t, body, "/api/time/date")
}

func TestPathTraversal(t *testing.T) {
	// Attempting to escape apiDir must not execute an arbitrary script.
	status, _ := test.Get(t, "/api/../../../../etc/hostname", auth)
	assert.NotEqual(t, 200, status, "path traversal must be rejected")
}
