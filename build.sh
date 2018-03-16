#!/bin/bash

cd $(dirname $0)

echo "Building Windows binary"
GOARCH=386
GOOS=windows
export GOARCH GOOS

go build -o gitrc_windows-386.exe
ls -l gitrc_windows-386.exe

echo "Building linux binary"
GOARCH=amd64
GOOS=linux

go build -o gitrc_linux-amd64
ls -l gitrc_linux-amd64

# EOF
