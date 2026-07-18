# Repository Go Toolchain Specification

## Purpose
Define the repository's single Go module layout and its required Go toolchain.

## Requirements

### Requirement: Repository contains one Go application module
The repository SHALL contain the root application Go module as its only Go module and SHALL NOT retain the `inf/` Pulumi program, Pulumi-only Taskfile commands, or CI cache references to infrastructure module metadata.

#### Scenario: Pulumi module is removed
- **WHEN** the repository layout is inspected after the change
- **THEN** `inf/`, its Go module metadata, its Pulumi project file, and its program source MUST be absent

#### Scenario: Stale infrastructure integrations are removed
- **WHEN** Taskfile, CI, and repository guidance are inspected after the change
- **THEN** they MUST NOT contain commands, cache entries, or module-layout guidance that requires `inf/`

### Requirement: Go 1.26.0 is used consistently
The root application module and every Go-based GitHub Actions workflow SHALL declare Go 1.26.0.

#### Scenario: Module toolchain is upgraded
- **WHEN** a developer inspects the root `go.mod`
- **THEN** its `go` directive MUST be `1.26.0`

#### Scenario: CI toolchain matches the module
- **WHEN** build, lint, and test workflows install Go
- **THEN** each workflow MUST request Go 1.26.0
