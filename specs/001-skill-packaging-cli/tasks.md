# Tasks: Skill Packaging CLI

**Input**: Design documents from `/specs/001-skill-packaging-cli/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/cli-contract.md

**Tests**: TDD is MANDATORY per constitution (Principle II). Tests MUST be written FIRST and FAIL before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `cmd/psk/`, `internal/`, `tests/` at repository root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Go module initialization, project structure, build tooling

- [X] T001 Initialize Go module with `go mod init` and add `gopkg.in/yaml.v3` dependency in go.mod
- [X] T002 Create directory structure: cmd/psk/, internal/cli/, internal/skill/, internal/store/, internal/exitcode/, tests/integration/, tests/integration/testdata/, tests/unit/
- [X] T003 Create Makefile with targets: build (CGO_ENABLED=0 go build -ldflags="-s -w" -o psk ./cmd/psk/), test (go test ./...), lint (go vet + golangci-lint), clean
- [X] T004 [P] Configure golangci-lint with .golangci.yml at repository root

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

### Tests for Foundational

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T005 [P] Write unit tests for YAML frontmatter parsing in tests/unit/parse_test.go — test cases: valid frontmatter extraction, missing opening `---`, missing closing `---`, empty frontmatter, additional `---` in body content
- [X] T006 [P] Write unit tests for validation rules in tests/unit/validate_test.go — test cases: valid name, uppercase name, consecutive hyphens, leading/trailing hyphens, name >64 chars, name-directory mismatch, empty description, description >1024 chars, valid semver, major.minor normalization to major.minor.0, invalid version "abc", empty author

### Implementation for Foundational

- [X] T007 Define exit code constants in internal/exitcode/codes.go — Success=0, ErrGeneral=1, ErrValidation=2, ErrConflict=3, ErrIO=4
- [X] T008 Implement SKILL.md frontmatter parser in internal/skill/parse.go — split on `---` delimiters, yaml.Unmarshal into SkillFrontmatter struct with fields: Name, Description, License, Compatibility, Metadata (map with version, author), AllowedTools
- [X] T009 Implement validation rules in internal/skill/validate.go — name regex `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`, length 1-64, directory name match, description non-empty and <=1024, version semver/major.minor with normalization, author non-empty
- [X] T010 [P] Implement manifest.json struct and read/write functions in internal/store/manifest.go — Manifest struct with ManifestVersion, Name, Version, Description, Author, Maintainer, BuildTimestamp, Contents, SourceHash; MarshalJSON and UnmarshalJSON via encoding/json
- [X] T011 Create test fixture: valid skill directory in tests/integration/testdata/valid-skill/SKILL.md with name=valid-skill, description, metadata.version="1.0.0", metadata.author="test-author"
- [X] T012 [P] Create test fixture: missing-author skill in tests/integration/testdata/missing-author/SKILL.md with name and version but no metadata.author
- [X] T013 [P] Create test fixture: missing-version skill in tests/integration/testdata/missing-version/SKILL.md with name and author but no metadata.version
- [X] T014 [P] Create test fixture: bad-name skill in tests/integration/testdata/bad-name/SKILL.md with name="Bad-Name" (uppercase violation)
- [X] T015 [P] Create test fixture: malformed-yaml skill in tests/integration/testdata/malformed-yaml/SKILL.md with broken YAML syntax

**Checkpoint**: Foundation ready — parser, validator, manifest, exit codes, and test fixtures all in place. User story implementation can now begin.

---

## Phase 3: User Story 1 - Build a Skill from a Local Directory (Priority: P1) MVP

**Goal**: `psk build <path> --maintainer <identity>` reads a skill directory, validates it, and stores a Proven Skill Artifact in the local store.

**Independent Test**: Build a valid skill, verify artifact appears in store with correct manifest.json. Build invalid skills, verify specific error messages on stderr and correct exit codes.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T016 [US1] Write integration test for successful build in tests/integration/build_test.go — build psk binary, invoke `psk build testdata/valid-skill --maintainer "Test <test@example.com>"` with PSK_STORE set to t.TempDir(), assert exit code 0, assert stdout contains "Built skill: valid-skill@1.0.0", assert manifest.json exists in store at valid-skill/1.0.0/manifest.json, assert SKILL.md copied to store
- [X] T017 [P] [US1] Write integration test for missing author in tests/integration/build_test.go — invoke build against testdata/missing-author, assert exit code 2, assert stderr contains "metadata.author"
- [X] T018 [P] [US1] Write integration test for missing version in tests/integration/build_test.go — invoke build against testdata/missing-version, assert exit code 2, assert stderr contains "metadata.version"
- [X] T019 [P] [US1] Write integration test for missing SKILL.md in tests/integration/build_test.go — invoke build against nonexistent path, assert exit code 4, assert stderr contains "not found"
- [X] T020 [P] [US1] Write integration test for missing --maintainer in tests/integration/build_test.go — invoke build without --maintainer, assert exit code 2, assert stderr contains "--maintainer"
- [X] T021 [P] [US1] Write integration test for name-directory mismatch in tests/integration/build_test.go — create temp skill dir where SKILL.md name field differs from directory name, assert exit code 2
- [X] T022 [US1] Write integration test for duplicate build conflict in tests/integration/build_test.go — build same skill twice, assert second build exits with code 3, assert stderr contains "already exists"
- [X] T023 [US1] Write integration test for --force overwrite in tests/integration/build_test.go — build same skill twice with --force on second, assert exit code 0
- [X] T024 [US1] Write integration test for --json output in tests/integration/build_test.go — build with --json flag, assert stdout is valid JSON with name, version, author, maintainer, path fields

### Implementation for User Story 1

- [X] T025 [US1] Implement store operations in internal/store/store.go — New(storePath) constructor resolving PSK_STORE env or default ~/.psk/store/, Init() creates store dir and .store-version file, Exists(name, version) checks directory existence, Add(name, version, sourceDir, manifest) copies skill files atomically via tmp dir + os.Rename
- [X] T026 [US1] Implement build command in internal/cli/build.go — parse --maintainer, --force, --json flags; read and validate skill dir; create manifest; call store.Add; format output per cli-contract.md
- [X] T027 [US1] Implement root CLI dispatcher in internal/cli/root.go — parse os.Args[1] for subcommand (build, list, validate, --help, --version), dispatch to appropriate handler, handle unknown commands with usage error
- [X] T028 [US1] Implement main.go entrypoint in cmd/psk/main.go — call cli.Run(os.Args), os.Exit with returned exit code

**Checkpoint**: `psk build` fully functional. Can package valid skills, reject invalid ones with specific errors, handle conflicts and --force. MVP complete.

---

## Phase 4: User Story 2 - List Locally Stored Skills (Priority: P2)

**Goal**: `psk list` displays all Proven Skill Artifacts in the local store with name, version, author, and maintainer.

**Independent Test**: Build several skills via `psk build`, then run `psk list` and verify output contains all expected entries. Also verify empty store behavior.

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T029 [US2] Write integration test for list with skills in tests/integration/list_test.go — build two skills into temp store, invoke `psk list` with PSK_STORE, assert exit code 0, assert stdout contains both skill names, versions, authors, maintainers in tabular format
- [X] T030 [P] [US2] Write integration test for list with empty store in tests/integration/list_test.go — invoke `psk list` against empty temp store, assert exit code 0, assert stdout contains "No skills found"
- [X] T031 [P] [US2] Write integration test for list --json in tests/integration/list_test.go — build skills, invoke `psk list --json`, assert stdout is valid JSON array, assert each element has name, version, description, author, maintainer

### Implementation for User Story 2

- [X] T032 [US2] Implement store.List() in internal/store/store.go — walk store directory, read manifest.json from each {name}/{version}/ directory, return slice of Manifest structs sorted by name then version
- [X] T033 [US2] Implement list command in internal/cli/list.go — parse --json flag, call store.List(), format as aligned table (human) or JSON array (--json) per cli-contract.md, handle empty store message
- [X] T034 [US2] Register list subcommand in internal/cli/root.go dispatcher

**Checkpoint**: `psk list` fully functional. Users can build and list skills end-to-end.

---

## Phase 5: User Story 3 - Validate a Skill Directory Without Building (Priority: P3)

**Goal**: `psk validate <path>` checks all validation rules and reports pass/fail without storing anything.

**Independent Test**: Point at valid and invalid skill directories, verify validation output lists each checked field and correct exit codes.

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T035 [US3] Write integration test for validate with valid skill in tests/integration/validate_test.go — invoke `psk validate testdata/valid-skill`, assert exit code 0, assert stdout contains "Validation passed" and lists checked fields (name, description, version, author, dir match)
- [X] T036 [P] [US3] Write integration test for validate with invalid skill in tests/integration/validate_test.go — invoke `psk validate testdata/bad-name`, assert exit code 2, assert stderr lists specific validation errors
- [X] T037 [P] [US3] Write integration test for validate with malformed YAML in tests/integration/validate_test.go — invoke `psk validate testdata/malformed-yaml`, assert exit code 2, assert stderr contains parse error

### Implementation for User Story 3

- [X] T038 [US3] Implement validate command in internal/cli/validate.go — parse path arg and --json flag, call skill.Parse then skill.Validate, format success output listing each checked field per cli-contract.md, format errors listing each violation
- [X] T039 [US3] Register validate subcommand in internal/cli/root.go dispatcher

**Checkpoint**: All three user stories independently functional. Full CLI contract implemented.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T040 Implement --help output in internal/cli/root.go matching cli-contract.md help text exactly
- [X] T041 [P] Implement --version output in internal/cli/root.go printing "psk version 0.1.0"
- [X] T042 [P] Add input sanitization for path traversal attacks in internal/cli/build.go and internal/cli/validate.go — reject paths containing ".." segments
- [X] T043 Run quickstart.md validation — execute each command from specs/001-skill-packaging-cli/quickstart.md against the built binary and verify expected outputs match
- [X] T044 Run full test suite (`go test ./...`) and verify all tests pass
- [X] T045 [P] Run linter (`golangci-lint run`) and fix any findings

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational phase completion
- **User Story 2 (Phase 4)**: Depends on Foundational phase. Depends on US1 store.Add for building test fixtures (integration tests build skills then list them)
- **User Story 3 (Phase 5)**: Depends on Foundational phase only (validate does not use the store)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) — no dependencies on other stories
- **User Story 2 (P2)**: Integration tests depend on `psk build` being functional (US1 must be complete first)
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) — independent of US1 and US2. Can run in parallel with US1 if desired.

### Within Each User Story

- Tests MUST be written and FAIL before implementation (Constitution Principle II)
- Parser/validator before store operations
- Store operations before CLI commands
- CLI commands before integration
- Story complete before moving to next priority

### Parallel Opportunities

- T003 and T004 (setup tasks) can run in parallel
- T005 and T006 (foundational tests) can run in parallel
- T010, T012, T013, T014, T015 (manifest + test fixtures) can run in parallel
- T017, T018, T019, T020, T021 (US1 error-path tests) can run in parallel
- T030 and T031 (US2 empty/json tests) can run in parallel
- T036 and T037 (US3 error-path tests) can run in parallel
- T041, T042, T045 (polish tasks) can run in parallel

---

## Parallel Example: User Story 1 Tests

```bash
# Launch all error-path integration tests together (different test cases, same file):
Task: "Integration test for missing author in tests/integration/build_test.go"
Task: "Integration test for missing version in tests/integration/build_test.go"
Task: "Integration test for missing SKILL.md in tests/integration/build_test.go"
Task: "Integration test for missing --maintainer in tests/integration/build_test.go"
Task: "Integration test for name-directory mismatch in tests/integration/build_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run `psk build` against real skill directories, verify artifacts in store
5. MVP deliverable: users can build and inspect skill artifacts

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → MVP!
3. Add User Story 2 → Test independently → Build + List workflow
4. Add User Story 3 → Test independently → Full validation story
5. Polish → Help text, version, security hardening, quickstart validation

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (RED phase of TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All integration tests invoke the real compiled `psk` binary — no in-process testing
