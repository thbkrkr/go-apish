GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date '+%Y%m%d-%H%M%S')

build: build-binary build-image

build-binary:
	docker run --rm \
		-w /go/src/github.com/thbkrkr/go-apish \
		-v $(shell pwd):/go/src/github.com/thbkrkr/go-apish \
		-e CGO_ENABLED=0 -e GOOS=linux \
		-ti golang:1.6.2 \
			go build -a -installsuffix cgo \
				-ldflags "-X=main.gitCommit=$(GIT_COMMIT) -X=main.buildDate=$(BUILD_DATE)"

build-image:
	@docker build --rm -t krkr/apish .

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
