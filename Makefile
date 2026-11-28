.DEFAULT_GOAL := help

GO ?= go
NPM ?= npm
PYTHON ?= python3
CLI := bin/cli-silly
GO_PACKAGES := ./cmd/... ./course ./internal/...

.PHONY: help setup cli doctor run-dev web-deps web-dev web-check go-check course-check terminal-check docker-check check package clean

help:
	@printf '%s\n' \
		'make setup          Fetch locked Go and website dependencies' \
		'make cli            Build bin/cli-silly' \
		'make doctor         Show available tools by workshop' \
		'make run-dev        Open the CLI with checkout-local workspaces and progress' \
		'make web-deps       Install website dependencies from package-lock.json' \
		'make web-dev        Start the companion site on localhost' \
		'make check          Run Go checks, the terminal smoke test, and website checks' \
		'make docker-check   Grade the Compose labs against a running daemon' \
		'make package VERSION=0.1.0  Build Linux archives in artifacts/releases' \
		'make clean          Remove bin, dist, dist-ssr, and artifacts'

setup: web-deps
	$(GO) mod download

cli:
	$(GO) build -o $(CLI) ./cmd/cli-silly

doctor: cli
	$(CLI) doctor

run-dev: cli
	XDG_DATA_HOME="$(CURDIR)/.local/data" XDG_STATE_HOME="$(CURDIR)/.local/state" $(CLI)

web-deps:
	$(NPM) ci

web-dev:
	$(NPM) run dev

web-check:
	$(NPM) run lint
	$(NPM) run build

course-check:
	$(PYTHON) scripts/course-quality.py
	$(PYTHON) scripts/test-course-quality.py

go-check: course-check
	$(GO) test $(GO_PACKAGES)
	$(GO) vet $(GO_PACKAGES)

terminal-check: cli
	$(PYTHON) scripts/test-terminal.py

docker-check:
	docker compose version
	docker info --format '{{.ServerVersion}}'
	CLI_SILLY_TEST_DOCKER=1 $(GO) test ./internal/lab -run 'TestCatalogLabs/docker' -count=1 -v

check: go-check terminal-check web-check

package:
	@if [ -z "$(VERSION)" ]; then printf '%s\n' 'Set VERSION, for example make package VERSION=0.1.0' >&2; exit 2; fi
	GO="$(GO)" bash scripts/package-linux.sh "$(VERSION)"

clean:
	rm -rf -- bin dist dist-ssr artifacts
