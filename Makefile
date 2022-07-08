#====================
AUTHOR         ?= The sacloud/go-template Authors
COPYRIGHT_YEAR ?= 2022

BIN            ?= packer-plugin-sakuracloud

include includes/go/common.mk
include includes/go/single.mk
#====================

HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

default: fmt generate go-licenses-check goimports lint test

.PHONY: tools
tools: dev-tools install-packer-sdc

install-packer-sdc:
	go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}

# CI環境向けにpackerをセットアップ
.PHONY: install-packer
install-packer:
	apt update; apt install -y curl zip
	curl -LO https://releases.hashicorp.com/packer/1.8.0/packer_1.8.0_linux_amd64.zip
	unzip -o packer_1.8.0_linux_amd64.zip
	install packer /usr/local/bin/

# CI環境向けにpackerプラグイン(sakuracloud)をセットアップ
.PHONY: install-plugin
install-plugin: dev

dev: build
	@mkdir -p ~/.packer.d/plugins/
	@mv ${BIN} ~/.packer.d/plugins/${BIN}

generate:
	go generate ./...

ci-release-docs: install-packer-sdc
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check ${BIN}