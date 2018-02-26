#! /usr/bin/make

pb-compile:
	protoc --go_out=${GOPATH}/src protofiles/uniqueid.proto

run:
	go run main.go generator.go connector.go foundation

install-tools:
	go get -u github.com/goreleaser/goreleaser
	go get -u github.com/golang/dep/cmd/dep

build:
	goreleaser 
