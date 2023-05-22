
BUILD_TARGET := $(dir $(wildcard cmd/**/main.go))
CLEAN_PATH := tmp dist
MAKE_CMD := make --no-print-directory

.PHONY: all
all: $(BUILD_TARGET)
	# all build ok

.PHONY: $(BUILD_TARGET)
$(BUILD_TARGET): tidy
	## build $@ ...
	@cd $@ && $(MAKE_CMD)

.PHONY: tidy
tidy:
	# tidy workspace ...
	@if [[ -n "$$(go env GOWORK)" ]]; then go work sync; else go mod tidy;fi

CLEAN_TARGET := $(addprefix clean-,$(BUILD_TARGET))
.PHONY: clean $(CLEAN_TARGET)
clean: $(CLEAN_TARGET)
	# clean root ...
	rm -rf $(CLEAN_PATH)
$(CLEAN_TARGET):
	## clean $@ ...
	@cd $(patsubst clean-%,%,$@) && $(MAKE_CMD) clean

TEST_TARGET := $(addprefix test-,$(BUILD_TARGET))
.PHONY: test $(TEST_TARGET)
test: $(TEST_TARGET)
	# unit test ...
	@go test -v ./...
$(TEST_TARGET):
	# integration test $(patsubst test-%,%,$@) ...
	@cd $(patsubst test-%,%,$@) && $(MAKE_CMD) test

.PHONY: generate
generate:
	go generate ./...

.PHONY: update
update:
	git pull --recurse-submodules
