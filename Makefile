.EXPORT_ALL_VARIABLES:
VERSION := $(shell git describe --always --tags)
DATE := $(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS := -ldflags "-X=main.version=$(VERSION)-$(DATE)"
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
PROJECT = mysuperproject
...

GO=GO111MODULE=on go

all: update generate bench test

update:
	cd confluent && $(GO) mod vendor

generate:
	cd confluent && $(GO) generate -mod=vendor

bench: 
	cd confluent && $(GO) test -mod=vendor -bench=.

test:
	cd confluent && $(GO) test
