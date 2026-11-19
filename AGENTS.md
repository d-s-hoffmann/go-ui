# AGENTS.md

## Goal
A self-written graphical app in Go, compiled to WASM, that runs well in
desktop and mobile browsers. Deployment comes later; structure matters now.

## How to work
- Small, incremental PRs: at most 500 changed lines, ideally fewer.
  Generated files and toolchain bumps don't count toward the limit.
- One concern per commit; renames and formatting go in their own commits.
- Reviews are strict: plan before coding and explain choices in the PR.
- Follow Go conventions (https://go.dev/doc/modules/layout). Code must pass
  gofmt and go vet, on the latest Go (see go.mod).

## Layout
- `apps/<name>/`: one deployable program each (`gallery` = WASM client,
  `server` = dev server). Keep them thin.
- `internal/`: shared packages. Rendering must not import `syscall/js`.
- `static/`: hand-written web files. `build/`: generated, never committed.

## Direction
- The menu is separate from the view and collapses to a small button.
- The view fills the whole screen. Scenes (plasma, julia, lissajous) scale
  to any size and are interchangeable.
- Testing (later PRs): test the whole WASM binary against inputs and
  outputs. Keep rendering deterministic so frames are reproducible.

## Commands
    make build   # -> build/
    make run     # build, then serve on :8080
    make clean
