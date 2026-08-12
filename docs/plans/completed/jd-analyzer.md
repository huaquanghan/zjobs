---
id: 01KZTD66588TH9WA6FV7S0DGME
type: plan
intake_id: 01KZTD682G1RN5SH4CKAPXXCCC
lane: normal
status: completed
created: 2026-08-12
updated: 2026-08-12
---

# Plan: Local-first Go JD Analyzer with Claude Skills

## Outcome
- result: A local-first Go CLI plus Claude skill files that, given one profile + one CV variant and one JD from URL/file/paste, emits a Markdown + JSON report with verdict `Apply / Stretch / Skip`, a 6-axis rubric score, evidence-backed gaps, and a lite prep pack (interview questions, verified references, next actions); the tool is read-only toward profile/CV source files.
- success_signals:
  - A cold reviewer can run `analyze jd` against golden fixtures and see Markdown + JSON output matching the expected schema.
  - Golden fixture verdicts and gaps match stored expectations, and `go test ./...` plus the golden suite pass.
  - One smoke run on a real public JD produces a valid report without login or CAPTCHA bypass.
  - Profile/CV source files are never modified by the tool.

## Authority and Requirements
- authority:
  - Locked spec `/home/tinhpt/Personal/zjobs/SPEC.md` (interviewed 2026-08-12, accepted by owner)
  - Upstream inspiration reviewed: https://github.com/MadsLorentzen/ai-job-search
- requirements:
  - R1 [accepted]: CLI command `analyze jd` accepts exactly one JD via URL, file path, or pasted text, plus one selected profile and one CV variant | source: SPEC Outcome
  - R2 [accepted]: Report emits both Markdown and JSON from a single stable schema, with verdict `Apply / Stretch / Skip` | source: SPEC Outcome
  - R3 [accepted]: Scoring uses a 6-axis rubric (hard constraints, must-have skills, nice-to-have skills, seniority/scope, domain/context, evidence strength), with hard gates on work constraints that can block `Apply` | source: SPEC Key Decisions 4
  - R4 [accepted]: Gaps are evidence-backed: each gap cites JD/CV evidence, impact on verdict, and an action (learn, rewrite CV, prepare story, or skip) | source: SPEC Key Decisions 6
  - R5 [accepted]: Lite prep pack includes top 3 gaps, 3–5 interview questions, verified references (titles/urls/source, never invented), and 1–2 next actions | source: SPEC Key Decisions 6
  - R6 [accepted]: JD acquisition uses public fetch plus a search-provider interface; Exa is the first real provider; no login, CAPTCHA bypass, or aggressive scraping in MVP | source: SPEC Key Decisions 3
  - R7 [accepted]: One profile per run only; batch/multi-profile ranking is deferred | source: SPEC Key Decisions 2
  - R8 [accepted]: Tool is read-only for profile/CV source files; suggestions are diff/proposal text, not auto-applied edits | source: SPEC Key Decisions 5
  - R9 [accepted]: Runs append Markdown + JSON reports and an index (JSONL/CSV) for history and dedupe | source: SPEC Key Decisions 5
  - R10 [accepted]: Verification gate = golden fixtures + unit tests + one real public JD smoke run | source: SPEC Key Decisions 7

## Non-goals
- NG1: Multi-profile ranking, batch analyze, and profile auto-choosing in a single run
- NG2: Automated cover letter / CV drafting, PDF compilation, or ATS validation
- NG3: Login-based scraping, CAPTCHA bypass, aggressive rotation, or browser automation in MVP
- NG4: Full multi-week learning roadmap generation per JD (only lite prep pack)
- NG5: SQLite tracker or rich query layer for history (reports + JSONL/CSV index only)
- NG6: Tavily/Brave providers running in MVP (adapter interface may exist; Exa is the only real provider)
- NG7: Mutation of profile/CV files, including auto-apply or apply-with-prompt

## Approach and Risks
- approach: Go core + Claude skill orchestration. Cobra CLI at `cmd/jdctl`; internal packages `domain`, `ingest`, `analysis`, `reporting`; Claude skill `analyze-jd` emits a structured JSON analysis contract from JD + CV, Go validates the schema and computes deterministic rubric/verdict; search provider interface with Exa as first adapter; reports to `data/reports` with a JSONL index; no write path to profile/CV files by design.
- constraints:
  - One profile per run; JD from URL, file, or paste; public fetch only, no login/CAPTCHA bypass
  - Stable CLI command names, flags, and Markdown/JSON report schema across the MVP
  - Read-only toward profile/CV source files; suggestions are proposals only
- risks:
  - Exa API key missing, rate-limited, or costly | mitigation: provider interface, env-key config, mock adapter in tests, clear error when key absent
  - Public JD pages block fetch or change HTML | mitigation: file/paste fallback contract, public-fetch-only rule, no bypass
  - Claude analysis JSON drifts from Go schema or hallucinates references | mitigation: strict schema validation in Go, verified-URL-only reference rule, evidence citations required
  - Contract drift between skill files and Go | mitigation: golden fixtures plus JSON schema tests as the gate
- rejected_alternatives:
  - Direct Claude API calls from Go: adds key/cost/retry scope beyond MVP orchestration boundary
  - SQLite tracker/history: deferred to NG5; reports + JSONL index is enough for MVP proof
  - Browser automation or login-based scraping: violates NG3
  - Boards-first job sources (LinkedIn/Indeed): anti-bot and ToS risk; ATS + manual import is the MVP source

## Phases and Verification
<!-- Phase and task definitions are immutable after to-plan. Do not add task status fields. Append-only Progress is the sole task execution-status source. Only each phase lifecycle status changes to mirror DB transitions: to-plan=planned; work after run create=in-progress; clean durable check=checked; closing handoff=done. Each planned phase records phase_slug, story_id, status, goal, depends_on, waves, tasks, and checks. -->
- planning_status: planned
- phases:
  - phase_slug: scaffold-cli
    story_id: 01KZTD73KTYGDT9NVXR3TFJ7X8
    status: checked
    goal: Scaffold Go module, cobra CLI entrypoint and internal package layout with a stable command surface
    depends_on: []
    waves:
      - wave: 1
        tasks:
          - Init Go module and directory layout `cmd/jdctl` + `internal/{domain,ingest,analysis,reporting}`
          - Add cobra CLI with root command and stable `analyze jd` command surface (flags: --profile, --cv, --input/--url/--file/--paste, --out)
          - Wire config loading (YAML) for rubric weights and provider settings
          - Add `--version` and `--help` with expected text
        checks:
          - `go build ./...` compiles
          - `go run ./cmd/jdctl --help` lists `analyze jd` with locked flags
          - `go vet ./...` passes
  - phase_slug: domain-ingest
    story_id: 01KZTD73M4FH3GPBFGZTC166RQ
    status: checked
    goal: Define domain models (profile, CV variant, JD), YAML/MD loaders, and JD ingestion from URL/file/paste plus a search provider interface with an Exa adapter
    depends_on: [scaffold-cli]
    waves:
      - wave: 1
        tasks:
          - Define `Profile`, `CVVariant`, and `JobDescription` domain models with JSON tags
          - Implement YAML profile loader with validation (role targets, work constraints, must-have/nice-to-have skills)
          - Implement CV variant loader from Markdown with frontmatter
          - Implement JD ingestion from file, pasted text, and public URL fetch with timeout and robots/ToS respect
          - Define `Provider` interface and implement `ExaProvider` with env-key config and mockable HTTP client
        checks:
          - Unit tests for loaders and validators on valid and invalid fixtures pass
          - Unit test for Exa adapter uses `httptest` and passes without network
          - URL fetch test against a local `httptest` server passes
  - phase_slug: analysis
    story_id: 01KZTD73MAS89BDWEAC84G3JJ2
    status: checked
    goal: Implement 6-axis rubric scoring, work-constraint hard gates, evidence-backed gaps, and Apply/Stretch/Skip verdict from a Claude skill structured JSON contract validated by Go
    depends_on: [domain-ingest]
    waves:
      - wave: 1
        tasks:
          - Define the Claude analysis JSON contract (axes, gates, gaps, evidence citations) with JSON schema
          - Implement Go validation of analysis payloads with clear error messages
          - Implement hard-gate evaluation on work constraints (location/remote, work authorization, language, must-have years/seniority)
          - Implement weighted 6-axis scoring and `Apply / Stretch / Skip` threshold mapping
          - Normalize evidence-backed gaps into the shared model with impact and action (learn, rewrite CV, prepare story, skip)
        checks:
          - Schema validation rejects malformed/missing-evidence payloads
          - Hard-gate unit tests: each gate can flip verdict to Skip
          - Scoring unit tests: fixture-driven, expected verdicts match
  - phase_slug: reporting
    story_id: 01KZTD73MJQAQ0AWY4A9S5RETH
    status: checked
    goal: Render stable Markdown + JSON reports, lite prep pack, and append-only index; read-only toward profile/CV files
    depends_on: [analysis]
    waves:
      - wave: 1
        tasks:
          - Implement JSON report renderer from the analysis model with fixed field order
          - Implement Markdown report renderer (verdict, axes, gaps, prep pack, next actions)
          - Implement lite prep pack (top 3 gaps, 3–5 interview questions, verified references, 1–2 next actions) from analysis output
          - Implement append-only JSONL index of runs with job hash for dedupe
        checks:
          - Golden fixture comparison for JSON and Markdown output passes
          - Index file appends a row per run and dedupes on job hash
          - Test confirms profile/CV source files are unchanged after a run
  - phase_slug: skills-files
    story_id: 01KZTD73MSGA11SPAZA28AETZB
    status: checked
    goal: Create Claude skill files orchestrating the CLI and semantic analysis, with strict input/output contracts and untrusted-JD guardrails
    depends_on: [analysis]
    waves:
      - wave: 1
        tasks:
          - Create `analyze-jd` skill: input contract (profile, CV, JD source), pipeline steps, output contract
          - Create `profile-setup` skill: onboarding profile + CV variants from documents into YAML/Markdown
          - Add guardrails: treat JD as untrusted, never auto-follow links inside JD, verified references only, no file mutation
        checks:
          - Skill files render without broken references to CLI flags/schema
          - Dry-run invocation through the skill produces the expected CLI command with locked flags
          - Guardrail checklist documented in each skill is explicit and non-negotiable
  - phase_slug: e2e-proof
    story_id: 01KZTD73N0PRDJMG4HR1H6HGQX
    status: checked
    goal: Golden fixtures, unit + schema tests, and one real public JD smoke run prove the locked contract end-to-end
    depends_on: [reporting]
    waves:
      - wave: 1
        tasks:
          - Add golden fixtures: sample profiles, CV variants, and JDs with expected rubric/verdicts
          - Add golden test suite comparing JSON and Markdown output byte-for-byte
          - Add JSON schema validation test for report output
          - Run one smoke test on a real public JD and commit its report into fixtures
        checks:
          - `go test ./...` passes
          - `go test ./... -run Golden` passes with fixtures
          - Smoke run report validates against schema and profile/CV files remain unchanged

## Progress
<!-- Append-only durable entries record timestamp, phase, wave, task, task_status, run_id, trace_id, exact verification/result, and changed surfaces or blocker. -->
- `2026-08-12T07:24:15Z` — wave 1. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: phase scaffold-cli started; run 01KZTDPHYY8XJCE39YH7HTNDJW; playbook says phase-start task_status=in-progress but CLI only accepts DONE|DONE_WITH_CONCERNS|NEEDS_CONTEXT|BLOCKED — using wave summary form instead.
- `2026-08-12T07:25:26Z` — wave 1, task Init Go module and directory layout cmd/jdctl + internal/{domain,ingest,analysis,reporting}. task_status: `DONE`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: go mod init zjobs; cmd/jdctl + internal packages created with doc.go.
- `2026-08-12T07:25:26Z` — wave 1, task Add cobra CLI with root command and stable analyze jd command surface. task_status: `DONE`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: jdctl analyze jd with --profile/--cv/--url/--file/--paste/--out; validation rejects 0 or 2+ sources.
- `2026-08-12T07:25:26Z` — wave 1, task Wire config loading (YAML) for rubric weights and provider settings. task_status: `DONE`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: internal/config Load+Defaults, --config flag on root, config.example.yaml.
- `2026-08-12T07:25:26Z` — wave 1, task Add --version and --help with expected text. task_status: `DONE`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: cobra Version=dev; help lists analyze jd with locked flags.
- `2026-08-12T07:25:29Z` — wave 1. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. summary: scaffold-cli wave 1 done: go build, go vet pass; jdctl analyze jd --help lists locked flags; --version and source validation verified.
- `2026-08-12T07:27:28Z` — wave 1, task Define Profile, CVVariant, and JobDescription domain models with JSON tags. task_status: `DONE`. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: models.go with JSON+YAML tags; Profile/WorkConstraints/CVVariant/JobDescription/JDSource.
- `2026-08-12T07:27:28Z` — wave 1, task Implement YAML profile loader with validation. task_status: `DONE`. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: LoadProfile/ParseProfile + Validate: name, role targets, must-have skills, locations/remote, work auth, languages, min_years.
- `2026-08-12T07:27:28Z` — wave 1, task Implement CV variant loader from Markdown with frontmatter. task_status: `DONE`. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: LoadCVVariant/ParseCVVariant splitFrontmatter + Validate; fixtures valid/invalid.
- `2026-08-12T07:27:28Z` — wave 1, task Implement JD ingestion from file, pasted text, and public URL fetch with timeout and robots/ToS respect. task_status: `DONE`. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: FromFile/FromPaste/FromURL with 15s timeout, UA header, best-effort robots.txt Disallow parser, sha256 hash.
- `2026-08-12T07:27:28Z` — wave 1, task Define Provider interface and implement ExaProvider with env-key config and mockable HTTP client. task_status: `DONE`. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: Provider.Search interface, JobHit, ExaProvider with EXA_API_KEY, injectable client+base URL.
- `2026-08-12T07:27:28Z` — wave 1. run: `01KZTDTK24KTX43JSP98ATKV3F`. summary: domain-ingest wave 1 done: 8 unit tests pass across domain+ingest; httptest covers URL fetch, robots disallow, Exa adapter.
- `2026-08-12T07:29:22Z` — wave 1, task Define the Claude analysis JSON contract (axes, gates, gaps, evidence citations) with JSON schema. task_status: `DONE`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: contract.go + schema.json (draft 2020-12) embedded via go:embed; axes/gates/gaps/prep with enums and evidence minLength.
- `2026-08-12T07:29:22Z` — wave 1, task Implement Go validation of analysis payloads with clear error messages. task_status: `DONE`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: Validate(): strict decode DisallowUnknownFields + jsonschema/v5 compile; tests reject unknown fields, bad axis range, bad gate/action, non-http reference.
- `2026-08-12T07:29:22Z` — wave 1, task Implement hard-gate evaluation on work constraints. task_status: `DONE`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: 5 gates (location, work_authorization, language, experience_years, seniority); any fail blocks to Skip; TestEachGateCanFlipVerdict proves each flips.
- `2026-08-12T07:29:22Z` — wave 1, task Implement weighted 6-axis scoring and Apply/Stretch/Skip threshold mapping. task_status: `DONE`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: Evaluate(): weighted score over config weights; Apply>=0.8, Stretch>=0.6, else Skip; missing weight/axis errors.
- `2026-08-12T07:29:22Z` — wave 1, task Normalize evidence-backed gaps into the shared model with impact and action. task_status: `DONE`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: Gap model with requirement/jd_evidence/cv_evidence/impact/action; action enum learn|rewrite_cv|prepare_story|skip enforced.
- `2026-08-12T07:29:22Z` — wave 1. run: `01KZTDXY8S28D29B7AXKQH6JK8`. summary: analysis wave 1 done: 12 tests pass incl. schema rejections and per-gate flip; fixtures apply/gate-fail/stretch/malformed.
- `2026-08-12T07:30:58Z` — wave 1, task Implement JSON report renderer from the analysis model with fixed field order. task_status: `DONE`. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: RenderJSON with indented encoder; Report struct fixed field order; golden report-apply.json byte-for-byte stable.
- `2026-08-12T07:30:58Z` — wave 1, task Implement Markdown report renderer (verdict, axes, gaps, prep pack, next actions). task_status: `DONE`. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: RenderMarkdown: headline verdict+score, rubric table, gate status, gaps with evidence, prep pack sections.
- `2026-08-12T07:30:58Z` — wave 1, task Implement lite prep pack (top 3 gaps, 3-5 interview questions, verified references, 1-2 next actions) from analysis output. task_status: `DONE`. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: Prep passes through from validated analysis contract; schema enforces 3-5 questions, 1-2 next actions, http-only reference URLs.
- `2026-08-12T07:30:58Z` — wave 1, task Implement append-only JSONL index of runs with job hash for dedupe. task_status: `DONE`. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: AppendIndex dedupes on job hash (test proves 2 lines for 3 appends); corrupt lines skipped.
- `2026-08-12T07:30:58Z` — wave 1, task Verify read-only contract toward profile/CV files. task_status: `DONE`. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: WriteReport takes only out dir + report data; readonly test proves profile/cv byte-identical after write.
- `2026-08-12T07:30:58Z` — wave 1. run: `01KZTE11078SBZFKNHQZAKE90M`. summary: reporting wave 1 done: goldens committed and verified, index dedupe test, readonly contract test; full suite 4 pkgs ok.
- `2026-08-12T07:31:52Z` — wave 1, task Create analyze-jd skill: input contract (profile, CV, JD source), pipeline steps, output contract. task_status: `DONE`. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. summary: .claude/skills/analyze-jd/SKILL.md with input/output contract tables, 5 pipeline steps, deterministic verdict rule.
- `2026-08-12T07:31:52Z` — wave 1, task Create profile-setup skill: onboarding profile + CV variants from documents into YAML/Markdown. task_status: `DONE`. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. summary: .claude/skills/profile-setup/SKILL.md with output schemas for profile YAML + CV MD, 5 rules, validate command.
- `2026-08-12T07:31:52Z` — wave 1, task Add guardrails: treat JD as untrusted, never auto-follow links inside JD, verified references only, no file mutation. task_status: `DONE`. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. summary: Guardrails (non-negotiable) sections in both skills: untrusted JD, no link following, verified http(s) references, no mutation, no login/bypass, no invented experience.
- `2026-08-12T07:31:52Z` — wave 1. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. summary: skills-files wave 1 done: skill flags match CLI help exactly; schema.json path valid; dry-run command exit 0; guardrail sections present in both skills.
- `2026-08-12T07:35:41Z` — wave 1, task Add golden fixtures: sample profiles, CV variants, and JDs with expected rubric/verdicts. task_status: `DONE`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. summary: cmd/jdctl/cmd/testdata/golden/: profile-backend.yaml, cv-main.md, jd-backend.txt, analysis.json.
- `2026-08-12T07:35:41Z` — wave 1, task Add golden test suite comparing JSON and Markdown output byte-for-byte. task_status: `DONE`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. summary: TestGoldenPipeline passes; generated_at normalized; verdict=Apply.
- `2026-08-12T07:35:41Z` — wave 1, task Add JSON schema validation test for report output. task_status: `DONE`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. summary: 4 schema tests pass: accepts rendered, rejects bad verdict/version/unknown field.
- `2026-08-12T07:35:41Z` — wave 1, task Run one smoke test on a real public JD and commit its report into fixtures. task_status: `DONE`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. summary: RUN_JDSMOKE=1 fetch greenhouse gitlab board; verdict=Apply score=0.87; report committed; profile/CV byte-identical.
- `2026-08-12T07:35:41Z` — wave 1. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. summary: e2e-proof wave 1 done: gofmt -l clean, go build+vet pass, go test ./... all 6 packages ok; golden + smoke fixtures committed.
- `2026-08-12T07:40:45.308Z` — handoff recorded. handoff: `01KZTEN8BWK0RJG955ZKER86CF`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. check: `01KZTDSY4XQC2XJQ5YYTZG08VB`. phase closed.
- `2026-08-12T07:40:45.319Z` — handoff recorded. handoff: `01KZTEN8C7XFQD14JXC2VA3YXV`. run: `01KZTDTK24KTX43JSP98ATKV3F`. check: `01KZTDX5MAFHCFST7G61WVXVVF`. phase closed.
- `2026-08-12T07:40:45.329Z` — handoff recorded. handoff: `01KZTEN8CH021QF9BTFASNQ7H2`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. check: `01KZTE0HB2N5DNWTC3NHSPCS9Q`. phase closed.
- `2026-08-12T07:40:50.576Z` — handoff recorded. handoff: `01KZTENDGGHQ37KCTN7BC8GJCZ`. run: `01KZTE11078SBZFKNHQZAKE90M`. check: `01KZTE3F93YW5J97M6RS83GP9E`. phase closed.
- `2026-08-12T07:40:50.585Z` — handoff recorded. handoff: `01KZTENDGSJX9X8RF4FVFD06BX`. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. check: `01KZTE54FBN25A3PJKWKGKCWV4`. phase closed.
- `2026-08-12T07:40:50.593Z` — handoff recorded. handoff: `01KZTENDH19VDZFAZA16N79JXW`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. check: `01KZTECMCHPA5MTH0NKF0WG8A0`. phase closed.

## Decisions
<!-- Append-only durable entries record timestamp, phase/task, decision, and rationale. -->
- `2026-08-12T07:33:04Z` — analyze jd gains --analysis <path> flag for the Claude analysis JSON (phase: `e2e-proof`), task: Wire the full pipeline into analyze jd. rationale: The locked architecture routes semantic analysis through the Claude skill and deterministic verdict/report through Go; the CLI must accept that JSON to run end-to-end. All scaffold-phase flags are preserved; --analysis is an addition required by e2e-proof..

## Validation
<!-- Append-only durable entries record timestamp, phase, exact command/result/output, run_id, check_id, verdict, and proof_gaps. -->
- `2026-08-12T07:25:50.109Z` — check. verdict: `APPROVED`. check: `01KZTDSY4XQC2XJQ5YYTZG08VB`. run: `01KZTDPHYY8XJCE39YH7HTNDJW`. phase: `scaffold-cli`. judge: `same-session` (deepseek-v4-flash).
  - `go build ./...` → Validation 2026-08-12: exit 0, no output
  - `go vet ./...` → Validation 2026-08-12: exit 0, no findings
  - `go test ./...` → Validation 2026-08-12: 7 packages, no failures
  - `go run ./cmd/jdctl analyze jd --help` → Validation 2026-08-12: lists analyze jd with locked flags
  - `go run ./cmd/jdctl --version` → Validation 2026-08-12: prints jdctl version dev
- `2026-08-12T07:27:36.074Z` — check. verdict: `APPROVED`. check: `01KZTDX5MAFHCFST7G61WVXVVF`. run: `01KZTDTK24KTX43JSP98ATKV3F`. phase: `domain-ingest`. judge: `same-session` (deepseek-v4-flash).
  - `go build ./...` → Validation 2026-08-12: exit 0, no output
  - `go vet ./...` → Validation 2026-08-12: exit 0, no findings
  - `go test ./...` → Validation 2026-08-12: domain + ingest packages pass, 8 tests
- `2026-08-12T07:29:26.370Z` — check. verdict: `APPROVED`. check: `01KZTE0HB2N5DNWTC3NHSPCS9Q`. run: `01KZTDXY8S28D29B7AXKQH6JK8`. phase: `analysis`. judge: `same-session` (deepseek-v4-flash).
  - `go build ./...` → Validation 2026-08-12: exit 0
  - `go vet ./...` → Validation 2026-08-12: exit 0
  - `go test ./...` → Validation 2026-08-12: analysis+domain+ingest pass, 20 tests
- `2026-08-12T07:31:02.563Z` — check. verdict: `APPROVED`. check: `01KZTE3F93YW5J97M6RS83GP9E`. run: `01KZTE11078SBZFKNHQZAKE90M`. phase: `reporting`. judge: `same-session` (deepseek-v4-flash).
  - `go build ./...` → Validation 2026-08-12: exit 0
  - `go vet ./...` → Validation 2026-08-12: exit 0
  - `go test ./...` → Validation 2026-08-12: 4 packages ok incl. golden byte-compare, index dedupe, readonly contract
- `2026-08-12T07:31:57.035Z` — check. verdict: `APPROVED`. check: `01KZTE54FBN25A3PJKWKGKCWV4`. run: `01KZTE43YN1FXMSP79EA0D1Z0Y`. phase: `skills-files`. judge: `same-session` (deepseek-v4-flash).
  - `go build ./...` → Validation 2026-08-12: exit 0
  - `go vet ./...` → Validation 2026-08-12: exit 0
  - `go run ./cmd/jdctl analyze jd --profile internal/domain/testdata/profile-valid.yaml --cv internal/domain/testdata/cv-valid.md --file internal/ingest/testdata/jd.txt --out /tmp/jdctl-dryrun` → Validation 2026-08-12: dry-run exit 0 with locked flags
- `2026-08-12T07:36:02.705Z` — check. verdict: `APPROVED`. check: `01KZTECMCHPA5MTH0NKF0WG8A0`. run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. phase: `e2e-proof`. judge: `same-session` (deepseek-v4-flash).
  - `gofmt -l .` → empty output = clean
  - `go build ./... && go vet ./...` → pass
  - `go test ./...` → 6 packages ok; TestGoldenPipeline byte-for-byte; 4 schema tests; fixture smoke validates
  - `RUN_JDSMOKE=1 go test ./internal/e2e/ -run TestSmokeRealJD` → PASS verdict=Apply score=0.87; smoke report committed
- `2026-08-12T07:40:05Z` — check full (complete manual review). verdict: `APPROVED`. check: `01KZTECMCHPA5MTH0NKF0WG8A0` (final-phase gate check; harness `check record` requires story in-progress, so no second DB row — mode evidence per handoff.md step 6 is this session's record of invoking `/check full`). run: `01KZTE5JYAG75FDCY2R7ZQ7RM2`. phase: `e2e-proof`. judge: `same-session` (deepseek-v4-flash).
  - Security: JD untrusted as data end-to-end, no link-following from JD, no login/CAPTCHA bypass, strict schema+DisallowUnknownFields on analysis JSON, Exa key never logged, robots best-effort — pass
  - Performance: 4MB/1MB read caps, 15s timeouts, schema compiled once, index O(n) append — pass
  - Architecture: clean layering, deterministic verdict in Go over validated Claude JSON, R8 read-only structurally enforced (WriteReport takes only out dir) — pass
  - Code quality: gofmt clean, go vet clean, unit + golden byte-for-byte + schema + hermetic smoke + live smoke — pass
  - Not independently verified: live Exa API path (httptest mocks only) and analysis contract against real Claude output (committed fixture only) — named per same-session judge rule

## Current State and Next Action
- active_phase: e2e-proof
- lifecycle_status: checked
- latest_run_id: 01KZTE5JYAG75FDCY2R7ZQ7RM2
- latest_trace_ids: []
- latest_check_id: 01KZTECMCHPA5MTH0NKF0WG8A0
- latest_handoff_id: none
- blockers: none
- open_items: [close e2e-proof — closing handoff]
- exact_next_action: closing handoff
