# Contributing Themes

You don't need to be a developer to contribute a theme to figma-kit. If you can pick colors, you can contribute.

## Quick Start

1. **Create your theme** using one of these methods:
   - Use the [web theme builder](https://dop-amine.github.io/figma-kit/theme-builder.html) (recommended for designers)
   - Run `figma-kit theme init --name "My Theme" --bg "#0A1628" --primary "#2196F3" --accent "#00BCD4"`
   - Copy the template below and edit the colors

2. **Save the file** as `assets/themes/community/your-theme-name.json`

3. **Validate** your theme: `figma-kit validate theme assets/themes/community/your-theme-name.json`

4. **Open a pull request** on GitHub adding your single JSON file

That's it. CI will validate the theme automatically.

## Minimal Template

Copy this, change the colors, and save it:

```json
{
  "name": "My Theme",
  "description": "One sentence describing your theme's aesthetic.",

  "colors": {
    "BG":      { "r": 0.05, "g": 0.06, "b": 0.09 },
    "CARD":    { "r": 0.09, "g": 0.10, "b": 0.15 },
    "WT":      { "r": 0.96, "g": 0.97, "b": 0.98 },
    "BD":      { "r": 0.78, "g": 0.81, "b": 0.85 },
    "MT":      { "r": 0.45, "g": 0.48, "b": 0.55 },
    "BL":      { "r": 0.20, "g": 0.40, "b": 1.00 },
    "TL":      { "r": 0.08, "g": 0.72, "b": 0.65 },
    "STK":     { "r": 0.14, "g": 0.16, "b": 0.22 },
    "WARN":    { "r": 1.00, "g": 0.60, "b": 0.20 },
    "ERR":     { "r": 1.00, "g": 0.35, "b": 0.35 },
    "SUCCESS": { "r": 0.20, "g": 0.80, "b": 0.50 }
  }
}
```

## Color Token Reference

| Token | Purpose | Typical use |
|-------|---------|-------------|
| `BG` | Page background | Darkest color in your palette |
| `CARD` | Card / panel background | Slightly lighter than BG |
| `WT` | Primary text | Near-white for dark themes, near-black for light |
| `BD` | Secondary text | Slightly muted from WT |
| `MT` | Muted text | Captions, labels, hints |
| `BL` | Primary accent | Buttons, links, focus rings |
| `TL` | Secondary accent | Badges, code highlights, teal by convention |
| `STK` | Stroke / border | Subtle dividers and card borders |
| `WARN` | Warning | Amber/orange for caution states |
| `ERR` | Error | Red for errors and destructive actions |
| `SUCCESS` | Success | Green for confirmations |

Colors use Figma's 0-1 range (not 0-255). To convert hex: divide each channel by 255.

**Example:** `#3366FF` = `{ "r": 0.20, "g": 0.40, "b": 1.00 }`

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

## Using the Web Theme Builder

Visit [dop-amine.github.io/figma-kit/theme-builder.html](https://dop-amine.github.io/figma-kit/theme-builder.html) to:

1. Pick your background, primary, and accent colors visually
2. See a live preview of how your theme will look
3. Download the ready-to-use JSON file
4. Submit it as a PR to this repo
