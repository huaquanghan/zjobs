# zjobs — Claude Code guide

Local-first Go CLI (`jdctl`) that analyzes one JD against one profile + one CV variant and emits a Markdown + JSON fit report (verdict `Apply` / `Stretch` / `Skip`, 6-axis rubric score, evidence gaps, lite prep pack).

## Commands

```sh
go build ./cmd/jdctl      # build CLI
go test ./...             # unit + golden + e2e smoke suites
go vet ./...
```

## Hard contracts (do not break without a plan change)

- **CLI surface is locked** — `jdctl analyze jd` and its flags (`--profile`, `--cv`, `--url|--file|--paste`, `--analysis`, `--out`, `--config`). See `cmd/jdctl/cmd/root.go` and `analyze.go`.
- **Exactly one profile + one CV variant per run**, and exactly one JD source.
- **Read-only toward profile/CV files** — never write or mutate them.
- **Report contract** — always emit Markdown + JSON in the output dir (`data/reports/` by default); schemas live in `internal/analysis/schema.json` and `internal/reporting/schema.json`.
- **Go owns determinism** — the Claude skill produces the semantic analysis JSON; Go validates it (`internal/analysis.Validate`) and computes the deterministic verdict/score (`internal/analysis.Evaluate` with weights from config).

## Architecture

- `cmd/jdctl/` — Cobra entrypoint; `version` is injected via `-ldflags "-X zjobs/cmd/jdctl/cmd.version=vX.Y.Z"` (default `dev`).
- `internal/ingest/` — JD ingestion (`FromURL` / `FromFile` / `FromPaste`) + `Provider` search interface; Exa is the only real provider (config key or `EXA_API_KEY` env var, never a committed key).
- `internal/domain/` — `Profile` (YAML), `CVVariant` (Markdown), `JobDescription` models with validation.
- `internal/analysis/` — 6-axis rubric weights: hard_constraints, must_have_skills, nice_to_have_skills, seniority_scope, domain_context, evidence_strength.
- `internal/reporting/` — `Build` + `WriteReport` (JSON + Markdown, filename-stamped).

## Conventions

- Match existing style: table-driven tests, golden fixtures under `testdata/`, doc comments on exported symbols.
- Any change to scoring weights, schemas, or the CLI surface updates `config.example.yaml` and the golden fixtures, and proves itself with `go test ./...`.

## Workflow

`AGENTS.md` is the entrypoint. This repo runs the zharness workflow — see `docs/WORKFLOW.md` and the stage playbooks under `docs/playbooks/`; the spec lives in `SPEC.md`.
