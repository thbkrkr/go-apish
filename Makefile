GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date '+%Y%m%d-%H%M%S')

build: build-binary build-image

build-binary:
	docker run --rm \
		-w /go/src/github.com/thbkrkr/go-apish \
		-v $(shell pwd)/../../../..:/go \
		-e CGO_ENABLED=0 -e GOOS=linux \
		-ti golang:1.5.1 \
			go build -a -installsuffix cgo \
				-ldflags "-X=main.gitCommit=$(GIT_COMMIT) -X=main.buildDate=$(BUILD_DATE)"

build-image:
	@docker build --rm -t krkr/apish .

push:
	docker push krkr/apish

run:
	docker run -d -p 80:4242 krkr/apish
