# Implementation Plan: Skill Packaging CLI

**Branch**: `001-skill-packaging-cli` | **Date**: 2026-02-11 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-skill-packaging-cli/spec.md`

## Summary

Build `psk`, a CLI tool that packages AI Agent Skills (following
the agentskills.io format) into Proven Skill Artifacts stored in
a local filesystem-based store. The tool validates skill metadata,
enforces the agentskills.io naming spec, and tracks maintainer
provenance. Go is the implementation language, chosen for its
single-dependency YAML parsing, stdlib JSON/testing/exec support,
and trivial static binary compilation — directly aligned with the
constitution's minimal dependency and small auditable binary
mandates.

## Technical Context

**Language/Version**: Go 1.23+
**Primary Dependencies**: `gopkg.in/yaml.v3` (sole external dep)
**Storage**: Local filesystem — directory-per-artifact under
`~/.psk/store/` (overridable via `PSK_STORE` env var)
**Testing**: `go test` (stdlib) + `os/exec` for CLI integration
tests invoking the real `psk` binary with real filesystem stores
**Target Platform**: macOS (arm64, amd64), Linux (amd64, arm64)
**Project Type**: Single project
**Performance Goals**: Build a skill in <5s for <1MB content;
list 100 skills in <2s
**Constraints**: Single static binary, zero runtime deps, sole
external Go module is `yaml.v3`
**Scale/Scope**: Local-only store, hundreds of skills

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Evidence |
|-----------|--------|----------|
| I. CLI-First | PASS | `psk` is a CLI binary. All commands (`build`, `list`, `validate`) are CLI subcommands. stdout/stderr/exit code contract defined. `--json` flag for machine-readable output. |
| II. TDD (NON-NEGOTIABLE) | PASS | Integration tests invoke the real binary, assert on stdout/stderr/exit codes. Tests written before implementation per constitution. |
| III. Integration Testing over Mocks | PASS | Tests use real filesystem stores via `t.TempDir()`, invoke real `psk` binary via `os/exec`, parse real YAML skill directories. Zero mocks. |
| IV. Security-First | PASS | Input validation at CLI boundary (arg parsing). YAML frontmatter validated before processing. No secrets in output. Dependency (`yaml.v3`) audited: Google-maintained, zero transitive deps, strong security track record. |
| V. Minimal Dependency | PASS | Single external dependency: `gopkg.in/yaml.v3`. JSON, file I/O, testing, process execution all from stdlib. Static binary via `CGO_ENABLED=0`. |
| Code Quality Standards | PASS | Single-responsibility file structure. Public interfaces documented. `go vet` + `golangci-lint` enforced. |
| Development Workflow | PASS | Feature branch `001-skill-packaging-cli`. CI gates: lint, test, security scan. |

No violations. No complexity tracking entries needed.

## Project Structure

### Documentation (this feature)

```text
specs/001-skill-packaging-cli/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI contract)
│   └── cli-contract.md
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
cmd/
└── psk/
    └── main.go              # Entrypoint, subcommand dispatch

internal/
├── cli/
│   ├── build.go             # psk build command
│   ├── list.go              # psk list command
│   ├── validate.go          # psk validate command
│   └── root.go              # Top-level arg parsing, help, version
├── skill/
│   ├── parse.go             # SKILL.md frontmatter parser
│   └── validate.go          # Validation rules (agentskills.io)
├── store/
│   ├── store.go             # Store operations (add, list, exists)
│   └── manifest.go          # manifest.json read/write
└── exitcode/
    └── codes.go             # Exit code constants (0-4)

tests/
├── integration/
│   ├── build_test.go        # psk build integration tests
│   ├── list_test.go         # psk list integration tests
│   ├── validate_test.go     # psk validate integration tests
│   └── testdata/
│       ├── valid-skill/     # Valid agentskills.io skill fixture
│       │   └── SKILL.md
│       ├── missing-author/  # Fixture missing metadata.author
│       │   └── SKILL.md
│       ├── missing-version/ # Fixture missing metadata.version
│       │   └── SKILL.md
│       ├── bad-name/        # Fixture with naming rule violation
│       │   └── SKILL.md
│       └── malformed-yaml/  # Fixture with broken YAML
│           └── SKILL.md
└── unit/
    ├── parse_test.go        # Frontmatter parsing unit tests
    └── validate_test.go     # Validation rule unit tests

go.mod
go.sum
Makefile
```

**Structure Decision**: Single project layout. `cmd/psk/` for the
binary entrypoint, `internal/` for unexported packages (prevents
external import). Test fixtures in `tests/integration/testdata/`
as real skill directories — integration tests invoke the compiled
binary against these real fixtures and real temp-dir stores.

## Complexity Tracking

> No violations. No entries needed.
