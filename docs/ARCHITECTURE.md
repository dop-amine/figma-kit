# Architecture

## Overview

**figma-kit** is a Go CLI that **generates JavaScript** intended to run inside Figma (via the Plugin API), typically delivered through the **use_figma** MCP tool. The binary never talks to Figma's servers directly for design operations: it composes strings of JS that create or inspect nodes when executed in the plugin runtime.

## AI Agent Workflow

figma-kit is designed as an **AI agent tool** — each command is a composable unit that an LLM can select, sequence, and execute without human intervention.

```
User prompt -> AI selects figma-kit commands -> CLI generates JS -> Figma MCP executes -> design appears
```

The AI agent (via Cursor, Claude Code, or any MCP-compatible client) acts as the orchestrator: it reads the user's natural-language request, picks the appropriate figma-kit commands, runs them to produce JS, and feeds the output to the Figma MCP server's `use_figma` tool. Designers can describe what they want in plain English and get production Figma designs without touching the CLI directly.

## Why Go and JavaScript

- **JavaScript** is required because Figma's Plugin API is JS-only (async fonts, node creation, `figma.*` APIs).
- **Go** implements the CLI, configuration, theme parsing, tests, and release pipeline. Keeping generation logic in Go yields a single static binary, strong typing for theme/config structs, and straightforward embedding of assets.

## Package layout

| Path | Role |
| --- | --- |
| `cmd/figma-kit/main.go` | Entry point; calls `cli.Execute()`. |
| `internal/cli/` | Cobra commands, theme/page resolution, orchestration; includes **batch** recipe execution (`batch.go`). |
| `internal/codegen/` | Fluent **Builder** for composing JS (`New`, `Line`, `Raw`, `Preamble`, `PreambleWithPage`, etc.). |
| `internal/theme/` | Theme structs, JSON parse/validate, `Load` / `LoadFile`, embedded theme registry. |
| `internal/config/` | `.figmarc.json` load/save and `config set` / `get` helpers. |
| `internal/mcpclient/` | Embedded MCP client for direct `exec` execution; OAuth flow, token caching, `use_figma` calls. |
| `internal/restapi/` | Optional Figma REST API client; auto-enabled when a PAT is set (`FIGMA_TOKEN`). File metadata, image exports. |
| `assets/` | `go:embed` for themes, templates, helper JS; re-exported into `internal/theme` / codegen where needed. |

YAML batch recipes are parsed in **`internal/cli/batch.go`** (not a separate `internal/batch` package).

## Codegen `Builder`

`internal/codegen/builder.go` exposes a small fluent API for appending JS:

- **`New()`** — empty buffer.
- **`Comment` / `Line` / `Linef` / `Blank` / `Raw`** — structure and raw fragments.
- **`PageSetup(pageIndex)`** — select `figma.root.children[i]` and `setCurrentPageAsync`.
- **`FontLoading()`** — standard Inter + Geist Mono `loadFontAsync` loop.
- **`ThemeColors` / `ThemeColorsOrdered`** — emit `const TOKEN={r,g,b};` lines.
- **`ReturnIDs` / `ReturnDone` / `ReturnExpr`** — tail return shapes for MCP-style results.

`internal/codegen/preamble.go` defines **`Preamble`** (fonts + sorted theme colors) and **`PreambleWithPage`** (page setup then preamble) — the pattern most commands use after resolving theme and page index.

## Theme system

- Built-in themes are embedded from `assets/themes/*.json` via `assets/embed.go` and wired in `internal/theme/embed.go` (`embeddedThemes` map).
- **`theme.Load(name)`** resolution: embedded name -> `~/.config/figma-kit/themes/<name>.json` -> `./themes/<name>.json`.
- **`theme.List()`** only enumerates **embedded** themes (for `figma-kit themes`).
- **`theme.LoadFile(path)`** is available for explicit paths (tests, future CLI extensions).

## Command pattern

Most commands follow the same shape:

1. **Resolve theme** — `resolveTheme(cmd)` (`-t` -> `.figmarc.json` `theme` -> `"default"`), then `theme.Load`.
2. **Resolve page** — `resolvePage()` (`-p` if >= 0, else config `page`, else `0`).
3. **Build** — `b := codegen.New()`, then `PreambleWithPage` or specialized sections, template `Raw`, deliverable-specific lines.
4. **Emit** — `output(b.String())` writes JS to stdout.

Commands that only export tokens from Go (e.g. `figma-kit export tokens`) skip the Builder and write JSON/CSS directly after `resolveTheme`.

## Compose engine

`internal/cli/compose.go` implements the `compose` command — the primary workflow for batching N figma-kit commands into a single `use_figma` call.

**How it works:**

1. **Resolve steps** — from positional args (`"ui hero --title X" "card glass"`) and/or `--recipe` YAML.
2. **Validate** — each step is resolved against the command tree; only commands annotated `composable: true` are allowed.
3. **Capture** — each step is executed as a child invocation of the root command, capturing its full JS stdout.
4. **Merge** — a superset preamble (page setup, **theme-aware font loading from `theme.Fonts`**, theme colors, type scale) is emitted once, plus **only helper functions referenced across steps** (`detectNeededHelpers` tree-shakes the embedded helper bundle). Each step's body is stripped of its individual preamble/return and wrapped in `{ }` for scope isolation.
5. **Cross-step results** — compose emits `const _results = [];` and pushes each step’s primary node id so later steps can use `--parent _results[N]`, `text create --parent _results[0]`, or `fx … --last` to target the previous step’s output without hard-coded ids.
6. **Emit** — the merged JS is printed to stdout, ending with `return { createdNodeIds: _ids }`.

The compose recipe YAML format uses `theme`, `page`, and `steps` (list of command strings) — distinct from the older `batch` format which stores raw JS.

## MCP client (`internal/mcpclient/`)

`internal/mcpclient/client.go` provides a `Session` that connects to `https://mcp.figma.com/mcp` over MCP Streamable HTTP. Used by `exec` and `new-file` commands.

- **OAuth flow** — `auth login` opens a browser for Figma authorization using PKCE. Tokens are cached at `~/.config/figma-kit/token.json`.
- **Direct token** — `auth login --token <token>` saves a pre-existing access token, skipping the browser flow.
- **Session lifecycle** — `Connect(ctx)` creates a session, `CallUseFigma(fileKey, js)` invokes the `use_figma` tool, `GetScreenshot(fileKey)` captures a frame.

## REST API client (`internal/restapi/`)

`internal/restapi/client.go` wraps the Figma REST API (`https://api.figma.com`). Created via `NewClient()`, which returns `nil` when no PAT is available — callers check for nil before using. Response types live in `internal/restapi/types.go`.

Token resolution order: `FIGMA_TOKEN` → `FIGMA_PAT` → `FIGMA_PERSONAL_ACCESS_TOKEN`.

The REST API is optional and complements the MCP path. It enables:
- **File operations**: metadata retrieval, image exports, node subtree fetching
- **Library discovery**: listing published components, component sets, and styles from team or file libraries (used by `ds library list` and `ds library info`)

### Library endpoints

Team-level endpoints (`/v1/teams/:id/components`, `/v1/teams/:id/component_sets`, `/v1/teams/:id/styles`) support pagination via `page_size` and `after` cursor. File-level endpoints (`/v1/files/:key/components`, etc.) return all assets in a single response. Single-asset lookups (`/v1/components/:key`, `/v1/component_sets/:key`, `/v1/styles/:key`) return detailed metadata for a specific published asset.

### Library imports vs. discovery

Discovery (`ds library list`, `ds library info`) uses the REST API and requires a PAT. Import commands (`ds library import`, `ds library import-set`, `ds library import-style`) generate Plugin API JavaScript using `figma.importComponentByKeyAsync()`, `figma.importComponentSetByKeyAsync()`, and `figma.importStyleByKeyAsync()` — these run in the Figma context via `use_figma` and do **not** require a PAT.

The CLI helper `resolveRESTClient()` (in `internal/cli/root.go`) returns a client or a descriptive error pointing the user to PAT setup.

## Image pipeline

`internal/cli/image.go` provides three subcommands:

- **`image place`** — creates a new image frame from a local file or URL. Local files ≤ 33 KB are base64-encoded inline; larger files require a URL.
- **`image fill`** — replaces an existing node's fill with an image.
- **`image serve`** — starts a local HTTP file server for images too large for inline embedding.

Both `image place` and `image fill` are annotated `composable: true` and work inside `compose` recipes.

See [STANDARDS.md](STANDARDS.md) for design conventions and coding standards.

## Adding a new command

1. Add a constructor in the appropriate `internal/cli/*.go` file (or a new file) returning `*cobra.Command`.
2. In `RunE`, call `resolveTheme` / `resolvePage` if the script needs them.
3. Use `codegen.New()`, `PreambleWithPage` when the snippet runs in a page context, then add JS via `Line` / `Raw` (embed strings from `assets/` if large).
4. Register the command in `internal/cli/root.go` with `cmd.AddCommand(...)`.
5. Add a **table-driven test** case in `internal/cli/cli_test.go` (or a focused `_test.go`) that checks stdout contains expected substrings (valid structure, key API calls).
6. Run `go test ./...` and lint.

## Testing strategy

- **Table-driven tests** in `internal/cli/cli_test.go`: each row runs the root command with args, captures stdout, asserts required fragments (e.g. `figma.loadFontAsync`, `createFrame`).
- **Stdout capture**: `executeCmd` combines `cmd.SetOut` with a pipe-based `captureOutput` wrapper so code paths that write to `os.Stdout` (via `output()`) are still asserted.
- Package-level tests also exist under `internal/theme`, `internal/config`, `internal/codegen` for parsing, colors, and builder behavior.

## Build and distribution

- **GoReleaser** (`.goreleaser.yml`): builds `figma-kit` from `./cmd/figma-kit` for darwin/linux/windows (amd64 + arm64 where applicable), stamps version via ldflags (`internal/cli.Version`), archives tarballs/zip, checksums, and changelog filters.
- **Homebrew**: `brews` publishes to tap repo `dop-amine/homebrew-tap` for `brew install`.
- **`install.sh`**: downloads the latest GitHub release archive for the current OS/arch, extracts, and `install`s the binary to `/usr/local/bin` (or `INSTALL_DIR`).

Local development typically uses `go build -o figma-kit ./cmd/figma-kit` or a project `Makefile` target.
