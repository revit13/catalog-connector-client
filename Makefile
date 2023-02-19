GOARCH = amd64
OS = linux

.DEFAULT_GOAL := all

all: source-build

.PHONY: source-build
source-build:
	CGO_ENABLED=0 GOOS="$(OS)" GOARCH="$(GOARCH)" go build 

.PHONY: run-read
run-read:
	./catalog-fake-client --request resources/request.json --operation "read" --creds "qqq" --url "http://localhost:888" --dataset "asset"

include hack/make-rules/verify.mk

