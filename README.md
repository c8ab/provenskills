
# Proven Skills

[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/c8ab/provenskills/badge)](https://scorecard.dev/viewer/?uri=github.com/c8ab/provenskills)

> **Conceptual — we are exploring whether this approach can work.** Not yet ready for general use.

Package AI skills as OCI artifacts — with provenance, signing, and standard registry distribution.

## What is Proven Skills?

Proven Skills (`psk`) packages AI skill definitions into OCI-compliant artifacts. Skills become versioned, signable, and distributable through any OCI registry — the same infrastructure already used for container images.

## Why OCI?

- **Immutable, versioned artifacts** — every skill version is a distinct, content-addressable artifact
- **Sign with existing tooling** — use cosign, Notation, or any OCI signing tool for authenticity
- **Provenance built-in** — supply-chain integrity for AI skills using established standards
- **Use any OCI registry** — Docker Hub, GitHub Container Registry, AWS ECR, or your own

## Architecture

Proven Skills is split into two components:

- **`psk` CLI** — builds skill artifacts, pushes to and pulls from OCI registries, and manages a local store
- **Tool plugins** — integrate with AI coding tools (OpenCode, Goose, and others) to manage skills in tool-specific locations (e.g., `.opencode/skills/`). The CLI handles artifacts; plugins handle how each tool consumes them.

## Current Status

> **Work in progress** — under active development and not yet ready for general use.

Today `psk` can validate skill definitions, build packages to a local store, and list stored artifacts. OCI registry push/pull and tool plugins are on the roadmap.

## Quick Start

### Requirements

- Go 1.23+

### Install

```sh
go install github.com/c8ab/provenskills/cmd/psk@latest
```

This places the `psk` binary in your `$GOBIN` (usually `~/go/bin`). Make sure it's on your `PATH`.

### Build

```sh
go build -o psk ./cmd/psk/
```

### Usage

```sh
# Validate a skill definition
psk validate ./path/to/skill-dir

# Build a skill package
psk build ./path/to/skill-dir --maintainer "Name <email>"

# List skills in the local store
psk list
```

## Development

```sh
# Run tests
go test ./...

# Run linter (requires golangci-lint v2)
golangci-lint run ./...
```

## License

See [LICENSE](LICENSE) for details.
