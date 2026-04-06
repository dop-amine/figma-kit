# Contributing Themes

You don't need to be a developer to contribute a theme to figma-kit. If you can pick colors, you can contribute.

> **See it in action:** Run `figma-kit cookbook match-a-website` to see how AI creates themes from references, or browse the [Prompt Cookbook](docs/COOKBOOK.md).

## Quick Start

1. **Create your theme** using one of these methods:

   - **Ask your AI agent** (fastest): Paste a screenshot or URL into Cursor or Claude Code and say: *"Create a figma-kit theme matching this website"* or *"Create a dark figma-kit theme with orange accents."* The AI extracts colors and runs `figma-kit theme init` for you.
   - Use the [web theme builder](https://dop-amine.github.io/figma-kit/theme-builder.html) (visual, no install needed)
   - Run `figma-kit theme init --name "My Theme" --bg "#0A1628" --primary "#2196F3" --accent "#00BCD4"`
   - Copy the template below and edit the colors

2. **Save the file** as `assets/themes/community/your-theme-name.json`

3. **Validate** your theme: `figma-kit validate theme assets/themes/community/your-theme-name.json`

4. **Open a pull request** on GitHub adding your single JSON file

That's it. CI will validate the theme automatically.

## Minimal Template

Copy this, change the colors, and save it. Hex strings are the recommended format:

```json
{
  "name": "My Theme",
  "description": "One sentence describing your theme's aesthetic.",

  "colors": {
    "BG":      "#0D0F17",
    "CARD":    "#161A25",
    "WT":      "#F5F7FA",
    "BD":      "#C7CFD9",
    "MT":      "#737A8C",
    "BL":      "#3366FF",
    "TL":      "#14B8A6",
    "STK":     "#242938",
    "WARN":    "#FF9933",
    "ERR":     "#FF5959",
    "SUCCESS": "#33CC80"
  }
}
```

Both hex strings (`"#3366FF"`) and Figma 0-1 objects (`{ "r": 0.20, "g": 0.40, "b": 1.00 }`) are accepted.

## Color Token Reference

| Token | Purpose | Typical use |
|-------|---------|-------------|
| `BG` | Page background | Darkest color in your palette |
| `CARD` | Card / panel background | Slightly lighter than BG |
| `CARD2` | Second elevation surface | Slightly lighter than CARD |
| `WT` | Primary text | Near-white for dark themes, near-black for light |
| `BD` | Secondary text | Slightly muted from WT |
| `MT` | Muted text | Captions, labels, hints |
| `BL` | Primary accent | Buttons, links, focus rings |
| `TL` | Secondary accent | Badges, code highlights, teal by convention |
| `LINK` | Link color | Defaults to BL |
| `STK` | Stroke / border | Subtle dividers and card borders |
| `HOVER` | Interactive highlight | Subtle highlight on hover |
| `WARN` | Warning | Amber/orange for caution states |
| `ERR` | Error | Red for errors and destructive actions |
| `SUCCESS` | Success | Green for confirmations |

## Optional Sections

A minimal theme only needs `name` and `colors`. For a complete theme, you can also include:

- **`type`** — Typography scale (h1-h4, body, small, label, mono)
- **`fonts`** — Font families (heading, body, mono)
- **`effects`** — Glass and shadow presets
- **`spacing`** — Page, card, and slide padding/gap
- **`gradients`** — Named gradient paints
- **`brand`** — Logo, tagline, and marketing metadata

See `assets/themes/noir.json` for a complete example with all sections.

## Guidelines

- **File naming**: Use lowercase kebab-case: `ocean-breeze.json`, `sunset-warm.json`
- **Unique name**: The `"name"` field should not conflict with existing themes
- **Test visually**: Run `figma-kit theme preview -t community/your-name` to see it in Figma
- **Keep it focused**: One cohesive aesthetic per theme
- **No offensive content**: Theme names and descriptions should be appropriate for all audiences

## Full `theme init` flags

For fine-grained control:

```bash
figma-kit theme init \
  --name "My Theme" \
  --bg "#0A1628" --primary "#2196F3" --accent "#00BCD4" \
  --font-heading "Poppins" --font-body "DM Sans" --font-mono "JetBrains Mono" \
  --warn "#FFAA00" --error "#FF4444" --success "#22CC66" \
  --spacing compact \
  -o assets/themes/community/my-theme.json
```

Or extend an existing theme: `--from assets/themes/community/ocean.json`

## Using the Web Theme Builder

Visit [dop-amine.github.io/figma-kit/theme-builder.html](https://dop-amine.github.io/figma-kit/theme-builder.html) to:

1. Pick your background, primary, and accent colors visually
2. Customize fonts, status colors, spacing, and brand info in the Advanced panel
3. See a live preview of how your theme will look
4. Download the ready-to-use JSON file
5. Submit it as a PR to this repo
