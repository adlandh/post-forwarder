## ADDED Requirements

### Requirement: Echo v5 is the single HTTP runtime
The service SHALL use Echo v5 as the only Echo runtime for generated code, handwritten HTTP server code, and HTTP-facing tests. The repository SHALL not retain direct Echo v4 runtime dependencies in the main application module after the migration is complete.

#### Scenario: Generated and handwritten HTTP code target the same Echo major version
- **WHEN** code generation is run and the application is built
- **THEN** generated server interfaces and handwritten HTTP handlers MUST both compile against Echo v5 types

#### Scenario: Echo v4 runtime dependencies are removed from the main application path
- **WHEN** the migration is complete
- **THEN** the root module MUST depend on Echo v5-compatible runtime packages for the application path instead of Echo v4-specific equivalents

### Requirement: Code generation stays on Echo v5
The repository SHALL pin `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen` to the requested commit and SHALL configure `.codegen.yml` to generate Echo v5 server code so future `task generate` runs preserve the migrated runtime.

#### Scenario: Generation uses the requested tool pin
- **WHEN** a developer inspects the root `go.mod` tool block after the migration
- **THEN** the `oapi-codegen` tool entry MUST resolve to commit `434cc16b427dcc1d63fc3e7e6b9c2f5337418be1`

#### Scenario: Generation emits Echo v5 handlers
- **WHEN** `task generate` is run after the migration
- **THEN** generated HTTP server code MUST be produced using the Echo v5 server target configured through `.codegen.yml`

### Requirement: Echo-specific middleware stays behaviorally equivalent
The service SHALL use Echo v5-compatible versions of the Swagger, Sentry, and Zap middleware packages while preserving the current middleware responsibilities: Swagger asset serving, Sentry request instrumentation, request ID propagation, and request/response logging.

#### Scenario: Middleware dependencies are upgraded to Echo v5 variants
- **WHEN** the migration updates module dependencies
- **THEN** `github.com/adlandh/echo-oapi-middleware`, `github.com/adlandh/echo-sentry-middleware`, and `github.com/adlandh/echo-zap-middleware` MUST be consumed via their `/v2` module paths

#### Scenario: Existing middleware behavior is preserved
- **WHEN** the service starts after the migration
- **THEN** Swagger serving, Sentry instrumentation, request logging, and request ID propagation MUST remain enabled with behavior equivalent to the pre-migration configuration

### Requirement: Public HTTP behavior remains stable across the migration
The migration SHALL preserve the existing externally visible webhook and message endpoints, including route patterns, successful responses, authorization failure behavior, and message retrieval semantics.

#### Scenario: Existing routes remain available
- **WHEN** the migrated service registers its HTTP handlers
- **THEN** it MUST continue to expose `GET /`, `POST /api/{token}/{service}`, `GET /api/{token}/{service}`, and `GET /api/message/{id}`

#### Scenario: Existing error semantics remain available
- **WHEN** a request uses an invalid auth token or requests a missing message
- **THEN** the migrated service MUST continue to return the same HTTP status class and response semantics as before the migration

### Requirement: Repository verification passes after migration
The migrated repository SHALL complete the standard repository verification flow successfully after regeneration.

#### Scenario: Repository-native checks succeed
- **WHEN** the migration implementation is finished
- **THEN** `task generate`, `task test`, and `task lint` MUST succeed without requiring a fallback to Echo v4 code paths
