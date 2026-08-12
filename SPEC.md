---
title: Go JD Analyzer for Multi-Profile Job Search
status: validated
interviewed: 2026-07-31
---

## Outcome
A local-first Go CLI plus Claude skill files that let one user select **1 profile + 1 CV variant**, feed **1 JD** from URL/file/paste, and receive a **Markdown + JSON report** with:
- verdict `Apply / Stretch / Skip`
- rubric score across 6 axes
- evidence-backed gaps
- a light prep pack: a few interview questions, verified references, and next actions

The tool must be **read-only** for profile/CV files.

## Success Condition
A reviewer who did not take part in the interview can:
1. run `analyze jd` on golden fixtures
2. see Markdown + JSON output that matches the expected schema
3. verify fixture verdicts and gaps against stored expectations
4. run `go test ./...` plus the golden suite and see them pass
5. run one smoke test on a real public JD and get a valid report without login or CAPTCHA bypass
6. confirm the tool does not modify profile/CV source files

## Scope
**May change:** new repo scaffolding under `zjobs`, Go CLI, parsing/scoring/reporting code, search provider adapters, Claude skill files, fixtures, generated reports, and an index file for history.

**Must not change:** profile/CV source-of-truth files, one-profile-per-run contract, stable CLI command names and primary flags once locked, Markdown/JSON report contract, read-only behavior toward user profile/CV files, and public-fetch-only automation with no login/CAPTCHA bypass.

## Context to Read First
- this validated spec
- `https://github.com/MadsLorentzen/ai-job-search` README and any flow/docs that explain their discovery or apply loop
- the local profile/CV schema docs and golden fixtures once they exist in the repo

## Key Decisions
1. **Go core + Claude skills orchestration**: Go owns parsing, scoring, storage, and report generation; Claude skill files orchestrate the workflow and semantic analysis.
2. **One profile per run**: each analysis uses exactly one profile and one CV variant so the MVP stays narrow and predictable.
3. **Input contract**: JD input comes from URL, file, or pasted text; automation uses public fetch plus a search-provider interface, with Exa as the first real provider.
4. **Scoring model**: use a 6-axis rubric, apply hard gates for work constraints, then a weighted score to produce `Apply / Stretch / Skip`.
5. **Output contract**: always emit Markdown + JSON reports and append history to an index; the tool is read-only for profile/CV files.
6. **Prep pack lite**: include only top gaps, a few interview questions, verified references, and next actions.
7. **Verification strategy**: prove the contract with golden fixtures, unit tests, and one smoke run on a real public JD.

## Validation Loop
**During work:** run `go test ./...`, golden fixture assertions for rubric/verdict/schema, and a short smoke run on one JD sample after parser or scoring changes.

**Final proof:** full golden suite passes, JSON schema validates, a real public JD smoke run produces the expected report, and profile/CV files remain untouched.

## Stop / Pause
**Done when:** one command can analyze one JD against one chosen profile/CV and emit the stable report + prep pack, with golden fixtures proving the contract.

**Pause if:** the JD input is missing, the user has not chosen a profile/CV, or the target source requires login, CAPTCHA bypass, or another access decision that is outside the locked MVP.