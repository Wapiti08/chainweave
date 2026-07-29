# gosecure

`gosecure` is an experimental, Go-specific dependency upgrade reviewer. It aims
to explain security-relevant behavior introduced by a new module release before
that behavior is associated with a published vulnerability or malware advisory.

The project focuses on one question:

> What security-sensitive behavior did this Go module upgrade introduce, and
> can the consuming application reach it?

## Project status

The repository is being rebuilt around this narrower direction. It currently
contains the intended package structure and a placeholder CLI, but no working
analysis pipeline. The examples below describe the target behavior, not a
released feature set.

## Intended workflow

Given two versions of a Go module and, optionally, a consuming project,
`gosecure` should:

1. acquire and normalize both module releases;
2. identify security-sensitive behavior in each release;
3. report only behavior introduced by the newer release;
4. determine whether the consumer can reach that behavior;
5. verify that published module artifacts match the expected source revision;
6. render evidence in terminal, JSON, or SARIF form.

An intended command could look like:

```bash
gosecure review example.com/module v1.4.2 v1.4.3 --consumer ./my-project
```

## Initial behavior coverage

The first implementation should remain deliberately small and concentrate on:

- process and shell execution;
- outbound network access and newly referenced domains;
- reads of credentials, sensitive environment variables, and sensitive files;
- use of CGO, `unsafe`, reflection, or embedded executables;
- changes to `go:generate`, build scripts, build tags, and platform-specific files.

Findings should include source locations, changed code, applicable build
constraints, and call paths where those paths can be established. A single
unexplained risk score is not considered sufficient evidence.

## Example report

The intended output is evidence-oriented:

```text
example.com/module v1.4.2 -> v1.4.3

HIGH    introduced process execution via os/exec.Command
        internal/update/install.go:48
        reachable from consumer command/update.Apply

MEDIUM  introduced outbound connection to telemetry.example.com
        internal/client/report.go:27
        linux/amd64 only

WARN    module proxy archive contains a file absent from the expected Git tag

Decision: manual review required
```

## Positioning

`gosecure` is not another general-purpose software composition analysis tool.
Dependency Review, Dependabot, OSV-Scanner, Trivy, Grype, and similar projects
already cover dependency inventory and known vulnerabilities well.

Known-vulnerability data may be included later as supporting context. The core
value must come from reproducible behavior differences, Go-aware reachability,
and source/release provenance evidence.

## Repository structure

- `cmd/gosecure/`: CLI entry point.
- `internal/source/`: acquire and normalize module releases.
- `internal/behavior/`: identify security-sensitive Go behavior.
- `internal/compare/`: compare behavior, source, and artifacts between releases.
- `internal/reachability/`: calculate consumer-to-dependency call paths.
- `internal/provenance/`: verify module artifacts against source revisions.
- `internal/consumer/`: discover the consumer's resolved dependency graph and
  build context.
- `internal/advisory/`: optional known-vulnerability context.
- `internal/report/`: terminal, JSON, and SARIF output.

These are intentionally empty package boundaries. Public APIs should be added
only after an end-to-end review workflow establishes the required data model.

## MVP order

1. Compare the source of two versions of one Go module.
2. Detect a small set of newly introduced sensitive API uses.
3. Report exact source evidence without a numeric risk score.
4. Add build-tag and platform awareness.
5. Add consumer reachability.
6. Add release provenance checks and SARIF output.

Multi-ecosystem support, generic lockfile scanning, policy engines, dashboards,
and vulnerability-database breadth are outside the initial MVP.

## Development

```bash
go test ./...
```

The CLI is not yet ready for general use.
