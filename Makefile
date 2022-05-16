NAME=sakuracloud
BINARY=packer-plugin-${NAME}
HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

TEST?=$$(go list ./... | grep -v vendor)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: lint build

.PHONY: tools
tools: install-packer-sdc
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get github.com/client9/misspell/cmd/misspell
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.38.0/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.38.0

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

.PHONY: clean
clean:
	rm -f ${BINARY}

build:
	@go build -o ${BINARY}

dev: build
	@mkdir -p ~/.packer.d/plugins/
	@mv ${BINARY} ~/.packer.d/plugins/${BINARY}

generate:
	go generate ./...

.PHONY: test testacc
test:
	go test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

testacc:
	@echo "WARN: Acceptance tests will take a long time to run and may cost money. Ctrl-C if you want to cancel."
	PACKER_ACC=1 go test -v $(TEST) $(TESTARGS) -timeout=45m

.PHONY: lint fmt golint goimports
lint: fmt goimports golangci-lint

fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

golangci-lint: fmt
	golangci-lint run ./...

goimports:
	find . -name '*.go' | grep -v vendor | xargs goimports -l -w

ci-release-docs: install-packer-sdc
	@packer-sdc renderdocs -src docs -partials docs-partials/ -dst docs/
	@/bin/sh -c "[ -d docs ] && zip -r docs.zip docs/"

plugin-check: install-packer-sdc build
	@packer-sdc plugin-check $(GOPATH)/bin/packer-plugin-sakuracloud