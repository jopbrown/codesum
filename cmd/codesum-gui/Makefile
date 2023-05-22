WORKSPACE := $(dir $(shell go env GOMOD))

BUILD_NAME := $(notdir $(shell pwd))

TARGET := windows-amd64

RESOURCE := config.yml

CLEAN_PATH := tmp frontend/src/lib/wailsjs

ENABLE_UPX := true

ENABLE_ARCHIVE := false

include $(WORKSPACE)/script/build.mk

################ special configure for wails start ################

BUILD_BIN_DIR := $(PWD)/build/bin
CLEAN_PATH += $(BUILD_BIN_DIR)

PRECMD = rm -rf $(DIST_DIR) && \
	mkdir -p $(BUILD_BIN_DIR) $(WORKSPACE)/dist

ifdef OS
	PRECMD += && cmd <<< "mklink /D \"$(DIST_DIR)\" \"$(BUILD_BIN_DIR)\"" > /dev/null
else
	PRECMD += && ln -s $(BUILD_BIN_DIR) $(DIST_DIR)
endif

BUILD_CMD = wails build -v 0 -trimpath -ldflags="-s -w $(BUILD_LDFLAGS)" -o $(BUILD_NAME)_$(VERSION)_$(BUILD_TIME)/$(TARGET_EXE) .

################ special configure for wails end ################

.PHONY: test
test:
	## $(BUILD_NAME) test pass

.PHONY: gen
gen:
	wails generate module

.PHONY: dev
dev:
	wails dev -appargs "-c config-test.yml"
