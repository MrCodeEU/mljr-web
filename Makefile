# mljr-web Makefile — PROJECT-parameterized.

# --- pinned versions (bump + `make upgrade-deps` to refresh) -------------------
TAILWIND_VERSION ?= v4.1.17
DATASTAR_VERSION ?= 1.0.2
ALTCHA_VERSION   ?= 3.0.11
LEAFLET_VERSION  ?= 1.9.4
LEAFLET_MARKERCLUSTER_VERSION ?= 1.5.3

# --- paths ---------------------------------------------------------------------
PROJECT ?= homepage
PKG      := ./projects/$(PROJECT)
CSS_IN   := projects/$(PROJECT)/assets/css/input.css
CSS_OUT  := projects/$(PROJECT)/assets/static/app.css
PORT     ?= 8090
BUILD_TAGS := $(if $(filter showcase,$(PROJECT)),showcase)
BUILD_TAG_FLAGS := $(if $(BUILD_TAGS),-tags=$(BUILD_TAGS))

TAILWIND := bin/tailwindcss
AIR      := bin/air

# Pretty URL line printed by `make dev` / `make dev-showcase`.
URL_LINE = printf "\n  \033[1;36m▸ http://localhost:%s\033[0m\n\n"

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)
ifeq ($(UNAME_S),Linux)
  ifeq ($(UNAME_M),x86_64)
    TW_BIN := tailwindcss-linux-x64
  else ifeq ($(UNAME_M),aarch64)
    TW_BIN := tailwindcss-linux-arm64
  endif
endif
ifeq ($(UNAME_S),Darwin)
  ifeq ($(UNAME_M),arm64)
    TW_BIN := tailwindcss-macos-arm64
  else
    TW_BIN := tailwindcss-macos-x64
  endif
endif

.PHONY: help setup tailwind vendor-js air icons dev dev-showcase build check fmt vet lint test test-showcase vuln guard-classes docker upgrade-deps clean sync-assets data-update data-pull

help:           ## list targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN{FS=":.*?## "}; {printf "  %-22s %s\n", $$1, $$2}'

# --- setup ---------------------------------------------------------------------
setup: tailwind vendor-js air ## fetch tailwind binary + datastar.js + altcha.js + air + go mod tidy
	go mod tidy

air: $(AIR) ## install air for Go hot-reload into ./bin
$(AIR):
	@mkdir -p bin
	GOBIN=$(PWD)/bin go install github.com/air-verse/air@latest

tailwind: $(TAILWIND) ## download tailwind standalone binary
$(TAILWIND):
	@mkdir -p bin
	curl -fsSL https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/$(TW_BIN) -o $(TAILWIND)
	chmod +x $(TAILWIND)

vendor-js: ## vendor datastar.js + altcha.js (ESM browser bundle) into every project's assets/static
	@mkdir -p projects/homepage/assets/static/altcha-workers projects/showcase/assets/static
	curl -fsSL https://cdn.jsdelivr.net/gh/starfederation/datastar@v$(DATASTAR_VERSION)/bundles/datastar.js \
		-o projects/homepage/assets/static/datastar.js
	cp projects/homepage/assets/static/datastar.js projects/showcase/assets/static/datastar.js
	# altcha: external ESM bundle plus the local SHA worker used by altcha-loader.js
	curl -fsSL https://cdn.jsdelivr.net/npm/altcha@$(ALTCHA_VERSION)/dist/external/altcha.min.js \
		-o projects/homepage/assets/static/altcha.js
	curl -fsSL https://cdn.jsdelivr.net/npm/altcha@$(ALTCHA_VERSION)/dist/workers/sha.js \
		-o projects/homepage/assets/static/altcha-workers/sha.js
	cp projects/homepage/assets/static/altcha.js projects/showcase/assets/static/altcha.js
	cp -R projects/homepage/assets/static/altcha-workers projects/showcase/assets/static/altcha-workers
	# leaflet + markercluster: vendored map JS/CSS for the open-map component
	curl -fsSL https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.js -o projects/homepage/assets/static/leaflet.js
	curl -fsSL https://unpkg.com/leaflet@$(LEAFLET_VERSION)/dist/leaflet.css -o projects/homepage/assets/static/leaflet.css
	curl -fsSL https://unpkg.com/leaflet.markercluster@$(LEAFLET_MARKERCLUSTER_VERSION)/dist/leaflet.markercluster.js -o projects/homepage/assets/static/leaflet.markercluster.js
	curl -fsSL https://unpkg.com/leaflet.markercluster@$(LEAFLET_MARKERCLUSTER_VERSION)/dist/MarkerCluster.css -o projects/homepage/assets/static/MarkerCluster.css
	curl -fsSL https://unpkg.com/leaflet.markercluster@$(LEAFLET_MARKERCLUSTER_VERSION)/dist/MarkerCluster.Default.css -o projects/homepage/assets/static/MarkerCluster.Default.css
	for f in leaflet.js leaflet.css leaflet.markercluster.js MarkerCluster.css MarkerCluster.Default.css; do \
	  cp projects/homepage/assets/static/$$f projects/showcase/assets/static/$$f; \
	done
	# strip sourceMappingURL comments: we don't ship .map files, and the
	# resulting browser devtools 404s are just noise
	for f in projects/homepage/assets/static/datastar.js projects/showcase/assets/static/datastar.js \
	         projects/homepage/assets/static/leaflet.js projects/showcase/assets/static/leaflet.js \
	         projects/homepage/assets/static/leaflet.markercluster.js projects/showcase/assets/static/leaflet.markercluster.js; do \
	  sed -i -e '/sourceMappingURL/d' $$f; \
	done

upgrade-deps: ## re-fetch tailwind + datastar + altcha at the version pins
	rm -f $(TAILWIND)
	$(MAKE) tailwind vendor-js

# --- code-gen ------------------------------------------------------------------
icons: ## regenerate ui/icon/icons_gen.go from tools/icongen/icons.txt
	go run ./tools/icongen

# --- dev / build ---------------------------------------------------------------
dev: $(TAILWIND) $(AIR) sync-assets data-update ## hot-reload a project:  make dev PROJECT=homepage
	$(TAILWIND) -i $(CSS_IN) -o $(CSS_OUT)
	@$(URL_LINE) $(PORT)
	@$(TAILWIND) -i $(CSS_IN) -o $(CSS_OUT) --watch & TW=$$!; \
	  trap "kill $$TW 2>/dev/null" EXIT INT TERM; \
	  if [ -f projects/$(PROJECT)/.env ]; then set -a; . projects/$(PROJECT)/.env; set +a; fi; \
	  MLJR_ENV=dev PORT=$(PORT) AIR_MAIN_PKG=$(PKG) AIR_BUILD_TAGS= \
	    $(AIR) -c .air.toml

dev-showcase: $(TAILWIND) $(AIR) ## run the catalogue (build tag: showcase)
	$(TAILWIND) -i projects/showcase/assets/css/input.css -o projects/showcase/assets/static/app.css
	@$(URL_LINE) 8091
	@$(TAILWIND) -i projects/showcase/assets/css/input.css -o projects/showcase/assets/static/app.css --watch & TW=$$!; \
	  trap "kill $$TW 2>/dev/null" EXIT INT TERM; \
	  if [ -f projects/showcase/.env ]; then set -a; . projects/showcase/.env; set +a; fi; \
	  MLJR_ENV=dev PORT=8091 AIR_MAIN_PKG=./projects/showcase AIR_BUILD_TAGS=showcase \
	    $(AIR) -c .air.toml

build: $(TAILWIND) sync-assets ## static binary -> bin/$(PROJECT)
	$(TAILWIND) -i $(CSS_IN) -o $(CSS_OUT) --minify
	CGO_ENABLED=0 go build $(BUILD_TAG_FLAGS) -buildvcs=false -trimpath -ldflags="-s -w -X mljr-web/internal/version.Tag=$$(cat VERSION 2>/dev/null || echo dev)" -o bin/$(PROJECT) $(PKG)

docker: ## per-project image:  make docker PROJECT=homepage TAG=v1
	docker build -f projects/$(PROJECT)/Dockerfile -t mljr/$(PROJECT):$(TAG) .

# --- data ------------------------------------------------------------------
data-update: ## regenerate mljr-data/generated/site-data.json + update seed-cache (needs GITHUB_TOKEN)
	@if [ -f projects/homepage/.env ]; then set -a; . projects/homepage/.env; set +a; fi; \
	  cd mljr-data/generator && \
	  GITHUB_TOKEN="$$GITHUB_TOKEN" GITHUB_USER="$${GITHUB_USERNAME:-MrCodeEU}" \
	  STRAVA_CLIENT_ID="$$STRAVA_CLIENT_ID" STRAVA_CLIENT_SECRET="$$STRAVA_CLIENT_SECRET" STRAVA_REFRESH_TOKEN="$$STRAVA_REFRESH_TOKEN" \
	  go run ./cmd/generate
	@cp mljr-data/generated/site-data.json projects/homepage/data/seed-cache.json
	@echo "seed-cache.json updated"

data-pull: ## fetch latest mljr-data submodule from origin
	git submodule update --remote --merge mljr-data

sync-assets: ## copy mljr-data portfolio images into projects/homepage/assets/static/portfolio
	@rm -rf projects/homepage/assets/static/portfolio
	@mkdir -p projects/homepage/assets/static/portfolio
	@cp -r mljr-data/assets/portfolio/. projects/homepage/assets/static/portfolio/

# --- quality -------------------------------------------------------------------
check: fmt vet lint guard-classes test test-showcase vuln ## full local gate
fmt:
	gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	elif [ -x bin/goimports ]; then \
		bin/goimports -w .; \
	fi
vet:  ; go vet ./...
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi
vuln:
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	elif [ -x bin/govulncheck ]; then \
		bin/govulncheck ./...; \
	else \
		echo "govulncheck not installed, skipping"; \
	fi
test: ; go test ./... -race -cover
test-showcase: ; go test -tags showcase ./... -race -cover

guard-classes: ## enforce data-* contract: no class= / Class( in ui/**.go
	@! grep -rnE '(Class\(|h\.Class|[^-]class=)' ui --include='*.go' \
	   || (echo "✗ ui components must not use classes — use data-* + CSS" && exit 1)
	@echo "✓ guard-classes ok"

clean:
	rm -rf bin/*.css bin/$(PROJECT)
