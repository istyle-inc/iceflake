#! /usr/bin/make

pb-compile:
	protoc --go_out=${GOPATH}/src protofiles/uniqueid.proto

run:
	go run main.go generator.go logger.go connector.go
