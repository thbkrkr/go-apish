GIT_COMMIT = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date '+%Y%m%d-%H%M%S')
NAME 		= apish
REPO 		= krkr/$(NAME)
GOPATH 		= $(shell pwd)/../../../..

build: build-binary build-image

build-in-docker:
	@echo "--> Build $(NAME) go binary inside Docker"
	docker run --rm \
		-w /go/src/github.com/thbkrkr/go-apish \
		-v ${GOPATH}:/go \
		-e CGO_ENABLED=0 -e GOOS=linux \
		-ti golang:1.5.1 \
			go build -a -installsuffix cgo -o $(NAME) \
				-ldflags "-X=main.gitCommit=$(GIT_COMMIT) -X=main.buildDate=$(BUILD_DATE)"

build-image:
	@echo "-->  Package $(NAME) go binary in Docker"
	@docker build --rm -t $(REPO) .

push:
	docker push $(REPO)

run:
	docker run -d -p 80:4242 $(REPO)