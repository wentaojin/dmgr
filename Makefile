.PHONY: build run gotool clean help

CMDPATH="./cmd"
BINARYPATH="bin/dmgr"
CONFIGPATH="./conf/dmgr.toml"

REPO    := github.com/wentaojin/dmgr

GOOS    := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH  := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
GOENV   := GO111MODULE=on CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH)
GO      := $(GOENV) go
GOBUILD := $(GO) build
GORUN   := $(GO) run
SHELL   := /usr/bin/env bash

COMMIT  := $(shell git describe --no-match --always --dirty)
BUILDTS := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GITHASH := $(shell git rev-parse HEAD)
GITREF  := $(shell git rev-parse --abbrev-ref HEAD)


LDFLAGS := -w -s
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.Version=$(COMMIT)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.BuildTS=$(BUILDTS)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.GitHash=$(GITHASH)"
LDFLAGS += -X "$(REPO)/pkg/dmgrutil.GitBranch=$(GITREF)"


build: clean gotool
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(BINARYPATH) $(CMDPATH)

run: gotool
	$(GORUN) -race $(CMDPATH) --config $(CONFIGPATH)

gotool:
	$(GO) mod tidy

clean:
	@if [ -f ${BINARYPATH} ] ; then rm ${BINARYPATH} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'mod tidy'"