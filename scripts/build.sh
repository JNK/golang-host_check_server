#!/bin/bash

function build {
    OUTPUT_NAME="bin/${1}-${2}-v${3}-host_check_server"
    export GOOS=$1
    export GOARCH=$2
    /usr/local/go/bin/go build -o $OUTPUT_NAME .
}

HCS_VERSION=$1

echo "Building ${HCS_VERSION}..."

set -x

build linux amd64 $HCS_VERSION

build linux arm64 $HCS_VERSION

build darwin amd64 $HCS_VERSION

build windows amd64 $HCS_VERSION

chmod +x bin/*-host_check_server