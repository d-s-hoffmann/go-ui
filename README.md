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

Debian/Ubuntu:

```sh
sudo apt install gcc golang rustup
echo 'export PATH="$HOME/.cargo/bin:$PATH"' >> ~/.profile && . ~/.profile
rustup default stable
cargo install wasm-bindgen-cli --version 0.2.130 --locked
```

macOS (Homebrew):

```sh
brew install go rustup
echo 'export PATH="$(brew --prefix rustup)/bin:$HOME/.cargo/bin:$PATH"' >> ~/.zprofile && . ~/.zprofile
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
