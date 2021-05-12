#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)

BinaryName="edgex-gateway"

echo "$CURDIR/bin/${BinaryName} -conf=${CONF_FILE_PATH}"

exec $CURDIR/bin/${BinaryName} -conf=${CONF_FILE_PATH}