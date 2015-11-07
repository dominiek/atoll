BUILD_PARAMS=--no-cache=false

all:
	exec boot2docker shellinit
	docker build $(BUILD_PARAMS) -t atoll-server/api api
	docker build $(BUILD_PARAMS) -t atoll-server/app app
	docker build $(BUILD_PARAMS) -t atoll-server/mongodb mongodb
	docker build $(BUILD_PARAMS) -t atoll-server/elasticsearch elasticsearch
	docker build $(BUILD_PARAMS) -t atoll-server/worker worker


clean:
	-docker rmi `docker images -q`

stop:
	-docker kill `docker ps -a -q`
	-docker rm "0.mongodb"
	-docker rm "1.mongodb"
	-docker rm "0.api"
	-docker rm "1.api"
	-docker rm "0.app"
	-docker rm "1.app"
	-docker rm "0.elasticsearch"
	-docker rm "1.elasticsearch"
	-docker rm "0.worker"
	-docker rm "1.worker"

run:
	make stop
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.mongodb" -t atoll-server/mongodb
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "1.mongodb" -t atoll-server/mongodb
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.elasticsearch" -t atoll-server/elasticsearch
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "1.elasticsearch" -t atoll-server/elasticsearch
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.worker" -t atoll-server/worker
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "1.worker" -t atoll-server/worker
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.api" --link 0.mongodb:0.mongodb --link 1.mongodb:1.mongodb --link 0.worker:0.worker --link 1.worker:1.worker -t atoll-server/api
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "1.api" --link 0.mongodb:0.mongodb --link 1.mongodb:1.mongodb --link 0.worker:0.worker --link 1.worker:1.worker -t atoll-server/api
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "0.app" --link 0.mongodb:0.mongodb --link 1.mongodb:1.mongodb --link 0.api:0.api --link 1.api:1.api -t atoll-server/app
	docker run -d --add-host=api.atoll.io:192.168.59.3 --name "1.app" --link 0.mongodb:0.mongodb --link 1.mongodb:1.mongodb --link 0.api:0.api --link 1.api:1.api -t atoll-server/app
