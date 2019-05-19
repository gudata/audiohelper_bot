.PHONY: clean

all: build checks install

checks:
	go vet ./...
	gocyclo -top 10 ./commands/main.go ./packages/
	-golint ./packages/... ./commands/
	ineffassign ./commands/main.go

build: test build-linux checks

build-linux:
	go build -o bin/audio-helper commands/main.go

install: build
	upx bin/audio-helper


clean:
	rm -f bin/audio-helper

clean:

test:
	go test ./commands/... ./commands/

install-packages:
	go get -u golang.org/x/lint/golint
	go get thub.com/gordonklaus/ineffassign
	go get golang.org/x/tools/cmd/guru

