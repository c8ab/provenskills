# Research: Skill Packaging CLI

**Feature**: 001-skill-packaging-cli
**Date**: 2026-02-11

## Decision 1: Implementation Language

**Decision**: Go 1.23+

**Rationale**: Go produces static binaries trivially
(`CGO_ENABLED=0 go build`), has stdlib JSON (`encoding/json`),
stdlib testing with `t.TempDir()` for filesystem tests, and
stdlib `os/exec` for integration tests that invoke the real
binary. The sole external dependency is `gopkg.in/yaml.v3`
(Google-maintained, zero transitive deps). This directly
satisfies Constitution Principle V (Minimal Dependency
Footprint) — the entire `go.mod` has one `require` line.

**Alternatives considered**:
- **Rust**: Produces smaller binaries (1-3 MB vs 5-8 MB) but
  requires `serde` + `serde_yaml` (15-25 transitive deps),
  slower build times (15-45s vs 1-3s), and more complex
  cross-compilation setup. The dependency footprint directly
  violates our minimal dependency principle for marginal
  binary size gain.
- **Python/Node**: Ruled out immediately — cannot produce a
  single static binary without bundlers. Violates SC-005.

## Decision 2: CLI Argument Parsing

**Decision**: Go `flag` stdlib package with manual subcommand
dispatch.

**Rationale**: The tool has 3 subcommands (`build`, `list`,
`validate`) with a handful of flags each. The stdlib `flag`
package handles flag parsing. Subcommand dispatch is a simple
switch on `os.Args[1]`. Adding `cobra` or `urfave/cli` for
this scope would violate the minimal dependency principle.

**Alternatives considered**:
- **cobra**: Popular but pulls in `pflag`, `viper`, and
  transitive deps. Massive overkill for 3 subcommands.
- **urfave/cli**: Lighter than cobra but still an unnecessary
  dependency for this scope.

## Decision 3: Artifact Store Format

**Decision**: Directory-per-artifact layout with `manifest.json`
per version.

**Rationale**: The filesystem IS the index — no derived state
to corrupt. Layout: `{store}/{name}/{version}/manifest.json`
plus the original skill files. Human-auditable with `ls` and
`cat`. Conflict detection is a single `stat()` call. Atomic
writes via rename pattern: write to `{version}.tmp.{pid}/`,
then `os.Rename()` to `{version}/`. A `.store-version` file
at root enables future layout migrations.

**Layout**:
```
~/.psk/store/
  .store-version          # "1"
  {name}/
    {version}/
      manifest.json       # Artifact metadata
      SKILL.md            # Original skill file
      scripts/            # Optional, copied from source
      references/         # Optional
      assets/             # Optional
```

**Alternatives considered**:
- **Tar/archive-based**: Fails the human-auditability
  requirement — must extract to inspect. Adds tar dependency.
- **Flat dir + index.json**: Index file is a coordination
  point that can drift from filesystem reality. Crash between
  dir-write and index-update = inconsistent state. Two-phase
  write is not atomic.

## Decision 4: Manifest Schema

**Decision**: JSON manifest with `manifestVersion: 1` for
forward compatibility.

**Rationale**: JSON is readable, Go stdlib can marshal/unmarshal
it, and a version integer allows future schema evolution without
breaking existing stores.

**Required fields**: `manifestVersion`, `name`, `version`,
`description`, `author`, `maintainer`, `buildTimestamp`.
**Optional fields**: `contents` (file listing), `sourceHash`
(SHA-256 of source tree).

See `data-model.md` for the full schema.

## Decision 5: YAML Frontmatter Parsing Strategy

**Decision**: Split `SKILL.md` on `---` delimiters, unmarshal
the YAML block between the first and second `---` into a Go
struct.

**Rationale**: The agentskills.io spec requires YAML frontmatter
delimited by `---`. This is a well-established convention
(Jekyll, Hugo, etc.). Simple string splitting + `yaml.Unmarshal`
is sufficient — no Markdown parsing library needed.

**Edge cases handled**:
- Missing opening `---`: validation error
- Missing closing `---`: validation error
- Empty frontmatter: validation error (required fields missing)
- Additional `---` in body content: only first two delimiters
  used for frontmatter extraction

## Decision 6: Version Format Acceptance

**Decision**: Accept both semver (1.0.0) and major.minor (1.0)
formats, normalizing major.minor to major.minor.0 for storage.

**Rationale**: The agentskills.io spec shows `version: "1.0"`
in examples, but semver is the industry standard. Accepting both
avoids rejecting valid skills while storing a consistent format.
The store directory uses the normalized semver string.
