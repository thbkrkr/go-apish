package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func ExecScript(w http.ResponseWriter, r *http.Request) { //(int, error) {

	url := r.URL.String()
	log.Printf("URL %s", url)

	script := fmt.Sprintf("./scripts%v", r.URL.String())
	log.Printf("Exec script %s", script)

	// '/' => OK
	if script == "./scripts/" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		io.WriteString(w, "{\"ok\": true, \"status\": 200}")
		return
	}

	// Script execution
	stdout, err := exec.Command(script).Output()

	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\":\"%v\"}", err.Error()), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	io.WriteString(w, string(stdout))
}
