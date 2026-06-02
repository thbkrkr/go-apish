package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

var (
	ServerURL string
	PrefixURL string
)

type BasicAuth struct {
	Username string
	Password string
}

func do(t *testing.T, verb, path, body string, auth *BasicAuth, apiKey string) (int, string) {
	url := fmt.Sprintf("%s%s%s", ServerURL, path, PrefixURL)
	req, err := http.NewRequest(verb, url, strings.NewReader(body))
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

func Get(t *testing.T, path string, auth *BasicAuth) (int, string) {
	return do(t, "GET", path, "", auth, "")
}

func GetWithKey(t *testing.T, path, apiKey string) (int, string) {
	return do(t, "GET", path, "", nil, apiKey)
}

func Post(t *testing.T, path, body string, auth *BasicAuth) (int, string) {
	return do(t, "POST", path, body, auth, "")
}
