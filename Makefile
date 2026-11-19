.RECIPEPREFIX := >
GO        ?= go
GOROOT    := $(shell $(GO) env GOROOT)
ADDR      ?= :8080
WASM_EXEC := $(GOROOT)/lib/wasm/wasm_exec.js

PUBLIC    := public
APP_WASM  := $(PUBLIC)/app.wasm
WASM_JS   := $(PUBLIC)/wasm_exec.js

.PHONY: run build clean help

## Build everything, then serve.
run: build
>$(GO) run ./apps/server -addr $(ADDR)

## Compile the WASM client.
$(APP_WASM): $(wildcard cmd/gallery/*.go) go.mod
>mkdir -p $(PUBLIC)
>GOOS=js GOARCH=wasm $(GO) build -o $(APP_WASM) ./apps/gallery

## Ship the toolchain JS glue (copied only if missing).
$(WASM_JS):
>cp "$(WASM_EXEC)" $(WASM_JS)

## Build outputs only.
build: $(APP_WASM) $(WASM_JS)
>@echo build complete: $(APP_WASM)

## Remove generated artifacts (keeps index.html).
clean:
>rm -f $(APP_WASM) $(WASM_JS)

help:
>@grep -E "^## " $(MAKEFILE_LIST) | sed "s/^## //"
