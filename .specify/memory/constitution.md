<!--
  Sync Impact Report
  ==================
  Version change: 0.0.0 (template) → 1.0.0
  Modified principles:
    - [PRINCIPLE_1_NAME] → I. CLI-First Development
    - [PRINCIPLE_2_NAME] → II. Test-Driven Development (NON-NEGOTIABLE)
    - [PRINCIPLE_3_NAME] → III. Integration Testing over Mocks
    - [PRINCIPLE_4_NAME] → IV. Security-First Posture
    - [PRINCIPLE_5_NAME] → V. Minimal Dependency Footprint
  Added sections:
    - Code Quality Standards (Section 2)
    - Development Workflow (Section 3)
    - Governance (filled from template placeholder)
  Removed sections: none
  Templates requiring updates:
    - .specify/templates/plan-template.md ✅ no update needed
      (Constitution Check is dynamically resolved at plan time)
    - .specify/templates/spec-template.md ✅ no update needed
      (User scenarios and requirements structure is compatible)
    - .specify/templates/tasks-template.md ✅ no update needed
      (TDD task ordering already enforced; security tasks in
       Phase N polish section)
    - .opencode/command/*.md ✅ no update needed
      (No outdated agent-specific references found)
  Deferred items:
    - RATIFICATION_DATE set to today (first adoption)
-->

# ProvenSkills Constitution

## Core Principles

### I. CLI-First Development

Every feature MUST be exposed as a CLI command before any other
interface is considered.

- The CLI binary is the primary artifact. All functionality MUST
  be invocable via command-line arguments, flags, and stdin.
- Text in/out protocol: arguments and stdin for input, stdout for
  results, stderr for diagnostics. Exit codes MUST follow POSIX
  conventions (0 = success, non-zero = specific error category).
- Support both human-readable (default) and machine-readable
  (JSON via `--json` or equivalent flag) output formats.
- No GUI, web UI, or API server layer MUST be introduced unless
  the CLI equivalent exists first and the addition is justified
  in a plan document.

### II. Test-Driven Development (NON-NEGOTIABLE)

Tests MUST be written before implementation code. No exceptions.

- Red-Green-Refactor cycle is strictly enforced:
  1. Write a failing test that defines the expected behavior.
  2. Confirm the test fails for the right reason.
  3. Write the minimum implementation to make the test pass.
  4. Refactor while keeping tests green.
- Test files MUST exist in the repository before or in the same
  commit as the implementation they cover.
- Every pull request MUST include tests for new or changed
  behavior. PRs without tests for functional changes MUST be
  rejected.
- CLI output format and exit codes are the contract. Tests MUST
  assert on stdout content, stderr content, and exit codes as
  the primary verification mechanism.

### III. Integration Testing over Mocks

Integration tests against real artifacts and real registries
MUST be preferred over mocks and stubs.

- Tests MUST exercise real file systems, real container
  registries, real binaries, and real network calls where the
  system under test interacts with them.
- Mocks are permitted ONLY when the real dependency is
  non-deterministic, prohibitively slow (>30s), or requires
  paid third-party credentials not available in CI. Every mock
  MUST include a comment justifying its existence.
- Contract tests MUST validate the interface between components
  using actual CLI invocations, not in-process function calls.
- CI pipelines MUST run the full integration suite. A passing
  unit-only run does NOT qualify as a green build.

### IV. Security-First Posture

Security is a design constraint, not a post-hoc review step.

- All inputs MUST be validated and sanitized at the boundary
  (CLI argument parsing layer). Reject unknown or malformed
  input before processing begins.
- Dependencies MUST be audited before adoption. Every new
  dependency MUST include a justification comment in the
  dependency manifest explaining why it is necessary and what
  alternatives were evaluated.
- The binary MUST be compiled with hardening flags appropriate
  to the target platform (e.g., PIE, stack canaries, stripped
  symbols for release builds).
- Secrets MUST NOT appear in CLI output, logs, or error
  messages. Credential handling MUST use environment variables
  or secure credential stores, never command-line arguments.
- Vulnerability scanning MUST run in CI on every build. Known
  vulnerabilities with available patches MUST be resolved
  before release.

### V. Minimal Dependency Footprint

Every dependency MUST be justified. The default answer to
"should we add this library?" is NO.

- The project MUST produce a small, statically-linked (where
  the toolchain permits), auditable binary. Binary size MUST
  be tracked and regressions investigated.
- Before adding a dependency, the author MUST document:
  (a) what it provides, (b) the cost of implementing it
  in-house, (c) its transitive dependency count, (d) its
  maintenance status and security track record.
- Vendored or copied code (with license compliance) is
  preferred over adding a dependency for small, stable
  functionality.
- Dependency updates MUST be reviewed for changelog impact,
  not blindly merged. Automated update PRs MUST still pass
  the full integration suite.

## Code Quality Standards

- Every source file MUST have a clear, single responsibility.
  Files exceeding 500 lines MUST be split with justification.
- Public interfaces MUST be documented with usage examples.
- Linting and formatting MUST be enforced in CI. Code that
  does not pass lint MUST NOT be merged.
- Complexity MUST be justified. If a simpler approach exists,
  the simpler approach MUST be chosen unless a plan document
  explicitly justifies the complexity with measurable criteria.

## Development Workflow

- All work MUST happen on feature branches. Direct commits to
  the main branch are prohibited.
- Every feature branch MUST pass CI (lint, unit tests,
  integration tests, security scan) before merge.
- Code review is mandatory. The reviewer MUST verify
  constitution compliance as part of the review checklist.
- Commits MUST be atomic and focused. Each commit MUST compile
  and pass tests independently.

## Governance

This constitution is the highest-authority document for the
ProvenSkills project. All development practices, code reviews,
and architectural decisions MUST comply with its principles.

- **Amendment procedure**: Any change to this constitution
  MUST be proposed as a pull request with a rationale section.
  The PR MUST include an updated Sync Impact Report and version
  bump. All active contributors MUST be notified.
- **Versioning policy**: The constitution follows semantic
  versioning. MAJOR for principle removals or incompatible
  redefinitions, MINOR for new principles or material
  expansions, PATCH for clarifications and typo fixes.
- **Compliance review**: Every PR review MUST include a
  constitution compliance check. Violations MUST be resolved
  before merge or explicitly granted an exception with a
  documented justification and expiration date.
- **Conflict resolution**: If a practice conflicts with this
  constitution, the constitution wins. Update the practice or
  amend the constitution — do not ignore the conflict.

**Version**: 1.0.0 | **Ratified**: 2026-02-11 | **Last Amended**: 2026-02-11
