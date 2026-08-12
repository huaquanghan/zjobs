<!-- ZHARNESS:BEGIN -->
## Harness

Run `zharness --version`, then `zharness preflight <stage> [--mode <mode>] --json` for every workflow skill invocation. Follow a returned stop and recovery exactly.

Read `docs/WORKFLOW.md`, then only the returned stage playbook and repository material relevant to the requested outcome. Repository docs, code, tests, and observable behavior are authoritative; the database is a lifecycle ledger and recovery index.

Read-only and bounded work may use reduced mode and must not mutate harness state. Durable planning, full execution, full checks, and durable handoffs require an initialized database. Claim completion only with executable or observable evidence.
<!-- ZHARNESS:END -->

## Project

`zjobs` is a local-first Go CLI (`jdctl`) that analyzes one job description (JD) against one profile + one CV variant and emits a Markdown + JSON fit report: verdict `Apply` / `Stretch` / `Skip`, a 6-axis rubric score, evidence-backed gaps, and a lite prep pack. Spec: `SPEC.md`. Full agent guide: `CLAUDE.md`.

## Commands

- Build: `go build ./cmd/jdctl`
- Test: `go test ./...` (unit + golden fixtures + e2e smoke)

## Contracts

- CLI surface is locked: `jdctl analyze jd --profile X --cv Y --url|--file|--paste Z --analysis A [--out DIR] [--config PATH]`.
- Read-only toward profile/CV source files — never mutate them.
- Exactly one profile, one CV variant, one JD source per run.
- Go computes the deterministic verdict/score from the Claude-produced analysis JSON; weights come from config (see `config.example.yaml`).
- Never commit an Exa API key — use the `EXA_API_KEY` env var or an empty config placeholder.
