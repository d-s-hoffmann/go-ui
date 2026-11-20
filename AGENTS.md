# AGENTS.md

## Project
A self-written graphical app in Go, compiled to WASM, for desktop and
mobile browsers. Good structure is the top priority.

## Working rules
- Keep PRs small and incremental: at most 500 changed lines, ideally fewer.
  Generated files and toolchain bumps don't count toward the limit.
- One concern per commit; renames and formatting go in their own commits.
- Reviews are strict: plan before coding and explain choices in the PR.
- Follow Go conventions (https://go.dev/doc/modules/layout). Code must pass
  gofmt and go vet, on the Go version in go.mod.

## Layout
- `apps/<name>/`: one deployable program each. Keep `main` thin; packages
  only that app uses go in subfolders (e.g. `apps/gallery/scene/`).
- Rendering code must not import `syscall/js`.
- `static/`: hand-written web files. `build/`: generated, never committed.

## UI
- The menu is separate from the view and collapses to a small button.
- The view fills the whole screen.
- Scenes scale to any size and are interchangeable.

## Testing
- Prefer tests that run the whole WASM binary and check its outputs for
  given inputs.
- Keep rendering deterministic so frames are reproducible.

## Commands
    make build   # -> build/
    make run     # build, then serve on :8080
    make clean
