#! /usr/bin/make

ICEFLAKE_SOCKETFILE_PATH ?= /tmp/iceflake-worker-1.sock

pb-compile:
	protoc --go_out=${GOPATH}/src protofiles/uniqueid.proto

run:
	go run main.go generator.go connector.go -w 1 -s $(ICEFLAKE_SOCKETFILE_PATH)

install-tools:
	go get -u github.com/goreleaser/goreleaser
	go get -u github.com/golang/dep/cmd/dep

build:
	goreleaser 

bench: 
	go test -bench Benchmark
