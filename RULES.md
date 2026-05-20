# LLM Resume Tailoring Rules

These rules apply when an LLM edits this repo, creates a tailored resume, or converts `verboseResume.md` into upload-ready JSON.

## Truth first

- Never lie.
- Never invent employers, dates, titles, tools, certifications, degrees, citizenship, clearance, metrics, or outcomes.
- If a job description asks for something not present in the source material, do not claim it. Either omit it or phrase nearby truthful experience carefully.
- It is acceptable to reword, reorder, shorten, and emphasize true experience for a specific job.
- Prefer direct, plain language over buzzwords when both would mean the same thing.

## Source and output

- Treat `verboseResume.md` as the detailed source of truth (see `docs/VERBOSE-RESUME.md`).
- Encourage monthly updates to `verboseResume.md` while actively working in a role so details are captured before they are forgotten.
- It is fine for `verboseResume.md` to contain too much information. The tailoring step is responsible for cutting it down.
- Do not upload the verbose Markdown file to the web app—only tailored JSON.
- Treat job-specific JSON/YAML files as tailored outputs derived from the source and target job description.
- The upload format is the JSON shape used by `resume.json` and documented in `README.md`.
- Keep skill entries short and machine-readable, for example `Kubernetes`, `EKS`, `Terraform`, `Grafana`.
- Put explanations and context in `Summary` or `Experience[].Details`, not inside long skill labels.

## Tailoring strategy

- Optimize for the smallest number of pages that still makes the candidate credible for the role.
- Prefer fewer, stronger bullets over many broad bullets.
- Put the most relevant experience, tools, and outcomes first.
- Tune the summary and the first two bullets of the most recent role especially carefully.
- Remove details that do not help the target job, even if they are impressive in general.
- Preserve important keywords from the job description when they are truthful.

## Tone

- Keep the resume human and direct.
- Avoid robotic phrasing, inflated claims, and generic filler.
- Avoid em dashes in tailored output unless the user explicitly asks for them.
- Do not over-explain company-specific internal tool names that a recruiter will not recognize.

## Page layout

- Optimize for fewer pages whenever possible.
- Keep summary paragraphs concise.
- Keep bullets focused enough that browser print can avoid awkward page breaks.
- If a section spans pages badly, shorten or remove lower-value bullets before shrinking everything aggressively.

## Templates

- Multiple HTML templates are supported under `templates/`.
- Reusable built-in templates must be registered in `templateOptions` in `main.go` before they appear in the upload form.
- One-off custom templates can be uploaded alongside the JSON resume and do not need to be committed or registered.
- Resume templates should be content fragments wrapped by `templates/resume-page.html`; do not duplicate the WYSIWYG editor, print controls, app header, or app footer in each template.
- Templates may change layout, typography, spacing, section order, and print rules.
- Templates must not invent, hide, or distort resume facts.
