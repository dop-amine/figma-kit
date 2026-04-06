# Naming Conventions

Standardized naming rules for Figma frames, pages, and nodes.

## Page Names

Use numbered prefixes for ordering:
```
1. Rationale & Visual Concept
2. One-Pager & LinkedIn Carousel
3. Storyboard & Styleframes
```

## Frame Names

### Slides (Carousel)
```
Slide 1 - Cover
Slide 2 - Problem
Slide 3 - Pivot
Slide 4 - Product Hero
Slide 5 - Capabilities
Slide 6 - Social Proof
Slide 7 - CTA
```

### One-Pagers
```
One-Pager Print (EN)
One-Pager Digital (EN)
One-Pager Print (ES)
```

### Storyboard
```
SF 01 — El problema
SF 02 — La transición
SF 03 — El mundo real
SF 04 — El producto
SF 05 — El cierre
Storyboard Panorámico — Teaser 12s
```

### General Frames
```
Rationale
Visual Concept
```

## Patterns

- Use sentence case for frame names
- Include language code `(EN)` or `(ES)` when bilingual
- Use em dash ` — ` (not hyphen) for scene separators in storyboards
- Zero-pad scene numbers: `SF 01`, `Slide 1`
- Include descriptive suffix after dash: `Slide 3 - Pivot`

## Slide Positioning

Slides within a page should be laid out horizontally:
```
x = startX + index * (slideWidth + gap)
y = 0
```

Default values:
- LinkedIn slides: startX=1500, gap=60, w=1080, h=1350
- Styleframes: startX=0, gap=120, w=1920, h=1080
- One-pagers: x=0, y varies (print at y=0, digital at y=1700)

## Color Variable Names

In theme configs, use short uppercase keys:
| Key | Meaning |
|-----|---------|
| BG | Background |
| CARD | Card surface |
| WT | White / primary text |
| TH | Text heading (light mode) |
| BD | Body text |
| MT | Muted text / labels |
| BL | Blue accent |
| TL | Teal accent |
| STK | Stroke / border |
| RL | Rule line (light mode) |
| WARN | Warning orange |
| ERR | Error red |
| SUCCESS | Success green |
