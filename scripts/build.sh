#!/bin/bash

CGO_ENABLED=1

system=$(uname -s)
case $system in
Linux)
    GOOS="linux"
    ;;
*)
    echo "Unknown system: ${system}"
    exit 1
    ;;
esac

arch=$(uname -m)
case $arch in
x86_64)
    GOARCH="amd64"
    ;;
*)
    echo "Unknown arch: ${arch}"
    exit 1
    ;;
esac

export CGO_ENABLED GOOS GOARCH
go build -v -buildmode c-shared -o "$1" ./lib
