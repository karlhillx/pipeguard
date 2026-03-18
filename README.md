# PipeGuard 🛡️

**PipeGuard** is a Go-based policy-as-code validator for Bitbucket Pipelines. It enforces organizational standards, security best practices, and CI/CD governance beyond basic YAML validity.

## Staff-Level Positioning

In modern engineering organizations, platform teams must balance developer autonomy with organizational safety. PipeGuard provides the "guardrails" necessary to:

- **Enforce Security**: Prevent the use of untrusted Bitbucket Pipes or unversioned/`:latest` Docker images.
- **Ensure Stability**: Mandate manual triggers for production deployments to prevent accidental releases.
- **Maintain Standards**: Require specific pipeline steps (like `test` or `security-scan`) to be present in all pipelines.
- **Governance at Scale**: Shift policy enforcement left into the development workflow, providing immediate feedback to engineers before code is merged.

## Features

- **Fast & Lightweight**: Written in Go for near-instant execution in local pre-commit hooks or CI runners.
- **Rich Rules Engine**: Supports complex logic across all pipeline triggers (branches, custom, tags, etc.).
- **Flexible Output**: Human-readable text for developers and JSON for machine-to-machine integration.
- **Severity Levels**: Support for `ERROR`, `WARNING`, and `INFO` to allow for phased policy rollouts.

## Installation

```bash
git clone github.com/karlhill/pipeguard
cd pipeguard
go build ./cmd/pipeguard
```

## Usage

Validate your `bitbucket-pipelines.yml`:

```bash
./pipeguard --config bitbucket-pipelines.yml
```

### Options

- `--config`: Path to the Bitbucket Pipelines YAML file (default: `bitbucket-pipelines.yml`).
- `--format`: Output format, either `text` (default) or `json`.

## Rules (MVP V1)

1.  **require-step**: Ensures a specific step name exists in the pipeline (e.g., `test`).
2.  **forbid-image-tag**: Denies usage of specific tags (e.g., `latest`) to ensure reproducible builds.
3.  **require-manual-trigger**: Mandates `trigger: manual` for steps with specific `deployment` environments (e.g., `production`).
4.  **allow-pipe-list**: Restricts Bitbucket Pipe usage to an approved allowlist.

## Example Output

```text
[ERROR  ] (forbid-image-tag) Forbidden image tag 'latest' used in global configuration: node:latest
[ERROR  ] (require-manual-trigger) Deployment 'production' in branch 'master' must have 'trigger: manual'.
[WARNING] (allow-pipe-list) Forbidden pipe used in branch 'master': docker/push-image:latest
```

## Project Structure

- `cmd/pipeguard/`: CLI entry point and configuration.
- `internal/parser/`: YAML parsing logic mapping Bitbucket's schema to Go structs.
- `internal/rules/`: The core rules engine and individual policy implementations.
- `internal/report/`: Logic for formatting and presenting findings.
