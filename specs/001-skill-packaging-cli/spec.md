# Feature Specification: Skill Packaging CLI

**Feature Branch**: `001-skill-packaging-cli`
**Created**: 2026-02-11
**Status**: Draft
**Input**: User description: "Build a CLI tool called 'psk' that packages and manages AI Agent Skills as Proven Skill Artifacts."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build a Skill from a Local Directory (Priority: P1)

A skill author has a local directory containing a valid
agentskills.io skill (with a `SKILL.md` file and optional
`scripts/`, `references/`, `assets/` directories). They run
`psk build <path>` to package it into a Proven Skill Artifact.
The tool reads the skill's `SKILL.md` frontmatter, validates
that all mandatory metadata is present (`name`, `description`,
and `version` and `author` from the `metadata` field), prompts
for or accepts a `--maintainer` flag identifying who is
packaging the skill, and produces a stored artifact in the
local skill store. The user sees a success message with the
skill name, version, and maintainer on stdout. If validation
fails, the tool prints the specific missing fields to stderr
and exits with a non-zero code.

**Why this priority**: This is the foundational capability.
Without building skills, no other feature has meaning. It
establishes the packaging contract, validation rules, and
local store format.

**Independent Test**: Can be fully tested by creating a
minimal skill directory with a valid `SKILL.md`, running
`psk build`, and verifying the artifact appears in the
local store. Delivers the core value proposition of
packaging skills as verifiable artifacts.

**Acceptance Scenarios**:

1. **Given** a directory containing a valid `SKILL.md` with
   all required frontmatter fields (`name`, `description`,
   `metadata.version`, `metadata.author`),
   **When** the user runs `psk build ./my-skill --maintainer "Jane Doe <jane@example.com>"`,
   **Then** the tool produces a Proven Skill Artifact in the
   local store, prints the skill name, version, author, and
   maintainer to stdout, and exits with code 0.

2. **Given** a directory containing a `SKILL.md` that is
   missing the required `metadata.author` field,
   **When** the user runs `psk build ./my-skill --maintainer "Jane Doe <jane@example.com>"`,
   **Then** the tool prints an error to stderr listing
   `metadata.author` as missing, does NOT create an artifact,
   and exits with a non-zero code.

3. **Given** a directory containing a `SKILL.md` that is
   missing the required `metadata.version` field,
   **When** the user runs `psk build ./my-skill --maintainer "Jane Doe <jane@example.com>"`,
   **Then** the tool prints an error to stderr listing
   `metadata.version` as missing, does NOT create an artifact,
   and exits with a non-zero code.

4. **Given** a directory that does not contain a `SKILL.md`
   file,
   **When** the user runs `psk build ./not-a-skill --maintainer "Jane Doe <jane@example.com>"`,
   **Then** the tool prints an error to stderr indicating
   no `SKILL.md` found, and exits with a non-zero code.

5. **Given** a valid skill directory,
   **When** the user runs `psk build ./my-skill` without
   providing `--maintainer`,
   **Then** the tool prints an error to stderr indicating
   the maintainer is required, and exits with a non-zero code.

6. **Given** a valid skill directory where the directory name
   does not match the `name` field in `SKILL.md` frontmatter,
   **When** the user runs `psk build ./wrong-name --maintainer "Jane Doe <jane@example.com>"`,
   **Then** the tool prints an error to stderr indicating the
   name mismatch (per agentskills.io spec: name must match
   parent directory), and exits with a non-zero code.

---

### User Story 2 - List Locally Stored Skills (Priority: P2)

A user wants to see all Proven Skill Artifacts stored locally.
They run `psk list` and see a table of all skills with their
name, version, author, and maintainer. If the store is empty,
they see a message indicating no skills are stored. The command
supports `--json` for machine-readable output.

**Why this priority**: Listing is the natural complement to
building. Users need to verify what they have stored. This
also validates the store format is readable and complete.

**Independent Test**: Can be tested by building one or more
skills, then running `psk list` and verifying the output
contains the expected skill metadata. Also testable with
an empty store.

**Acceptance Scenarios**:

1. **Given** the local store contains two packaged skills,
   **When** the user runs `psk list`,
   **Then** stdout displays a human-readable table showing
   each skill's name, version, author, and maintainer,
   and exits with code 0.

2. **Given** the local store is empty,
   **When** the user runs `psk list`,
   **Then** stdout displays a message indicating no skills
   are stored, and exits with code 0.

3. **Given** the local store contains skills,
   **When** the user runs `psk list --json`,
   **Then** stdout displays a JSON array of skill objects
   each containing name, version, description, author, and
   maintainer fields, and exits with code 0.

---

### User Story 3 - Validate a Skill Directory Without Building (Priority: P3)

A skill author wants to check if their skill directory meets
all requirements before building. They run `psk validate <path>`
to get a pass/fail report without producing an artifact.

**Why this priority**: Validation is useful during skill
development iteration but is not essential for the core
build-and-store workflow. The build command already validates,
so this is a convenience feature.

**Independent Test**: Can be tested by pointing at valid and
invalid skill directories and verifying the validation output
and exit codes without checking the store.

**Acceptance Scenarios**:

1. **Given** a directory with a fully valid `SKILL.md`,
   **When** the user runs `psk validate ./my-skill`,
   **Then** stdout displays a validation success message
   listing all checked fields, and exits with code 0.

2. **Given** a directory with an invalid `SKILL.md` (missing
   required fields),
   **When** the user runs `psk validate ./my-skill`,
   **Then** stderr lists each validation error with the
   specific field name and reason, and exits with a non-zero
   code.

### Edge Cases

- What happens when the skill directory path does not exist?
  The tool MUST print "directory not found" to stderr and
  exit with a non-zero code.
- What happens when `SKILL.md` contains invalid YAML
  frontmatter (malformed syntax)? The tool MUST print a
  parse error to stderr and exit with a non-zero code.
- What happens when a skill with the same name and version
  already exists in the local store? The tool MUST reject
  the build by default, printing a conflict message to stderr.
  A `--force` flag MAY be provided to overwrite.
- What happens when the `name` field in `SKILL.md` contains
  invalid characters (uppercase, consecutive hyphens, etc.)?
  The tool MUST validate against agentskills.io naming rules
  and report specific violations.
- What happens when the skill directory contains files outside
  the expected structure (no `SKILL.md` at root)? The tool
  MUST report "SKILL.md not found at directory root."
- What happens when `metadata.version` is not a valid semver
  string? The tool MUST report the version format requirement.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tool MUST be a single CLI binary named `psk`.
- **FR-002**: `psk build <path> --maintainer <identity>` MUST
  read a skill directory, validate it, and store a Proven Skill
  Artifact in the local store.
- **FR-003**: The tool MUST validate the following fields exist
  and are non-empty before accepting a skill: `name`,
  `description` (from SKILL.md frontmatter), `metadata.version`,
  and `metadata.author` (from SKILL.md metadata section).
- **FR-004**: The tool MUST require a `--maintainer` flag on
  `psk build` identifying who packaged the skill. The maintainer
  is distinct from the skill author.
- **FR-005**: The tool MUST validate `name` conforms to
  agentskills.io naming rules (lowercase alphanumeric + hyphens,
  1-64 chars, no leading/trailing/consecutive hyphens, must
  match parent directory name).
- **FR-006**: `psk list` MUST display all locally stored skills
  with name, version, author, and maintainer.
- **FR-007**: `psk list --json` MUST output skill data as a
  JSON array to stdout.
- **FR-008**: All errors MUST be written to stderr. All
  successful output MUST be written to stdout.
- **FR-009**: Exit codes MUST follow POSIX conventions: 0 for
  success, non-zero for failure. Fixed exit code taxonomy:
  1=general/unexpected error, 2=validation failure (missing
  or invalid metadata fields, naming rule violations),
  3=store conflict (duplicate name+version already exists),
  4=IO/filesystem error (path not found, permission denied,
  read/write failure).
- **FR-010**: The local store MUST only track skills. It MUST
  have no concept of agents, runtimes, or deployment targets.
- **FR-011**: `psk validate <path>` MUST check all validation
  rules without producing an artifact.
- **FR-012**: The tool MUST reject building a skill if an
  artifact with the same name and version already exists in
  the store, unless `--force` is provided.

### Key Entities

- **Skill Source**: A local directory following the agentskills.io
  format, containing at minimum a `SKILL.md` file with YAML
  frontmatter defining `name`, `description`, and `metadata`
  fields (`version`, `author`).
- **Proven Skill Artifact**: A packaged skill stored in the local
  store, consisting of the original skill content plus packaging
  metadata (maintainer identity, build timestamp).
- **Local Skill Store**: A local directory structure where Proven
  Skill Artifacts are persisted. Organized by skill name and
  version. Multiple versions of the same skill coexist
  simultaneously. Contains no agent or runtime information.
- **Maintainer**: The identity of the person or organization that
  packaged the skill. Distinct from the skill's author. Provided
  via `--maintainer` flag at build time.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can build a valid skill from a local
  directory and see it in the store in under 5 seconds for a
  typical skill (< 1MB of content).
- **SC-002**: 100% of validation errors produce specific,
  actionable error messages naming the exact field and rule
  that failed.
- **SC-003**: The `psk list` output for a store with 100 skills
  completes in under 2 seconds.
- **SC-004**: A user unfamiliar with the tool can successfully
  build and list a skill within 5 minutes using only `psk --help`
  and `psk build --help` output.
- **SC-005**: The tool produces a single binary with zero
  runtime dependencies (no interpreters, no shared libraries
  beyond the OS standard).

## Clarifications

### Session 2026-02-11

- Q: Can multiple versions of the same skill coexist in the store? → A: Yes, multiple versions coexist (e.g., my-skill@1.0.0 and my-skill@2.0.0 both stored).
- Q: What are the distinct exit codes for error categories? → A: Small fixed set: 1=general error, 2=validation failure, 3=store conflict, 4=IO/filesystem error.

## Assumptions

- The `metadata.version` field follows semantic versioning
  (e.g., "1.0.0", "0.2.1"). The agentskills.io spec shows
  `version: "1.0"` in examples, so the tool accepts both
  semver and major.minor formats.
- The `metadata.author` field is a free-form string (e.g.,
  "example-org", "Jane Doe <jane@example.com>").
- The `--maintainer` flag accepts a free-form string identity
  (e.g., "Jane Doe <jane@example.com>", "acme-corp").
- The local store location defaults to `~/.psk/store/` unless
  overridden by a `PSK_STORE` environment variable.
- The Proven Skill Artifact format and local store structure
  are implementation details to be determined during planning.
