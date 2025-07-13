SHELL=/usr/bin/bash

gopath := $(shell go env GOPATH)
goos := $(shell go env GOOS)
goarch := $(shell go env GOARCH)
goflags := GOOS=${goos} GOARCH=${goarch}
gocmd := go
verbose ?= 0

supported_oses = windows linux darwin
supported_archs = arm64 amd64

puml_targets = ./pkg/{go,po}/parse ./pkg/po/compile ./pkg/po

local_prefix = $(HOME)/.local
root_prefix = /usr/local

override app_paths = $(shell ./scripts/app-paths.sh)
override app_names = $(shell ./scripts/app-names.sh)
override bin_ext := $(shell [[ $(goos) == "windows" ]] \
	&& echo ".exe")
override cmd = $(goflags) $(gocmd)
override v_flag := $(shell [ $(verbose) -eq 1 ] && echo "-v")

.PHONY: $(shell grep -E '^[a-zA-Z0-9_-]+:' Makefile | cut -d ':' -f 1)

default:
	for app in $(app_names);do \
		$(MAKE) build-$$app;\
	done

benchmark:
	$(MAKE) bench path="./..."

bench:
	$(cmd) test $(v_flag) -bench=. $(path)

test:
	$(cmd) clean -testcache
	$(cmd) test $(v_flag) ./...

test-all:
	$(MAKE) test
	$(MAKE) gocmd=go1.18 test

run-%:
	$(cmd) run -C ./cli $(v_flag) ./$* $(args)

clean:
	rm -rf $$(find . \( -name "*.po" -o -name "*.mo" -o -name "*.pot" -o -name "*.log" \)) \
		builds cli/dist completions

diff:
	git diff --staged > log.diff

_gen_releasemsg:
	echo '## Changelog' > changelog.md
	echo >> changelog.md

	latest_tag=$$(git describe --tags --abbrev=0);\
	penultimate_tag=$$(git describe --tags --abbrev=0 "$$latest_tag^");\
	\
	git log --pretty=format:'- [%h](https://github.com/Tom5521/gotext-tools/commit/%H): %s' \
		$$penultimate_tag..$$latest_tag >> changelog.md

build-%:
	$(cmd) build -C ./cli $(v_flag) -o \
		../builds/$*-$(goos)-$(goarch)$(bin_ext) \
		-ldflags '-s -w' $(v_flag) \
		./$*

build-tool-%:
	valid=$$($(cmd) tool dist list);\
	for os in $(supported_oses); do \
		for arch in $(supported_archs); do \
			if ! echo $$valid | grep -qw "$$os/$$arch"; then \
				continue;\
			fi;\
			$(MAKE) goos=$$os goarch=$$arch build-$*;\
		done;\
	done

build-all: clean
	for app in $(app_names); do \
		$(MAKE) build-tool-$$app;\
	done

cli-docs:
	rm -rf cli/docs
	$(cmd) run -C cli ./gotext-tools doc-tree -D docs

release: clean _gen_releasemsg build-all
	gh release create $$(git describe --tags --abbrev=0) \
		--notes-file ./changelog.md --fail-on-no-commits builds/*

puml:
	for path in $(puml_targets); do \
		structure="$$path/structure.puml";\
		goplantuml "$$path" > "$$structure";\
		plantuml -theme spacelab "$$structure";\
	done

completions:
	mkdir -p completions

	for app in $(app_names);do \
		$(MAKE) completion-$$app;\
	done

completion-%:
	go run -C cli ./$* completion bash > ./completions/$*.bash
	go run -C cli ./$* completion fish > ./completions/$*.fish
	go run -C cli ./$* completion zsh > ./completions/$*.zsh

_%-prefix:
	@if [[ "$*" == "local" ]]; then \
		echo $(local_prefix);\
	else \
		echo $(root_prefix);\
	fi;

go-install-all:
	for app in $(app_names);do \
		$(MAKE) go-install-$$app;\
	done

go-uninstall-all:
	for app in $(app_names);do \
		$(MAKE) go-uninstall-$$app;\
	done

%-install-all:
	for app in $(app_names);do \
		$(MAKE) $*-install app=$$app;\
	done

%-uninstall-all:
	for app in $(app_names);do \
		$(MAKE) $*-uninstall app=$$app;\
	done

%-install: build-$(app)
	install -D "./builds/$(app)-$(goos)-$(goarch)" \
		"$$($(MAKE) -s _$*-prefix)/bin/$(app)"
	$(MAKE) $*-completions-install

%-uninstall:
	rm -f "$$($(MAKE) -s _$*-prefix)/bin/$(app)"
	$(MAKE) $*-completions-uninstall

%-completions-install: completions
	install -D "./completions/$(app).fish" \
		"$$($(MAKE) -s _$*-prefix)/share/fish/vendor_completions.d/$(app).fish"
	install -D "./completions/$(app).bash" \
		"$$($(MAKE) -s _$*-prefix)/share/bash-completion/completions/$(app)"
	install -D "./completions/$(app).zsh" \
		"$$($(MAKE) -s _$*-prefix)/share/zsh/site-functions/_$(app)"

%-completions-uninstall:
	prefix=$$($(MAKE) -s _$*-prefix);\
	rm -f "$$prefix/share/fish/vendor_completions.d/$(app).fish" \
	"$$prefix/share/bash-completion/completions/$(app)" \
	"$$prefix/share/zsh/site-functions/_$(app)"

go-install-%:
	$(cmd) install -C cli $(v_flag) ./$*
	$(MAKE) local-completions-install app=$*

go-uninstall-%:
	rm -f $(gopath)/bin/$*
	$(MAKE) local-completions-uninstall app=$*
