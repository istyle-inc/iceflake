#! /usr/bin/make

ICEFLAKE_SOCKETFILE_PATH ?= /tmp/iceflake-worker-1.sock

pb-compile:
	protoc --go_out=${GOPATH}/src protofiles/uniqueid.proto

run:
	go run cmd/iceflake/*.go -w 1 -s $(ICEFLAKE_SOCKETFILE_PATH)

install-tools:
	go get -u github.com/goreleaser/goreleaser
	go get -u github.com/golang/dep/cmd/dep

build:
	goreleaser 

bench: 
	go test -benchmem -bench Benchmark github.com/istyle-inc/iceflake/app
