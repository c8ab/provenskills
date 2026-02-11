# CLI Contract: psk

**Feature**: 001-skill-packaging-cli
**Date**: 2026-02-11

This document defines the CLI interface contract. Tests MUST
assert against these exact behaviors (stdout content, stderr
content, exit codes).

## Exit Codes

| Code | Name | Meaning |
|------|------|---------|
| 0 | Success | Operation completed successfully |
| 1 | ErrGeneral | Unexpected/internal error |
| 2 | ErrValidation | Metadata validation failure |
| 3 | ErrConflict | Store conflict (name+version exists) |
| 4 | ErrIO | Filesystem/IO error |

## Commands

### `psk build <path> --maintainer <identity> [--force] [--json]`

Packages a skill directory into a Proven Skill Artifact.

**Arguments**:
- `<path>` (required): Path to skill source directory

**Flags**:
- `--maintainer <identity>` (required): Who is packaging
- `--force` (optional): Overwrite existing artifact
- `--json` (optional): Output result as JSON

**Success output** (exit 0):

Human-readable (default):
```
Built skill: code-review@2.0.0
  author:     example-org
  maintainer: Jane Doe <jane@example.com>
  stored:     ~/.psk/store/code-review/2.0.0/
```

JSON (`--json`):
```json
{
  "name": "code-review",
  "version": "2.0.0",
  "author": "example-org",
  "maintainer": "Jane Doe <jane@example.com>",
  "path": "/home/user/.psk/store/code-review/2.0.0/"
}
```

**Error output** (exit 2, validation failure, stderr):
```
error: validation failed for ./my-skill

  - metadata.author: required field is missing
  - metadata.version: required field is missing
```

**Error output** (exit 3, conflict, stderr):
```
error: skill code-review@2.0.0 already exists in store

Use --force to overwrite.
```

**Error output** (exit 4, IO error, stderr):
```
error: directory not found: ./nonexistent-path
```

**Error output** (exit 2, missing --maintainer, stderr):
```
error: --maintainer flag is required

Usage: psk build <path> --maintainer <identity>
```

### `psk list [--json]`

Lists all skills in the local store.

**Flags**:
- `--json` (optional): Output as JSON array

**Success output** (exit 0, skills present):

Human-readable (default):
```
NAME            VERSION  AUTHOR       MAINTAINER
code-review     2.0.0    example-org  Jane Doe <jane@example.com>
code-review     1.0.0    example-org  Jane Doe <jane@example.com>
terraform-plan  1.0.0    acme-corp    Bob Smith <bob@acme.com>
```

JSON (`--json`):
```json
[
  {
    "name": "code-review",
    "version": "2.0.0",
    "description": "Guides agents to perform code reviews",
    "author": "example-org",
    "maintainer": "Jane Doe <jane@example.com>"
  },
  {
    "name": "terraform-plan",
    "version": "1.0.0",
    "description": "Reviews Terraform plans for best practices",
    "author": "acme-corp",
    "maintainer": "Bob Smith <bob@acme.com>"
  }
]
```

**Success output** (exit 0, empty store):
```
No skills found in store.
```

JSON (`--json`, empty store):
```json
[]
```

### `psk validate <path> [--json]`

Validates a skill directory without storing it.

**Arguments**:
- `<path>` (required): Path to skill source directory

**Flags**:
- `--json` (optional): Output result as JSON

**Success output** (exit 0):

Human-readable:
```
Validation passed: ./my-skill

  name:        code-review (valid)
  description: present (145 chars)
  version:     2.0.0 (valid semver)
  author:      example-org (present)
  dir match:   code-review == code-review (ok)
```

**Error output** (exit 2, stderr):
```
error: validation failed for ./my-skill

  - name: contains uppercase characters (must be lowercase)
  - metadata.version: "abc" is not valid semver
```

### `psk --help`

```
psk - Package and manage Proven Skill Artifacts

Usage:
  psk <command> [flags]

Commands:
  build     Package a skill directory into an artifact
  list      List all skills in the local store
  validate  Validate a skill directory

Flags:
  --help      Show this help message
  --version   Show psk version

Environment:
  PSK_STORE   Override default store location (~/.psk/store/)
```

### `psk --version`

```
psk version 0.1.0
```
