# figma-kit — Cursor Skill for Programmatic Figma Design

## Overview

`figma-kit` is a Go CLI that generates `use_figma`-compatible JavaScript for the Figma MCP server. This skill teaches you how to use it to create, inspect, and audit Figma designs from Cursor.

## Prerequisites

1. **figma-kit binary** in PATH — verify: `figma-kit --version`
2. **Figma MCP server** configured in `.cursor/mcp.json`:
   ```json
   { "mcpServers": { "figma": { "url": "https://mcp.figma.com/mcp" } } }
   ```
3. A Figma file open in the browser with Dev Mode enabled.

## Core Workflow

Every interaction follows the same 3-step loop:

1. **Generate JS** — run a `figma-kit` command to produce use_figma JavaScript
2. **Execute** — feed the JS output to the `use_figma` MCP tool
3. **Verify** — call `get_screenshot` to visually confirm the result

```
User: "Create a hero section"
→ Run: figma-kit node create frame --name "Hero" -w 1440 --height 800
→ Feed output to use_figma
→ get_screenshot to verify
```

## Command Reference

### Layer 0 — File & Session

| Command | Description |
|---------|-------------|
| `figma-kit init [name]` | Create `.figmarc.json` project config |
| `figma-kit config set\|get\|list` | Manage config (fileKey, theme, page, exportDir) |
| `figma-kit whoami` | Check Figma identity (wraps MCP) |
| `figma-kit open` | Open current file in browser |
| `figma-kit status` | JS to inspect pages, frames, node counts |

### Layer 1 — Primitives

**Node CRUD** — `figma-kit node <verb>`
| Verb | Key Flags |
|------|-----------|
| `create <type>` | `--name`, `-w`, `--height`, `--x`, `--y` |
| `clone <id>` | `--dx`, `--dy` |
| `delete <id>` | — |
| `move <id>` | `--x`, `--y` |
| `resize <id>` | `-w`, `--height` |
| `rename <id>` | `--name` |
| `reparent <id> <parentId>` | — |
| `lock <id>` | `--unlock` |
| `visible <id>` | `--hide` |

**Styling** — `figma-kit style <verb>`
| Verb | Key Flags |
|------|-----------|
| `fill <id>` | `--solid "#HEX"`, `--opacity` |
| `stroke <id>` | `--color`, `--weight`, `--align` |
| `effect <id>` | `--shadow`, `--blur`, `--blur-type` |
| `corner <id>` | `--radius` or `--tl/--tr/--br/--bl` |
| `blend <id>` | `--mode`, `--opacity` |
| `gradient <id>` | `--type`, `--angle`, `--stops` |
| `clip <id>` | `--off` |

**Text** — `figma-kit text <verb>`
| Verb | Key Flags |
|------|-----------|
| `create` | `--content`, `--font`, `--weight`, `--size`, `--color`, `--parent` |
| `edit <id>` | `--content` |
| `style <id>` | `--size`, `--lh`, `--ls`, `--align` |
| `range <id>` | `--start`, `--end`, `--weight`, `--color` |
| `list-fonts` | — |
| `load-fonts` | `--families "Inter,Geist Mono"` |

**Layout** — `figma-kit layout <verb>`
| Verb | Key Flags |
|------|-----------|
| `auto <id>` | `--dir`, `--gap`, `--pad`, `--align`, `--wrap` |
| `grid <id>` | `--columns`, `--gutter`, `--margin` |
| `constraints <id>` | `--h`, `--v` |
| `sizing <id>` | `-w`, `--height` |
| `align <id>` | `--primary`, `--counter` |
| `distribute <ids>` | `--axis`, `--gap` |

### Layer 2 — Design Patterns

**Cards** — `figma-kit card <type>`
- `glass` — glassmorphism with presets (subtle, default, strong, pill)
- `solid` — flat card with bg, border, shadow, radius
- `gradient` — gradient fill card
- `image` — image fill with overlay
- `bento` — grid layout of cards

**UI Primitives** — `figma-kit ui <component>`
- `button`, `input`, `badge`, `avatar`, `divider`, `icon`, `progress`, `toggle`, `tooltip`, `stat`, `table`, `nav`, `footer`

**Visual Effects** — `figma-kit fx <effect>`
- `glow`, `mesh`, `noise`, `vignette`, `grain`, `blur-bg`, `accent-bar`, `shadow`, `parallax-layer`

### Layer 3 — Deliverables

`figma-kit make <deliverable>` — generates complete production designs.

**Marketing & Social:** `carousel`, `instagram-post`, `instagram-story`, `twitter-card`, `facebook-cover`, `youtube-thumb`, `og-image`, `banner`, `email-header`, `ad-set`

**Sales & Business:** `one-pager`, `pitch-deck`, `case-study`, `proposal`, `invoice`, `business-card`, `letterhead`, `contract`

**Motion:** `storyboard`, `styleframe`, `animatic`, `transition-spec`

**UI/UX:** `wireframe`, `screen`, `dashboard`, `form`, `modal`, `empty-state`, `error-page`, `onboarding`, `settings`

**Print:** `poster`, `brochure`, `packaging`, `signage`, `menu`

Many accept `--content <file.yml>` for data-driven generation.

### Layer 4 — Design System

`figma-kit ds <verb>` — create, colors, type-scale, spacing, elevation, radius, icons, component, variables, search, import, sync-tokens, audit

### Layer 5 — Inspect & QA

| Command | Description |
|---------|-------------|
| `figma-kit inspect <id>` | Dump node properties |
| `figma-kit screenshot <id>` | Capture via MCP |
| `figma-kit tree` | Hierarchical node tree |
| `figma-kit find --name "..."` | Search nodes by name |
| `figma-kit measure <a> <b>` | Distance between nodes |
| `figma-kit diff <a> <b>` | Compare node properties |
| `figma-kit qa <check>` | contrast, touch-targets, orphans, fonts, colors, spacing, naming, responsive, checklist |

### Layer 6 — Export & Handoff

**Export:** `figma-kit export png\|svg\|pdf\|page\|sprites\|tokens`
**Handoff:** `figma-kit handoff spec\|redline\|css\|react\|assets`

### Layer 7 — Batch Orchestration

```bash
figma-kit batch recipe.yml
```

Recipe format:
```yaml
name: "Q2 Campaign"
steps:
  - title: "Carousel"
    js: |
      // JS from figma-kit make carousel
  - title: "One-pager"
    js: |
      // JS from figma-kit make one-pager
```

## Theme System

Three built-in themes, selectable via `-t` flag or `.figmarc.json`:

| Theme | Description |
|-------|-------------|
| `default` | Dark theme for tech/SaaS. Blue-teal accents. |
| `light` | Light mode for print-friendly deliverables. |
| `noir` | Noir dark premium theme. Primary blue #3366FF. |

Custom themes: place JSON in `~/.config/figma-kit/themes/` or `./themes/`.

Export tokens: `figma-kit export tokens --format css`

## Workflow Patterns

### Creating a Carousel
```bash
figma-kit make carousel --content slides.yml -t noir
# → produces JS that creates all slides in one use_figma call
```

### Building a Component from Scratch
```bash
figma-kit node create frame --name "Card" -w 400 --height 300
# → get the node ID from use_figma response
figma-kit style fill <id> --solid "#1A1C2B"
figma-kit style corner <id> --radius 16
figma-kit layout auto <id> --dir VERTICAL --gap 16 --pad 24
figma-kit text create --content "Title" --size 24 --weight Bold --parent <id>
```

### QA Audit
```bash
figma-kit qa checklist --page 0
# → generates JS that runs all QA checks and returns scored report
```

### Design System Setup
```bash
figma-kit ds create -t noir
# → creates full DS page with swatches, type specimens, spacing scale
```

## Reference-Driven Workflow

When a user shares a URL, screenshot, brand guide, or mood description, follow this sequence:

### 1. Extract colors from the reference

Analyze the image or website. Identify three key colors:
- **Background** — the dominant dark/light surface color
- **Primary** — the main brand/accent color
- **Accent** — a secondary highlight color

### 2. Create a theme

```bash
figma-kit theme init \
  --name "Brand" \
  --bg "#0A2540" \
  --primary "#635BFF" \
  --accent "#00D4AA" \
  --font-heading "Inter" \
  --font-body "Inter" \
  -o themes/brand.json
```

Optional flags: `--font-mono`, `--warn`, `--error`, `--success`, `--spacing compact|spacious`, `--from` (extend existing theme).

### 3. Preview the theme

```bash
figma-kit theme preview -t brand
# → execute via use_figma, then get_screenshot to verify
```

### 4. Build the preamble

```bash
figma-kit preamble -t brand
# → execute via use_figma (sets up colors + fonts on the page)
```

### 5. Compose the design

Sequence commands based on what the user wants:

```bash
# Landing page
figma-kit make screen --type landing --sections "hero,features,pricing,cta" -t brand
figma-kit card glass -t brand --title "Feature 1" --desc "Description"
figma-kit ui button --variant primary -t brand
figma-kit fx mesh <heroId> -t brand

# Pitch deck
figma-kit make pitch-deck --slides 7 --template saas -t brand

# Design system
figma-kit ds create -t brand

# Social assets
figma-kit make og-image --title "Product" --description "Tagline" -t brand
figma-kit make carousel --content slides.yml -t brand
```

### 6. QA and export

```bash
figma-kit qa checklist --page 0
figma-kit export tokens -t brand --format css
```

## Project Templates

Common prompt patterns with recommended command sequences:

### Landing Page
```
preamble → node create frame "Hero" → fx mesh → text create (headline) →
text create (subtitle) → ui button → card glass ×3 → make screen --type pricing →
ui footer → qa checklist
```

### Pitch Deck
```
preamble → make pitch-deck --slides 7 --content deck.yml
```
Or with carousel format: `make carousel --content deck.yml`

### Design System
```
preamble → ds create → ds variables-create → export tokens --format css
```

### Social Campaign
```
preamble → make carousel --content slides.yml → make og-image → make twitter-card →
make instagram-post
```

### Full Project (multiple pages)
```
page create "Landing" → page create "Dashboard" → page create "Design System" →
(switch to each page and build with above patterns)
```

## Prompt Cookbook

Run `figma-kit cookbook` to browse 15 complete prompt-to-design sessions, or see [docs/COOKBOOK.md](../../docs/COOKBOOK.md).

## MCP Tool Routing

| figma-kit output | MCP tool to use |
|------------------|-----------------|
| JavaScript code | `use_figma` |
| Screenshot instructions | `get_screenshot` |
| Search instructions | `search_design_system` |
| React/handoff | `get_design_context` |
| Plain text (themes, info) | Display directly |

## Error Recovery

- **Font not loaded**: Always call `figma-kit preamble` or use `scaffold` which includes font loading
- **Node not found**: Use `figma-kit tree` to list available nodes and get correct IDs
- **Theme not found**: Check available themes with `figma-kit themes`
- **use_figma fails**: Check that the JS is a single top-level async expression (figma-kit handles this)
- **Visual issues**: Always verify with `get_screenshot` after mutations

## Global Flags

| Flag | Description |
|------|-------------|
| `-t, --theme` | Theme name or path |
| `-p, --page` | Page index (0-based) |
| `-v, --version` | Show version |
