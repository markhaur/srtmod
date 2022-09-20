.PHONY: clean test build

default: build

clean:
	rm -f output/*
	rm -f app

build:
	go build -o app main.go

test:
	go test -cover
