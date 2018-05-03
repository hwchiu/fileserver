SHELL := /bin/bash

GO_VENDOR := $(shell which govendor)
VENDOR_JSON_FILE := vendor/vendor.json

PUBLIC_DOCKER_REGISTRY = docker.io
DOCKER_PROJECT = linkernetworks

BUILD_DATE := $(shell date +%Y.%m.%d.%H:%M:%S)

GIT_SYMREF := $(shell git rev-parse --abbrev-ref HEAD | sed -e 's![^A-Za-z0-9]!-!g')
GIT_REV_SHORT := $(shell git rev-parse --short HEAD)
GIT_DESCRIBE := $(shell git describe --all --long)
BUILD_REVISION := $(GIT_REV_SHORT)

# container image definitions
# IMAGE_TAG := latest
# IMAGE_TAG := $(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(IMAGE_TAG),)
IMAGE_TAG := $(GIT_SYMREF)-$(GIT_REV_SHORT)
endif

# image anchor tag should refers to "latest" or "develop"
ifeq ($(IMAGE_ANCHOR_TAG),)
IMAGE_ANCHOR_TAG := $(GIT_SYMREF)
endif

all: build build-image push-image

build: vendor/.deps
	go build .
clean:
	@rm -rf fileserver

tool-govendor:
	if [ "$(GO_VENDOR)" == "" ] ; then go get github.com/kardianos/govendor ; fi

vendor/.deps: tool-govendor $(VENDOR_JSON_FILE)
	govendor sync
	touch $@

build-image:
	time docker build $(DOCKER_BUILD_FLAGS) \
		--tag $(PUBLIC_DOCKER_REGISTRY)/$(DOCKER_PROJECT)/fileserver:$(IMAGE_TAG) \
		--tag $(PUBLIC_DOCKER_REGISTRY)/$(DOCKER_PROJECT)/fileserver:$(IMAGE_ANCHOR_TAG) \
		.
push-image:
	docker push $(PUBLIC_DOCKER_REGISTRY)/$(DOCKER_PROJECT)/fileserver:$(IMAGE_ANCHOR_TAG)
