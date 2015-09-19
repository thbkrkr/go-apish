build:
	docker build --rm -t api-example .

run:
	docker run -d \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v /usr/bin/docker:/usr/bin/docker \
		-p 80:4242 \
		api-example
