TEST?=$$(go list ./... | grep -v vendor)
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
CURRENT_VERSION = $(shell cat sakuracloud/version.go | perl -ne 'if(/^var Version = "([0-9\.]+)"/){print $$1;exit}')

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

install: build
	cp -f $(CURDIR)/bin/packer-builder-sakuracloud $(GOPATH)/bin/packer-builder-sakuracloud

build: clean vet
	govendor build -ldflags "-s -w" -o $(CURDIR)/bin/packer-builder-sakuracloud $(CURDIR)/main.go

build-x: clean vet
	sh -c "'$(CURDIR)/scripts/build.sh'"

test: vet
	govendor test $(TEST) $(TESTARGS) -v -timeout=30m -parallel=4 ;

vet: fmt
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -s -l -w $(GOFMT_FILES)

docker-shell:
	docker-compose run --rm packer

docker-test:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"

prepare-homebrew:
	rm -rf homebrew-packer-builder-sakuracloud/; \
	sh -c "'$(CURDIR)/scripts/update_homebrew_formula.sh' '$(CURRENT_VERSION)'"

.PHONY: default test vet fmt lint
