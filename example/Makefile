build:
	docker build --rm -t krkr/apish:example .

push:
	docker --config ~/.docker-config-krkr push krkr/apish:example

run:
	docker run --rm -ti \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /usr/bin/docker:/usr/bin/docker \
		-p 4242:4242 \
		apish:example

port-forward:
	kubectl port-forward svc/apish 4242:80
