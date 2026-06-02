# Code Review — go-apish

Review date: 2026-06-02

`go-apish` is a small Gin-based HTTP server that exposes shell scripts (and static
files) as a REST API. The idea is neat and the codebase is small and readable. The
notes below are grouped by severity. The most important section is **Security** — by
its very nature this project runs arbitrary shell commands, so the attack surface
deserves the most care.

---

## 🔴 Security (high priority)

### 1. Path traversal → arbitrary script execution — `handlers/exec.go:23-24`
```go
path   := c.Param("path")
script := fmt.Sprintf("%s%s%s", *h.ApiDir, path, ".sh")
```
The wildcard `*path` is concatenated straight onto `ApiDir` with no validation. A
request such as `GET /api/../../../../tmp/evil` resolves to `./api/../../../../tmp/evil.sh`,
letting a caller execute any `.sh` file on the filesystem outside `apiDir`. Mitigate by:
- `filepath.Clean` the joined path and verify it still has `apiDir` as a prefix
  (use `filepath.Abs` on both and check `strings.HasPrefix`), or
- reject any path containing `..`.

This applies to both `ExecScript` and `PostExecScript` (duplicated logic).

### 2. `/docker` runs arbitrary `docker run` — `handlers/docker.go:31-32`
```go
args := append([]string{"run"}, strings.Split(form.Cmd, " ")...)
output, err := exec.Command("docker", args...).CombinedOutput()
```
Any authenticated caller can run any container with any flags
(`-v /:/host`, `--privileged`, `--pid=host`, …), which is effectively root on the
host. This may be intentional for the tool's purpose, but it should be:
- gated behind an explicit opt-in flag (off by default), and
- clearly documented as "grants full host access."

Also note `strings.Split(cmd, " ")` breaks on quoted arguments and multiple spaces;
a real shell-style tokenizer (e.g. `shellwords`) would be more correct if this stays.

### 3. CORS allows any origin *with* credentials — `middlewares/cors.go:8,13`
```go
c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
```
`Allow-Origin: *` together with `Allow-Credentials: true` is rejected by browsers and
is a misconfiguration. If credentials are needed, reflect a vetted origin instead of
`*`. If not, drop the credentials header. For an unauthenticated, open API the wildcard
is fine — but then the basic-auth story below conflicts with it.

### 4. Auth middleware comparisons are not constant-time — `middlewares/auth.go:12`
```go
if c.Request.Header.Get(AuthHeaderKey) == apiKey {
```
The API-key check uses `==`, which is vulnerable to timing attacks. Use
`crypto/subtle.ConstantTimeCompare`. Minor for a hobby tool, but trivial to fix.

### 5. Auth is opt-in — `router.go:26`
Authentication is only enabled when `-password` is set. The default is a fully open
server that can execute scripts. Consider failing closed (require a password, or print
a loud warning at startup when running unauthenticated).

### 6. Example scripts demonstrate shell injection — `example/api/test/param.sh`
```sh
echo '{ "param": "'$1'" }'
```
Because Go's `exec.Command` does not invoke a shell, the Go side is safe from command
injection — but `$1` unquoted/unescaped inside the script means a value like `"} ...`
breaks the JSON, and any script that does `eval`/backticks on `$1` would be exploitable.
Since these are the canonical examples users copy, they should model safe quoting
(e.g. build JSON with `jq -n --arg p "$1" '{param:$p}'`).

---

## 🟠 Correctness / bugs

### 7. `PostExecScript` ignores all execution errors — `handlers/exec.go:88-105`
```go
c1 := exec.Command(script)
...
_ = c1.Start()
_ = c1.Wait()
if err != nil {   // err is still nil here — declared at top, never assigned
```
The `err` checked on line 98 is the zero-value from line 71; the real errors from
`Start()`/`Wait()` are discarded with `_ =`. A script that fails will fall through and
likely produce an "Invalid JSON" 400 instead of a 500 with the real error. Capture and
check those errors. Also `c1.Stderr` is never set, so stderr is lost.

### 8. `favicon` returns JSON `null` with 200 — `router.go:81-83`
`c.JSON(200, nil)` sends `null` as a favicon. Harmless but odd; returning `204 No Content`
would be cleaner.

### 9. `indexExists` swallows non-NotExist errors — `router.go:60-67`
If `os.Stat` fails for a reason other than "not exist" (e.g. permissions), the function
returns `true` and a redirect is issued for a file that can't be read. Treat any error as
"not present."

### 10. Server restart loop hides crashes — `main.go:46-48`
```go
for {
    s.ListenAndServe()
}
```
`ListenAndServe` only returns on error; the bare `for` loop silently restarts it (busy-loop
if the port is unavailable) and the returned error is never logged. Log the error and exit,
or restart with backoff. Also, the "API started" log on line 44 prints *before*
`ListenAndServe` is called, so it's measuring `Router()` build time, not real startup.

### 11. `version` route is inside the authorized group — `router.go:36`
Minor: a version/health endpoint is usually fine to expose unauthenticated; currently it's
behind auth when a password is set, which complicates health checks.

---

## 🟡 Code quality / maintainability

### 12. `ExecScript` and `PostExecScript` are ~90% duplicated — `handlers/exec.go`
The path-building, existence check, JSON-unmarshal, and error responses are copy-pasted.
Extract a helper like `runScript(c, script, stdin io.Reader)` and a shared
`resolveScript(path) (string, bool)`. This also means the path-traversal fix (item 1) only
has to be made once.

### 13. Inconsistent logging — `fmt.Printf` vs `log.Printf` vs `logrus`
The project uses three logging styles: `fmt.Printf` (`exec.go:31,82`), `log.Printf`
(`exec.go:49`), and `logrus` (`docker.go`). Pick one (logrus is already a dependency) and
use structured, leveled logging consistently. The `fmt.Printf` calls also omit trailing
`\n`.

### 14. Resolving `*string` flags via pointers passed into handlers — `router.go:38-39`
`LsHandler{ApiDir: apiDir}` stores a `*string` to a global flag. It works, but passing the
resolved `string` value (after `flag.Parse`) is simpler and avoids handlers depending on
global mutable state — which is exactly what makes the test on `main_test.go:23` have to
poke `*apiDir` directly.

### 15. `ListResources` walks `_static` twice — `handlers/ls.go:48-69`
Two `filepath.Walk` passes over the same `htmlDir`, one for `.html` and one for everything
else. A single walk with a branch would halve the I/O and the code. The `err` from the
first scripts walk (line 30) can never be non-nil because the walk func always returns
`nil` — so the error checks are dead code unless the walk func propagates `err`.

### 16. `fileToUrl` hardcodes `http://` — `handlers/ls.go:89`
Generated URLs are always `http://`, so behind TLS/a proxy the listed links are wrong.
Derive the scheme from the request (`X-Forwarded-Proto` / `c.Request.TLS`).

### 17. Magic values / globals — `router.go:12`
`basicAuthUser = "zuperadmin"` is a hardcoded global; the API key default `"42"`
(`main.go:20`) is a weak, shipped default. Make the username configurable and avoid a
guessable default key (or require it to be set).

---

## 🟢 Build / tooling / docs

### 18. `go.mod` says `go 1.25.1` but Makefile builds with `golang:1.6.2` — `Makefile:11`
These are wildly out of sync. `golang:1.6.2` predates modules entirely and cannot build a
`go 1.25` module. The Dockerfile (`alpine:3.7`) is also from 2018 and has known CVEs. Update
to a current Go builder image and a maintained base (e.g. multi-stage build on `golang:1.25`
+ `alpine:3.20` or distroless).

### 19. `release.sh` referenced but missing — `Makefile:19`
`make release` calls `./release.sh` which isn't in the repo. Either add it or drop the target.

### 20. README is thin
No mention of the auth model, the `-apiKey`/`-password`/`-port`/`-apiDir` flags, the
`/docker` endpoint's risks, or the JSON contract (scripts must emit valid JSON or callers
get a 400). A short "Endpoints" and "Security" section would help a lot.

### 21. Tests don't cover the risky paths
`main_test.go` is a good start but there are no tests for: POST script execution, the
`/docker` endpoint, invalid-JSON handling (a 400 case), or — most importantly — path
traversal. Add a test asserting `/api/../../something` is rejected; it will fail today and
guard the fix for item 1.

### 22. Struct literal without field names — `main_test.go:19`
```go
var auth = &test.BasicAuth{"zuperadmin", "42"}
```
`go vet` flags unkeyed struct literals. Use `&test.BasicAuth{Username: ..., Password: ...}`.

---

## Summary

The architecture is clean and the intent is clear. Priorities:

1. **Fix path traversal** (item 1) — this is the one outright vulnerability that isn't
   "by design."
2. **Gate / document `/docker`** (item 2) and **fix the swallowed POST errors** (item 7).
3. **Modernize the build** (item 18) so the project actually builds reproducibly.
4. Then the deduplication (item 12) and logging cleanup (item 13) for maintainability.
