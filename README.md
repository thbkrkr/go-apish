# apish — REST API for shell scripts

Write shell scripts that return JSON ([example](example/time/date.sh)).

Serve static files from [_static](example/_static) directory.

## Run

```sh
./go-apish \
  -port=4242 \          # HTTP port
  -apiDir=example/api \ # directory of .sh scripts and _static files
  -user=zuperadmin \    # basic-auth username
  -password=secret \    # basic-auth password (empty = no auth)
  -apiKey=mykey         # X-Auth header key (empty = disabled)
```

```sh
make binary   # build ./go-apish
make build    # build krkr/apish Docker image
make run      # run image, mounting ./example as /api on port 80
```

## Example

### Layout

```
<apiDir>/
  time/date.sh        →  GET  /api/time/date
  test/param.sh       →  GET  /api/test/param?q=<value>
  test/post.sh        →  POST /api/test/post
  _static/index.html  →  GET  /s/
```

### Endpoints

```sh
# build info (no auth)
❯ curl localhost:4242/version
{"build_date":"20260602-233836","git_commit":"cacfab6"}

# list available API URLs
❯ curl -s localhost:4242/ls | jq '.api[]' -r
http://localhost:4242/api/test/invalid-json
http://localhost:4242/api/test/param
http://localhost:4242/api/test/post
http://localhost:4242/api/time/date

# run a script (GET)
❯ curl localhost:4242/api/time/date
{"date":1780435972,"human_date":"Tue Jun  2 23:32:52 CEST 2026"}

# run a script (GET, optional ?q= passed as $1)
❯ curl localhost:4242/api/test/param?q=hello
{"param":"hello"}

# run a script (POST, request body piped to stdin)
❯ curl -d '{"key":"42"}' localhost:4242/api/test/post
{"jackpot": "42"}

# serve static files from <apiDir>/_static
❯ curl localhost:4242/s/ -s | head -1
<!doctype html>
```

Invalid JSON from a script → `400`. Non-zero exit → `500`.

```sh
❯ curl localhost:4242/api/test/invalid-json
< HTTP/1.1 400 Bad Request
{"error":"Invalid JSON"}

❯ curl localhost:4242/api/test/fail
< HTTP/1.1 500 Internal Server Error
{"error":"exit status 1"}
```