# Batch recipes

Batch recipes are YAML files that group several **pre-generated JavaScript snippets** into a single workflow. figma-kit does not run Figma itself: it prints numbered blocks meant for an AI or operator to paste into the **use_figma** MCP tool (or equivalent) **one after another**.

## YAML format

```yaml
title: Optional recipe heading
steps:
  - title: Short label for this step
    js: |
      // multi-line JavaScript
      // (output from a figma-kit command)
  - title: Next step
    js: |
      ...
```

| Field | Required | Description |
| --- | --- | --- |
| `title` (root) | No | Printed as `// Recipe: …` at the top of the output. |
| `steps` | Yes | Non-empty list of steps. |
| `steps[].title` | No | Defaults to `Step 1`, `Step 2`, … |
| `steps[].js` | Yes | Non-empty string; trimmed and wrapped in a labeled block. |

The parser lives in `internal/cli/batch.go` (`batchRecipe` struct). Root-level `name` is **not** read by the CLI—use `title` for the recipe heading.

## Generating step bodies

Each step’s `js` should be the **exact stdout** of a figma-kit command that emits plugin JS:

```bash
figma-kit make carousel --content slides.yml -t noir > /tmp/step1.js
# Paste contents into the YAML under js: |
```

Tips:

- Re-run the same command whenever you change flags, theme, or content files.
- Use the same `-t` / `.figmarc.json` theme across steps so colors and type stay aligned.
- Use `-p` / `--page` when you need a specific Figma page index consistently.

## How `figma-kit batch` works

```bash
figma-kit batch path/to/recipe.yaml
```

The command:

1. Reads and parses the YAML.
2. Validates that there is at least one step and that each step has non-empty `js`.
3. Prints to stdout:
   - Optional `// Recipe: <title>`
   - For each step: `// --- Block N: <title> ---` followed by the JS and a blank line.

An AI assistant executing the recipe should run **Block 1**, wait for completion, then **Block 2**, and so on. Each block is independent plugin code (often including preamble + page setup when generated from subcommands that use `PreambleWithPage`).

## Example: campaign-style recipe

Below, replace the `js` placeholders with real command output (heredocs or saved files). The structure shows carousel + one-pager + OG image + QA checklist in sequence.

```yaml
title: Q2 launch — Figma pass
steps:
  - title: Carousel slides
    js: |
      PASTE_STDOUT_OF: figma-kit make carousel -t noir --content slides.yml

  - title: B2B one-pager
    js: |
      PASTE_STDOUT_OF: figma-kit make one-pager -t noir --content one-pager.yml

  - title: OG image
    js: |
      PASTE_STDOUT_OF: figma-kit make og-image -t noir --title "Q2 Launch"

  - title: QA checklist
    js: |
      PASTE_STDOUT_OF: figma-kit qa checklist -t noir
```

After filling in real JS:

```bash
figma-kit batch campaign.yml
```

## Compose recipes (command-based format)

The newer `compose` command uses a simpler recipe format where steps are **command strings** instead of raw JS. The compose engine runs each command internally, merges the output into one JS payload with a shared preamble, and emits it as a single `use_figma` call.

### Compose YAML format

```yaml
theme: noir       # optional — overrides -t flag
page: 0           # optional — overrides -p flag
steps:
  - "ui hero --title 'Ship Faster' --cta 'Start'"
  - "card glass --title 'Feature 1'"
  - "card glass --title 'Feature 2'"
  - "ui pricing --tiers '[{\"name\":\"Pro\",\"price\":\"$29\"}]'"
```

| Field | Required | Description |
| --- | --- | --- |
| `theme` | No | Theme name; overridden by `-t` flag if both set. |
| `page` | No | Page index; overridden by `-p` flag. |
| `steps` | Yes | List of figma-kit command strings (without the `figma-kit` prefix). |

Only commands annotated as `composable` (most Layer 1–2 commands that emit Plugin API JS) are valid as steps.

### Running compose recipes

```bash
# Generate merged JS to stdout
figma-kit compose --recipe landing.yml

# Generate + execute directly in Figma
figma-kit exec compose --recipe landing.yml

# Mix recipe steps with inline commands
figma-kit compose --recipe base.yml "ui footer --cols 4"
```

### Example: landing page compose recipe

```yaml
# landing.yml
theme: noir
steps:
  - "ui hero --title 'Build Faster' --subtitle 'AI-powered design' --cta 'Get Started'"
  - "card glass --title 'Speed' --preset strong"
  - "card glass --title 'Scale' --preset strong"
  - "card glass --title 'Security' --preset strong"
  - "image place ./brand/logo.png --name 'Logo' --width 200 --height 60"
  - "ui footer --cols 3 --copyright '© 2026 Acme'"
```

### Compose vs batch

| | `batch` | `compose` |
|---|---|---|
| Step format | Raw JS (pre-generated output) | Command strings |
| Preamble | Repeated per step | Shared (emitted once) |
| Output | Concatenated JS blocks | Single merged JS payload |
| Execution | One `use_figma` call per block | One `use_figma` call total |
| Best for | Pre-built JS workflows | Multi-command designs |

Both formats are valid. Use `compose` when you want the engine to handle preamble merging and scope isolation automatically. Use `batch` when you have pre-generated JS snippets you want to sequence.

## Tips

- **Validate between steps**: Use the Figma MCP **get_screenshot** tool on key frames after each block; `figma-kit screenshot` prints guidance for that workflow.
- **Theme consistency**: Prefer one theme for the whole recipe (`-t` on every generating command or a single `.figmarc.json`).
- **Page index**: If everything should land on the same page, use the same `-p` value when generating each step.
- **Content files**: Keep YAML content paths stable relative to where you run `figma-kit`; the batch file only stores JS text, not paths to your slide/one-pager YAML.
