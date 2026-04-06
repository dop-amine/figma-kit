# figma-kit

**CLI for programmatic Figma design via the MCP server.**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/dop-amine/figma-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/dop-amine/figma-kit/actions/workflows/ci.yml)

A single Go binary that generates `use_figma`-compatible JavaScript for the official [Figma MCP server](https://mcp.figma.com). 120+ commands across 8 layers of abstraction — from `node create frame` to `make carousel --content slides.yml`.

<img width="547" height="191" alt="image" src="https://github.com/user-attachments/assets/36c55805-a5c9-4ef4-8ff8-558620c414a1" />

## Install

```bash
# Homebrew (macOS / Linux)
brew install dop-amine/tap/figma-kit

# Go install
go install github.com/dop-amine/figma-kit@latest

# curl one-liner
curl -fsSL https://raw.githubusercontent.com/amine/figma-kit/main/install.sh | sh

# Or download from GitHub Releases
```

## Quick Start

```bash
# 1. Initialize a project
figma-kit init my-project

# 2. Link your Figma file
figma-kit config set fileKey YOUR_FILE_KEY

# 3. Generate a carousel with the Noir theme
figma-kit make carousel --content slides.yml -t noir
# → outputs use_figma JS — feed it to the MCP tool

# 4. Create a hero frame
figma-kit node create frame --name "Hero" -w 1440 --height 800

# 5. Run a QA audit
figma-kit qa checklist --page 0
```

## How It Works

```
figma-kit CLI → generates JavaScript → use_figma MCP tool → Figma Plugin API → Figma file
```

The CLI doesn't connect to Figma directly. It generates JavaScript compatible with Figma's `use_figma` MCP tool, which executes it inside the Plugin API sandbox. The workflow is:

1. Run a `figma-kit` command to produce JS
2. Feed the JS to `use_figma` via the MCP server
3. Verify the result with `get_screenshot`

## Command Layers

| Layer | Group | Commands | Description |
|-------|-------|----------|-------------|
| 0 | Session | `init`, `config`, `whoami`, `open`, `status` | File & project management |
| 1 | Primitives | `node`, `style`, `text`, `layout` | Low-level Figma node operations |
| 2 | Patterns | `card`, `ui`, `fx` | Mid-level design components & effects |
| 3 | Deliverables | `make` | Full production designs (36 templates) |
| 4 | Design System | `ds` | Token management, specimens, audit |
| 5 | Inspect & QA | `inspect`, `tree`, `find`, `qa` | Quality checks & inspection |
| 6 | Export | `export`, `handoff` | PNG/SVG/PDF, CSS, React specs |
| 7 | Orchestration | `batch` | YAML-driven multi-step recipes |

Run `figma-kit --help` or `figma-kit <command> --help` for full details.

## Themes

Three built-in themes, switchable with `-t`:

| Theme | Description |
|-------|-------------|
| `default` | Dark theme for tech/SaaS. Blue-teal accents on `rgb(13,15,23)`. |
| `light` | Light mode for print-friendly deliverables. |
| `noir` | Brand theme. Primary blue `#3366FF`, dark premium aesthetic. |

Custom themes: place a JSON file in `~/.config/figma-kit/themes/` or `./themes/`.

```bash
# List themes
figma-kit themes

# Use a specific theme
figma-kit make og-image --title "Hello" -t noir

# Export theme as CSS variables
figma-kit export tokens -t default --format css
```

## Content Specs

Data-driven templates accept YAML content files via `--content`:

```yaml
# slides.yml
theme: noir
total: 7
slides:
  - name: Cover
    glow: topRight
    headline: "How many\ndecisions\ndid you miss?"
    subtitle: "Because you didn't have the data."
  - name: Problem
    glow: subtle
    headline: "The real cost\nof blind spots"
    chips: ["Revenue", "Speed", "Trust"]
  - name: CTA
    glow: cta
    headline: "Start now."
    cta:
      text: "Book a demo"
      centered: true
```

## Batch Recipes

Orchestrate multi-step workflows with YAML recipes:

```yaml
# campaign.yml
name: "Q2 Campaign"
steps:
  - title: Carousel
    js: |
      // output from: figma-kit make carousel --content slides.yml
  - title: OG Image
    js: |
      // output from: figma-kit make og-image --title "Launch"
```

```bash
figma-kit batch campaign.yml
```

## Architecture

The binary is pure Go. JavaScript files in `assets/` are embedded at compile time via `go:embed` — they are the **output format** for Figma's Plugin API, not implementation code.

```
cmd/figma-kit/main.go          # Entry point
internal/cli/                   # Cobra command definitions
internal/codegen/               # Fluent JS code builder
internal/theme/                 # Theme loading & validation
internal/config/                # .figmarc.json management
internal/batch/                 # YAML recipe parser
assets/                         # Embedded JS helpers, themes, templates
```

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for details.

## Contributing

```bash
# Clone
git clone https://github.com/dop-amine/figma-kit && cd figma-kit

# Build
make build

# Test
make test

# Lint
make lint

# Install locally
make install
```

To add a new command, create a function in `internal/cli/` following the existing patterns — use the codegen builder, resolve theme and page, generate JS, and call `output()`.

## Acknowledgments

- [Figma MCP Server](https://mcp.figma.com) — the official Figma integration that makes this possible
- [figma-use](https://github.com/nicekid1/figma-use) / [figma-cli](https://github.com/nicekid1/Figma-CLI) — prior art in the Figma CLI space
- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)

## License

[MIT](LICENSE)
