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

bench: iceflakebench.pid

iceflakebench.pid: 
	@rm -f $(ICEFLAKE_SOCKETFILE_PATH)
	go run main.go generator.go connector.go -w 1 -s $(ICEFLAKE_SOCKETFILE_PATH) & echo $$! > $@;
	@sleep 5
	cd tool/icebench; ICEFLAKE_SOCKETFILE_PATH=$(ICEFLAKE_SOCKETFILE_PATH) go test -v -bench .
	kill -KILL `cat $@` && rm $@
	@rm -f $(ICEFLAKE_SOCKETFILE_PATH)

