
GOCMD ?= go
TARGET ?= linux-amd64 windows-amd64
CLEAN_PATH ?= tmp
ENABLE_UPX ?= false
ENABLE_ARCHIVE ?= false

ifdef OS
WORKSPACE := $(shell cygpath -m $(WORKSPACE))
endif

DIST_DIR := $(WORKSPACE)/dist/$(BUILD_NAME)

VERSION := $(shell cat ver.txt 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
BUILD_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date +%Y%m%d%H%M%S)

BUILD_LDFLAGS := -X main.BuildName=$(BUILD_NAME) -X main.BuildVersion=$(VERSION) -X main.BuildHash=$(BUILD_HASH) -X main.BuildTime=$(BUILD_TIME)

export GOOS GOARCH
TARGET_DIR := $(DIST_DIR)/$(BUILD_NAME)_$(VERSION)_$(BUILD_TIME)
TARGET_NAME = $(BUILD_NAME)_$(GOOS)_$(GOARCH)
TARGET_EXE = $(TARGET_NAME)$(TARGET_EXT)
TARGET_PATH = $(TARGET_DIR)/$(TARGET_EXE)
BUILD_CMD = go build -trimpath -ldflags="-s -w $(BUILD_LDFLAGS)" -buildvcs=false -o $(TARGET_PATH) .

.PHONY: all $(TARGET)
all: $(TARGET)
	## $(BUILD_NAME) build ok
$(TARGET): resource
	### build $@ ...
	$(BUILD_CMD)
ifeq ($(ENABLE_UPX),true)
	### upx $@ ...
	@upx $(TARGET_PATH) 1>/dev/null
endif
	$(POSTCMD)
ifeq ($(ENABLE_ARCHIVE),true)
	### archive $@ ...
	@cd $(TARGET_DIR) && tar zcf $(TARGET_NAME).tar.gz $(TARGET_EXE) $(notdir $(RESOURCE))
endif

linux-amd64: GOOS=linux
linux-amd64: GOARCH=amd64

windows-amd64: GOOS=windows
windows-amd64: GOARCH=amd64
windows-amd64: TARGET_EXT=.exe

.PHONY: mkdir
mkdir:
	@mkdir -p $(TARGET_DIR)
	@rm -rf $(DIST_DIR)/latest
ifdef OS
	@cmd <<< "mklink /D \"$(DIST_DIR)/latest\" \"$(TARGET_DIR)\"" > /dev/null
else
	@ln -s $(TARGET_DIR) $(DIST_DIR)/latest
endif

.PHONY: resource
resource: mkdir
	$(PRECMD)
ifneq ($(RESOURCE),)
	## copy resource ...
	@cp -rf $(RESOURCE) $(TARGET_DIR)/
endif

.PHONY: clean
clean:
	rm -rf $(CLEAN_PATH) $(DIST_DIR)

echo:
	mklink /D $(DIST_DIR)/latest $(TARGET_DIR)
