GIT_REVISION=$(shell git rev-parse --short=8 HEAD)
# PLATFORM=linux
# ARCH=amd64
VERSION=$(shell head ./VERSION)

LDFLAGS=-ldflags "-X github.com/superioz/artemis/appversion.Version=${VERSION} \
-X github.com/superioz/artemis/appversion.Build=${GIT_REVISION}"

all: test build install

build:
	go build ${LDFLAGS} -o artemis ./cmd/artemiscli/artemiscli.go

install:
	go install ${LDFLAGS}

test:
	${GOTEST} ./... -v
