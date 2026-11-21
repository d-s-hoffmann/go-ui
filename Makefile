.RECIPEPREFIX := >
GO        ?= go
CARGO     ?= cargo
BINDGEN   ?= wasm-bindgen
ADDR      ?= :8080
WASM_EXEC := $(shell $(GO) env GOROOT)/lib/wasm/wasm_exec.js

## Which gallery client to build and serve: go or rust.
GALLERY   ?= go
GALLERIES := go rust

STATIC    := static
BUILD     := build
OUT       := $(BUILD)/$(GALLERY)
GO_SRC    := $(shell find . -name '*.go' -not -path './$(BUILD)/*') go.mod

RS_DIR    := apps/gallery-rs
RS_WASM   := $(RS_DIR)/target/wasm32-unknown-unknown/release/gallery_rs.wasm
RS_SRC    := $(shell find $(RS_DIR) -path $(RS_DIR)/target -prune -o -type f -print)

# Shared files in static/, plus the ones in static/$(GALLERY)/.
SHARED    := $(filter-out $(addprefix $(STATIC)/,$(GALLERIES)),$(wildcard $(STATIC)/*))
ASSETS    := $(patsubst $(STATIC)/%,$(OUT)/%,$(SHARED)) \
             $(patsubst $(STATIC)/$(GALLERY)/%,$(OUT)/%,$(wildcard $(STATIC)/$(GALLERY)/*))

ifeq ($(GALLERY),go)
CLIENT    := $(OUT)/app.wasm $(OUT)/wasm_exec.js
else ifeq ($(GALLERY),rust)
CLIENT    := $(OUT)/gallery_rs.js
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

## Compile the Rust client. Cargo runs inside the crate so that its
## rust-toolchain.toml picks the pinned toolchain.
$(RS_WASM): $(RS_SRC)
>cd $(RS_DIR) && $(CARGO) build --locked --release --target wasm32-unknown-unknown

## Generate the JS glue (gallery_rs.js) and final gallery_rs_bg.wasm.
$(BUILD)/rust/gallery_rs.js: $(RS_WASM)
>$(BINDGEN) --target web --no-typescript --out-dir $(@D) $<

## Copy hand-written static files; per-gallery files win over shared ones.
$(OUT)/%: $(STATIC)/$(GALLERY)/%
>@mkdir -p $(@D)
>cp "$<" $@

$(OUT)/%: $(STATIC)/%
>@mkdir -p $(@D)
>cp "$<" $@

## Vet all Go packages (natively and for js/wasm), then lint the Rust crate.
vet:
>$(GO) vet ./...
>GOOS=js GOARCH=wasm $(GO) vet ./...
>cd $(RS_DIR) && $(CARGO) fmt --check
>cd $(RS_DIR) && $(CARGO) clippy --locked --release --target wasm32-unknown-unknown -- -D warnings

## Remove all generated output.
clean:
>rm -rf $(BUILD) $(RS_DIR)/target

help:
>@grep -E "^## " $(MAKEFILE_LIST) | sed "s/^## //"
