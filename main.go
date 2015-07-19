package main // import "github.com/thbkrkr/go-apish"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/thbkrkr/go-apish/pkg/handlers"
)

var (
	buildTime    = "undefined"
	buildSha1    = "undefined"
	argPort      = flag.Int("port", 4242, "port to listen")
	apiKey       = flag.String("apiKey", "42", "API key to authenticate")
	apiKeyHeader = flag.String("apiKeyHeader", "X-apish-auth", "API key header name to authenticate")
)

const scriptsDir = "./scripts"

func listScripts(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer
	files, _ := ioutil.ReadDir(scriptsDir)

	hostname := "localhost"
	baseURL := fmt.Sprintf("http://%v:%v", hostname, *argPort)

	buffer.WriteString("[")
	for _, f := range files {
		url := fmt.Sprintf("\"%v/%v\",", baseURL, f.Name())
		buffer.WriteString(url)
	}
	buffer.WriteString("\"/ls\"]")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	io.WriteString(w, buffer.String())
	return
}

func version(w http.ResponseWriter, r *http.Request) {
	v := fmt.Sprintf("{\"version\":\"%v.%v\"}", buildSha1, buildTime)
	io.WriteString(w, string(v))
}

//ype appHandler func(http.ResponseWriter, *http.Request) (int, error)

type appHandler func(http.ResponseWriter, *http.Request)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(*apiKeyHeader) != *apiKey {
		http.Error(w, "{\"error\":\"Invalid API key\"}", 401)
		return
	}
	fn(w, r)

	/*if status, err := fn(w, r); err != nil {
		// We could also log our errors centrally:
		// i.e. log.Printf("HTTP %d: %v", err)
		switch status {
		// We can have cases as granular as we like, if we wanted to
		// return custom errors for specific status codes.
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}*/
}

func main() {
	flag.Parse()
	startTime := time.Now()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		uptime := fmt.Sprintf("{\"uptime\":\"%v\"}", time.Since(startTime))
		io.WriteString(w, string(uptime))
	})

	mux.HandleFunc("/version", version)
	mux.HandleFunc("/ls", listScripts)

	mux.HandleFunc("/", appHandler(handlers.ExecScript))

	/*	mux.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(*apiKeyHeader) != *apiKey {
				http.Error(w, "{\"error\":\"Invalid API key\"}", 401)
				return
			}
			handlers.ExecScript(w, r)
		})
	*/

	http.Handle("/", mux)
	addr := fmt.Sprintf(":%v", *argPort)
	log.Printf("Magic server started on port %d", *argPort)
	log.Fatal(http.ListenAndServe(addr, nil))
}
