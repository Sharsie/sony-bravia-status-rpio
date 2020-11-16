BUILD_DIR 		:= $(shell pwd)
export GOBIN 		:= $(BUILD_DIR)/bin
export GO111MODULE 	:= on


GIT_TAG		:= $(shell git describe --tags --exact-match 2>/dev/null)
GIT_HASH	:= $(shell git rev-parse HEAD)
VERSION		:= $(or $(GIT_TAG),latest)
BUILD_TIME	:= $(shell date '+%s')

# Runtime environment variables for docker mounts
DOCKER_USER	:= $(shell id -u):$(shell id -g)
# TODO: Does not work, sets to empty string
#HOME_CACHE_DIR  := $(or $${XDG_CACHE_HOME},$$HOME/.cache)
HOME_CACHE_DIR  := $$HOME/.cache
GO_CACHE_DIR	:= $(or $(GOCACHE), $(HOME_CACHE_DIR)/go-build)

REGISTRY	:= docksee
APP_NAME 	:= tv-status-rpio

-include .env

all: clean build
.PHONY: all


clean:
	@go clean -mod readonly ./...
	@rm -rf $(GOBIN)/*
.PHONY: clean

docker-clean/%:
	@docker run -u $(DOCKER_USER) -v $(GO_CACHE_DIR):/.cache/go-build  -v $(BUILD_DIR):/go/src/github.com/Sharsie/$(APP_NAME) -v ~/go/pkg/mod:/go/pkg/mod --rm $*/golang:1.15.5 \
		sh -c "cd /go/src/github.com/Sharsie/$(APP_NAME) && go clean -mod readonly ./..."
	@rm -rf $(GOBIN)/*
.PHONY: clean


test:
	@go test -mod readonly ./...
.PHONY: test


build/%:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -mod readonly -ldflags "-s -w \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.Tag=$(VERSION) \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.Commit=$(GIT_HASH) \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.buildTime=$(BUILD_TIME)" -o ./bin/amd64/linux/$* ./cmd/$*
build-arm32v7/%:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -mod readonly -ldflags "-s -w \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.Tag=$(VERSION) \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.Commit=$(GIT_HASH) \
		-X github.com/Sharsie/$(APP_NAME)/cmd/$*/version.buildTime=$(BUILD_TIME)" -o ./bin/arm32v7/linux/$* ./cmd/$*

build: build/is-on
build-arm32v7: build-arm32v7/is-on
.PHONY: build build-arm32v7

docker-build:
	@docker run -u $(DOCKER_USER) -v $(GO_CACHE_DIR):/.cache/go-build -v $(BUILD_DIR):/go/src/github.com/Sharsie/$(APP_NAME) -v ~/go/pkg/mod:/go/pkg/mod --rm golang:1.15.5 \
		sh -c "cd /go/src/github.com/Sharsie/$(APP_NAME) && make build"
docker-build-arm32v7:
	@docker run -u $(DOCKER_USER) -v $(GO_CACHE_DIR):/.cache/go-build  -v $(BUILD_DIR):/go/src/github.com/Sharsie/$(APP_NAME) -v ~/go/pkg/mod:/go/pkg/mod --rm arm32v7/golang:1.15.5 \
		sh -c "cd /go/src/github.com/Sharsie/$(APP_NAME) && make build-arm32v7"
.PHONY: docker-build docker-build-arm32v7


release/%:
	@docker build --build-arg FROM_IMAGE=alpine:3.12 --build-arg ARCH=amd64 --build-arg OS=linux --build-arg APP_NAME=$* -t $(REGISTRY)/$(APP_NAME)-$*:$(VERSION) -f ./docker/Dockerfile ./bin
	@docker push $(REGISTRY)/$(APP_NAME)-$*:$(VERSION)
release-arm32v7/%:
	@docker build --build-arg FROM_IMAGE=arm32v7/alpine:3.12 --build-arg ARCH=arm32v7 --build-arg OS=linux --build-arg APP_NAME=$* -t $(REGISTRY)/$(APP_NAME)-$*:$(VERSION)-arm32v7 -f ./docker/Dockerfile ./bin
	@docker push $(REGISTRY)/$(APP_NAME)-$*:$(VERSION)-arm32v7


release: clean docker-build release/is-on
release-arm32v7: docker-clean/arm32v7 docker-build-arm32v7 release-arm32v7/is-on
.PHONY: release release-arm32v7


develop/%:
	@go run --tags dev -race ./cmd/$*

develop: develop/is-on
.PHONY: develop
