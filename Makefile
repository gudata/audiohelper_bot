# all: build
install: build
# 	upx bin/audio-helper

build-linux:
	go build -o bin/audio-helper commands/main.go

build: build-linux

clean:
	rm -f bin/audio-helper

test:
	go test ./commands/... ./commands/