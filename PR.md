## Summary

- **Security:** prevent path traversal in script execution; fix wildcard+credentials CORS combo; compare API key in constant time; gate `/docker` behind a flag then remove it entirely; configurable basic-auth user with no default API key; warn when running unauthenticated
- **Fixes:** handle errors in POST handler; return 204 for favicon; treat any stat error as missing index; log and exit on server error; derive `/ls` URL scheme from the request
- **Refactor:** expose `/version` without auth; deduplicate GET/POST exec handlers; unify logging on logrus; single static walk with real error propagation; drop favicon route and slim CORS headers
- **Build:** multi-stage Dockerfile with Go 1.25; drop dead `release` target
- **Docs:** KISS README rewrite with inline flag comments, real curl output, layout tree; example scripts build JSON safely with `jq`
- **Tests:** add coverage for path traversal, POST, invalid JSON, `/ls`, script failure, wrong credentials; fix broken `apiDir` path after restructuring; simplify HTTP helpers
- **Helm:** minimal chart (Deployment + Service) with configurable image, port, auth flags; `helm-*` Makefile targets using `apish` as release name
- **Build:** bump Go 1.25 → 1.26.4 and alpine 3.20 → 3.23; tag example image as `krkr/apish:example`; add `push` and `port-forward` targets to example Makefile

## Test plan

- [ ] `go test ./app/...` passes
- [ ] `helm template` renders without errors (`make helm-render`)
- [ ] `./go-apish -apiDir=example` serves scripts and static files
- [ ] Auth is required when `-password` is set; server warns when it is not
- [ ] `make -C example build && make -C example run` serves the example image on port 80
