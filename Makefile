.PHONY: clean test build

default: build

clean:
	rm -f srtmod*

build:
	go build

all: windows darwin_x86 darwin_arm linux

windows:
	GOOS=windows GOARCH=amd64 go build -o srtmod_windows_x86_64.exe

darwin_x86:
	GOOS=darwin GOARCH=amd64 go build -o srtmod_darwin_x86_64

darwin_arm:
	GOOS=darwin GOARCH=arm64 go build -o srtmod_darwin_arm_64

linux:
	GOOS=linux GOARCH=amd64 go build -o srtmod_linux_x86_64

test:
	go test -cover
