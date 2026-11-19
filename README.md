# Go × WebAssembly Gallery

Server, renderer and UI written entirely in Go. The only JavaScript is the
toolchain-generated wasm_exec.js and the 8-line loader in index.html.

## Quick start

    make run          # builds into build/ if stale, serves on :8080
    make run ADDR=:9090
    make clean        # removes build/

Three animated 2D scenes: Plasma, Julia set, Lissajous.

## Layout

    apps/gallery/   WASM client (GOOS=js GOARCH=wasm)
    apps/server/    development file server
    internal/scene/ scenes, rendered into image.RGBA; no browser APIs
    internal/web/   browser glue: canvas, animation loop, controls
    static/         hand-written files copied into build/ (index.html)
    build/          generated output, not committed
