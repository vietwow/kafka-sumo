#!/bin/bash

export ROOT_DIR=$PWD
#export BUILD_NUMBER=`[ ${VERSION_APP} ] && cat ${VERSION_APP} || echo "0.0.0"`
#export SHA1=`[ ${COMMIT_SHA1} ] && cat ${COMMIT_SHA1} || echo "000000000"`
export BUILD_NUMBER="1.0"
export SHA1=`git rev-parse HEAD`
export DATE=$(date +%Y_%m_%d)

#CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build '-mod=vendor'  -ldflags  " -X main.version=$BUILD_NUMBER -X main.commit_sha1=$SHA1 -X main.builddate=$DATE " -a -installsuffix cgo -o /go/bin/kafka-sumo ./cmd/forwarder/
CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -o /go/bin/kafka-sumo ./cmd/forwarder/