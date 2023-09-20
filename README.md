<div align="center">
    <img src=".github/grove_logo.png" alt="Grove logo" width="600"/>
    <h1>Utils Go</h1>
    <big>The Utils Go package provides common used utility funcs, an HTTP client and more to other Go backend packages.</big>
</div>
<br/>

## Pre-Commit Installation

Run the command `make init-pre-commit` from the repository root.

Once this is done, the following commands will be performed on every commit to the repo and must pass before the commit is allowed:

- **go-fmt** - Runs `gofmt`
- **go-imports** - Runs `goimports`
- **golangci-lint** - run `golangci-lint run ./...`
- **go-critic** - run `gocritic check ./...`
- **go-build** - run `go build`
- **go-mod-tidy** - run `go mod tidy -v`
