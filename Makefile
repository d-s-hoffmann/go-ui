.RECIPEPREFIX := >
GO        ?= go
ADDR      ?= :8080
WASM_EXEC := $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js

STATIC    := static
BUILD     := build
APP_WASM  := $(BUILD)/app.wasm
WASM_JS   := $(BUILD)/wasm_exec.js
ASSETS    := $(patsubst $(STATIC)/%,$(BUILD)/%,$(wildcard $(STATIC)/*))
GO_SRC    := $(shell find . -name '*.go' -not -path './$(BUILD)/*') go.mod

.PHONY: run build vet clean help

## Build everything, then serve.
run: build
>$(GO) run ./apps/server -addr $(ADDR) -dir $(BUILD)

## Build outputs only.
build: $(APP_WASM) $(WASM_JS) $(ASSETS)
>@echo build complete: $(BUILD)/

## Compile the WASM client.
$(APP_WASM): $(GO_SRC)
>@mkdir -p $(BUILD)
>GOOS=js GOARCH=wasm $(GO) build -o $@ ./apps/gallery

## Copy the JS glue from the active toolchain; it must match the Go version.
$(WASM_JS): $(WASM_EXEC)
>@mkdir -p $(BUILD)
>cp "$<" $@

## Copy hand-written static files.
$(BUILD)/%: $(STATIC)/%
>@mkdir -p $(BUILD)
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
