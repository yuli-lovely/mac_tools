#!/bin/bash

set -e

rm -rf bin

function go_build(){
    # CGO_ENABLED=0 GOOS=windows GOARCH=amd64
    # CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64
   CGO_ENABLED=0 GOOS=$1 GOARCH=$2  go build -v -ldflags "-s -w" -o bin/"$1_$2"/mac_tools main.go
}

go_build windows amd64
go_build darwin amd64
go_build darwin arm64
go_build linux amd64