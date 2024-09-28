#!/bin/bash

set -e

# Note: BY ZWF;

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

rm -f ./filebeat
go build -ldflags "-s -w"

IMAGE_TAG=8.14.4-hurricanezwf0.2.7
docker build -f ./Dockerfile -t ops-huanle.tencentcloudcr.com/public/filebeat:${IMAGE_TAG} .
docker push ops-huanle.tencentcloudcr.com/public/filebeat:${IMAGE_TAG}
