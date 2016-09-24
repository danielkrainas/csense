VERSION_FILE=VERSION
SRC_PKGS=$(shell go list ./... | grep -v vendor)
REV=$(shell git rev-parse --short HEAD)

ifeq ($(BUILD_VERSION),)
	BUILD_VERSION=$(shell cat $(VERSION_FILE))-$(REV)
endif

IMAGE_REPO=dakr/csense
ifeq ($(IMAGE),)
	IMAGE_NAME=$(IMAGE_REPO):$(BUILD_VERSION)
endif 

.PHONY: clean image test

all: compile

clean:
	go clean ./...

compile:
	go build -ldflags "-X main.appVersion=$(BUILD_VERSION)" .

dist:
	GOOS=linux go build -ldflags "-X main.appVersion=$(BUILD_VERSION)" -o dist .

image:
	docker build -t $(IMAGE_NAME) .

test:
	set -e; 
	for pkg in $(SRC_PKGS); \
	do \
		go test -v $$pkg; \
	done 