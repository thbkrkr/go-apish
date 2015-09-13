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

func MakeHttp(t *testing.T, verb string, path string, json string, auth *BasicAuth) (int, string) {
	reader := strings.NewReader(json)

	url := fmt.Sprintf("%s%s%s", ServerURL, path, PrefixURL)
	req, err := http.NewRequest(verb, url, reader)

	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
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

/*func Get(t *testing.T, path string) (int, string) {
	return MakeHttp(t, "GET", path, "", nil)
}

func Post(t *testing.T, path string, json string) (int, string) {
	return MakeHttp(t, "POST", path, json, nil)
}*/

func Get(t *testing.T, path string, auth *BasicAuth) (int, string) {
	return MakeHttp(t, "GET", path, "", auth)
}

func Post(t *testing.T, path string, json string, auth *BasicAuth) (int, string) {
	return MakeHttp(t, "POST", path, json, auth)
}
