BINARY=gg_ad_auction_system

VERSION=1.0.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
GOFILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

LDFLAGS=-ldflags "-X main.Version=${VERSION}-${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}"

OUTPUT_DIR=.
OUTPUT_NAME=${BINARY}

GO_PKG=bitbucket.org/greedygames/ad_request_auction_system

build:
	go build ${LDFLAGS} -o ${BINARY}

fmt:
	gofmt -s -w ${GOFILES}

install:
	go install ${LDFLAGS}

dist: clean
	GOOS=linux go build ${LDFLAGS} -o ${BINARY}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install

run:
	./gg_ad_auction_system
