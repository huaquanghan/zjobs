---
name: profile-setup
description: Build or update a job-search profile (YAML) and CV variants (Markdown with frontmatter) from user documents such as CVs, LinkedIn exports, or written summaries. Use when the user wants to onboard or refresh their profile for jdctl.
---

# profile-setup

Turns user documents into the two source-of-truth file types the analyzer
reads: one profile YAML and one or more CV variant Markdown files.

## Inputs

- Any of: existing CV (PDF/DOCX/Markdown), LinkedIn export, GitHub profile,
  portfolio, or the user's own written summary.
- The user's answer to: what roles are you targeting? (drives `role_targets`
  and `seniority`)

## Outputs

`profiles/<name>.yaml`:

```yaml
name: <handle>
role_targets: [<target role 1>, <target role 2>]
work_constraints:
  locations: [<city>, Remote]        # at least one, or remote_ok: true
  remote_ok: <true|false>
  work_authorization: [<country>]    # required
  languages: [<language>]            # required
  min_years: <int>                   # >= 0
  seniority: <junior|mid|senior|staff>
must_have_skills: [<skill>, ...]     # required, non-empty
nice_to_have_skills: [<skill>, ...]
```

`cvs/<variant>.md`:

```markdown
---
name: <variant name>
profile: <profile name>
---
<CV body: experience, skills, projects — plain Markdown>
```

## Rules

1. **Only facts from provided documents.** Never invent skills, years, or
   projects. If a document is ambiguous, ask the user.
2. **Map skills honestly.** An item belongs in `must_have_skills` only if the
   documents back it. If unsure between must-have and nice-to-have, put it in
   nice-to-have and note it.
3. **One profile, many variants.** A profile holds the intent; CV variants
   are different angles on the same career record. Create a new variant, not
   a new profile, for a different emphasis.
4. **Write, then validate.** After writing files, run:

```bash
go test ./internal/domain/ -run TestLoadProfile -v
```

   and fix any validation error before handing back.
5. **Never delete or rewrite existing profile/CV content** the user did not
   ask to change; propose changes as a diff instead.

## Guardrails (non-negotiable)

1. Facts must trace to a provided document or an explicit user statement.
2. No fabricated skills, titles, dates, or project outcomes.
3. Existing profile/CV files are only edited at the user's explicit request.
4. `must_have_skills` stays minimal — it drives the hard scoring axes.
