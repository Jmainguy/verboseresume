#!/usr/bin/env python3
"""One-off helper: convert verboseResume.yaml to verboseResume.md."""

import re
import sys
from pathlib import Path

import yaml


def para(text: str) -> str:
    text = (text or "").strip()
    if not text:
        return ""
    lines = [line.strip() for line in text.splitlines()]
    chunks: list[str] = []
    buf: list[str] = []
    for line in lines:
        if not line:
            if buf:
                chunks.append(" ".join(buf))
                buf = []
            continue
        buf.append(line)
    if buf:
        chunks.append(" ".join(buf))
    return "\n\n".join(chunks)


def md_escape_inline(text: str) -> str:
    return text.replace("\r", "").strip()


def convert(data: dict) -> str:
    personal = data.get("personal") or {}
    out: list[str] = [
        "<!--",
        "  Verbose Resume — source of truth for tailoring.",
        "  Maintain this file in Markdown for humans and agents.",
        "  Do not upload this file to the site; produce upload JSON instead.",
        "  See docs/VERBOSE-RESUME.md for structure and maintenance rules.",
        "-->",
        "",
        f"# {personal.get('name', 'Your Name')} — Verbose Resume",
        "",
        "## Contact",
        "",
    ]

    for label, key in [
        ("Location", "location"),
        ("Phone", "phone"),
        ("Email", "email"),
        ("Website", "website"),
        ("GitHub", "github"),
    ]:
        value = personal.get(key)
        if value:
            out.append(f"- **{label}:** {value}")

    out.extend(["", "## Summary", "", para(data.get("summary", "")), "", "## Skills", ""])

    skill_order = [
        "languages",
        "cloud_and_containers",
        "ci_cd",
        "iac_and_automation",
        "observability",
        "security",
        "Kafka",
    ]
    skill_labels = {
        "languages": "Languages",
        "cloud_and_containers": "Cloud & Containers",
        "ci_cd": "CI / CD",
        "iac_and_automation": "IaC & Automation",
        "observability": "Observability",
        "security": "Security",
        "Kafka": "Kafka",
    }
    skills = data.get("skills") or {}
    groups = [g for g in skill_order if g in skills] + [
        g for g in skills.keys() if g not in skill_order
    ]
    for group in groups:
        items = skills[group] or []
        if not items:
            continue
        title = skill_labels.get(group, group.replace("_", " ").title())
        out.append(f"### {title}")
        out.append("")
        for item in items:
            out.append(f"- {item}")
        out.append("")

    out.extend(["## Experience", ""])

    for role in data.get("experience") or []:
        title = role.get("title", "")
        company = role.get("company", "")
        out.append(f"### {company} — {title}")
        out.append("")
        meta = []
        if role.get("date"):
            meta.append(f"**Dates:** {role['date']}")
        if role.get("location"):
            meta.append(f"**Location:** {role['location']}")
        if meta:
            out.append(" | ".join(meta))
            out.append("")

        overview = para(role.get("description", ""))
        if overview:
            out.extend(["#### Overview", "", overview, ""])

        teams = role.get("teams") or []
        if teams:
            out.extend(["#### Teams", ""])
            for team in teams:
                name = team.get("name") if isinstance(team, dict) else team
                if name:
                    out.append(f"- {name}")
            out.append("")

        projects = role.get("notable_projects") or []
        if projects:
            out.extend(["#### Notable work", ""])
            for project in projects:
                pname = project.get("name", "Project")
                details = para(project.get("details", ""))
                out.append(f"##### {pname}")
                out.append("")
                if details:
                    out.append(details)
                    out.append("")

    out.extend(["## Education", ""])
    for edu in data.get("education") or []:
        parts = [
            edu.get("degree"),
            edu.get("school"),
            edu.get("location"),
            edu.get("date"),
        ]
        parts = [p for p in parts if p]
        out.append(f"- {' | '.join(parts)}")

    out.extend(["", "## Certifications", ""])
    for cert in data.get("certifications") or []:
        out.append(f"- {cert}")

    out.append("")
    return "\n".join(out)


def main() -> int:
    root = Path(__file__).resolve().parents[1]
    src = root / "verboseResume.yaml"
    dst = root / "verboseResume.md"
    if len(sys.argv) >= 2:
        src = Path(sys.argv[1])
    if len(sys.argv) >= 3:
        dst = Path(sys.argv[2])

    raw = src.read_text(encoding="utf-8")
    data = yaml.safe_load(re.sub(r"^#.*\n", "", raw, count=1))
    dst.write_text(convert(data), encoding="utf-8")
    print(f"Wrote {dst}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
