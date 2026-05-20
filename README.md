# Verbose Resume

[Verbose Resume](https://verboseresume.com) turns structured resume data into an HTML resume you can review in the browser, edit lightly, and print to PDF.

The workflow is built around two formats:

- `verboseResume.md`: the long, Markdown source-of-truth resume (for humans and agents).
- `resume.json`: the compact upload format consumed by the Go web app.

Job-specific files such as `cava.json`, `cava.yaml`, or `tetra*` files are intentionally local tailoring artifacts and are ignored by git.

Full Markdown structure and maintenance rules: [`docs/VERBOSE-RESUME.md`](docs/VERBOSE-RESUME.md).

Brand assets and colors: [`BRAND.md`](BRAND.md) (`static/brand/`).

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
make build-prod && ./verboseresume
```

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

## CI and release

GitHub Actions (see `.github/workflows/`):

- **CI** on push/PR: lint, test, production build, smoke-test embedded server
- **Release** on main via release-please: GoReleaser binaries plus `ko` images to GHCR

Do not commit personal resume JSON or `verboseResume.md`—see `.gitignore`. Use `static/sample-upload.json` as a public example only.
