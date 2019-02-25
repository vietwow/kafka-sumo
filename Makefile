.EXPORT_ALL_VARIABLES:
BUILD_NUMBER = "1.0"
SHA1 := $(shell git rev-parse HEAD)
DATE := $(shell date +%Y_%m_%d)
LDFLAGS := -ldflags  " -X main.version=$(BUILD_NUMBER) -X main.commit_sha1=$(SHA1) -X main.builddate=$(DATE) "
GOOS=linux
GO111MODULE=on

all: update build test

update:
	go mod vendor

build:
	go build '-mod=vendor' $(LDFLAGS) -a -installsuffix cgo -o /go/bin/kafka-sumo ./cmd/forwarder/

test:
	go test
