.PHONY: build dep

APP := manual-gcp-autoscaler
VERSION := "0.0.1"
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
        -X 'main.revision=$(REVISION)'

dep:
	go get -u github.com/golang/dep/...
	dep ensure

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(PWD)/bin/darwin_amd64/$(APP)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(PWD)/bin/linux_amd64/$(APP)
