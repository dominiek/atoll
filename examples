
run:
	docker run -d \
		-p 8400:8400 \
		-p 8500:8500 \
		-p 53:53/udp \
		-h consul-server-node \
		--name consul \
		progrium/consul \
			-server \
			-bootstrap
	docker run -d \
		--link consul:consul \
		--name registrator \
		-h "Triggerfish.local" \
		-v /var/run/docker.sock:/tmp/docker.sock \
		gliderlabs/registrator:latest \
			consul://consul:8500
