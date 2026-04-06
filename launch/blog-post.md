# I Built a CLI That Generates Figma Designs from YAML Specs

## The Problem

Every designer knows the drill: open Figma, create a frame, set up auto-layout, pick colors from the brand palette, load fonts, position elements, repeat. For recurring deliverables — social media carousels, pitch decks, OG images — it's the same setup dance every time.

What if you could just describe what you want and have it appear in your Figma file?

## The Tool

**figma-kit** is a Go CLI with 120+ commands that generates production-ready Figma designs from your terminal. Instead of clicking through the Figma UI, you run commands like:

```bash
figma-kit make carousel --content slides.yml -t noir
```

This outputs JavaScript that the official Figma MCP server executes directly in your file. No plugins to install, no browser automation, no CDP hacking — just the official API.

## How It Works

The architecture is intentionally simple:

```
Your Terminal → figma-kit (Go binary) → JavaScript → Figma MCP Server → Your Figma File
```

The Go binary is a sophisticated code generator. It takes your command, loads a theme (design tokens, colors, typography), and composes a JavaScript string compatible with Figma's Plugin API. The JavaScript files aren't implementation code — they're the output format.

When you run `figma-kit make carousel`, here's what happens:

1. The binary loads the `noir` theme (colors, fonts, spacing)
2. It generates font loading code for Inter and Geist Mono
3. It emits theme color constants (`const BL = {r:0.2, g:0.4, b:1}`)
4. It injects helper functions (glass cards, gradient glows, text factories)
5. It outputs the carousel template with your YAML content interpolated
6. The AI feeds this entire JS blob to `use_figma` via MCP
7. Seven perfectly themed carousel slides appear in your Figma file

## The 8 Layers

figma-kit is organized into 8 layers of increasing abstraction:

**Layer 0 — Session:** `init`, `config`, `open`, `status`
**Layer 1 — Primitives:** Node CRUD, styling, text, auto-layout (28 commands)
**Layer 2 — Patterns:** Cards, UI components, visual effects (28 commands)
**Layer 3 — Deliverables:** Carousels, pitch decks, wireframes, storyboards (36 templates)
**Layer 4 — Design System:** Token specimens, palette generation, auditing (14 commands)
**Layer 5 — QA:** Contrast checking, touch targets, orphan detection (16 commands)
**Layer 6 — Export:** PNG/SVG/PDF export, CSS/React handoff (11 commands)
**Layer 7 — Orchestration:** YAML batch recipes for multi-step workflows

You can work at whatever level you need. Create a single rectangle, or generate an entire marketing campaign.

## The Theme System

Every command is theme-aware. Themes are JSON files that define design tokens:

```json
{
  "name": "Noir Studio",
  "colors": {
    "BG":   { "r": 0.05, "g": 0.06, "b": 0.09 },
    "BL":   { "r": 0.2,  "g": 0.4,  "b": 1.0  },
    "CARD": { "r": 0.086, "g": 0.1,  "b": 0.145 }
  },
  "type": {
    "h1": { "fontSize": 48, "style": "Bold", "lineHeight": 56 },
    "body": { "fontSize": 16, "style": "Regular", "lineHeight": 24 }
  },
  "fonts": { "heading": "Inter", "body": "Inter", "mono": "Geist Mono" }
}
```

Three themes ship built-in: a dark tech aesthetic, a print-friendly light theme, and the Noir dark premium theme. Custom themes go in `~/.config/figma-kit/themes/`.

## Data-Driven Templates

The most powerful commands accept YAML content specs:

```yaml
# slides.yml
theme: noir
slides:
  - name: Cover
    glow: topRight
    headline: "How many\ndecisions\ndid you miss?"
  - name: CTA
    glow: cta
    headline: "Start now."
    cta: { text: "Book a demo", centered: true }
```

One command, one YAML file, seven production slides. Change the content, re-run, done.

## Why Go?

Go gives us exactly what this tool needs:

- **Single binary** — no runtime, no node_modules, no Python venv
- **Cross-platform** — `goreleaser` builds for macOS, Linux, Windows with one config
- **`go:embed`** — themes, helpers, and templates compile into the binary
- **Homebrew** — `brew install dop-amine/tap/figma-kit` just works

The entire tool is ~5,000 lines of Go and ~700 lines of embedded JavaScript. It compiles in under 3 seconds.

## Try It

```bash
brew install dop-amine/tap/figma-kit
figma-kit init my-project
figma-kit themes
figma-kit make og-image --title "Hello HN" -t default
```

The project is MIT licensed and open source: [github.com/dop-amine/figma-kit](https://github.com/dop-amine/figma-kit)

---

*Tags: go, figma, cli, design-tools, mcp, open-source*
