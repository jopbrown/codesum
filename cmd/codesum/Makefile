
BUILD_NAME := $(notdir $(shell pwd))

WORKSPACE := $(dir $(shell go env GOMOD))

TARGET := windows-amd64 linux-amd64

RESOURCE :=

CLEAN_PATH := tmp

ENABLE_UPX := false

ENABLE_ARCHIVE := false

include $(WORKSPACE)/script/build.mk

.PHONY: test
test:
	rm -rf tmp
	go build -o tmp/codesum .
	cp config.yml tmp/
	cd tmp && ./codesum $(WORKSPACE)
