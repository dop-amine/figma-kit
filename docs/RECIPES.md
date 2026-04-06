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

## Tips

- **Validate between steps**: Use the Figma MCP **get_screenshot** tool on key frames after each block; `figma-kit screenshot` prints guidance for that workflow.
- **Theme consistency**: Prefer one theme for the whole recipe (`-t` on every generating command or a single `.figmarc.json`).
- **Page index**: If everything should land on the same page, use the same `-p` value when generating each step.
- **Content files**: Keep YAML content paths stable relative to where you run `figma-kit`; the batch file only stores JS text, not paths to your slide/one-pager YAML.
