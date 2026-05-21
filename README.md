# Verbose Resume

[Verbose Resume](https://verboseresume.com) turns structured resume data into an HTML resume you can review in the browser, edit lightly, and print to PDF.

The workflow is built around two formats:

- `verboseResume.md`: the long, Markdown source-of-truth resume (for humans and agents).
- `resume.json`: the compact upload format consumed by the Go web app.

Job-specific files such as `cava.json`, `cava.yaml`, or `tetra*` files are intentionally local tailoring artifacts and are ignored by git.

Full Markdown structure and maintenance rules: [`docs/VERBOSE-RESUME.md`](docs/VERBOSE-RESUME.md).

Brand assets and colors: [`BRAND.md`](BRAND.md) (`static/brand/`).

Contributing: fork the repo and open a PR from your fork—see [/contribute](https://verboseresume.com/contribute) for templates, docs, and workflow tips (no direct write access). Licensed under **GNU GPLv2** ([`LICENSE`](LICENSE)).

## The core idea

Write one very verbose resume that captures more context than any one application should include.

That verbose resume can then be fine-tuned for a specific job by an LLM using the job description. The LLM should select, reword, compress, and order the most relevant experience for that role, while preserving factual accuracy. The result is shorter job-specific JSON you upload to Verbose Resume and print as a polished PDF.

This keeps the full career history in one structured place, while making each application easier to tune for the employer's actual needs.

## `verboseResume.md`

`verboseResume.md` is intentionally detailed. It is meant for humans and LLMs to read—not for direct upload to the HTML renderer.

Suggested structure:

```markdown
# Your Name — Verbose Resume

## Contact
- **Location:** ...
- **Email:** ...

## Summary
Long-form career narrative.

## Skills
### Languages
- Go

## Experience
### Company — Title
**Dates:** Apr 2024 – Present | **Location:** Remote

#### Overview
Role scope and impact.

#### Teams
- Team name

#### Notable work
##### Project or theme
Specific deliverables, tools, outcomes.

## Education
- Degree | School | Location | Date

## Certifications
- Name
```

When adding to `verboseResume.md`, prefer truth and specificity over polish. It is fine for this file to be long, repetitive, and too detailed—it is raw material for tailoring.

**In this repository**

- [`docs/VERBOSE-RESUME.md`](docs/VERBOSE-RESUME.md) — format spec and maintenance
- [`docs/VERBOSE-RESUME-QUESTIONS.md`](docs/VERBOSE-RESUME-QUESTIONS.md) — sample interview questions (self or LLM interviewer)
- [`docs/example-verbose-resume.md`](docs/example-verbose-resume.md) — fictional, very verbose example (structure only)

On the site, open [Docs](https://verboseresume.com/docs#verbose-resume-questions) for the question guide and example.

### Keep it fresh

Update `verboseResume.md` at least once a month while actively working in a role.

Capture what actually happened before you forget it:

- Projects started or finished.
- Tools, platforms, frameworks, or services you touched.
- Incidents, outages, migrations, releases, or production problems you helped solve.
- Architecture decisions and tradeoffs.
- Team leadership, mentoring, planning, or cross-team coordination.
- Metrics, scale, business context, and constraints.

Agents can help by reviewing recent prompts, git history, Jira tickets, pull requests, and docs, then proposing additions under the correct `###` role and `#####` notable-work headings.

### Legacy YAML

`verboseResume.yaml` is deprecated. Convert with:

```bash
python3 scripts/yaml_to_verbose_md.py verboseResume.yaml verboseResume.md
```

## Upload JSON format

The Go app expects JSON matching the `Resume` struct in `main.go`.

Uploads are processed in memory for the current request only. The app does not save uploaded resume JSON or custom templates to disk.

Required top-level fields:

```json
{
  "Name": "Jonathan Mainguy",
  "Location": "Garner, NC",
  "Phone": "434-229-1127",
  "Email": "jon@example.com",
  "Website": "https://example.com",
  "Github": "https://github.com/example",
  "Summary": "Short tailored summary.",
  "Skills": {
    "Languages": ["Go", "Python", "Bash"],
    "Cloud & Containers": ["Kubernetes", "EKS", "GKE", "AWS"],
    "CI / CD": ["ArgoCD", "GitHub Actions"],
    "IaC & Automation": ["Terraform", "Helm"],
    "Observability": ["Prometheus", "Grafana"],
    "Security": ["Hashicorp Vault"],
    "Other": []
  },
  "Experience": [
    {
      "Title": "Role title",
      "Company": "Company name",
      "Location": "Remote",
      "Date": "Apr 2024 - Present",
      "Details": [
        "Tailored bullet focused on the target role."
      ]
    }
  ],
  "Education": [
    {
      "Degree": "Bachelor of Science, Accounting",
      "School": "Liberty University",
      "Location": "Lynchburg, Virginia",
      "Date": "May 2009"
    }
  ],
  "Certifications": [
    "Certification name"
  ]
}
```

### JSON notes

- `Summary` is a string. Keep it as a concise tailored summary.
- `Skills` is a map of display group names to arrays of short skill labels.
- `Experience[].Details` should be tailored bullets.
- Validate with `python3 -m json.tool your-resume.json > /dev/null` before upload if editing by hand.

## Tailoring with an LLM

Work in a JSON loop—one `verboseResume.md`, many tailored JSON files:

1. Give the LLM `verboseResume.md`, the job description, and the upload JSON spec.
2. Ask for **only** valid upload JSON. Save and upload to Verbose Resume.
3. Preview in the browser. Not right? Paste the **current JSON** back with feedback; get revised JSON; re-upload.
4. Print or save PDF when satisfied. Use WYSIWYG edit mode only for tiny fixes, or sync changes back into JSON in the chat.

Do not maintain a second Markdown copy of the tailored resume. The LLM should not invent experience, employers, dates, tools, credentials, citizenship status, clearance, or outcomes.

## Agent-native MCP

POST JSON-RPC to `/mcp`. Tools include:

- `get_resume_generator_guide`
- `get_verbose_resume_format` — Markdown structure for `verboseResume.md`
- `get_verbose_resume_questions` — interview prompts and LLM interviewer template
- `get_example_verbose_resume` — fictional verbose example for structure reference
- `get_upload_json_format`
- `get_llm_prompt_guide`
- `create_resume_artifact` — agent-side HTML/JSON artifacts and local PDF instructions from upload JSON

## Templates

The upload form supports multiple HTML templates. Pick the layout that best fits the audience and page count.

Included starter templates:

- `templates/resume-page.html`: shared web wrapper (edit toolbar, print controls, template switch).
- `templates/classic.html`: default traditional serif layout.
- `templates/clean.html`: clean layout with skill pills.
- `templates/compact.html`: tighter layout.

`clean.html`, `compact.html`, `classic.html`, and uploaded custom templates receive the same `Resume` data from `main.go`. They can use:

- `.Name`, `.Location`, `.Phone`, `.Email`, `.Website`, `.Github`
- `.Summary`, `.Skills`, `.Experience`, `.Education`, `.Certifications`
- `summaryParagraphs`, `hasSummaryRoleLead`, `summaryAfterRoleLead`

### Writing your own template

Custom templates should be resume fragments, not full pages. The app wraps them in `templates/resume-page.html`.

To add a reusable built-in layout:

1. Create `templates/your-template.html` with `{{ define "your-template.html" }}`.
2. Register it in `templateOptions` in `main.go`.

Templates should not invent, hide, or distort resume facts.

## Running locally

```bash
go run .
# or
make run
```

`make build-prod` writes **`verboseresume`** and copies it to **`resumeGen`** (legacy name). After pulling updates, rebuild before `./resumeGen`:

```bash
make build-prod && ./resumeGen
# or: go build -o resumeGen . && ./resumeGen
```

Kill any old server on port 8080 if the docs page looks wrong: `pkill -f resumeGen; fuser -k 8080/tcp`

Open `http://localhost:8080`, upload JSON, review HTML, and print to PDF from Chrome.

Templates, static assets, and `docs/VERBOSE-RESUME.md` are embedded in the binary—no files need to be mounted beside the executable.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port when `LISTEN_ADDR` is unset |
| `LISTEN_ADDR` | `:8080` | Full listen address (overrides `PORT`) |
| `MAX_UPLOAD_BYTES` | `5242880` | Max multipart upload size (5 MiB) |

## Container

```bash
docker build -t verboseresume:local .
docker run --rm -p 8080:8080 -e PORT=8080 verboseresume:local
```

Published images (on release): `zot.soh.re/verboseresume/verboseresume:latest`

## Contributing

See **[verboseresume.com/contribute](https://verboseresume.com/contribute)** for the fork-and-PR workflow, resume templates, and what we are looking for. We do not grant push access to the upstream repo. Contributions are licensed under **GNU GPLv2** ([`LICENSE`](LICENSE)).

### Repository rules

- **Public app code only** — templates, `main.go`, embedded docs, `static/sample-upload.json`, and fictional examples under `docs/`.
- **No personal resume data in git** — do not commit `verboseResume.md`, `verboseResume.yaml`, `resume.json`, or job-specific `cava*` / `tetra*` files (see [`.gitignore`](.gitignore)). Use [`static/sample-upload.json`](static/sample-upload.json) and [`docs/example-verbose-resume.md`](docs/example-verbose-resume.md) for public examples.
- **No secrets** — never commit `.env`, registry passwords, or API tokens.
- **Truth in tailored output** — when using LLMs on resume content, follow [`RULES.md`](RULES.md) (no invented employers, dates, tools, or metrics).
- **Templates must not distort facts** — layout and typography only; do not hide or fabricate resume content in HTML templates.

### Dependencies

Prefer the **latest stable** versions when changing the project:

- **Go** — match or exceed the version in [`go.mod`](go.mod) (currently Go 1.26). Update `go.mod`, workflows, and the [`Dockerfile`](Dockerfile) together.
- **GitHub Actions** — bump action pins in [`.github/workflows/`](.github/workflows/) to current major releases (checkout, setup-go, golangci-lint-action, release-please, GoReleaser, ko, docker/login-action).
- **golangci-lint** — keep the CI pin in [`golang-ci.yml`](.github/workflows/golang-ci.yml) aligned with the latest v2.x release.

[Renovate](https://docs.renovatebot.com/) is enabled for dependency PRs; still verify CI locally before merging.

This app has **no third-party Go modules** today—only the standard library—so `go.mod` is mostly the toolchain version.

### Commit messages

This repository follows [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/). [release-please](https://github.com/googleapis/release-please) on `main` maps commit types to semantic versions:

| Type | Typical release bump |
|------|----------------------|
| `fix:` | patch |
| `feat:` | minor |
| `feat!:` or `BREAKING CHANGE:` footer | major |
| `docs:`, `chore:`, `ci:`, `refactor:`, `test:`, … | no version bump unless paired with user-facing `feat`/`fix` |

Examples:

```text
feat: add classic template preview
fix: reject empty resume JSON uploads
docs: clarify verbose resume workflow
ci: publish images to zot.soh.re
chore: bump Go 1.26 and GitHub Actions
```

- Use **`feat:`** / **`fix:`** on `main` when you intend a versioned release.
- Prefer **one logical change per commit**; squash PRs to conventional messages when needed.
- Non-conventional subjects (e.g. `Fix CI lint`) are ignored by release-please.

### CI must pass

**All CI jobs must pass** on every push and pull request before merge. Run the same checks locally when you can:

```bash
# Same steps as CI (fmt, vet, test, production build)
make ci

# Lint (install v2.12.2 or newer to match CI)
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
make lint
```

CI workflow ([`golang-ci.yml`](.github/workflows/golang-ci.yml)) runs:

1. **golangci-lint** — `golangci-lint run --timeout=5m`
2. **fmt** — `gofmt` with no diffs
3. **vet** — `go vet ./...`
4. **test** — `go test ./...`
5. **build** — `make build-prod`
6. **smoke test** — start the binary and `curl` `/`, `/docs`, and a static asset

Optional manual smoke test after `make build-prod`:

```bash
./verboseresume &
curl -fsS http://127.0.0.1:8080/ http://127.0.0.1:8080/docs
kill %1
```

Open a PR only when local checks are green (or when CI is green on the PR branch). Fix failures on your branch; do not merge with failing checks.

### Pull requests

- Target **`main`**.
- Keep diffs focused; avoid unrelated refactors in the same PR.
- Update **docs/tests** when behavior or MCP tools change.
- Do not include personal resume files or generated PDFs in the PR.

## Release

On **`main`**, after conventional commits land:

1. **release-please** opens or updates a release PR (changelog + version bump).
2. Merging that PR creates a GitHub release tag.
3. The **release** workflow builds binaries with GoReleaser and publishes images to `zot.soh.re/verboseresume/verboseresume`.

See [`.github/workflows/golang-release.yml`](.github/workflows/golang-release.yml). Deploy config for production lives in the separate [clusters](https://github.com/Standouthost/clusters) repo (Argo CD).
