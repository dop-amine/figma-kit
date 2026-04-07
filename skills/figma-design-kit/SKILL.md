# figma-kit — Cursor Skill for Programmatic Figma Design

## Compose First

**ALWAYS use `figma-kit compose` when building 2+ elements.** Individual commands are for single tweaks only.

Compose batches N commands into 1 JS payload → 1 `use_figma` call → 1 round-trip. Shared preamble emitted once, each command body scope-isolated in `{ }`, all node IDs collected automatically.

Compose **tree-shakes** shared helpers into minimal bundles (roughly **14KB → ~3KB** payloads). Behavior for agents is unchanged; payloads are just smaller.

The generated script declares `const _results = [];` and, after each step, pushes that step’s main node. Later steps can reference **`_results[0]`**, **`_results[1]`**, and so on in **`--parent`** flags to nest under prior outputs.

### CLI args

```bash
figma-kit compose -t noir \
  "ui hero --title 'Ship Faster' --cta 'Start'" \
  "card glass --title 'Feature 1'" \
  "card glass --title 'Feature 2'" \
  "ui footer"
```

### Chaining with `_results[]`

```bash
figma-kit compose -t noir \
  "ui section --title Features --label FEATURES" \
  "card glass --parent _results[0] --title 'Feature 1'"
```

### `fx --last` in compose

All **`fx`** commands accept **`--last`** to target the previous step’s result (no manual index):

```bash
figma-kit compose -t noir \
  "card glass --title Hero" \
  "fx noise --last" \
  "fx glow --last --position center"
```

**`--parent` in compose**: Key **`ui`** and **`card`** commands accept **`--parent`** for nesting—use a concrete node ID, **`_results[n]`**, or (for **`fx`**) **`--last`**.

### Recipe YAML

```bash
figma-kit compose --recipe landing.yml
```

```yaml
# landing.yml
theme: noir
page: 0
steps:
  - "preamble"
  - "ui hero --title 'Ship Faster' --cta 'Get Started'"
  - "card glass --title 'Speed' --desc 'Sub-ms reads'"
  - "card glass --title 'Scale' --desc 'Planet-scale'"
  - "card glass --title 'Sync' --desc 'Real-time sync'"
  - "ui pricing --tiers '[{\"name\":\"Pro\",\"price\":\"$29\",\"highlighted\":true}]'"
  - "ui footer"
```

### Execution

**AI agent**: capture compose stdout → feed to `use_figma` once.

**Direct**: `figma-kit exec compose -t noir --recipe landing.yml` — generates JS and sends to Figma MCP in one shot.

### Images in compose

Local images (< 33KB) embed as base64 inline. Include `image place` steps in your compose:

```bash
figma-kit compose -t noir \
  "ui hero --title 'Launch'" \
  "image place ./logo.png --name Logo --width 200 --height 60" \
  "card glass --title 'Feature'"
```

For larger files, start `figma-kit image serve ./assets` first (separate process), then use the URL:

```bash
figma-kit compose -t noir \
  "image place http://127.0.0.1:8741/hero.jpg --width 1440 --height 900" \
  "ui hero --title 'Launch'"
```

### When to use individual commands

- Single property tweak: `figma-kit style fill <id> --solid "#FF0000"`
- One-off inspection: `figma-kit inspect <id>`
- Non-composable operations: `theme init`, `auth`, `config`, `export tokens`, `image serve`

## Rate Limits

| MCP Tool | Limit | Notes |
|----------|-------|-------|
| `use_figma` | **Unlimited** (beta) | Write freely. Batch via compose. |
| `get_screenshot` | 200–600/day | Call ONCE at the end, not between steps. |
| `get_metadata`, `get_design_context`, etc. | 200–600/day | Minimize reads. |

**Strategy**: maximize writes per `use_figma` call with compose. Call `get_screenshot` once after all mutations are done to verify the final result. Never screenshot between compose steps.

## Verification

After composing and executing:

1. Call `get_screenshot` **once** to verify the final result.
2. If something is off, fix with a targeted individual command or a small compose, then screenshot again.
3. Do NOT screenshot after every step — it wastes rate-limited reads.

## Authentication

Several paths, from zero-config to fully manual:

| Tier | Method | Setup |
|------|--------|-------|
| **AI agent** | OAuth (MCP) | Cursor/Claude Code handle OAuth transparently. Nothing to configure. |
| **Cached token** | Browser OAuth after `auth login` | Token stored at `~/.config/figma-kit/token.json` for reuse. |
| **REST / PAT (optional)** | Personal Access Token | Set `FIGMA_ACCESS_TOKEN` for `exec`, scripts, or REST-style calls without going through MCP OAuth. |
| **PAT bootstrap** | `FIGMA_TOKEN` + `auth login` | PAT can register the OAuth client; browser flow still caches token at `~/.config/figma-kit/token.json`. |

Use **OAuth + cached token** for interactive MCP workflows; use **PAT / `FIGMA_ACCESS_TOKEN`** when you need a portable token for automation or direct API access.

## Prerequisites

1. **figma-kit binary** in PATH — verify: `figma-kit --version`
2. **Figma MCP server** in `.cursor/mcp.json`:
   ```json
   { "mcpServers": { "figma": { "url": "https://mcp.figma.com/mcp" } } }
   ```
3. A Figma file open in the browser with Dev Mode enabled.

## Core Workflow

```
User prompt → AI picks commands → figma-kit compose → 1 JS payload → use_figma → get_screenshot (once)
```

For multi-element designs:

```bash
# 1. Create theme (if needed — not composable, run separately)
figma-kit theme init --name "Brand" --bg "#0A2540" --primary "#635BFF" --accent "#00D4AA" -o themes/brand.json

# 2. Compose everything into one call
figma-kit compose -t brand \
  "preamble" \
  "ui hero --title 'Title' --cta 'Start'" \
  "card glass --title 'Feature 1'" \
  "card glass --title 'Feature 2'" \
  "ui footer"
# → feed output to use_figma

# 3. Verify once
# → call get_screenshot
```

For single tweaks:

```bash
figma-kit style fill <id> --solid "#FF0000"
# → feed to use_figma
```

## Command Layers

See [docs/STANDARDS.md](../../docs/STANDARDS.md) for the canonical architecture, composable command contract, and JS generation rules.

### Layer 0 — File & Session

| Command | Description | Composable |
|---------|-------------|:----------:|
| `init [name]` | Create `.figmarc.json` project config | No |
| `config set\|get\|list` | Manage config (fileKey, theme, page) | No |
| `auth login\|logout\|status` | Manage Figma MCP authentication | No |
| `exec <command>` | Generate JS + send to MCP directly | No |
| `new-file <name>` | Create a new Figma file via MCP | No |
| `whoami` | Check Figma identity | No |
| `open` | Open current file in browser | No |
| `status` | Inspect pages, frames, node counts | Yes |

### Layer 1 — Primitives

All composable. Use inside `compose` for multi-element designs.

**Node CRUD** — `node <verb>`: `create`, `clone`, `delete`, `move`, `resize`, `rename`, `reparent`, `lock`, `visible`, `boolean`, `svg`, `variant-set`

**Styling** — `style <verb>`: `fill`, `stroke`, `effect`, `corner`, `blend`, `gradient`, `clip`

**Text** — `text <verb>`: `create`, `edit`, `style`, `range`, `list-fonts`, `load-fonts`. For **`text create`**, use **`--line-height`**, **`--letter-spacing`**, **`--align`**, and **`--auto-resize`** alongside existing flags.

**Layout** — `layout <verb>`: `auto`, `grid`, `constraints`, `sizing`, `align`, `distribute`

### Layer 2 — Design Patterns

All composable. These are the primary compose building blocks.

**Cards** — `card <type>`: `glass`, `solid`, `gradient`, `image`, `bento`, `neumorphic`, `clay`, `outline`

**UI Components** — `ui <component>`: `button`, `input`, `badge`, `avatar`, `divider`, `icon`, `progress`, `toggle`, `tooltip`, `stat`, `table`, `nav`, `footer`, `checkbox`, `radio`, `tabs`, `dropdown`, `breadcrumb`, `skeleton`, `chip`, `toast`, `modal`, `card-list`, `sidebar`, `avatar-group`, `rating`, `search`, `pagination`, `color-picker`, `section`, `hero`, `pricing`, `feature-grid`, `testimonial`, `timeline`, `stepper`, `accordion`

**`ui section`** — Recommended **section wrapper** for pages built with compose. Flags: **`--title`**, **`--label`**, **`--subtitle`**, **`--label-color`**, **`--width`**, **`--padding`**, **`--spacing`**, **`--divider`**, **`--parent`**.

**`ui stat`** and **`ui badge`** accept **`--items`** for **batch** creation (multiple stats/badges in one step).

Key **`ui`** and **`card`** commands accept **`--parent`** for compose chaining wherever nesting applies.

**Effects** — `fx <effect>`: `glow`, `mesh`, `noise`, `vignette`, `grain`, `blur-bg`, `accent-bar`, `shadow`, `parallax-layer`, `aurora`, `morph`, `gradient-border`, `spotlight`, `pattern`

**Images** — `image <action>`: `place` (composable), `fill` (composable), `serve` (not composable — starts HTTP server)

### Layer 3 — Deliverables

`make <deliverable>` — composable, generates complete production designs.

Marketing: `carousel`, `instagram-post`, `instagram-story`, `twitter-card`, `facebook-cover`, `youtube-thumb`, `og-image`, `banner`, `email-header`, `ad-set`

Business: `one-pager`, `pitch-deck`, `case-study`, `proposal`, `invoice`, `business-card`, `letterhead`, `contract`

Motion: `storyboard`, `styleframe`, `animatic`, `transition-spec`

UI/UX: `wireframe`, `screen`, `dashboard`, `form`, `modal`, `empty-state`, `error-page`, `onboarding`, `settings`

Print: `poster`, `brochure`, `packaging`, `signage`, `menu`

Meta: `changelog`

Many accept `--content <file.yml>` for data-driven generation.

### Layer 4 — Design System

`ds <verb>`: `create`, `colors`, `type-scale`, `spacing`, `elevation`, `radius`, `icons`, `component`, `component-sheet`, `variables`, `variables-create`, `search`, `import`, `sync-tokens`, `audit`, `tokens`

`ds create` and `ds variables-create` are composable. Search/import/audit are MCP-only (not composable).

### Layer 5 — Inspect & QA

| Command | Composable |
|---------|:----------:|
| `inspect <id>` | Yes |
| `screenshot <id>` | No (MCP read) |
| `tree` | Yes |
| `find --name "..."` | Yes |
| `measure <a> <b>` | Yes |
| `diff <a> <b>` | Yes |
| `qa <check>` | Yes |

### Layer 6 — Export & Handoff

**Export**: `export png|svg|pdf|page|sprites|tokens` — `tokens` is not composable (Go data output).

**Handoff**: `handoff spec|redline|css|react|assets`

### Layer 7 — Compose & Orchestration

| Command | Description |
|---------|-------------|
| `compose` | Batch N commands → 1 JS payload (this is the primary workflow) |
| `batch recipe.yml` | Legacy YAML orchestration (use compose --recipe instead) |

## Theme System

Three built-in themes via `-t` flag:

| Theme | Description |
|-------|-------------|
| `default` | Dark tech/SaaS. Blue-teal accents. |
| `light` | Light mode for print-friendly work. |
| `noir` | Noir premium. Primary blue `#3366FF`. |

Custom themes: `figma-kit theme init` → place JSON in `~/.config/figma-kit/themes/` or `./themes/`.

**Fonts**: Compose and standalone commands load fonts from the **active theme’s fonts spec** (not hardcoded Inter/Geist Mono). Ensure the theme defines the families you need; use a **`preamble`** compose step when you need upfront font loading.

## Workflow Patterns

### Landing Page (compose)

```bash
figma-kit theme init --name "Brand" --bg "#0A2540" --primary "#635BFF" --accent "#00D4AA" -o themes/brand.json

figma-kit compose -t brand \
  "preamble" \
  "ui hero --title 'Ship Faster' --subtitle 'The modern platform' --cta 'Get Started'" \
  "card glass --title 'Speed' --desc 'Sub-millisecond reads'" \
  "card glass --title 'Scale' --desc 'Planet-scale infra'" \
  "card glass --title 'Sync' --desc 'Real-time everywhere'" \
  "ui pricing --tiers '[{\"name\":\"Free\",\"price\":\"$0\"},{\"name\":\"Pro\",\"price\":\"$29\",\"highlighted\":true}]'" \
  "ui footer"
# → single use_figma call, then get_screenshot once
```

### Pitch Deck (recipe)

```yaml
# deck.yml
theme: brand
steps:
  - "make carousel --content deck-content.yml"
```

```bash
figma-kit compose --recipe deck.yml
```

### Design System

```bash
figma-kit compose -t brand \
  "ds create" \
  "ds variables-create"
# then separately: figma-kit export tokens -t brand --format css
```

### Reference-Driven (from URL/screenshot/brand guide)

1. Extract 3 colors: background, primary, accent
2. `figma-kit theme init --bg ... --primary ... --accent ... -o themes/ref.json`
3. `figma-kit compose -t ref "preamble" "ui hero ..." "card glass ..." ...`
4. `get_screenshot` once to verify

## MCP Tool Routing

| figma-kit output | MCP tool |
|------------------|----------|
| JavaScript code | `use_figma` |
| Screenshot requests | `get_screenshot` (rate limited — use sparingly) |
| Design system search | `search_design_system` |
| File creation | `create_new_file` |
| Plain text (themes, tokens) | Display directly |

## Error Recovery

- **Font not loaded**: Use `preamble` as a compose step — it loads all fonts.
- **Node not found**: `figma-kit tree` to list nodes and get IDs.
- **Theme not found**: `figma-kit themes` to list available themes.
- **use_figma fails**: Ensure JS is a single top-level async expression (figma-kit handles this).
- **Payload too large**: Split compose into multiple calls (50,000 char limit per `use_figma`).

## Global Flags

| Flag | Description |
|------|-------------|
| `-t, --theme` | Theme name or path |
| `-p, --page` | Page index (0-based) |
| `-v, --version` | Show version |
