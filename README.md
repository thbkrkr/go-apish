# apish — REST API for shell scripts

`apish` turns a directory of shell scripts into a JSON REST API and serves
static files alongside them. Each script must print **valid JSON** to stdout
([example](example/api/time/date.sh)); the server parses it and returns it to
the caller.

## Build & run

```sh
make binary   # build the ./go-apish binary locally
make build    # build the krkr/apish Docker image
make run      # run the image, mounting ./example as /api on port 80
```

Run the binary directly against the example API:

```sh
./go-apish -apiDir=example/api -password=secret
```

## Flags

| Flag            | Default     | Description                                              |
| --------------- | ----------- | -------------------------------------------------------- |
| `-port`         | `4242`      | HTTP port to listen on                                   |
| `-apiDir`       | `./api`     | Directory of `.sh` scripts and `_static` files           |
| `-user`         | `zuperadmin`| Basic-auth username                                      |
| `-password`     | *(empty)*   | Basic-auth password. **Empty disables all auth.**        |
| `-apiKey`       | *(empty)*   | Key for `X-Auth` header auth. Empty disables header auth.|

## Endpoints

| Method | Path         | Description                                                      |
| ------ | ------------ | ---------------------------------------------------------------- |
| GET    | `/`          | JSON status, or redirect to `/s` if `_static/index.html` exists  |
| GET    | `/version`   | Build commit and date (no auth)                                  |
| GET    | `/ls`        | List script, HTML and static resource URLs                       |
| GET    | `/api/*path` | Run `<apiDir>/<path>.sh`; `?q=value` is passed as `$1`           |
| POST   | `/api/*path` | Run `<apiDir>/<path>.sh` with the request body piped to stdin    |
| GET    | `/s/*`       | Serve static files from `<apiDir>/_static`                        |

Scripts must emit valid JSON; otherwise the caller receives `400 Invalid JSON`.
A script that exits non-zero yields `500` with its error (stderr is logged).

## Authentication

When `-password` is set, all endpoints except `/` and `/version` require
either:

- HTTP basic auth (`-user` / `-password`), or
- an `X-Auth: <apiKey>` header (when `-apiKey` is set).

If `-password` is empty, **the server is fully open** and logs a warning at
startup.

## Security

`apish` executes shell scripts — treat it as a privileged service:

- Always set `-password` (and ideally an `-apiKey`) in any non-local deployment.
- Scripts receive request input (`$1` / stdin). Build their JSON output with a
  tool like `jq` so values are safely escaped — see
  [param.sh](example/api/test/param.sh).
