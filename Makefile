.DEFAULT_GOAL := build
HAS_UPX := $(shell command -v upx 2> /dev/null)

.PHONY: build
build:
	go build -o ./bin/lark_backup main.go
ifneq ($(and $(COMPRESS),$(HAS_UPX)),)
	upx -9 ./bin/lark_backup
endif

.PHONY: clean
clean:  ## Clean build bundles
	rm -rf ./bin

.PHONY: format
format:
	gofmt -l -w .
