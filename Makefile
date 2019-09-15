BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%S%z')
BUILD_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --long HEAD)
LDFLAGS := "-X main.BuildTimestamp=$(BUILD_TIME) -X main.BuildCommit=$(BUILD_COMMIT) -X main.BuildVersion=$(VERSION)"

all: porcelain

porcelain:
	go build -ldflags $(LDFLAGS) -o ./$@ ./cmd/$@

.PHONY: porcelain
