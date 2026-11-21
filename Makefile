.RECIPEPREFIX := >
GO        ?= go
ADDR      ?= :8080
WASM_EXEC := $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js

## Which gallery client to build and serve: go or rust.
GALLERY   ?= go
GALLERIES := go rust

STATIC    := static
BUILD     := build
OUT       := $(BUILD)/$(GALLERY)
GO_SRC    := $(shell find . -name '*.go' -not -path './$(BUILD)/*') go.mod

# Shared files in static/, plus the ones in static/$(GALLERY)/.
SHARED    := $(filter-out $(addprefix $(STATIC)/,$(GALLERIES)),$(wildcard $(STATIC)/*))
ASSETS    := $(patsubst $(STATIC)/%,$(OUT)/%,$(SHARED)) \
             $(patsubst $(STATIC)/$(GALLERY)/%,$(OUT)/%,$(wildcard $(STATIC)/$(GALLERY)/*))

ifeq ($(GALLERY),go)
CLIENT    := $(OUT)/app.wasm $(OUT)/wasm_exec.js
else
$(error GALLERY must be one of: $(GALLERIES))
endif

.PHONY: run build vet clean help

## Build everything, then serve.
run: build
>$(GO) run ./apps/server -addr $(ADDR) -dir $(OUT)

## Build outputs only.
build: $(CLIENT) $(ASSETS)
>@echo build complete: $(OUT)/

## Compile the Go WASM client.
$(BUILD)/go/app.wasm: $(GO_SRC)
>@mkdir -p $(@D)
>GOOS=js GOARCH=wasm $(GO) build -o $@ ./apps/gallery

## Copy the JS glue from the active toolchain; it must match the Go version.
$(BUILD)/go/wasm_exec.js: $(WASM_EXEC)
>@mkdir -p $(@D)
>cp "$<" $@

## Copy hand-written static files; per-gallery files win over shared ones.
$(OUT)/%: $(STATIC)/$(GALLERY)/%
>@mkdir -p $(@D)
>cp "$<" $@

$(OUT)/%: $(STATIC)/%
>@mkdir -p $(@D)
>cp "$<" $@

## Vet all packages. A native run skips js-only code, so vet that too.
vet:
>$(GO) vet ./...
>GOOS=js GOARCH=wasm $(GO) vet ./...

## Remove all generated output.
clean:
>rm -rf $(BUILD)

help:
>@grep -E "^## " $(MAKEFILE_LIST) | sed "s/^## //"
