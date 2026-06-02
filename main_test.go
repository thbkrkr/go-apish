package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

var auth = &basicAuth{Username: "zuperadmin", Password: "42"}

func init() {
	gin.SetMode(gin.TestMode)
	*apiDir = "example/api"
	*password = "42"
	*apiKey = "42"
	server = httptest.NewServer(Router())
}

type basicAuth struct {
	Username string
	Password string
}

func do(t *testing.T, verb, path, body string, auth *basicAuth, apiKey string) (int, string) {
	req, err := http.NewRequest(verb, server.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	if apiKey != "" {
		req.Header.Set("X-Auth", apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	out, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, string(out)
}

func get(t *testing.T, path string, auth *basicAuth) (int, string) {
	return do(t, "GET", path, "", auth, "")
}

func getWithKey(t *testing.T, path, apiKey string) (int, string) {
	return do(t, "GET", path, "", nil, apiKey)
}

func post(t *testing.T, path, body string, auth *basicAuth) (int, string) {
	return do(t, "POST", path, body, auth, "")
}

func TestBase(t *testing.T) {
	// / redirects to /s (index.html exists in the example), which is behind
	// auth, so the request must carry credentials to follow through to 200.
	status, _ := get(t, "/", auth)
	assert.Equal(t, 200, status, "should get a 200")

	// /version is exposed without auth.
	status, _ = get(t, "/version", nil)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = get(t, "/blablabla", nil)
	assert.Equal(t, 404, status, "should get a 404")
}

func TestAuthentication(t *testing.T) {
	status, _ := get(t, "/api/time/date", nil)
	assert.Equal(t, 401, status, "should get a 401")

	status, _ = get(t, "/api/time/date", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = getWithKey(t, "/api/time/date", "42")
	assert.Equal(t, 200, status, "should get a 200")
}

func TestScripts(t *testing.T) {
	status, _ := get(t, "/api/nothing", auth)
	assert.Equal(t, 404, status, "should get a 404")

	status, _ = get(t, "/api/time/date", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = get(t, "/api/test/param?q=hello", auth)
	assert.Equal(t, 200, status, "should get a 200")
}

func TestPages(t *testing.T) {
	status, _ := get(t, "/s/", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = get(t, "/s/date.html", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = get(t, "/s/css/styles.css", auth)
	assert.Equal(t, 200, status, "should get a 200")

	status, _ = get(t, "/s/js/script.js", auth)
	assert.Equal(t, 200, status, "should get a 200")
}

func TestPost(t *testing.T) {
	status, body := post(t, "/api/test/post", `{"o": 42}`, auth)
	assert.Equal(t, 200, status, "should get a 200")
	assert.Contains(t, body, "jackpot")
}

func TestInvalidJSON(t *testing.T) {
	status, _ := get(t, "/api/test/invalid-json", auth)
	assert.Equal(t, 400, status, "invalid JSON output should yield a 400")
}

func TestListResources(t *testing.T) {
	status, body := get(t, "/ls", auth)
	assert.Equal(t, 200, status, "should get a 200")
	assert.Contains(t, body, "/api/time/date")
}

func TestPathTraversal(t *testing.T) {
	// Attempting to escape apiDir must not execute an arbitrary script.
	status, _ := get(t, "/api/../../../../etc/hostname", auth)
	assert.NotEqual(t, 200, status, "path traversal must be rejected")
}
