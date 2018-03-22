BINARY=gitrc
VERSION=$(shell git describe --tag)
BUILD_DIR=${GOPATH}/src/github.com/wpueschel/${BINARY}

GOARCH=amd64
LDFLAGS= -ldflags "-s -X main.version=${VERSION}"

all: windows darwin linux

linux: 
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; \

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \

clean:
	-rm -f ${BINARY}-*
