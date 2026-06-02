GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date '+%Y%m%d-%H%M%S')
LDFLAGS = -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)

build:
	docker build --rm \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t krkr/apish .

binary:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o go-apish .

test:
	go test ./...

release:
	./release.sh $(GIT_COMMIT)

push:
	docker push krkr/apish

run:
	docker run -d \
		-v $$(pwd)/example:/api \
		-p 80:4242 \
		krkr/apish

golive:
	gohere
	golive -apiDir=example/api
