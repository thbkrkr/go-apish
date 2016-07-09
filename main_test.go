package main

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	test "github.com/thbkrkr/go-apish/test"
)

var (
	server *httptest.Server
	reader io.Reader //Ignore this for now
)

var auth = &test.BasicAuth{"zuperadmin", "42"}

func init() {
	gin.SetMode(gin.TestMode)
	*apiDir = "example/api"
	*password = "42"
	server = httptest.NewServer(Router())

	test.ServerURL = server.URL
}

func TestBase(t *testing.T) {
	test.PrefixURL = ""
	status, _ := test.Get(t, "/", nil)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/favicon.ico", nil)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = test.Get(t, "/blablabla", nil)
	assert.Equal(t, 404, status, "should get a 404")
}

func TestAuthentication(t *testing.T) {
	status, _ := test.Get(t, "/api/time/date", nil)
	assert.Equal(t, 401, status, "should get a 401")

	auth := &test.BasicAuth{"zuperadmin", "42"}
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
