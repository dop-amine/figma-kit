# Theme system

figma-kit themes are JSON files that drive colors, typography, spacing, effects, and gradients for generated Figma plugin JavaScript. The CLI loads a theme for almost every command so the output stays visually consistent.

## Theme JSON schema

Top-level object:

| Field | Type | Description |
| --- | --- | --- |
| `name` | string | **Required.** Display name (also validated on load). |
| `description` | string | Human-readable summary; defaults to `name` if omitted. |
| `website` | string | Optional URL (e.g. brand site). |
| `colors` | object | **Required.** Map of token name → `{ "r", "g", "b" }` in **0–1** (Figma RGB). |
| `type` | object | Map of style key → typography preset (see below). |
| `fonts` | object | `heading`, `body`, `mono` (strings), `weights` (array of style names like `"Semi Bold"`). |
| `effects` | object | `glass` and `shadow` maps (see below). |
| `spacing` | object | Presets: `page`, `card`, `slide`, `frame16`, `letter` — each may include `width`, `height`, `padding`, `margin`, `gap` as needed. |
| `gradients` | object | Named gradients; each has `type`, `gradientTransform`, `gradientStops` (Figma-compatible paint shape). |
| `brand` | object | Optional brand metadata: `primary`, `logo`, `logoFull`, `tagline`, `url`, `clients`, `product`, `features`, etc. |

### `type` entries

Each key maps to an object used as `theme.type.<key>` in templates:

| Field | Type | Description |
| --- | --- | --- |
| `fontSize` | number | Size in px. |
| `style` | string | Figma font style (e.g. `"Bold"`, `"Regular"`). |
| `lineHeight` | number or null | Line height in px; may be `null`. |
| `family` | string | Optional override (e.g. mono family). |

**Conventional keys** include `hero`, `h1`, `h2`, `h3`, `h4`, `body`, `small`, `label`, `code`, and `mono`. Built-in themes ship a subset (typically `h1`–`h4`, `body`, `small`, `label`, `mono`). You may add keys such as `hero` or `code` for custom templates; unused keys are ignored until referenced.

### `effects.glass`

Named presets (e.g. `subtle`, `default`, `strong`, `pill`). Each preset:

| Field | Description |
| --- | --- |
| `r` | Corner radius. |
| `f` | Fill opacity-related parameter. |
| `s` | Stroke-related parameter. |
| `ga` | Glass / overlay strength. |
| `bl` | Background blur. |

### `effects.shadow`

Named Figma-style shadows. Each preset includes `type`, `color` (`r`,`g`,`b`,`a`), `offset` (`x`,`y`), `radius`, `spread`, `visible`, `blendMode`.

## Built-in themes

Embedded in the binary (`assets/themes/*.json`):

| Key | Name / intent |
| --- | --- |
| `default` | Dark tech / SaaS baseline. Blue-teal accents on a near-black background (`rgb(13,15,23)`-class). |
| `light` | Print-friendly light mode: light greys, white cards; includes extra tokens such as `TH` and `RL` not present in all themes. |
| `noir` | Brand-forward dark theme with primary blue **#3366FF** (`brand.primary`), optional `website`, and `brand` content for deliverables. |

## Community themes

Community-contributed themes ship bundled in `assets/themes/community/`. Anyone can contribute a theme — see [CONTRIBUTING-THEMES.md](../CONTRIBUTING-THEMES.md).

| Key | Name / intent |
| --- | --- |
| `ocean` | Deep blue maritime palette. Navy backgrounds, aqua accents. |

`figma-kit themes` lists all themes grouped by source (built-in, community, user, local).

### Hex strings in JSON

Colors can be specified as hex strings instead of `{r, g, b}` objects — this is the recommended format for AI agents:

```json
{
  "colors": {
    "BG": "#0D0F17",
    "BL": "#3366FF",
    "TL": "#14B8A6"
  }
}
```

Both formats are accepted. Hex is preferred for readability and because AI agents extract hex from websites and images, not 0-1 floats.

## Color tokens

Themes use short names as **JavaScript identifiers** in generated code (`const BG = { r, g, b };`). Common tokens:

| Token | Typical role |
| --- | --- |
| `BG` | Page / canvas background |
| `CARD` | Elevated surfaces |
| `WT` | Primary light text on dark |
| `BD` | Body / secondary text |
| `MT` | Muted text |
| `STK` | Strokes / dividers |
| `BL` | Brand / link accent |
| `SUCCESS`, `ERR`, `WARN` | Semantic states |
| `TL` | Teal or secondary accent |
| `AC` | Accent (conventional; not all built-in themes define it) |
| `LINK` | Link color (defaults to BL) |
| `CARD2` | Second elevation surface |
| `HOVER` | Subtle interactive highlight |

**Not every theme defines every token.** Light adds `TH` (titles on light backgrounds) and `RL`; your templates or YAML should match the active theme or you should extend the JSON.

## Using themes

### CLI flag

```bash
figma-kit preamble -t noir
figma-kit make carousel -t light --content slides.yml
```

`-t` / `--theme` is a persistent flag on the root command.

### `.figmarc.json`

In the project directory (or a parent), set defaults:

```json
{
  "theme": "noir",
  "page": 0
}
```

Use `figma-kit init` and `figma-kit config set theme <name>` to manage this file.

### Resolution order

When a command needs a theme, `resolveTheme` applies:

1. **`-t` / `--theme` if non-empty** — theme name passed on the CLI.
2. **Else `theme` from `.figmarc.json`** (via `config.Load()`).
3. **Else `"default"`** — then `theme.Load` runs.

`theme.Load` search order for the **name**:

1. Embedded built-ins (`default`, `light`, `noir`).
2. Bundled community themes (`assets/themes/community/`).
3. `~/.config/figma-kit/themes/<name>.json` (OS user config dir).
4. `./themes/<name>.json` relative to the current working directory.

To use a one-off file path from the shell, place or symlink it under those directories with a stable `<name>.json`, or extend the CLI to call `theme.LoadFile` (Go API today: `theme.LoadFile(path)`).

## AI-Driven Theme Creation

The recommended workflow uses an AI agent to create themes from visual references:

```
1. Paste screenshot or URL into Cursor / Claude Code
2. AI extracts dominant colors (background, primary, accent)
3. AI runs: figma-kit theme init --name "Brand" --bg "#1a1a2e" --primary "#e94560" --accent "#0f3460"
4. AI runs: figma-kit theme preview -t "Brand"
5. Verify the preview in Figma
6. AI uses the theme to build: figma-kit make landing -t "Brand"
```

### Example prompts for your AI agent

```
"Create a figma-kit theme matching this website: https://stripe.com"

"Look at this screenshot and create a figma-kit theme from the dominant colors. Use Poppins for headings."

"Create a dark, warm figma-kit theme with orange accents. Then build a landing page with it."
```

### Example AI session

```
User: Create a figma-kit theme inspired by the Linear app, then build a feature showcase.

AI: I'll extract Linear's design language and create a matching theme.

$ figma-kit theme init \
    --name "Linear Dark" \
    --bg "#1B1B25" \
    --primary "#5E6AD2" \
    --accent "#26B5CE" \
    --font-heading "Inter" \
    --font-body "Inter" \
    --spacing compact \
    -o themes/linear-dark.json

Theme written to themes/linear-dark.json

$ figma-kit theme preview -t themes/linear-dark.json
# → Preview page appears in Figma

$ figma-kit make landing -t themes/linear-dark.json --title "Issue Tracking" --subtitle "Reimagined"
# → Full landing page with Linear-style colors
```

figma-kit does not do the AI extraction itself — the AI agent already can parse images and websites. figma-kit provides the **perfect API** for that agent: hex input everywhere, rich flags, clear output.

### Every way to create a theme

| Entry point | Who does the extraction? | What you say / do | figma-kit command |
|---|---|---|---|
| **Website URL** | AI reads the live site, extracts colors and fonts | *"Create a figma-kit theme matching stripe.com"* | `theme init --name "Stripe" --bg "#0A2540" --primary "#635BFF" --accent "#00D4AA" --font-heading "Sohne"` |
| **Screenshot / image** | AI analyzes pixels, picks dominant palette | *"Create a figma-kit theme from the dominant colors in this screenshot"* | `theme init --bg "#..." --primary "#..." --accent "#..."` |
| **Photoshop / PSD** | AI reads layer styles, color swatches, type | *"Use the colors and fonts from this PSD to create a figma-kit theme"* | `theme init --bg "#..." --primary "#..." --font-heading "Montserrat"` |
| **Mood / description** | AI interprets the aesthetic, chooses colors | *"Create a warm, minimal figma-kit theme with earth tones"* | `theme init --name "Earth" --bg "#2C2418" --primary "#D4956A" --accent "#8B9E6B"` |
| **Existing theme** | You specify overrides on a base | `--from ocean.json --bg "#F8F9FA"` | `theme init --from themes/ocean.json --name "Ocean Light" --bg "#F8F9FA"` |
| **Web builder** | You pick colors visually | Visit [theme-builder.html](https://dop-amine.github.io/figma-kit/theme-builder.html) | Download JSON, place in `themes/` |
| **Manual JSON** | You write it by hand | Copy a built-in theme and edit hex values | Place in `~/.config/figma-kit/themes/` or `./themes/` |

### What you can customize

Everything `theme init` (and the web builder) can control:

| Category | Flags | What it affects |
|---|---|---|
| **Seed colors** | `--bg`, `--primary`, `--accent` | Derives 14 tokens: BG, CARD, CARD2, WT, BD, MT, BL, TL, LINK, STK, HOVER, WARN, ERR, SUCCESS |
| **Fonts** | `--font-heading`, `--font-body`, `--font-mono` | Typography across all type scale entries; used in theme preview and templates |
| **Status colors** | `--warn`, `--error`, `--success` | Override the auto-derived status colors |
| **Spacing** | `--spacing compact` or `--spacing spacious` | Page padding, card padding, slide margins — affects `make` layout templates |
| **Base theme** | `--from <path>` | Start from any existing theme and override specific values |
| **Brand** | (JSON only) `brand.tagline`, `brand.url` | Shown in theme preview; available to templates |

### Architectural limits

These are intentional design choices:

- **figma-kit does not fetch URLs or parse images** — the AI agent (Cursor, Claude Code, etc.) does all visual perception. figma-kit is the execution API.
- **figma-kit does not read .psd / .sketch / .ai files** — but AI agents like Claude can analyze screenshots of those files and extract colors.
- **Font availability depends on Figma** — the font must be available in Figma (Google Fonts, uploaded fonts) for the generated JS to work.
- **Effects use fixed presets** — glass and shadow effects have named presets (subtle/default/strong), not arbitrary CSS.
- **`theme preview` is a static Figma frame** — not an interactive component.

## Creating themes

There are four ways to create a custom theme:

### 1. Web Theme Builder (recommended for designers)

Visit [dop-amine.github.io/figma-kit/theme-builder.html](https://dop-amine.github.io/figma-kit/theme-builder.html) to pick colors visually, see a live preview, and download ready-to-use JSON.

### 2. AI agent (let your AI create it)

Paste a screenshot, URL, or color description into Cursor or Claude Code and ask it to create a theme. See [AI-Driven Theme Creation](#ai-driven-theme-creation) above.

### 3. `figma-kit theme init` (CLI)

```bash
# Basic: 3 colors
figma-kit theme init --name "Ocean" --bg "#0A1628" --primary "#2196F3" --accent "#00BCD4" -o themes/ocean.json

# Full: custom fonts, status colors, compact spacing
figma-kit theme init \
  --name "Brand Kit" \
  --bg "#1a1a2e" --primary "#e94560" --accent "#0f3460" \
  --font-heading "Poppins" --font-body "DM Sans" --font-mono "JetBrains Mono" \
  --warn "#FFAA00" --error "#FF4444" --success "#22CC66" \
  --spacing compact \
  -o themes/brand.json

# Extend an existing theme
figma-kit theme init --from themes/brand.json --name "Brand Light" --bg "#F8F9FA" -o themes/brand-light.json
```

Available flags: `--name`, `--desc`, `--bg`, `--primary`, `--accent`, `--font-heading`, `--font-body`, `--font-mono`, `--warn`, `--error`, `--success`, `--spacing` (compact/spacious), `--from` (base theme), `-o` (output path).

Run with no flags to print a starter template.

### 4. Copy and edit

1. Copy `assets/themes/default.json` (or `light.json`) as a starting point.
2. Save as either:
   - `~/.config/figma-kit/themes/<mytheme>.json`, or
   - `./themes/<mytheme>.json` in your repo.
3. Run with `-t mytheme` (no `.json` suffix).

### Previewing themes in Figma

```bash
figma-kit theme preview -t mytheme
```

Generates `use_figma` JS that creates a compact preview page in Figma with color swatches, type scale specimens, and sample components.

### Validation

```bash
figma-kit validate theme themes/mytheme.json
```

`name` must be non-empty; `colors` must be non-empty. For deliverables (`make carousel`, `one-pager`, etc.), keep `type`, `fonts`, `effects`, `spacing`, and `gradients` aligned with the templates you use—dropping sections may break generated JS that expects keys like `theme.effects.glass.strong` or `theme.spacing.slide`.

### Minimal custom theme

Smallest JSON that still **loads** (enough for e.g. color-only preamble / token export):

```json
{
  "name": "My Brand",
  "description": "Minimal smoke-test theme",
  "colors": {
    "BG":   { "r": 0.1, "g": 0.1, "b": 0.12 },
    "WT":   { "r": 0.95, "g": 0.96, "b": 0.98 },
    "BL":   { "r": 0.2, "g": 0.4, "b": 1.0 }
  }
}
```

For real work, start from a full built-in file and edit values.

## Exporting tokens (no plugin)

Theme tokens can be exported from Go without running Figma:

```bash
figma-kit export tokens --format json
figma-kit export tokens --format css -t noir
```

- **`json`**: Full theme object as JSON (default).
- **`css`**: CSS custom properties for colors, fonts, typography scale, and spacing:

```css
:root {
  /* Colors */
  --fk-BG: #0D0F17;
  --fk-BL: #3366FF;
  /* ... all color tokens ... */

  /* Fonts */
  --fk-font-heading: 'Inter', sans-serif;
  --fk-font-body: 'Inter', sans-serif;
  --fk-font-mono: 'Geist Mono', monospace;

  /* Typography */
  --fk-h1-size: 72px;
  --fk-h1-lh: 86px;
  --fk-body-size: 16px;
  /* ... all type scale entries ... */

  /* Spacing */
  --fk-page-padding: 80px;
  --fk-card-padding: 24px;
  /* ... */
}
```

Both commands respect `-t` / `.figmarc.json` / `default` the same way as other commands.
