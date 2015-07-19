GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date '+%Y%m%d-%H%M%S')

build:
	go build -ldflags "-X main.gitCommit $(GIT_COMMIT) -X main.buildDate $(BUILD_DATE)"

build-in-docker:
	docker run --rm -v $(shell pwd):/src -t -i centurylink/golang-builder

run: build
	./go-apish -port=1234 -apiKeyHeader=X-pof-auth -apiKey=E7D9EJJD87EH7ED87H

define exec =
	@echo "\033[0;33mGET \033[0m/$(1)"
	@curl -sH 'X-pof-auth:E7D9EJJD87EH7ED87H' localhost:1234/$(1) | jq .
endef

test:
	$(call exec,version)
	$(call exec,ls)
	$(call exec,date.sh)
	$(call exec,ping)

stop:
	pkill go-apish
