package handlers

import "testing"

func TestScriptPathStaysWithinApiDir(t *testing.T) {
	h := &ExecHandler{ApiDir: "example/api"}

	cases := map[string]bool{
		"/time/date":            true,  // normal script
		"/test/param":           true,  // nested script
		"/../secret":            false, // climbs out of ApiDir
		"/../../../etc/passwd":  false, // deep traversal
		"/time/../../../tmp/x":  false, // traversal after a valid segment
	}

	for reqPath, wantOK := range cases {
		if _, ok := h.scriptPath(reqPath); ok != wantOK {
			t.Errorf("scriptPath(%q) ok = %v, want %v", reqPath, ok, wantOK)
		}
	}
}
