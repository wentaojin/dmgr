.PHONY: all build run gotool clean help

CMDPATH="./cmd"
BINARYPATH="bin/dmgr"

REPO    := github.com/wentaojin/dmgr

GOOS    := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH  := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
GOENV   := GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO      := $(GOENV) go
GOBUILD := $(GO) build $(BUILD_FLAGS)
GORUN   := $(GO) run
SHELL   := /usr/bin/env bash

COMMIT  := $(shell git describe --no-match --always --dirty)
BUILDTS := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GITHASH := $(shell git rev-parse HEAD)
GITREF  := $(shell git rev-parse --abbrev-ref HEAD)
GOVERSION := $(shell go version)


LDFLAGS := -w -s
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.ReleaseVersion=$(COMMIT)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.BuildTS=$(BuildTS)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.GitHash=$(GITHASH)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.GitBranch=$(GITREF)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.GoVersion=$(GOVERSION)"
LDFLAGS += $(EXTRA_LDFLAGS)


FILES   := $$(find . -name "*.go")

all: clean gotool build

build: clean gotool
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(BINARYPATH) $(CMDPATH)

run:
	$(GORUN) -race $(CMDPATH)

gotool:
	$(GO) mod tidy

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi


help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'mod tidy'"