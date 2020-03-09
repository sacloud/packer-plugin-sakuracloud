TEST?=$$(go list ./... | grep -v vendor)
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
CURRENT_VERSION = $(gobump show -r version/)

default: test vet

.PHONY: tools
tools:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get github.com/tcnksm/ghr
	GO111MODULE=off go get github.com/client9/misspell/cmd/misspell
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.23.8/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.23.8

.PHONY: clean
clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: install build build-x
install: build
	cp -f $(CURDIR)/bin/packer-builder-sakuracloud $(GOPATH)/bin/packer-builder-sakuracloud

build: clean 
	go build -mod vendor -ldflags "-s -w" -o $(CURDIR)/bin/packer-builder-sakuracloud $(CURDIR)/main.go

build-x: clean vet
	sh -c "'$(CURDIR)/scripts/build.sh'"

.PHONY: test testacc
test: vet
	go test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

# testacc runs acceptance tests
testacc:
	@echo "WARN: Acceptance tests will take a long time to run and may cost money. Ctrl-C if you want to cancel."
	PACKER_ACC=1 go test -v $(TEST) $(TESTARGS) -timeout=45m

.PHONY: lint vet fmt golint goimports
lint: fmt golangci-lint goimports

fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

golangci-lint: fmt
	golangci-lint run ./...

goimports:
	find . -name '*.go' | grep -v vendor | xargs goimports -l -w

.PHONY: docker-shell docker-test docker-testacc docker-build
docker-shell:
	docker-compose run --rm packer

docker-test:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-testacc:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"

.PHONY: prepare-homebrew
prepare-homebrew:
	rm -rf homebrew-packer-builder-sakuracloud/; \
	sh -c "'$(CURDIR)/scripts/update_homebrew_formula.sh' '$(CURRENT_VERSION)'"

