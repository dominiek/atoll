
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
