

all:
	exec boot2docker shellinit
	docker build -t atoll-server/api api
	docker build -t atoll-server/app app
	docker build -t atoll-server/mongodb mongodb

clean:
	-docker kill `docker ps -a -q`
	-docker rm "0.mongodb"
	-docker rm "0.api"
	-docker rm "0.app"

run:
	make clean
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.mongodb" -t atoll-server/mongodb
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.api" --link 0.mongodb:0.mongodb  -t atoll-server/api
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.app" --link 0.mongodb:0.mongodb  --link 0.api:0.api -t atoll-server/app
