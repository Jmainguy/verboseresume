# Verbose Resume — Brand Guide

Verbose Resume is a community tool for keeping one detailed career record and tailoring it into job-ready resumes. The brand should feel precise, calm, and technical—like good infrastructure documentation, not flashy recruiting marketing.

**Domain:** verboseresume.com  
**Product name:** Verbose Resume (two words, title case in UI)  
**Short name / mark:** VR (fallback monogram only when the SVG mark cannot load)

---

## Brand idea

| Concept | Expression |
|--------|------------|
| **Verbose** | Many lines of detail—your long memory |
| **Resume** | One focused output for a specific role |
| **Mark** | Three horizontal lines (detail) converging to a single bold stroke (tailored result) inside a rounded square |

The logo mark is **not** a document icon or a printer. It reads as “lots of input → one clear output.”

---

## Voice & tone

- Direct, plain language; avoid buzzwords and fake enthusiasm.
- Truth-first: never imply features or experience that are not in the source material.
- Community tool energy: helpful, open, engineer-friendly.
- Pair with Harborvane only in the site footer (consulting), not in product copy on every page.

**Tagline (default):** *Keep one detailed career record. Tailor it for every role.*

**One-liner:** Turn structured career JSON into a polished HTML resume you can print to PDF.

---

## Logo

### Files

| File | Use |
|------|-----|
| `static/brand/logo-mark.svg` | Nav icon, app chrome, favicon source |
| `static/brand/logo.svg` | Horizontal lockup (mark + wordmark) |
| `static/brand/favicon.svg` | Modern browsers (`rel="icon"`) |
| `static/brand/favicon-32.png` | Legacy / small raster |
| `static/brand/apple-touch-icon.png` | iOS home screen (180×180) |

Regenerate PNGs after SVG edits:

```bash
python3 scripts/generate_brand_pngs.py
```

### Mark (symbol)

- Rounded square, 12px corner radius at 38×38 (scale proportionally).
- Background: navy gradient `#152a48` → `#08111f`.
- Three cyan lines on the left, decreasing length and opacity (verbose detail).
- One vertical blue/cyan gradient stroke on the right (focused output).

### Wordmark

- **Verbose Resume** set in **Inter** (800 weight), tight tracking (`-0.02em`).
- Light text `#f4f8fb` on dark backgrounds; use `#102033` on light panels if needed.

### Clear space

Keep at least **½× mark height** of empty space on all sides of the mark (e.g. 19px when the mark is 38px).

### Don’t

- Stretch, rotate, or recolor the mark gradients arbitrarily.
- Place the mark on busy photography without a panel behind it.
- Replace the mark with only “VR” text in customer-facing UI when the SVG is available.
- Use the gold accent inside the logo mark (gold is for page atmosphere only).

---

## Color

### Core palette

| Name | Hex | Role |
|------|-----|------|
| Navy Deep | `#08111f` | Page background start |
| Navy Mid | `#101d32` | Background mid-stop |
| Navy Slate | `#152a48` | Background end, mark gradient |
| Primary Blue | `#1684e8` | Buttons, links on light panels |
| Bright Blue | `#46a3ff` | Highlights, mark stroke end |
| Cyan | `#6ee7f9` | Accent glow, detail lines in mark |
| Gold | `#f8c85c` | **Sparingly**—hero atmosphere only |
| Ink | `#102033` | Body text on light surfaces |
| Muted | `#5b667a` | Secondary text |
| Panel | `rgba(255,255,255,0.94)` | Cards, footer bar |
| Panel border | `rgba(198,217,237,0.72)` | Card and footer outlines |
| Link on panel | `#0f67b7` | Footer links on white panel |
| Link hover | `#1684e8` | Interactive hover |

### CSS variables (app shell)

Use these names in site templates for consistency:

```css
:root {
  --vr-navy-deep: #08111f;
  --vr-navy-mid: #101d32;
  --vr-navy-slate: #152a48;
  --vr-blue: #1684e8;
  --vr-blue-bright: #46a3ff;
  --vr-cyan: #6ee7f9;
  --vr-gold: #f8c85c;
  --vr-ink: #102033;
  --vr-muted: #5b667a;
  --vr-panel: rgba(255, 255, 255, 0.94);
  --vr-border: rgba(198, 217, 237, 0.72);
}
```

### Background

Default page background is a **diagonal navy gradient** with soft **cyan** (top-left) and **gold** (top-right) radial glows, plus a subtle grid overlay at ~4% white. Do not flatten to a single solid color—the depth is intentional, but **content and footer sit on opaque panels** so text stays readable.

### Resume paper tones

Print preview paper options stay neutral (white, soft white, cool linen, etc.). Do not apply the navy marketing gradient inside the printable resume area.

---

## Typography

| Context | Font | Notes |
|---------|------|--------|
| App UI | Inter, system-ui stack | 400 body, 700–800 UI labels and headings |
| Resume templates | Per template (sans or serif) | Brand does not force Inter inside PDF output |
| Code / JSON | ui-monospace, monospace | Docs and MCP examples |

**Scale (app):**

- Hero H1: `clamp(42px, 8vw, 82px)`, weight 900, tight letter-spacing
- Section H2: ~1.25–1.5rem, weight 800
- Body: 1rem / line-height 1.5–1.6
- Kicker: 12px, uppercase, letter-spacing `0.12em`, color `#0f67b7`

---

## UI components

- **Panels:** White frosted cards, 28px radius (hero/docs), 20px radius (footer bar).
- **Nav pills:** Semi-transparent white on navy, 999px radius.
- **Primary button:** Blue gradient `#46a3ff` → `#1676d2`, white text, soft shadow.
- **Footer:** Always on opaque panel (not floating text on the gradient).

---

## Favicon & app icons

```html
<link rel="icon" href="/static/brand/favicon.svg" type="image/svg+xml">
<link rel="icon" href="/static/brand/favicon-32.png" sizes="32x32" type="image/png">
<link rel="apple-touch-icon" href="/static/brand/apple-touch-icon.png">
```

Serve from `/static/` (see `main.go`). Prefer SVG for crisp tabs; keep PNG for older clients.

---

## MCP & agent artifacts

- MCP server name: `verboseresume`
- Artifact filenames: neutral (`resume.md`, `resume.html`)—no need to embed logo in generated files.
- Guides may mention the product name in prose; avoid large logo ASCII art in machine-readable outputs.

---

## File checklist for new surfaces

- [ ] Use `logo-mark.svg` or full `logo.svg` from `static/brand/`
- [ ] Link favicon trio in `<head>`
- [ ] Navy gradient + panel surfaces per this doc
- [ ] Inter for chrome; resume body follows template
- [ ] Footer: Jonathan Mainguy + Harborvane only (site chrome partial)

---

## Related docs

- Product workflow: `README.md`, `docs/VERBOSE-RESUME.md`
- Site chrome implementation: `templates/partials/site-chrome.html`
