# figma-kit

**AI-powered Figma design.**

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: BSL 1.1](https://img.shields.io/badge/License-BSL_1.1-blue.svg)](LICENSE)
[![CI](https://github.com/dop-amine/figma-kit/actions/workflows/ci.yml/badge.svg)](https://github.com/dop-amine/figma-kit/actions/workflows/ci.yml)

Describe what you want in natural language — your AI agent picks the right figma-kit commands and the [Figma MCP server](https://mcp.figma.com) executes them. A single Go binary with 150+ commands across 8 layers of abstraction, from `node create frame` to `make carousel --content slides.yml`.

**Works with [Cursor](https://cursor.com), [Claude Code](https://docs.anthropic.com/en/docs/claude-code), or any MCP-compatible AI agent.** Also runs standalone with `figma-kit exec` for direct MCP execution.

<img width="547" height="191" alt="image" src="https://github.com/user-attachments/assets/36c55805-a5c9-4ef4-8ff8-558620c414a1" />

## Install

```bash
# Homebrew (macOS / Linux)
brew install dop-amine/tap/figma-kit

# Go install
go install github.com/dop-amine/figma-kit@latest

# curl one-liner
curl -fsSL https://raw.githubusercontent.com/dop-amine/figma-kit/main/install.sh | sh

# Or download from GitHub Releases
```

## Quick Start

### With AI (recommended)

1. **Install figma-kit** and add the [Figma MCP server](https://mcp.figma.com) to your AI tool (Cursor, Claude Code, etc.)
2. **Prompt in natural language** — here are three real examples:

**From a reference website:**

> "Create a figma-kit theme that matches stripe.com, then build a landing page with a hero, 3 feature cards, and pricing."

The AI extracts Stripe's colors, runs `figma-kit theme init --bg "#0A2540" --primary "#635BFF" --accent "#00D4AA"`, creates a theme, then sequences `make screen`, `card glass`, `ui button` commands. Your Figma file fills with a complete, branded landing page.

**From a description:**

> "Build a 7-slide pitch deck for a climate-tech startup. Dark green theme, clean typography. Problem, solution, market, traction, team, ask."

**From a brand guide:**

> "Here's our brand guide [screenshot]. Extract the colors and fonts, create a theme, then generate a full design system page."

Run `figma-kit cookbook` to see 15 complete prompt sessions — from reference to finished design. See the full [Prompt Cookbook](docs/COOKBOOK.md).

### Standalone CLI

```bash
# 1. Initialize a project
figma-kit init my-project

# 2. Link your Figma file
figma-kit config set fileKey YOUR_FILE_KEY

# 3. Browse prompt examples (embedded in binary)
figma-kit cookbook --list

# 4. Dump starter content YAML files
figma-kit examples --dump

# 5. Generate a carousel with the Noir theme
figma-kit make carousel --content examples/saas-landing.yml -t noir
# → outputs use_figma JS — feed it to the MCP tool

# 5b. Or compose several steps into one payload (tree-shaken helpers)
figma-kit compose -t noir "ui section --title 'Demo'" "ui stat --parent _results[0] --value '150+' --label 'Commands'"

# 6. Run a QA audit
figma-kit qa checklist --page 0
```

### Compose (batch N commands → 1 call)

`compose` is the primary workflow for multi-step designs. It batches N commands into a single `use_figma` call — shared preamble, one JS payload, one round-trip. The engine **tree-shakes** embedded helpers: only functions your steps reference are emitted (often **~3KB** of helpers vs **~14KB** for the full bundle), which keeps recipes closer to the 50KB `use_figma` limit.

Chain steps with **`_results[]`**: compose emits `const _results = [];` and pushes each step’s main node. Reference prior steps with `--parent _results[0]` (or `fx … --last` for the immediately previous result):

```bash
figma-kit compose -t noir \
  "ui section --title 'Features' --label 'PRODUCT'" \
  "card glass --parent _results[0] --title 'Feature 1'" \
  "card glass --parent _results[0] --title 'Feature 2'" \
  "fx glow --last --position subtle"
```

Or define a recipe YAML:

```yaml
# landing.yml
theme: noir
steps:
  - "ui section --title 'Features' --label 'PRODUCT'"
  - "card glass --parent _results[0] --title 'Feature 1'"
  - "card glass --parent _results[0] --title 'Feature 2'"
  - "ui pricing --tiers '[{\"name\":\"Pro\",\"price\":\"$29\",\"highlighted\":true}]'"
```

```bash
figma-kit compose --recipe landing.yml
# or execute directly:
figma-kit exec compose --recipe landing.yml
```

### Authentication

Three paths depending on your workflow:

| Who | Setup | Run |
|-----|-------|-----|
| **Designers** | Install in Cursor / Claude Code. Add the [Figma MCP server](https://mcp.figma.com). | Prompt in natural language. Zero config. |
| **Developers** | `export FIGMA_TOKEN=<pat>` + `figma-kit auth login` | `figma-kit exec compose -t noir --recipe landing.yml` |
| **CI / Scripts** | `export FIGMA_ACCESS_TOKEN=<oauth-token>` | `figma-kit exec compose --recipe landing.yml` |

When a Personal Access Token (`FIGMA_TOKEN`) is set, the optional **REST API** client activates automatically — enabling file metadata, image exports, and other operations that complement the MCP-based Plugin API path.

See [docs/STANDARDS.md](docs/STANDARDS.md) for design conventions and coding standards.

## Why figma-kit?

Figma's `use_figma` MCP tool executes raw Plugin API JavaScript. You *can* build anything with it — but every request means writing 20–80 lines of JS, converting hex colors to `{r,g,b}` objects, manually calculating coordinates, and hoping the AI doesn't hallucinate a non-existent API.

figma-kit solves this:

| Without figma-kit | With figma-kit |
|---|---|
| 45 lines of JS for a glass card | `card glass -t noir --title "Feature"` |
| Hardcode `{r:0.2, g:0.4, b:1.0}` every time | `-t noir` propagates theme everywhere |
| AI must know the Plugin API intimately | AI just picks from 150+ named commands |
| No design system, no QA, no audit | `ds create`, `qa checklist` built in |
| Each run is unique, fragile JS | Tested, deterministic output every time |
| Need an AI to run anything | `exec` runs directly — no AI middleman |

figma-kit turns `use_figma` from a low-level escape hatch into a practical design tool.

## How It Works

```
You prompt → AI picks commands → figma-kit generates JS → Figma MCP executes → design appears
```

figma-kit generates JavaScript compatible with Figma's `use_figma` MCP tool, which executes inside the Plugin API sandbox. Two ways to use it:

### AI Agent Workflow (recommended)

1. **You describe** what you want in your AI tool (Cursor, Claude Code, etc.)
2. **The AI reasons** about which figma-kit commands to use and in what order
3. **figma-kit generates** Figma Plugin API JavaScript for each command
4. **The MCP server executes** the code inside your Figma file
5. **Verify** the result with `get_screenshot`

### Direct Execution with `exec`

```bash
# One-time auth
figma-kit auth login

# Execute any command directly — no AI middleman needed
figma-kit exec make carousel -t noir --content slides.yml
figma-kit exec card glass -t noir --title "Feature" --screenshot
figma-kit exec ui hero -t noir --title "Ship Faster"
```

The `exec` command generates JS and sends it to the Figma MCP server in one shot.

## AI Integration

figma-kit is designed to be an **AI agent tool**. Each command is a self-contained, composable unit that an LLM can select, sequence, and execute without human intervention.

| AI Client | How it works |
|-----------|-------------|
| **Cursor** | AI calls figma-kit commands via the Figma MCP server automatically |
| **Claude Code** | Connect the Figma MCP server; Claude uses figma-kit output in `use_figma` |
| **Any MCP client** | Any MCP-compatible agent can invoke figma-kit -> `use_figma` |
| **Terminal** | Run commands directly — pipe output to the MCP tool manually |

The entire [marketing site](https://dop-amine.github.io/figma-kit/) and its [Figma design](https://www.figma.com/design/olwJu5ffDAxhUyMYS3Nb4f/figma-kit-v2) were built using figma-kit with AI.

## Command Layers

| Layer | Group | Commands | Description |
|-------|-------|----------|-------------|
| 0 | Session | `init`, `config`, `whoami`, `open`, `status`, `auth`, `exec`, `new-file` | File, auth & direct execution |
| 1 | Primitives | `node`, `style`, `text`, `layout` | Low-level node ops + boolean, svg, variant-set |
| 2 | Patterns | `card`, `ui` (incl. `ui section` wrapper), `fx`, `image` | 8 card types, UI components + section layout helper, 14 effects |
| 3 | Deliverables | `make` | 37 production templates + changelog |
| 4 | Design System | `ds`, `ds library` | Token management, specimens, component-sheet, audit, **library browsing & import** |
| 5 | Inspect & QA | `inspect`, `tree`, `find`, `screenshot`, `qa` | Quality checks & MCP-backed screenshot |
| 6 | Export | `export`, `handoff` | PNG/SVG/PDF, CSS, React specs |
| 7 | Orchestration | `batch`, `compose` | YAML-driven recipes; compose batches N commands into 1 call |

Run `figma-kit --help` or `figma-kit <command> --help` for full details.

## Themes

Three built-in themes, switchable with `-t`:

| Theme | Description |
|-------|-------------|
| `default` | Dark theme for tech/SaaS. Blue-teal accents on `rgb(13,15,23)`. |
| `light` | Light mode for print-friendly deliverables. |
| `noir` | Brand theme. Primary blue `#3366FF`, dark premium aesthetic. |

Create themes from colors, screenshots, or existing sites — AI agents handle the extraction:

```bash
# AI workflow: paste a screenshot into Cursor, ask it to create a theme
figma-kit theme init --name "Brand" --bg "#1a1a2e" --primary "#e94560" --accent "#0f3460" -o brand.json

# Preview in Figma
figma-kit theme preview -t brand

# Or use the visual web builder
# → dop-amine.github.io/figma-kit/theme-builder.html

# List all themes (built-in, community, user, local)
figma-kit themes

# Export as CSS variables
figma-kit export tokens -t default --format css
```

Custom themes: place a JSON file in `~/.config/figma-kit/themes/` or `./themes/`. See [Theme Docs](docs/THEMES.md) for full details, or run `figma-kit cookbook theme-from-a-mood` for a step-by-step example.

## External Library Integration

Import published components, styles, and variable collections from team or community Figma libraries:

```bash
# Discover — requires a PAT (FIGMA_TOKEN)
figma-kit ds library list --team 123456
figma-kit ds library list --file abc123 --type styles
figma-kit ds library info a1b2c3d4e5

# Import — generates Plugin API JS, no PAT needed
figma-kit ds library import a1b2c3d4e5 --name "Hero"
figma-kit ds library import-set x9y8z7 --variant "Size=Large"
figma-kit ds library import-style s1t2y3l4e5 --apply 123:456
figma-kit ds library variables

# Compose with library components
figma-kit compose -t noir \
  "ds library import a1b2 --name Hero" \
  "fx glow --last"
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
internal/cli/                   # Cobra command definitions (150+ commands)
internal/codegen/               # Fluent JS code builder
internal/mcpclient/             # Embedded MCP client for direct execution
internal/theme/                 # Theme loading & validation
internal/config/                # .figmarc.json management
assets/                         # Embedded JS helpers, themes, templates, examples, cookbook
```

Build from source:
```bash
go build ./cmd/figma-kit
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

[BSL 1.1](LICENSE) — free for individuals and non-commercial use. Converts to MPL 2.0 on April 6, 2030.
