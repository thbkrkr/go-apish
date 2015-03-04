
build:
	docker run --rm -v $(shell pwd):/src -t -i centurylink/golang-builder
