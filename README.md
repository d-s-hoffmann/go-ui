# WebAssembly Gallery

Animated 2D scenes (Plasma, Julia set, Lissajous) rendered in the browser
by a WebAssembly client. There are two clients: the original one in Go and
a Rust port in progress. The development server is Go in both cases. The
only JavaScript is toolchain-generated glue plus a few-line loader.

## Quick start

    make run                  # Go client: builds into build/go/, serves on :8080
    make run GALLERY=rust     # Rust client: builds into build/rust/
    make run ADDR=:9090
    make vet                  # go vet, cargo fmt --check, cargo clippy
    make clean                # removes build/ and Cargo's target/

## Installing the toolchains

Go: install the version named in `go.mod` from https://go.dev/dl/.

Rust (only needed for `GALLERY=rust`):

1. A C linker. Rust's build scripts and macros run natively and need one.
   On Debian/Ubuntu:

       sudo apt install gcc

2. rustup, the Rust toolchain installer:

       curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

   Restart the shell (or `. ~/.cargo/env`) afterwards. The exact Rust
   version, the wasm32 target, rustfmt and clippy are pinned in
   `apps/gallery-rs/rust-toolchain.toml`; rustup installs them on the
   first build.

3. The wasm-bindgen CLI. Its version must equal the `wasm-bindgen` version
   in `apps/gallery-rs/Cargo.toml` (currently 0.2.130):

       cargo install wasm-bindgen-cli --version 0.2.130 --locked

   or download a prebuilt binary from
   https://github.com/wasm-bindgen/wasm-bindgen/releases and put
   `wasm-bindgen` on your `PATH`.

## Layout

    apps/gallery/          Go WASM client (GOOS=js GOARCH=wasm)
    apps/gallery/scene/    scenes, rendered into image.RGBA; no browser APIs
    apps/gallery/ui/       browser glue: canvas, animation loop, controls
    apps/gallery-rs/       Rust WASM client (wasm32-unknown-unknown)
    apps/server/           development file server (Go)
    static/                hand-written page shared by both clients
    static/go/, static/rust/
                           per-client loaders, copied next to the page
    build/                 generated output, not committed
