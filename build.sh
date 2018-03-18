#!/bin/bash

cd $(dirname $0)

echo "Building Windows binary"
GOARCH=386
GOOS=windows
export GOARCH GOOS

if ! go build -ldflags "-s -X main.version=$(git describe --tag)" -o gitrc_windows-386.exe; then 
   echo "Build failed. Exiting."
   exit 1
fi
ls -l gitrc_windows-386.exe

echo "Building linux binary"
GOARCH=amd64
GOOS=linux

if ! go build -ldflags "-s -X main.version=$(git describe --tag)" -o gitrc_linux-amd64; then
   echo "Build failed. Exiting."
   exit 1
fi
ls -l gitrc_linux-amd64

# EOF
