#! /usr/bin/make

pb-compile:
	protoc --go_out=${GOPATH}/src protofiles/uniqueid.proto
