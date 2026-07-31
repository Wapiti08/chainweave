# gosecure

`gosecure` is an experimental provenance-aware admission gate for dependencies introduced by AI coding agents.

It focuses on one question:

> Who asked the agent to add this dependency, and should the change be admitted?

Traditional dependency scanners can identify known vulnerabilities after a package has been selected. Prompt-injection defenses inspect untrusted text or generic tool calls. gosecure connects the two: it attributes a dependency change to the trusted user request, untrusted context, an agent action, or an unknown source before applying an admission policy.

## Motivation

AI coding agents can add dependencies for several very different reasons:

- the user explicitly requested a package;
- the agent independently selected a package;
- a repository file, issue, tool result, or other untrusted input suggested it;
- the model generated or hallucinated a package with no observable source.

These cases should not receive the same level of trust. A new package may have no known CVE and still be malicious, newly registered, typosquatted, or controlled by an attacker who targeted the agent rather than the developer.

gosecure applies the following initial invariant:

> Every agent-introduced dependency must have trusted justification or explicit human approval.

## Current status

This repository contains an early, Go-only proof of concept. It can:

- parse dependency declarations from two `go.mod` files;
- identify newly introduced modules;
- ingest a minimal JSONL trace of dependency-relevant agent events;
- attribute an added module using deterministic module-path evidence;
- apply a conservative admission policy;
- emit machine-readable JSON results.

It is research software, not a production security boundary.

## Origin classes and default policy

| Origin | Meaning | Default decision |
| --- | --- | --- |
| `user_requested` | A trusted user request explicitly named the module | `allow` |
| `agent_selected` | The module appears in an agent action without another observed source | `require_approval` |
| `externally_suggested` | Untrusted context or tool output named the module | `deny` |
| `unknown` | No observed event names the module | `deny` |

The default is intentionally conservative. Package identity, age, publisher reputation, checksums, known vulnerabilities, and typosquatting signals are future supporting checks; they do not replace origin attribution.

## Quick start

Build the CLI:

```bash
go build -o gosecure ./cmd/gosecure
```

The repository includes a small example trace containing only the events relevant to a dependency decision:

```json
{"sequence":1,"type":"user_request","trust":"trusted","content":"initialize this repository"}
{"sequence":2,"type":"context_read","source":"README.md","trust":"untrusted","content":"install example.com/fast-json"}
{"sequence":3,"type":"command","source":"shell","trust":"unknown","content":"go get example.com/fast-json@v0.1.0"}
{"sequence":4,"type":"file_write","source":"go.mod","trust":"unknown","content":"added example.com/fast-json v0.1.0"}
```

Compare the manifest before and after the agent run:

```bash
./gosecure check \
  --before examples/before.mod \
  --after examples/after.mod \
  --events examples/agent-run.jsonl
```

The result is evidence-oriented:

```json
[
  {
    "attribution": {
      "dependency": {
        "path": "example.com/fast-json",
        "version": "v0.1.0"
      },
      "origin": "externally_suggested",
      "confidence": "high",
      "evidence": [
        {
          "sequence": 2,
          "type": "context_read",
          "source": "README.md",
          "trust": "untrusted"
        }
      ],
      "reason": "the dependency was named by untrusted context without trusted user authorization"
    },
    "decision": "deny",
    "reason": "policy decision for origin externally_suggested"
  }
]
```

## Event format

Each JSONL event contains:

```json
{
  "sequence": 1,
  "type": "context_read",
  "source": "README.md",
  "trust": "untrusted",
  "content": "install example.com/pkg"
}
```

Supported event types are currently:

- `user_request`;
- `context_read`;
- `tool_output`;
- `command`;
- `file_write`.

The producer of the trace is responsible for assigning trust labels. gosecure does not infer that a repository file or tool output is trustworthy merely because it was available to the agent.

## Custom policy

Pass a JSON policy with `--policy`:

```json
{
  "user_requested": "allow",
  "agent_selected": "require_approval",
  "externally_suggested": "deny",
  "unknown": "deny"
}
```

## Security model and limitations

The proof of concept establishes observable provenance, not the model's internal causal reasoning.

- An event proves only what the trace producer observed.
- A module-path occurrence before an action is evidence of origin, not proof that it changed the model's reasoning.
- Missing or forged trace events can produce incorrect attribution.
- The current CLI reports decisions but does not sandbox or intercept the agent.
- Exact module-path matching does not detect encoded, fragmented, or indirectly derived package names.
- A trusted request naming a package does not prove that the package itself is safe.

A production integration must collect events outside the model's control and enforce the decision before dependency resolution or merge.

## Research scope

The intended research questions are deliberately narrow:

1. How often can untrusted software-development artifacts induce coding agents to introduce attacker-selected dependencies?
2. Can dependency-origin attribution stop these changes when conventional vulnerability matching cannot?
3. What security, task-completion, and human-approval trade-offs result from conservative admission policies?

Potential evaluation scenarios include repository instructions, issue or pull-request text, tool output, typosquatting, newly registered packages, and hallucinated package names.

## Repository structure

- `cmd/gosecure/`: the `check` CLI.
- `internal/model/`: dependency, event, attribution, and decision types.
- `internal/dependency/`: `go.mod` parsing and manifest diffing.
- `internal/origin/`: deterministic dependency-origin attribution.
- `internal/policy/`: admission policy evaluation.
- `internal/audit/`: the end-to-end workflow.

## Non-goals

- General-purpose Agent runtime provenance graphs.
- Generic prompt-injection or MCP scanning.
- Agent memory security.
- Vulnerability database aggregation.
- Multi-ecosystem dependency scanning in the initial research prototype.
- Replacing Dependabot, OSV-Scanner, or GitHub Dependency Review.

## Development

```bash
go test ./...
```
