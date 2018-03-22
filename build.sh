#!/bin/bash

cd $(dirname $0)

echo "Building Windows binary"
GOARCH=386
GOOS=windows
export GOARCH GOOS
if ! go build -ldflags "-s -X main.version=$(git describe --tag)" -o gitrc_$GOOS-$GOARCH.exe; then 
   echo "Build failed. Exiting."
   exit 1
fi
ls -l gitrc_windows-386.exe

GOARCH=amd64
for GOOS in darwin linux; do
   echo "Building $GOOS binary"
   if ! go build -ldflags "-s -X main.version=$(git describe --tag)" -o gitrc_$GOOS-$GOARCH; then
      echo "Build failed. Exiting."
      exit 1
   fi
   ls -l gitrc_$GOOS-$GOARCH
done

# EOF
