SHELL=/usr/bin/bash

GOPATH := $(shell go env GOPATH)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOFLAGS := GOOS=$(GOOS) GOARCH=$(GOARCH)
GOCMD := go
verbose ?= 0

SUPPORTED_OSES := windows linux darwin
SUPPORTED_ARCHITECTURES := arm64 amd64

PUML_TARGETS := ./pkg/{go,po}/parse ./pkg/po/compile ./pkg/po

SUDO = $(if $(and $(filter root,$(1)),$(filter-out root,$(USER))),sudo)

LOCAL_PREFIX := $(HOME)/.local
ROOT_PREFIX := /usr/local
PREFIX=$(if $(filter local,$(1)),$(LOCAL_PREFIX),$(ROOT_PREFIX))

override APP_DIRS := $(foreach dir,$(wildcard ./cli/*/),\
	$(if $(wildcard ${dir}*.go),${dir}))
override APP_NAMES := $(foreach dir,$(APP_DIRS),\
	$(shell basename ${dir}))
override WIN_EXT := $(if $(filter windows,$(GOOS)),.exe)
override CMD = $(GOFLAGS) $(GOCMD)
override V_FLAG := $(if $(filter 1,$(verbose)),-v)

.PHONY: all default benchmark bench test test-all run-% clean \
	build-% build-tool-% build-all release puml completions \
	completion-% go-install-all go-uninstall-all %-install-all \
	%-uninstall-all %-install %-uninstall %-completions-install \
	%-completions-uninstall go-install-% go-uninstall-%
.ONESHELL: default changelog.md build-tool-% build-all puml \
	completions go-install-all go-uninstall-all %-install-all \
	%-uninstall-all %-completions-uninstall
.DEFAULT_GOAL := default

default:
	for app in $(APP_NAMES);do
		$(MAKE) build-$$app
	done

benchmark:
	$(MAKE) bench path="./..."

bench:
	$(CMD) test $(V_FLAG) -bench=. $(path)

test:
	$(CMD) clean -testcache
	$(CMD) test $(V_FLAG) ./...

test-all:
	$(MAKE) test
	$(MAKE) GOCMD=go1.18 test

run-%:
	$(CMD) run -C ./cli $(V_FLAG) ./$* $(args)

clean:
	rm -rf $$(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
		builds cli/dist completions changelog.md cli/builds

log.diff:
	git diff --staged > log.diff

changelog.md:
	echo '## Changelog' > changelog.md
	echo >> changelog.md

	latest_tag=$$(git describe --tags --abbrev=0)
	penultimate_tag=$$(git describe --tags --abbrev=0 "$$latest_tag^")
	
	git log --pretty=format:'- [%h](https://github.com/Tom5521/gotext-tools/commit/%H): %s' \
		$$penultimate_tag..$$latest_tag >> changelog.md

build-%:
	$(CMD) build -C ./cli $(V_FLAG) -o \
		../builds/$*-$(GOOS)-$(GOARCH)$(WIN_EXT) \
		-ldflags '-s -w' \
		./$*

build-tool-%:
	valid=$$($(CMD) tool dist list)
	for os in $(SUPPORTED_OSES); do
		for arch in $(SUPPORTED_ARCHITECTURES); do
			if ! echo $$valid | grep -qw "$$os/$$arch"; then
				continue
			fi
			$(MAKE) GOOS=$$os GOARCH=$$arch build-$*
		done
	done

build-all: clean
	for app in $(APP_NAMES); do
		$(MAKE) build-tool-$$app
	done

cli/docs:
	$(CMD) run -C cli ./gotext-tools doc-tree -D docs

release: clean changelog.md build-all
	gh release create $$(git describe --tags --abbrev=0) \
		--notes-file ./changelog.md --fail-on-no-commits builds/*

puml:
	for path in $(PUML_TARGETS); do
		structure="$$path/structure.puml"
		goplantuml "$$path" > "$$structure"
		plantuml -theme spacelab "$$structure"
	done

completions:
	mkdir -p completions

	for app in $(APP_NAMES);do
		$(MAKE) completion-$$app
	done

completion-%:
	mkdir -p completions
	go run -C cli ./$* completion bash > ./completions/$*.bash
	go run -C cli ./$* completion fish > ./completions/$*.fish
	go run -C cli ./$* completion zsh > ./completions/$*.zsh

go-install-all:
	for app in $(APP_NAMES);do
		$(MAKE) go-install-$$app
	done

go-uninstall-all:
	for app in $(APP_NAMES);do
		$(MAKE) go-uninstall-$$app
	done

%-install-all:
	for app in $(APP_NAMES);do
		$(MAKE) $*-install APP=$$app
	done

%-uninstall-all:
	for app in $(APP_NAMES);do
		$(MAKE) $*-uninstall APP=$$app
	done

%-install: build-$(APP)
	$(call SUDO,$*) install -D "./builds/$(APP)-$(GOOS)-$(GOARCH)" \
		$(call PREFIX,$*)/bin/$(APP)
	$(MAKE) $*-completions-install

%-uninstall:
	$(call SUDO,$*) rm -f $(call PREFIX,$*)/bin/$(APP)
	$(MAKE) $*-completions-uninstall

%-completions-install: completion-$(APP)
	$(call SUDO,$*) install -D "./completions/$(APP).fish" \
		$(call PREFIX,$*)/share/fish/vendor_completions.d/$(APP).fish
	$(call SUDO,$*) install -D "./completions/$(APP).bash" \
		$(call PREFIX,$*)/share/bash-completion/completions/$(APP)
	$(call SUDO,$*) install -D "./completions/$(APP).zsh" \
		$(call PREFIX,$*)/share/zsh/site-functions/_$(APP)

%-completions-uninstall:
	prefix=$(call PREFIX,$*)
	$(call SUDO,$*) rm -f "$$prefix/share/fish/vendor_completions.d/$(APP).fish" \
	"$$prefix/share/bash-completion/completions/$(APP)" \
	"$$prefix/share/zsh/site-functions/_$(APP)"

go-install-%:
	$(CMD) install -C cli $(V_FLAG) ./$*
	$(MAKE) local-completions-install APP=$*

go-uninstall-%:
	rm -f $(GOPATH)/bin/$*
	$(MAKE) local-completions-uninstall APP=$*
