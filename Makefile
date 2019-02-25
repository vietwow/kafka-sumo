.EXPORT_ALL_VARIABLES:
BUILD_NUMBER = "1.0"
SHA1=$(git rev-parse HEAD)
DATE=$(date +%Y_%m_%d)

all: update build test

update:
	go mod vendor

build:
	GO111MODULE=on GOOS=linux go build '-mod=vendor'  -ldflags  " -X main.version=$(BUILD_NUMBER) -X main.commit_sha1=$(SHA1) -X main.builddate=$(DATE) " -a -installsuffix cgo -o /go/bin/kafka-sumo ./cmd/forwarder/

test:
	go test
