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

golive:
	gohere
	golive -apiDir=example/api

deploy-example:
	@make -C example build
	docker tag -f apish-example sailabove.io/krkr/apish-demo
	docker push sailabove.io/krkr/apish-demo

sail-add:
	sail services add krkr/apish-demo -p 80:4242 apish-demo

sail-redeploy:
	sail services redeploy apish-demo

open:
	sensible-browser http://apish-demo.krkr.app.runabove.io
