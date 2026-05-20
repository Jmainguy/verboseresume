# Verbose Resume (Markdown)

The verbose resume is your private, agent-readable source of truth. It holds more detail than any job-specific resume should include. An LLM (or you) tailors it into upload JSON for [Verbose Resume](https://verboseresume.com).

This file is **not** uploaded to the web app. The app consumes concise `resume.json` only.

## Design goals

- Easy for humans to read and edit in any editor
- Easy for agents to parse, diff, and append to over time
- Stable headings so prompts can say “add a `#####` under the current role’s Notable work”
- Room for long-form context: teams, projects, incidents, metrics, and tools

## Suggested file

Keep one file in your personal notes or repo, for example:

`verboseResume.md`

## Document structure

```markdown
# Your Name — Verbose Resume

## Contact
- **Location:** ...
- **Phone:** ...
- **Email:** ...
- **Website:** ...
- **GitHub:** ...

## Summary
One or more paragraphs of career narrative. More detail than the final resume.

## Skills

### Languages
- Go
- Python

### Cloud & Containers
- Kubernetes

## Experience

### Company Name — Job Title
**Dates:** Apr 2024 – Present | **Location:** Remote

#### Overview
Long-form description of the role, scope, and impact.

#### Teams
- Team name

#### Notable work

##### Short project or theme title
Paragraphs with specifics: tools, outcomes, scale, constraints, links to systems.

## Education
- Degree | School | Location | Date

## Certifications
- Certification name
```

### Heading rules

| Level | Use for |
|-------|---------|
| `#` | Document title: `{Name} — Verbose Resume` |
| `##` | Major sections: Contact, Summary, Skills, Experience, Education, Certifications |
| `###` | Skill group names, or each employer role as `### {Company} — {Title}` |
| `####` | Role metadata blocks: Overview, Teams, Notable work |
| `#####` | Individual projects, initiatives, incidents, or themes |

### Role metadata line

Immediately under each `### Company — Title` heading, use a single line:

`**Dates:** … | **Location:** …`

Agents and humans can scan roles quickly without extra YAML fields.

### Notable work

Put the granular material here—not in Overview alone. Overview should explain the role; Notable work should capture deliverables you might forget in a month:

- migrations, outages, releases, automation
- tools and platforms touched
- leadership, mentoring, process changes
- metrics, scale, and business context when truthful

## Maintenance (agent-native)

Update the verbose resume about **once a month** while actively employed. You will forget details faster than you expect.

Ask your agent to help by reviewing:

- recent prompts and notes
- git history (commits, PR titles, release notes)
- Jira or Linear tickets you closed or led
- incident/postmortem docs
- architecture docs and design threads

The agent should propose **additions only from evidence**, grouped under the correct role and `#####` heading.

## Tailoring to upload JSON

When applying for a job:

1. Provide `verboseResume.md`, the job description, and the upload JSON spec (see site Docs or MCP `get_upload_json_format`).
2. Ask the LLM to output **only** valid upload JSON—no invented facts.
3. Upload that JSON to Verbose Resume, preview, revise the JSON with your LLM if needed, then print or save PDF.

## Sample questions

Interview yourself—or use an LLM as a one-question-at-a-time interviewer—to fill gaps before they fade.

See [`VERBOSE-RESUME-QUESTIONS.md`](VERBOSE-RESUME-QUESTIONS.md) for prompts about employers, first teams, team moves, cross-team work, projects, incidents, and metrics. Includes a copy-paste **LLM interviewer** system prompt.

## Example

See [`example-verbose-resume.md`](example-verbose-resume.md) for a **fictional**, very verbose sample (Meridian Analytics / Lattice Systems / Northwind Digital). Use it for heading structure and level of detail—not as real experience.

Your private `verboseResume.md` stays out of git; see `.gitignore`.

## Legacy YAML

`verboseResume.yaml` is deprecated. Use `scripts/yaml_to_verbose_md.py` if you still have YAML to convert:

```bash
python3 scripts/yaml_to_verbose_md.py verboseResume.yaml verboseResume.md
```
