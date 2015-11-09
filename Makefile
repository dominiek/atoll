
export GOPATH=$(shell pwd)
export GOBIN=$(shell pwd)/bin

all:
	mkdir -p $(GOBIN)
	go get
	go install

clean:
	rm -rf bin

test:
	go test *.go

test.verbose:
	go test *.go -v

test.linux:
	docker build -t atoll-test .
	docker run -e BUILD_COMMAND="make test" -t atoll-test

test.linux.verbose:
	docker build -t atoll-test .
	docker run -e BUILD_COMMAND="make test.verbose" -t atoll-test

linux.shell:
	docker build -t atoll-test .
	docker run -i -t atoll-test /bin/bash
