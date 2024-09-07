#!/bin/bash

set -e

# Note: BY ZWF;

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

rm -f ./filebeat
go build -ldflags "-s -w"

IMAGE_TAG=8.14.4-hurricanezwf0.0.1
docker build -f ./Dockerfile -t harbor.123u.com/public/filebeat:${IMAGE_TAG} .
