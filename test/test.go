package test

import (
	"fmt"
	"io/ioutil"
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

func MakeHttp(t *testing.T, verb string, path string, json string, auth *BasicAuth, apiKey *string) (int, string) {
	reader := strings.NewReader(json)

	url := fmt.Sprintf("%s%s%s", ServerURL, path, PrefixURL)
	req, err := http.NewRequest(verb, url, reader)

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	if apiKey != nil {
		req.Header.Set("X-Auth", *apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}

	if resp.Body == nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	return resp.StatusCode, string(body)
}

func Get(t *testing.T, path string, auth *BasicAuth) (int, string) {
	return MakeHttp(t, "GET", path, "", auth, nil)
}

func Get2(t *testing.T, path string, apiKey *string) (int, string) {
	return MakeHttp(t, "GET", path, "", nil, apiKey)
}

func Post(t *testing.T, path string, json string, auth *BasicAuth) (int, string) {
	return MakeHttp(t, "POST", path, json, auth, nil)
}
