TEST?=$$(go list ./... | grep -v vendor)
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: lint build

.PHONY: tools
tools:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get github.com/client9/misspell/cmd/misspell
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.38.0/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.38.0
	go install github.com/hashicorp/packer/cmd/mapstructure-to-hcl2@latest

.PHONY: clean
clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: install build build-x
install: build
	cp -f $(CURDIR)/bin/packer-plugin-sakuracloud $(GOPATH)/bin/packer-plugin-sakuracloud

build: generate clean
	CGO_ENABLED=0 go build -mod vendor -ldflags "-s -w -extldflags -static" -o $(CURDIR)/bin/packer-plugin-sakuracloud $(CURDIR)/main.go

generate:
	go generate ./...

.PHONY: test testacc
test:
	go test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

# testacc runs acceptance tests
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
