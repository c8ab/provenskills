# Data Model: Skill Packaging CLI

**Feature**: 001-skill-packaging-cli
**Date**: 2026-02-11

## Entities

### SkillFrontmatter

Parsed from `SKILL.md` YAML frontmatter. Represents the
agentskills.io skill definition as read from the source
directory.

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| name | string | Yes | 1-64 chars, lowercase alphanumeric + hyphens, no leading/trailing/consecutive hyphens, must match parent directory name |
| description | string | Yes | 1-1024 chars, non-empty |
| license | string | No | Free-form |
| compatibility | string | No | 1-500 chars |
| metadata | map | Partial | `version` and `author` keys required |
| metadata.version | string | Yes | Semver or major.minor format |
| metadata.author | string | Yes | Non-empty, free-form |
| allowed-tools | string | No | Space-delimited list |

### Manifest (manifest.json)

Stored per artifact version. Combines skill metadata with
packaging provenance.

| Field | Type | Required | Constraints |
|-------|------|----------|-------------|
| manifestVersion | integer | Yes | Must be `1` (current schema) |
| name | string | Yes | Same constraints as SkillFrontmatter.name |
| version | string | Yes | Normalized semver (e.g., "1.0.0") |
| description | string | Yes | From SKILL.md frontmatter |
| author | string | Yes | From SKILL.md metadata.author |
| maintainer | string | Yes | From `--maintainer` flag |
| buildTimestamp | string | Yes | ISO 8601 UTC (e.g., "2026-02-11T14:30:00Z") |
| contents | object | No | File listing for integrity |
| contents.skillFile | string | No | Always "SKILL.md" |
| contents.scripts | []string | No | Relative paths of scripts/ files |
| contents.references | []string | No | Relative paths of references/ files |
| contents.assets | []string | No | Relative paths of assets/ files |
| sourceHash | string | No | SHA-256 hex digest of source directory tree |

### StoreLayout

The local skill store is a directory hierarchy. Not a data
structure in code — it is the filesystem itself.

```
{PSK_STORE}/
  .store-version              # Contains "1" (integer)
  {name}/                     # One directory per skill
    {version}/                # One directory per version (normalized semver)
      manifest.json           # Artifact metadata
      SKILL.md                # Copied from source
      scripts/                # Copied from source (if present)
      references/             # Copied from source (if present)
      assets/                 # Copied from source (if present)
```

**Identity & Uniqueness**: A skill artifact is uniquely
identified by the tuple `(name, version)`. The store enforces
this via directory structure — attempting to create a directory
that already exists is a conflict (exit code 3).

**Lifecycle**: Artifacts are immutable once stored. `--force`
on `psk build` replaces the entire version directory atomically.

## Relationships

```
SkillFrontmatter  --(parsed from)-->  SKILL.md file
Manifest          --(derived from)-->  SkillFrontmatter + --maintainer flag
StoreLayout       --(contains)-->      Manifest + skill files
```

## Validation Rules

1. `name` MUST match regex `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`
   and be 1-64 characters.
2. `name` MUST match the parent directory name of the source
   skill directory.
3. `description` MUST be non-empty and <= 1024 characters.
4. `metadata.version` MUST be valid semver or major.minor.
   major.minor is normalized to major.minor.0.
5. `metadata.author` MUST be non-empty.
6. `--maintainer` MUST be non-empty (provided at build time).
7. No `name` + normalized `version` collision in the store
   unless `--force` is specified.
