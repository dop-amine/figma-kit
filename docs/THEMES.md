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

`figma-kit themes` lists embedded themes and descriptions.

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
2. `~/.config/figma-kit/themes/<name>.json` (OS user config dir).
3. `./themes/<name>.json` relative to the current working directory.

To use a one-off file path from the shell, place or symlink it under those directories with a stable `<name>.json`, or extend the CLI to call `theme.LoadFile` (Go API today: `theme.LoadFile(path)`).

## Custom themes

1. Copy `assets/themes/default.json` (or `light.json`) as a starting point.
2. Save as either:
   - `~/.config/figma-kit/themes/<mytheme>.json`, or
   - `./themes/<mytheme>.json` in your repo.
3. Run with `-t mytheme` (no `.json` suffix).

Validation (Go): `name` must be non-empty; `colors` must be non-empty. For deliverables (`make carousel`, `one-pager`, etc.), keep `type`, `fonts`, `effects`, `spacing`, and `gradients` aligned with the templates you use—dropping sections may break generated JS that expects keys like `theme.effects.glass.strong` or `theme.spacing.slide`.

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

Theme tokens can be dumped from Go without running Figma:

```bash
figma-kit export tokens --format json
figma-kit export tokens --format css
```

- **`json`**: Full theme object as JSON (default).
- **`css`**: `:root { --fk-<COLORNAME>: #hex; }` for each color key (sorted).

Both commands respect `-t` / `.figmarc.json` / `default` the same way as other commands.
