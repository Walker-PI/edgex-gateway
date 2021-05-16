#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)

BinaryName="edgex-gateway"

exec $CURDIR/bin/${BinaryName} -conf=${CONF_FILE_PATH}