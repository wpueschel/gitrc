GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOTEST=${GOCMD} test
GOGET=${GOCMD} get

BINARY_NAME=gitrc
VERSION=${shell git describe --tag}
BUILD_DIR=${GOPATH}/src/github.com/wpueschel/${BINARY}

GOARCH=amd64
LDFLAGS= -ldflags "-s -X main.version=${VERSION}"

all: dep windows darwin linux

linux: 
	GOOS=linux GOARCH=${GOARCH} ${GOBUILD} ${LDFLAGS} -o ${BINARY_NAME}-linux-${GOARCH} . ; \

darwin:
	GOOS=darwin GOARCH=${GOARCH} ${GOBUILD} ${LDFLAGS} -o ${BINARY_NAME}-darwin-${GOARCH} . ; \

windows:
	GOOS=windows GOARCH=${GOARCH} ${GOBUILD} ${LDFLAGS} -o ${BINARY_NAME}-windows-${GOARCH}.exe . ; \

dep:
	${GOGET} ./... 

clean:
	-rm -f ${BINARY_NAME}-*
