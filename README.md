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

## Setting up the environment

1. Install Go, Rust's installer `rustup`, and a C linker (Rust needs one):

   ```sh
   sudo apt install gcc golang rustup    # Debian/Ubuntu
   brew install gcc go rustup            # Homebrew, e.g. on Bazzite
   ```

2. Add the tools to your `PATH`. Cargo's folder holds `wasm-bindgen`.
   Homebrew also needs rustup's folder, because it doesn't link rustup
   itself (`brew info rustup`). Run the line for your setup:

   ```sh
   # apt + bash
   echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.bashrc && . ~/.bashrc
   # Homebrew + bash
   echo 'export PATH="$(brew --prefix rustup)/bin:$HOME/.cargo/bin:$PATH"' >> ~/.bashrc && . ~/.bashrc
   ```

   ```fish
   # apt + fish
   fish_add_path ~/.cargo/bin
   # Homebrew + fish
   fish_add_path (brew --prefix rustup)/bin ~/.cargo/bin
   ```

3. Set up Rust and install the wasm-bindgen CLI (same in every shell):

   ```sh
   rustup default stable
   cargo install wasm-bindgen-cli --version 0.2.130 --locked
   ```

Go downloads the exact version from `go.mod` by itself if the installed one
is older. Likewise, rustup installs the Rust version and wasm target pinned in
`apps/gallery-rs/rust-toolchain.toml` on the first build. The wasm-bindgen
CLI version must match `apps/gallery-rs/Cargo.toml`.

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
