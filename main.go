package main // import "github.com/thbkrkr/go-apish"

import (
    "bytes"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os/exec"
    "time"
)

var argPort = flag.Int("port", 4242, "port to listen")
var apiKey = flag.String("apiKey", "42", "API key to authenticate")
var apiKeyHeader = flag.String("apiKeyHeader", "X-apish-auth", "API key header name to authenticate")

const ok = "{\"ok\": true, \"status\": 200}"
const scriptsDir = "./scripts"

func ls(w http.ResponseWriter, r *http.Request) {
  var buffer bytes.Buffer
  files, _ := ioutil.ReadDir("./scripts")
  
  hostname := "localhost"
  baseUrl := fmt.Sprintf("\"http://%v:%v/", hostname, *argPort)

  buffer.WriteString("[")
  for _, f := range files {
    buffer.WriteString(baseUrl)
    buffer.WriteString(f.Name())
    buffer.WriteString("\",")
  }
  buffer.WriteString("\"/ls\"]")

  io.WriteString(w, buffer.String())
  return
}

func biim(w http.ResponseWriter, r *http.Request) {

  script := fmt.Sprintf("./scripts%v", r.URL.String())

  // '/' => OK  
  if script == "./scripts/" {
    io.WriteString(w, ok)
    return
  }

  if r.Header.Get(*apiKeyHeader) != *apiKey {
     http.Error(w, "{\"error\":\"Invalid API key\"}", 401)
    return
  }

  stdout, err := exec.Command(script).Output()

  if err != nil {
    http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err.Error()), 500)
    return
  }

  io.WriteString(w, string(stdout))
}

func main() {
  flag.Parse()
  startTime := time.Now()

  // Register a ping
  http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, time.Since(startTime).String())
  })

  http.HandleFunc("/ls", ls)

  http.HandleFunc("/", biim)

  addr := fmt.Sprintf(":%v", *argPort)
  log.Printf("Magic server started on port %s", addr)
  log.Fatal(http.ListenAndServe(addr, nil))
}
