# AGENTS.md

## Workflow

- Follow the repo-local OpenSpec workflow when it applies: `.opencode/commands/opsx-*.md` and `.opencode/skills/openspec-*` are the primary instruction sources. Matching `.cursor/` and `.codex/` copies exist.
- For non-trivial work, prefer the wired flow: `opsx-explore` -> `opsx-propose` -> `opsx-apply` -> `opsx-archive`.

## Repo Shape

- Main service entrypoint: `internal/post-forwarder/main.go`.
- Secondary binary: `internal/chat-id-checker/main.go`.
- HTTP API source of truth is `api/post-forwarder.yaml`; generated Echo server code lives under `internal/post-forwarder/driver`.
- Runtime wiring in `internal/post-forwarder/main.go` uses Echo v5, Fx DI, optional Sentry setup, and generated Sentry decorators from `internal/post-forwarder/domain/wrappers`.

## Commands

- Prefer `task` over guessing raw commands.
- `task generate` runs `go generate ./internal/post-forwarder/...`.
- `task test` depends on `generate`, then runs `go test -cover -race -v ./...` from `internal/`.
- `task lint` depends on `generate`, downloads `.golangci.yml` from `adlandh/golangci-lint-config`, then runs `golangci-lint` from `internal/`.
- CI does not use the task runner: it runs `go generate ./...` and `go test -race -coverprofile=coverage.txt -covermode=atomic ./...` from the repo root.
- Release builds use `ko build -B --tags <tag>,latest --platform linux/amd64,linux/arm64 ./internal/post-forwarder`.
- `task up` runs Tilt using `Tiltfile`, which builds `./internal/post-forwarder` and port-forwards the app to `localhost:8081`.

## Codegen

- Never hand-edit `*.gen.go`.
- If you change `api/post-forwarder.yaml`, `internal/post-forwarder/domain/interfaces.go`, `.codegen.yml`, `.mockery.yml`, or `gen-wraps.sh`, run `task generate`.
- Code generation is driven by `//go:generate` directives in `internal/post-forwarder/driver/generators.go` and `internal/post-forwarder/domain/generators.go`.
- Tooling is pinned in the root `go.mod` `tool` block; use `go tool ...` rather than `go install`.
- `task rm-generated` deletes committed generated files across the repo; use it only intentionally.

## Verification

- Focused package tests: `go test ./internal/post-forwarder/<package>`.
- `testify/suite` tests can be narrowed with `-run SuiteName/TestName`.
- `internal/post-forwarder/driven/redisstorage_test.go` uses Testcontainers Redis; Docker must be running for that test path.
- Lefthook pre-push runs `task generate`, `task lint`, and `task test` in parallel for Go changes.

## Gotchas

- Do not edit `.golangci.yml` manually; `task lint` and CI both overwrite it from the shared config repo.
- The service is on Echo v5 (`github.com/labstack/echo/v5` plus `github.com/getsentry/sentry-go/echo v0.45.1`); keep generated and handwritten HTTP code on the same Echo major version.
- For local `task up` on Docker Desktop, prefer the `kubeadm` Kubernetes backend over `kind`; with `kind`, this repo hit image pull/storage issues (`unexpected EOF`, local `ko` images not visible to pods) that disappeared after switching to `kubeadm`.
- Keep the existing sentinel name `domain.ErrorNotFound`; that spelling is used intentionally in the codebase.
- `app-secrets.yaml`, `k8s/secrets.yaml`, and rendered Doppler output are sensitive; `app-secrets.yaml.gotmpl` is the safe template.
