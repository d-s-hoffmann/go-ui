# Go × WebAssembly Gallery

Server, renderer and UI written entirely in Go. The only JavaScript is the
toolchain-generated wasm_exec.js and the 8-line loader in index.html.

## Quick start

    make run          # builds wasm if stale, serves on :8080
    make run ADDR=:9090
    make clean

Three animated 2D scenes: Plasma, Julia set, Lissajous.
