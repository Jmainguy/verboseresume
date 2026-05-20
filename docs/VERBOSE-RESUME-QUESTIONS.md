# Sample questions for a verbose resume

Use these prompts on yourself, in a journal, or with an LLM acting as an interviewer. Answer in plain language first; your agent can later file answers under the right `### Company — Title` and `#####` headings in `verboseResume.md`.

**Rules for LLM-assisted interviews**

- The model should ask **one question at a time** (or a short cluster of related follow-ups), then wait for your answer.
- It must **not invent** employers, dates, tools, or metrics—only record what you say.
- After each answer, it should propose **where** the note belongs (Overview, Teams, or a new `#####` under Notable work).
- When you are done, it should output a **draft Markdown patch** or a bullet list of additions—not a tailored JSON resume.

## Getting started

- What is your full name and how do you want it to appear on resumes?
- Where are you based, and are you open to remote, hybrid, or relocation?
- In one paragraph, how would you explain your career to a peer at a meetup?
- What kinds of roles are you targeting in the next two years?
- What do you want to be known for technically (e.g. platform, data, security, frontend)?

## Contact and links

- What is the best email and phone for recruiters?
- What personal site or blog should stay on the resume?
- Which GitHub, LinkedIn, or portfolio links are still current?

## Skills inventory

- Which languages do you use in production today? Which are rusty?
- Which clouds and Kubernetes flavors have you operated in anger?
- What CI/CD, IaC, and observability tools have you **owned** vs only used?
- What certifications are active, expired, or planned?

## Per employer (repeat for each job)

### Context

- What company did you work for, and what did they sell or do?
- What was your official title? Any title changes while you were there?
- What were your start and end dates? Full-time, contract, or acquisition?
- Where was the job located (office, hybrid, remote)?

### Team you joined first

- What team did you start on? Who was your manager and roughly how many people were on the team?
- Why were you hired—what problem were they trying to solve?
- What was broken, missing, or painful when you arrived?
- What did you ship or fix in the first 90 days?

### Team moves and scope changes

- Did you move to another team at the same company? When and why?
- What was the new team’s mission? How did your day-to-day change?
- Did you keep responsibilities from the old team, or hand them off cleanly?
- Any **dotted-line** or **loan** work to other groups (e.g. 20% to a program)?

### Cross-team and leadership

- Which other teams did you depend on (security, networking, data, product)?
- Where did you lead design reviews, RFCs, or standards other teams adopted?
- Did you mentor anyone, run onboarding, or facilitate incidents/postmortems?
- Any conflict or negotiation worth remembering (priorities, budget, headcount)?

### Delivery and impact

- Name three projects you are proud of. What was hard about each?
- What tools, languages, and platforms did each project touch?
- What broke in production? What did you learn from incidents?
- Any metrics you can cite truthfully (latency, cost, error rate, time-to-deploy, team size)?
- What would have gone wrong if you had not been there?

### Things people forget

- Certificates, keys, or compliance work (SOC2, HIPAA, PCI) you touched
- On-call rotations, game days, runbooks, or docs you wrote
- Conference talks, internal trainings, or interview loops you ran
- Reasons you left (neutral wording for your own memory)

## Projects and notable work (dig deeper)

For each project theme, ask:

- What was the user or business outcome?
- What was your specific contribution vs the team’s?
- What alternatives did you consider (build vs buy, cloud A vs B)?
- What is the **before / after** story in one sentence?
- What would you put in a `#####` heading title for this work?

## Education and certifications

- Degree, school, location, graduation date?
- Relevant coursework, thesis, or honors worth mentioning?
- Certifications: passed date, renewal, lapse?

## Maintenance cadence

- What did you do **this month** at work that might matter in six months?
- Any tickets, PRs, postmortems, or design docs to skim before they disappear?
- What did you learn that is not in git yet?

## LLM interviewer prompt (copy/paste)

```text
You are helping me build verboseResume.md — a private, detailed career memory.
Ask me one interview question at a time from a structured career interview
(company, first team, team moves, cross-team work, projects, tools, incidents,
metrics only if I provide them).

Rules:
- Do not invent facts.
- After each answer, tell me which section heading the note belongs under.
- When I say "done for this role", summarize bullets I should paste into Markdown.
- Do not output tailored resume JSON.

Start with my current or most recent role. First question:
"What company do you work at, what is your title, and what team did you start on?"
```

## Example verbose resume

See the fictional sample in [`example-verbose-resume.md`](example-verbose-resume.md) (also on the site Docs page). It shows multiple employers, team changes, `#####` project sections, and an appendix for easy-to-forget facts.
