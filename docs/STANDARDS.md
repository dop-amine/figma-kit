# figma-kit Standards

This document defines the canonical architecture, composable command contract, and conventions for figma-kit. All other documentation references this file as the single source of truth.

## 1. Architecture Overview

figma-kit is a Go CLI that generates JavaScript strings for Figma's Plugin API. The binary never talks to Figma's servers for design operations — it composes JS that creates or inspects nodes when executed in the plugin runtime via the `use_figma` MCP tool.

Three delivery paths:

- **AI agent** (Cursor, Claude Code): AI picks figma-kit commands, captures JS output, pipes it to `use_figma`.
- **Direct exec**: `figma-kit exec compose -t noir --recipe site.yml` sends JS to Figma via the embedded MCP client. No AI middleman.
- **Stdout**: `figma-kit card glass -t noir` prints JS to stdout. User or script pipes it wherever needed.

### The Compose Model

Multiple commands can be batched into a single `use_figma` call using `figma-kit compose`. The compose engine emits a shared preamble once, wraps each command body in `{ }` for scope isolation, collects all created node IDs, and outputs a single JS payload under 50,000 characters.

```
N figma-kit commands → 1 compose call → 1 JS payload → 1 use_figma call → 1 round-trip
```

## 2. The Composable Command Contract

Every command that emits Figma Plugin API JavaScript is **composable** — it can participate in a `figma-kit compose` batch.

### Requirements

A composable command MUST:

- Set `cmd.Annotations["composable"] = "true"` on its Cobra command.
- Emit creation or modification JS to `output()` (which writes to `os.Stdout`).
- Use `codegen.Builder` for JS generation.
- Respect the `bodyOnly` flag on the Builder: when `b.IsBodyOnly() == true`, the command skips its own preamble and return statement (compose provides these).
- Return node IDs via `ReturnIDs()` or `ReturnDone()` for modifications.

A composable command MUST NOT:

- Emit preamble when `bodyOnly == true` (compose provides the superset preamble).
- Call `os.Exit()`.
- Write data to stderr (stderr is for progress/log messages only).
- Assume it is the only command running (use `const` for locals — compose wraps bodies in `{ }` for scope isolation).

### Non-Composable Commands

Commands that do not emit Plugin API JS are not composable:

- Go HTTP servers (`image serve`)
- Go data output (`export tokens`, `ds sync-tokens`, `ds tokens`, `theme init`)
- Config/utility (`init`, `config`, `open`, `status`, `info`, `themes`, `cookbook`, `examples`, `docs`)
- MCP-only (`auth`, `whoami`, `new-file`, `ds search`, `ds import`, `ds audit`)
- Own composition (`batch`, `watch`, `compose`)
- Validation (`validate theme`, `validate recipe`)

Compose rejects non-composable commands with a clear error message.

### Compose orchestration contracts (v2.1)

- **`_results[]` array** — `figma-kit compose` emits `const _results = [];` and, after each step, pushes that step’s primary created node. Later steps can reference prior outputs in flags as `_results[0]`, `_results[1]`, etc. (zero-based, in step order).
- **`--last` on `fx`** — Every `fx` subcommand accepts `--last`. In compose, it targets the previous step’s main node (equivalent to the last entry in `_results`) instead of passing a positional node id.
- **`ui section`** — Recommended pattern for a centered block wrapper (label, heading, subtitle, optional divider) when building pages inside compose; nest child commands with `--parent _results[N]` (or `--last` on `fx`) instead of flat unrelated frames.
- **Tree-shaking helpers** — Compose includes only helper functions that steps actually reference (via `detectNeededHelpers`), not the full helper bundle every time — roughly **~14KB → ~3KB** of helper JS when few helpers are used.
- **Universal `--parent`** — Composable `ui` and `card` commands accept `--parent` for chaining. Values are either a Figma node id string or, when the value **starts with** `_results[`, emitted as a **JavaScript expression** (no `getNodeByIdAsync` wrapper) so compose can wire steps together.

## 3. Preamble Families

Commands use two preamble patterns. Compose emits the superset.

### Standard Preamble (`PreambleWithPage`)

Used by `make`, `ds`, `qa`, `handoff`, `export`:

1. Page setup: `const pg = figma.root.children[N]; await figma.setCurrentPageAsync(pg);`
2. Font loading: `FontLoadingFromTheme` — families and style list from **`theme.Fonts`** (heading, body, mono, weights); built-in themes default to Inter + Geist Mono with standard weights.
3. Theme color constants: `const BG={r:...,g:...,b:...};` etc.

### UI Preamble (`uiPreamble`)

Used by `ui`, `card`, `fx` commands:

1. Page setup
2. Font loading (same **`theme.Fonts`**-driven pattern as standard preamble)
3. Theme color constants
4. Type scale tokens: `TY_BODY`, `ST_BODY`, `TY_SMALL`, `ST_SMALL`, `TY_LABEL`, `ST_LABEL`, `TY_H4`, `ST_H4`
5. Spacing tokens

### Compose Superset

The compose engine emits the union of both families plus **only the helper functions each step needs** (tree-shaken via `detectNeededHelpers`), not necessarily every symbol in `helpers.js`:

1. Page setup (once)
2. Full font loading derived from the active theme (**`theme.Fonts`** and standard faces — preamble tracks theme-aware families)
3. All theme color constants (sorted)
4. Type scale tokens
5. Needed helpers (`G`, `T`, and any functions referenced by merged step bodies)
6. `const _ids = [];` for ID collection and `const _results = [];` for per-step main nodes

## 4. Image Pipeline

### Local Files (< 33KB)

Go reads the file with `os.ReadFile`, base64-encodes it, and emits JS that decodes and creates an image:

```
image place ./logo.png
  → Go: os.ReadFile() → base64.Encode
  → JS: atob(b64) → Uint8Array → figma.createImage(buf)
```

The 33KB limit exists because base64 expands the data ~33%, and the total `use_figma` payload must stay under 50,000 characters.

### Remote URLs

JS fetches the image directly:

```
image place https://example.com/hero.jpg
  → JS: await fetch(url) → arrayBuffer → figma.createImage(buf)
```

### Large Files

For files > 33KB, use `figma-kit image serve` to start a local HTTP server, then reference the URL:

```bash
figma-kit image serve ./assets &    # starts http://127.0.0.1:8741/
figma-kit image place http://127.0.0.1:8741/hero.png
```

`image serve` is NOT composable (it's a Go HTTP server, not JS output).

## 5. Authentication Tiers

### Tier 1: AI Agent (Zero Config)

When using figma-kit through Cursor, Claude Code, or another MCP-compatible AI tool, the AI handles OAuth transparently. No setup required.

### Tier 2: Direct Token

Set `FIGMA_ACCESS_TOKEN` environment variable with an OAuth access token. figma-kit uses it directly — no PAT, no registration, no browser flow.

```bash
export FIGMA_ACCESS_TOKEN=your_oauth_token
figma-kit exec compose -t noir --recipe site.yml
```

Also available via `figma-kit auth login --token <token>`.

### Tier 3: PAT-Bootstrapped OAuth

For first-time `exec` setup, figma-kit uses a Figma Personal Access Token (PAT) to register a dynamic OAuth client, then runs a PKCE browser flow:

```bash
export FIGMA_TOKEN=your_pat    # or FIGMA_PAT or FIGMA_PERSONAL_ACCESS_TOKEN
figma-kit auth login           # opens browser, caches OAuth token
```

The PAT is only used once for client registration. The resulting OAuth token is cached at `~/.config/figma-kit/token.json`.

### Tier 0: Optional REST API

When a PAT is available (`FIGMA_TOKEN` or `FIGMA_PAT` env var), figma-kit can use Figma's REST API for richer read operations:

- `GET /v1/files/:key` — full file JSON
- `GET /v1/images/:key` — render nodes as PNG/SVG/PDF at any scale
- `GET /v1/files/:key/nodes` — specific node subtrees

REST API calls use a separate rate limit budget and do not consume MCP quota. Every command that uses the REST API falls back gracefully when no PAT is set.

## 6. JS Generation Rules

- Top-level `await` works (Plugin API auto-wraps code in async context).
- Use `const` for local variables (compose wraps bodies in `{ }` for scope isolation).
- Never call `figma.closePlugin()`.
- Never use `setPluginData()`.
- Maximum payload: 50,000 characters per `use_figma` call, 20KB output.
- Supported node types: Rectangle, Frame, Component, Text, Ellipse, Star, Line, Vector, Polygon, BooleanOperation, Slice, Page, Section, TextPath.
- No image/asset upload in beta — use base64 encoding or URL fetch.
- Fonts: system fonts + Google Fonts already loaded in Figma. Use `figma.loadFontAsync()` before setting text.

## 7. Adding a New Command

1. Create a constructor in the appropriate `internal/cli/*.go` file returning `*cobra.Command`.
2. Set `cmd.Annotations = map[string]string{"composable": "true"}` if the command emits Plugin API JS.
3. In `RunE`, call `resolveTheme` / `resolvePage` if needed.
4. Use `codegen.New()` for JS generation. Call preamble functions only when `!b.IsBodyOnly()`.
5. End with `ReturnIDs()` or `ReturnDone()` only when `!b.IsBodyOnly()`.
6. Call `output(b.String())`.
7. Register the command in `internal/cli/root.go`.
8. Add a table-driven test case in `internal/cli/cli_test.go`.
9. Document in `docs/COMMANDS.md`.
10. Run `go test ./...` and lint.

## 8. Figma MCP Server Reference

Endpoint: `mcp.figma.com/mcp`

17 tools available:

| Tool | Purpose | Rate Limited |
|------|---------|-------------|
| `use_figma` | Execute Plugin API JS (50KB input, 20KB output) | No (beta) |
| `create_new_file` | Create a new Figma file | Yes |
| `get_screenshot` | Capture node screenshot | Yes |
| `get_metadata` | Node metadata | Yes |
| `get_design_context` | Design context for a node | Yes |
| `get_variable_defs` | Variable definitions | Yes |
| `search_design_system` | Search components/styles | Yes |
| `whoami` | Current user info | Yes |
| `get_code_connect_map` | Code Connect mappings | Yes |
| `add_code_connect_map` | Add Code Connect mapping | Yes |
| `get_code_connect_suggestions` | Code Connect suggestions | Yes |
| `send_code_connect_mappings` | Bulk send mappings | Yes |
| `get_context_for_code_connect` | Context for mapping | Yes |
| `create_design_system_rules` | Create DS rules | Yes |
| `get_figjam` | FigJam board data | Yes |
| `generate_diagram` | Generate FigJam diagram | Yes |
| `generate_figma_design` | AI-generated design | Yes |

Write operations via `use_figma` are rate-limit exempt during beta. Read tools are limited to 200-600 calls/day depending on plan, 10-20/min. Minimize reads, maximize writes per call.
