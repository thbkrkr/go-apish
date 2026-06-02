build:
	docker build --rm -t apish-example .

run:
	docker run --rm -ti \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /usr/bin/docker:/usr/bin/docker \
		-p 80:4242 \
		apish-example
