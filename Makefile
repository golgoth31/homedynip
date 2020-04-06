#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------

DOCKER_USER ?= $(DOCKER_USER)
DOCKER_PASS ?=

DOCKER_BUILD_ARGS := --build-arg HTTP_PROXY=$(http_proxy) --build-arg HTTPS_PROXY=$(https_proxy)
DOCKER_AMD64_ARGS := --build-arg GOARM= --build-arg GOARCH=amd64
DOCKER_ARM64_ARGS := --build-arg GOARM= --build-arg GOARCH=arm64
DOCKER_ARMV7_ARGS := --build-arg GOARM=7 --build-arg GOARCH=arm
PROTOC_IMAGE := thethingsindustries/protoc:3.1.21
APP_VERSION := latest
PACKAGE ?= $(shell go list ./... | grep configs)
VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

GOLINTER:=$(shell command -v golangci-lint 2> /dev/null)

override LDFLAGS += \
  -X ${PACKAGE}.Version=${VERSION} \
  -X ${PACKAGE}.BuildDate=${BUILD_DATE} \
  -X ${PACKAGE}.GitCommit=${GIT_COMMIT} \


#-----------------------------------------------------------------------------
# BUILD
#-----------------------------------------------------------------------------

.PHONY: default build test publish build_local lint artifact_linux artifact_darwin deploy
default:  test lint build swagger

test:
	go test -v ./...

build_assets:
	docker run --rm -v $$(pwd):$$(pwd) -w $$(pwd) $(PROTOC_IMAGE) \
		-Iproto \
		-I$(GOPATH)/src \
		--go_out=$(GOPATH)/src \
		./proto/*.proto

run:
	go mod tidy
	go run main.go

build_local:
	go build -ldflags '${LDFLAGS}' -o ./homedynip main.go
build_amd64:
	docker build $(DOCKER_BUILD_ARGS) $(DOCKER_AMD64_ARGS) -t localhost/homedynip:$(APP_VERSION) -f ./build/Dockerfile .
build_arm64:
	docker build $(DOCKER_BUILD_ARGS) $(DOCKER_ARM64_ARGS) -t localhost/homedynip:$(APP_VERSION) -f ./build/Dockerfile .
build_armv7:
	docker build $(DOCKER_BUILD_ARGS) $(DOCKER_ARMV7_ARGS) -t localhost/homedynip:$(APP_VERSION) -f ./build/Dockerfile .
lint:
ifndef GOLINTER
	GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.15
endif
	golangci-lint run

artifact_linux_amd64:
	GOARM=7 GOARCH=amd64 GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -ldflags '${LDFLAGS}' -o homedynip-linux-amd64
artifact_linux_arm64:
	GOARM=7 GOARCH=arm64 GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -ldflags '${LDFLAGS}' -o homedynip-linux-arm64
artifact_linux_armv7:
	GOARM=7 GOARCH=arm GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -ldflags '${LDFLAGS}' -o homedynip-linux-armv7

artifact_darwin:
	GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=darwin go build -ldflags '${LDFLAGS}' -o homedynip-darwin

#-----------------------------------------------------------------------------
# PUBLISH
#-----------------------------------------------------------------------------

.PHONY: publish

publish:
	docker push $(DOCKER_USER)/homedynip:$(APP_VERSION)

#-----------------------------------------------------------------------------
# CLEAN
#-----------------------------------------------------------------------------

.PHONY: clean

clean:
	rm -rf homedynip
