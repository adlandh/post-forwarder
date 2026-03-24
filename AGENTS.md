# AGENTS.md
Guide for coding agents working in `github.com/adlandh/post-forwarder`.

## Scope
- Go 1.25 webhook forwarding service.
- Main app code is under `internal/post-forwarder`.
- Infra code is a separate Go module in `inf/`.
- Generated files are committed; do not hand-edit `*.gen.go` unless changing generation itself.

## Layout
- `internal/post-forwarder/main.go`: `fx` wiring, Echo setup, middleware, logging, Sentry.
- `internal/post-forwarder/application`: application-layer logic.
- `internal/post-forwarder/domain`: interfaces, request ID helper, generators, mocks, wrappers.
- `internal/post-forwarder/driven`: Redis and notifier implementations.
- `internal/post-forwarder/driver`: HTTP handlers, OpenAPI output, generated wrappers.
- `internal/chat-id-checker`: small auxiliary program.
- `api/post-forwarder.yaml`: OpenAPI source of truth.
- `Taskfile.yaml`: main local task runner.
- `.github/workflows/`: CI behavior for lint/test/build.

## Tooling
- Primary app toolchain: Go 1.25.
- Common workflows use `task`.
- Linting uses `golangci-lint`.
- Generation uses `go generate`, `oapi-codegen`, `gowrap`, `mockery`, and `find-interfaces`.
- Some tests require Docker because they use Testcontainers.

## Build, Generate, Lint, Test

### Generate
- Preferred: `task generate`
- Direct: `go generate ./internal/post-forwarder/...`
- CI-style: `go generate ./...`
- Sources/config: `internal/post-forwarder/driver/generators.go`, `internal/post-forwarder/domain/generators.go`, `gen-wraps.sh`, `.codegen.yml`, `.mockery.yml`

### Build
- Compile everything: `go build ./...`
- Build app package: `go build ./internal/post-forwarder`
- Install app locally: `go install ./internal/post-forwarder`
- CI release image build: `ko build -B --tags <tag>,latest --platform linux/amd64,linux/arm64 ./internal/post-forwarder`

### Test
- Preferred full run: `task test`
- Direct full run: `go test ./...`
- CI command: `go test -race -coverprofile=coverage.txt -covermode=atomic ./...`
- Package run: `go test ./internal/post-forwarder/driven`
- Verbose package run: `go test -v ./internal/post-forwarder/driver`

### Run A Single Test
- Single test function: `go test ./internal/post-forwarder/application -run TestProcessRequest`
- Single suite test method: `go test ./internal/post-forwarder/driven -run TestRedisStorage/TestRead`
- Single test in main package: `go test ./internal/post-forwarder -run TestCreateService`
- Single test with race detector: `go test -race ./internal/post-forwarder/driver -run TestHttpServer`
- For `testify/suite`, use `SuiteName/TestName` after `-run`.

### Lint And Format
- Preferred lint run: `task lint`
- Task/CI lint flow:
- `go generate ./...`
- `curl -sS https://raw.githubusercontent.com/adlandh/golangci-lint-config/refs/heads/main/.golangci.yml -o .golangci.yml`
- `golangci-lint run`
- Formatters enabled by GolangCI: `gofmt`, `goimports`
- Useful local commands: `gofmt -w <file>.go`, `goimports -w <file>.go`

## Hooks And Quality Gates
- `.lefthook.yml` defines a `pre-push` hook.
- Pre-push runs `task generate`, `task lint`, and `task test` in parallel.
- Expect generated files, lint issues, and tests to be resolved before push.

## Infra Commands
- Pulumi module is in `inf/`.
- Show stack outputs: `task inf-show`
- Interactive update: `task inf`
- Non-interactive update: `task inf-y`

## Code Style

### General
- Follow normal Go style; this repo already aligns with `gofmt` and `goimports`.
- Keep package boundaries clean: `driver` handles transport, `application` handles use cases, `domain` defines contracts, `driven` implements integrations.
- Keep dependency wiring in `main.go`; do not move business logic into bootstrap code.
- Prefer constructors named `New<Type>`.
- Keep files focused and small when practical.
- Use compile-time interface assertions for implementations.

### Imports
- Let `goimports` manage order and grouping.
- Standard library first, then third-party, then repo-local imports.
- Alias imports only for clarity or name collisions.
- Existing alias style is acceptable when it improves readability, e.g. `contextlogger`, `echoSentryMiddleware`, `sentryecho`.
- Blank imports are fine only for side effects or tool/generator support.

### Formatting
- Run `gofmt` or `goimports` after every Go edit.
- Use the formatting produced by Go tools; do not hand-align spacing.
- Keep control flow readable instead of golfing expressions.
- Group related constants in `const (...)` blocks.

### Types And Interfaces
- Use concrete structs for implementations and keep cross-layer contracts in `domain` interfaces.
- Preserve existing public shapes unless a larger refactor is intended.
- Prefer standard library types such as `context.Context`, `time.Time`, and typed slices.
- Named return values are acceptable when they simplify multi-value flows; this repo uses that style in several places.
- Avoid `interface{}` unless absolutely necessary.

### Naming
- Exported identifiers: PascalCase.
- Unexported identifiers: camelCase.
- Constructors: `New<Type>`.
- Sentinel errors follow the current repo pattern, e.g. `ErrorNotFound`.
- Config/service constants are descriptive and singular, e.g. `TelegramService`.
- Tests use `TestXxx`; suite methods also use `TestXxx`.

### Error Handling
- Wrap returned errors with `%w` when callers may need the cause.
- Use `errors.Is` for sentinel handling.
- Use `errors.Join` for multiple validation/configuration failures.
- Log internal operational failures with structured fields before returning generalized external errors when appropriate.
- Avoid panics in application code.
- Return early on invalid input and setup failures.

### Context, Logging, Observability
- Thread `context.Context` through request, storage, and notifier paths.
- Preserve request-scoped values like `domain.RequestID`.
- Use structured Zap fields rather than interpolated log strings when logging context.
- Keep Sentry optional; the no-Sentry path must continue to work.
- Respect cancellation and timeouts in lifecycle hooks and outbound calls.

### HTTP Layer
- Echo handlers generally return `echo.NewHTTPError(...)` on request failures.
- Keep request parsing and HTTP response shaping inside `driver`.
- Delegate business behavior to `application`.
- Prefer OpenAPI-generated route/types glue over hand-written contract code.
- Preserve current health check behavior unless requirements explicitly change.

### Testing
- The repo uses `require` and `testify/suite`.
- Prefer `require` for stop-now assertions.
- Regenerate mocks with `mockery`; do not edit generated mocks directly.
- Use `t.Cleanup(...)` when overriding package globals or test hooks.
- Integration tests should set explicit timeouts when containers or network services are involved.
- `internal/post-forwarder/driven/redisstorage_test.go` uses Testcontainers Redis.

### Generated Code
- Never manually edit `*.gen.go` files.
- If interfaces, generator config, or `api/post-forwarder.yaml` change, run `task generate` and commit generated outputs.
- Keep generator sources and generated outputs in sync.

### Secrets And Safety
- Never commit secrets, tokens, copied credentials, or local-only values.
- Treat `app-secrets.yaml`, `k8s/secrets.yaml`, and personal notes as sensitive.
- Do not paste secret values into docs, tests, fixtures, issues, or commit messages.

## Cursor / Copilot Rules
- No `.cursor/rules/` directory was found.
- No `.cursorrules` file was found.
- No `.github/copilot-instructions.md` file was found.
- If any appear later, merge them into this document and treat repo-local agent rules as higher priority than generic guidance.

## Change Checklist
- Regenerate code if interfaces or OpenAPI changed.
- Format edited Go files.
- Run the narrowest relevant tests first.
- Run broader tests when changes cross package boundaries.
- Run lint before finishing any Go change.
- Do not hand-edit generated files.
