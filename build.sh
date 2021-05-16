#!/bin/bash

RUN_NAME="iot-gateway"

export GO111MODULE=on
go mod download
mkdir -p output/bin output/log
cp script/bootstrap.sh output 2>/dev/null
chmod +x output/bootstrap.sh

go build -o output/bin/${RUN_NAME}
