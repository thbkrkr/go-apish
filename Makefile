
build:
	docker run --rm -v $(shell pwd):/src -t -i centurylink/golang-builder

run:
	./go-apish -port=1234 -apiKeyHeader=X-pof-auth -apiKey=E7D9EJJD87EH7ED87H &

test:
	curl -H 'X-pof-auth:E7D9EJJD87EH7ED87H' localhost:1234/date.sh

stop:
	pkill go-apish
