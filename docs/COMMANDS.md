# figma-kit command reference

`figma-kit` is a CLI with **150+ commands** that generates JavaScript for Figma's `use_figma` MCP tool. It covers everything from low-level node primitives to complete design templates, all powered by a built-in theme system.

**Two ways to use commands:**
- **AI workflow:** Let your AI agent (Cursor, Claude Code) select and sequence commands. The JS output is piped to `use_figma` automatically.
- **Direct execution:** Run `figma-kit exec <command>` to generate JS and send it to Figma in one shot (requires `figma-kit auth login` first).
- **Standalone:** Run any command and pipe/paste the output into your Figma MCP execution path.

Commands that resolve a theme use **`--theme` / `-t`** and page index **`--page` / `-p`** when relevant.

---

## Global flags

Available on the root command and inherited by subcommands (where applicable).

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--theme` | `-t` | string | *(empty)* | Theme name: built-in (`default`, `light`, `noir`) or path to a theme file. Empty uses `.figmarc.json` `theme`, then `default`. |
| `--page` | `-p` | int | `-1` | Zero-based page index for generated JS. `-1` uses `.figmarc.json` `page`, then `0`. |
| `--version` | `-v` | bool | — | Print CLI version and exit. |
| `--help` | `-h` | bool | — | Help for the current command (Cobra). |

---

## Layer 0 — Session

### `init`

**Usage:** `figma-kit init [name]`

**Description:** Create `.figmarc.json` in the current directory (optional project `name`, default `Untitled`). Does not emit JS.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit init "My file"
```

---

### `config set`

**Usage:** `figma-kit config set <key> <value>`

**Description:** Set a key in `.figmarc.json`. Common keys: `fileKey`, `theme`, `page`, `exportDir`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | Positional only. |

```bash
figma-kit config set fileKey abcdefghijklmnop
figma-kit config set theme noir
figma-kit config set page 1
```

---

### `config get`

**Usage:** `figma-kit config get <key>`

**Description:** Print a single config value.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | Positional only. |

```bash
figma-kit config get fileKey
```

---

### `config list`

**Usage:** `figma-kit config list`

**Description:** Print all known config fields (`fileKey`, `theme`, `page`, `exportDir`).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit config list
```

---

### `whoami`

**Usage:** `figma-kit whoami`

**Description:** Prints guidance to use the Figma MCP **`whoami`** tool (not JS).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit whoami
```

---

### `open`

**Usage:** `figma-kit open`

**Description:** Opens `https://www.figma.com/file/<fileKey>` in the default browser. Requires `fileKey` in config.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit open
```

---

### `status`

**Usage:** `figma-kit status`

**Description:** Emit JS that walks all pages and returns frame summaries (counts, ids, sizes).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Use `-p` for page setup in other commands; `status` scans all pages in generated JS. |

```bash
figma-kit status
```

---

### `auth login`

**Usage:** `figma-kit auth login [flags]`

**Description:** Authenticate with Figma for direct command execution. Without flags, runs the full OAuth flow (opens browser). With `--token`, saves a pre-existing access token directly — useful for CI or headless environments.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--token` | string | `""` | Pre-existing OAuth access token (skips browser flow) |

```bash
figma-kit auth login                    # Full OAuth flow (opens browser)
figma-kit auth login --token <token>    # Direct token injection (CI / headless)
```

---

### `auth logout`

**Usage:** `figma-kit auth logout`

**Description:** Clear the cached OAuth token.

```bash
figma-kit auth logout
```

---

### `auth status`

**Usage:** `figma-kit auth status`

**Description:** Check whether a valid OAuth token exists.

```bash
figma-kit auth status
```

---

### `exec`

**Usage:** `figma-kit exec <sub-command> [flags]`

**Description:** Generate JS from any figma-kit command and execute it directly in Figma via the MCP server in one shot. No AI middleman required. Requires prior `auth login`. Works with any subcommand including `compose` for batched multi-step execution.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--file-key` | string | *(from config)* | Figma file key to target |
| `--screenshot` | bool | `false` | Take a screenshot after execution |
| `--timeout` | duration | `30s` | MCP call timeout |

```bash
figma-kit exec make carousel -t noir --content slides.yml
figma-kit exec card glass -t noir --title "Feature" --screenshot
figma-kit exec ui hero -t noir --title "Ship Faster"

# Execute a compose recipe directly
figma-kit exec compose -t noir --recipe landing.yml
```

---

### `new-file`

**Usage:** `figma-kit new-file <name>`

**Description:** Create a new Figma file via MCP. Requires prior `auth login`.

```bash
figma-kit new-file "My Landing Page"
```

---

## Layer 1 — Primitives

### `node create`

**Usage:** `figma-kit node create <type>`

**Description:** Emit JS to create a node (`frame`, `rect`, `rectangle`, `text`, `ellipse`, `line`, `polygon`, `star`, `vector`, `component`, `component-set`) and append it to the current page.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--name` | `-n` | string | `Untitled` | Node name. |
| `--width` | `-w` | int | `400` | Width (not applied to `text` / `line`). |
| `--height` | | int | `300` | Height (not applied to `text` / `line`). |
| `--x` | | int | `0` | X position. |
| `--y` | | int | `0` | Y position. |

```bash
figma-kit node create frame -n Hero -w 1440 -h 900 -t default -p 0
```

---

### `node clone`

**Usage:** `figma-kit node clone <nodeId>`

**Description:** Duplicate a node with offset.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dx` | int | `100` | X offset from original. |
| `--dy` | int | `0` | Y offset from original. |

```bash
figma-kit node clone "123:456" --dx 40 --dy 0
```

---

### `node delete`

**Usage:** `figma-kit node delete <nodeId>`

**Description:** Remove a node.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit node delete "123:456"
```

---

### `node move`

**Usage:** `figma-kit node move <nodeId>`

**Description:** Set absolute position.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--x` | int | `0` | X position. |
| `--y` | int | `0` | Y position. |

```bash
figma-kit node move "123:456" --x 100 --y 200 -p 0
```

---

### `node resize`

**Usage:** `figma-kit node resize <nodeId>`

**Description:** Resize node.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--width` | `-w` | int | `400` | Width. |
| `--height` | | int | `300` | Height. |

```bash
figma-kit node resize "123:456" -w 320 -h 240
```

---

### `node rename`

**Usage:** `figma-kit node rename <nodeId>`

**Description:** Rename a node.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--name` | `-n` | string | *(required)* | New name. |

```bash
figma-kit node rename "123:456" -n "Card / Primary"
```

---

### `node reparent`

**Usage:** `figma-kit node reparent <nodeId> <parentId>`

**Description:** Move node under another parent.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit node reparent "123:456" "789:012"
```

---

### `node lock`

**Usage:** `figma-kit node lock <nodeId>`

**Description:** Lock node (`locked = true`) unless `--unlock`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--unlock` | bool | `false` | Unlock instead of lock. |

```bash
figma-kit node lock "123:456"
figma-kit node lock "123:456" --unlock
```

---

### `node visible`

**Usage:** `figma-kit node visible <nodeId>`

**Description:** Show node unless `--hide`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--hide` | bool | `false` | Hide instead of show. |

```bash
figma-kit node visible "123:456" --hide
```

---

### `node order`

**Usage:** `figma-kit node order <direction> <nodeId>`

**Description:** Change layer order. Direction: `front`, `back`, `forward`, `backward`.

```bash
figma-kit node order front "123:456"
```

---

### `node group`

**Usage:** `figma-kit node group <nodeId> [nodeId...]`

**Description:** Group two or more nodes together.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` / `-n` | string | `Group` | Group name. |

```bash
figma-kit node group "1:2" "1:3" "1:4" -n "Header"
```

---

### `node ungroup`

**Usage:** `figma-kit node ungroup <nodeId>`

**Description:** Ungroup a group node, reparenting children to the group's parent.

```bash
figma-kit node ungroup "5:6"
```

---

### `node component`

**Usage:** `figma-kit node component <nodeId>`

**Description:** Convert an existing frame/node into a Figma component using `createComponentFromNode`.

```bash
figma-kit node component "2:3"
```

---

### `node flatten`

**Usage:** `figma-kit node flatten <nodeId>`

**Description:** Flatten a node subtree into a single vector. Useful for export prep.

```bash
figma-kit node flatten "4:5"
```

---

### `node boolean`

**Usage:** `figma-kit node boolean <operation> <nodeA> <nodeB>`

**Description:** Perform a boolean shape operation on two nodes. Operations: `union`, `subtract`, `intersect`, `exclude`.

```bash
figma-kit node boolean union "1:2" "1:3"
figma-kit node boolean subtract "1:2" "1:3"
figma-kit node boolean intersect "1:2" "1:3"
```

---

### `node svg`

**Usage:** `figma-kit node svg <path-data>`

**Description:** Create a vector node from SVG path data.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fill` | string | `""` | Fill color hex |
| `--stroke` | string | `""` | Stroke color hex |
| `--size` | int | `100` | Viewbox size |

```bash
figma-kit node svg "M10 10 L90 90 L10 90 Z" --fill "#3B82F6"
figma-kit node svg "M50 0 L100 100 L0 100 Z" --fill "#EF4444" --size 200
```

---

### `node variant-set`

**Usage:** `figma-kit node variant-set <componentId>`

**Description:** Create a component set with variant rows from a base component.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--variants` | string (JSON) | `[]` | Array of variant definitions |

```bash
figma-kit node variant-set "5:10" --variants '[{"name":"Size=Small","width":100},{"name":"Size=Large","width":200}]'
```

---

### `style fill`

**Usage:** `figma-kit style fill <nodeId>`

**Description:** Set solid fill.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--solid` | string | `#FFFFFF` | Hex color. |
| `--opacity` | float | `1.0` | Opacity 0–1. |

```bash
figma-kit style fill "123:456" --solid "#1A1D24" --opacity 0.9
```

---

### `style stroke`

**Usage:** `figma-kit style stroke <nodeId>`

**Description:** Set stroke color, weight, alignment.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--color` | string | `#FFFFFF` | Stroke hex. |
| `--weight` | int | `1` | Stroke weight. |
| `--align` | string | `inside` | `inside`, `outside`, or `center` (mapped to Figma enums). |

```bash
figma-kit style stroke "123:456" --color "#64748B" --weight 2 --align center
```

---

### `style effect`

**Usage:** `figma-kit style effect <nodeId>`

**Description:** Apply blur or generic drop shadow. If `--blur` > 0, uses layer or background blur; else if `--shadow` non-empty, applies a fixed drop shadow preset in generated JS.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--shadow` | string | `""` | Non-empty triggers drop shadow branch. |
| `--blur` | int | `0` | Blur radius. |
| `--blur-type` | string | `layer` | `layer` or `background`. |

```bash
figma-kit style effect "123:456" --blur 16 --blur-type background
```

---

### `style corner`

**Usage:** `figma-kit style corner <nodeId>`

**Description:** Uniform `--radius`, or per-corner if any of `--tl` `--tr` `--br` `--bl` is **changed** from default.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--radius` | int | `0` | Uniform corner radius. |
| `--tl` | int | `0` | Top-left. |
| `--tr` | int | `0` | Top-right. |
| `--br` | int | `0` | Bottom-right. |
| `--bl` | int | `0` | Bottom-left. |

```bash
figma-kit style corner "123:456" --radius 12
```

---

### `style blend`

**Usage:** `figma-kit style blend <nodeId>`

**Description:** Set `blendMode` and node opacity.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--mode` | string | `NORMAL` | Figma blend mode string. |
| `--opacity` | float | `1.0` | Opacity 0–1. |

```bash
figma-kit style blend "123:456" --mode MULTIPLY --opacity 0.85
```

---

### `style gradient`

**Usage:** `figma-kit style gradient <nodeId>`

**Description:** Linear or radial gradient fill.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `linear` | `linear` or `radial`. |
| `--angle` | int | `0` | Degrees (declared; transform uses fixed matrix in generator). |
| `--stops` | string | `0:#000000,1:#FFFFFF` | Stops as `pos:hex,...`. |

```bash
figma-kit style gradient "123:456" --type radial --stops "0:#3B82F6,1:#14B8A6"
```

---

### `style clip`

**Usage:** `figma-kit style clip <nodeId>`

**Description:** Enable `clipsContent` on a frame-like node unless `--off`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--off` | bool | `false` | Disable clipping. |

```bash
figma-kit style clip "123:456"
```

---

### `style apply`

**Usage:** `figma-kit style apply <nodeId>`

**Description:** Apply a named local style (paint, text, or effect) to a node by looking up the style name.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | *(required)* | Style name to apply. |
| `--type` | string | `fill` | Style type: `fill`, `text`, or `effect`. |

```bash
figma-kit style apply "123:456" --name "Primary Blue" --type fill
```

---

### `text create`

**Usage:** `figma-kit text create`

**Description:** Create a text layer with font load.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--content` | | string | `Text` | Text body. |
| `--font` | | string | `Inter` | Font family. |
| `--weight` | | string | `Regular` | Font style. |
| `--size` | | int | `16` | Font size. |
| `--color` | | string | `#FFFFFF` | Fill hex. |
| `--parent` | | string | `""` | Parent node id (optional). |
| `--x` | | int | `0` | X. |
| `--y` | | int | `0` | Y. |
| `--width` | `-w` | int | `0` | Text box width; `0` = auto. |
| `--line-height` | | int | `0` | Line height in px; `0` = default auto behavior. |
| `--letter-spacing` | | float | `0` | Letter spacing in pixels. |
| `--align` | | string | `""` | `LEFT`, `CENTER`, `RIGHT`, or `JUSTIFIED` when set. |
| `--auto-resize` | | string | `""` | `NONE`, `WIDTH_AND_HEIGHT`, or `HEIGHT` when set (with width, defaults to `HEIGHT` if unset). |

```bash
figma-kit text create --content "Hello" --font Inter --weight Bold --size 24 --color "#0F172A"
figma-kit text create --content "Subhead" --line-height 28 --letter-spacing 0.5 --align CENTER --auto-resize WIDTH_AND_HEIGHT
```

---

### `text edit`

**Usage:** `figma-kit text edit <nodeId>`

**Description:** Replace `characters` on a text node.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--content` | string | *(required)* | New string. |

```bash
figma-kit text edit "123:456" --content "Updated label"
```

---

### `text style`

**Usage:** `figma-kit text style <nodeId>`

**Description:** Update typography; only emits lines for flags explicitly **changed**.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--size` | int | `16` | Font size. |
| `--lh` | int | `0` | Line height (px). |
| `--ls` | int | `0` | Letter spacing (%). |
| `--align` | string | `""` | `LEFT`, `CENTER`, `RIGHT`, `JUSTIFIED`. |

```bash
figma-kit text style "123:456" --size 18 --align CENTER
```

---

### `text range`

**Usage:** `figma-kit text range <nodeId>`

**Description:** Mixed styles on substring (Inter family for weight range).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--start` | int | `0` | Start index. |
| `--end` | int | `0` | End index. |
| `--weight` | string | `""` | Font style for range. |
| `--color` | string | `""` | Hex for range fills. |

```bash
figma-kit text range "123:456" --start 0 --end 4 --weight Bold --color "#F59E0B"
```

---

### `text list-fonts`

**Usage:** `figma-kit text list-fonts`

**Description:** JS listing available font families in the file.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No command-specific flags. |

```bash
figma-kit text list-fonts
```

---

### `text load-fonts`

**Usage:** `figma-kit text load-fonts`

**Description:** Emit loops to load all styles for given families.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--families` | string | `Inter,Geist Mono` | Comma-separated families. |

```bash
figma-kit text load-fonts --families "Inter,Roboto"
```

---

### `layout auto`

**Usage:** `figma-kit layout auto <nodeId>`

**Description:** Configure auto-layout on a frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | `VERTICAL` | `HORIZONTAL` or `VERTICAL`. |
| `--gap` | int | `16` | Item spacing. |
| `--pad` | int | `0` | Uniform padding. |
| `--align` | string | `""` | Counter-axis align (`MIN`, `CENTER`, `MAX`, `BASELINE`). |
| `--wrap` | bool | `false` | Enable `WRAP`. |

```bash
figma-kit layout auto "123:456" --dir HORIZONTAL --gap 12 --pad 16 --wrap
```

---

### `layout grid`

**Usage:** `figma-kit layout grid <nodeId>`

**Description:** Column layout grid on frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--columns` | int | `12` | Column count. |
| `--gutter` | int | `24` | Gutter. |
| `--margin` | int | `80` | Offset. |

```bash
figma-kit layout grid "123:456" --columns 8 --gutter 32
```

---

### `layout constraints`

**Usage:** `figma-kit layout constraints <nodeId>`

**Description:** Set resizing constraints.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--h` | string | `MIN` | Horizontal: `MIN`, `CENTER`, `MAX`, `STRETCH`, `SCALE`. |
| `--v` | string | `MIN` | Vertical: same set. |

```bash
figma-kit layout constraints "123:456" --h STRETCH --v MIN
```

---

### `layout sizing`

**Usage:** `figma-kit layout sizing <nodeId>`

**Description:** Set `layoutSizingHorizontal` / `layoutSizingVertical` when flags non-empty.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--width` | `-w` | string | `""` | `FIXED`, `HUG`, `FILL`. |
| `--height` | | string | `""` | `FIXED`, `HUG`, `FILL`. |

```bash
figma-kit layout sizing "123:456" -w FILL --height HUG
```

---

### `layout align`

**Usage:** `figma-kit layout align <nodeId>`

**Description:** Auto-layout child alignment on primary/counter axis when set.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--primary` | string | `""` | `MIN`, `CENTER`, `MAX`, `SPACE_BETWEEN`. |
| `--counter` | string | `""` | `MIN`, `CENTER`, `MAX`, `BASELINE`. |

```bash
figma-kit layout align "123:456" --primary SPACE_BETWEEN --counter CENTER
```

---

### `layout distribute`

**Usage:** `figma-kit layout distribute <nodeIds>`

**Description:** Even spacing; `nodeIds` is a **single** comma-separated argument.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--axis` | string | `H` | `H` / `V` (or `VERTICAL`). |
| `--gap` | int | `24` | Gap between nodes. |

```bash
figma-kit layout distribute "id1,id2,id3" --axis H --gap 16
```

---

## Layer 2 — Patterns

Composable **`ui`** and **`card`** commands accept **`--parent`** with either a Figma node id or a compose expression **`_results[N]`** (emitted as JS, not wrapped in `getNodeByIdAsync`). Use this to nest outputs inside a prior compose step.

### `card glass`

**Usage:** `figma-kit card glass`

**Description:** Glassmorphism card using embedded `G()` helper + full helpers.

| Flag | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| `--preset` | | string | `default` | `subtle`, `default`, `strong`, `pill`. |
| `--width` | `-w` | int | `320` | Width. |
| `--height` | | int | `200` | Height. |
| `--parent` | | string | `""` | Optional parent id. |

```bash
figma-kit card glass --preset strong -w 400 -h 240 -t default -p 0
```

---

### `card solid`

**Usage:** `figma-kit card solid`

**Description:** Solid card frame with optional border and shadow preset.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--bg` | string | `#1A1D24` | Background hex. |
| `--border` | string | `""` | Border stroke hex (optional). |
| `--shadow` | string | `""` | `sm`, `md`, or `lg`. |
| `--radius` | int | `16` | Corner radius. |

```bash
figma-kit card solid --bg "#0B1020" --border "#334155" --shadow md
```

---

### `card gradient`

**Usage:** `figma-kit card gradient`

**Description:** Linear gradient card (fixed 320×200 frame in generator).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from` | string | `#3B5BFF` | Start hex. |
| `--to` | string | `#14B8A6` | End hex. |
| `--angle` | float | `135` | Angle in degrees. |

```bash
figma-kit card gradient --from "#6366F1" --to "#EC4899" --angle 90
```

---

### `card image`

**Usage:** `figma-kit card image`

**Description:** Fetch image by URL, image fill, optional overlay and title.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--url` | string | *(required)* | Image URL. |
| `--overlay` | string | `""` | `dark` or `light`. |
| `--title` | string | `""` | Title text (loads Inter Semi Bold). |

```bash
figma-kit card image --url "https://example.com/hero.jpg" --overlay dark --title "Launch"
```

---

### `card bento`

**Usage:** `figma-kit card bento`

**Description:** Grid of framed cells.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--cols` | int | `3` | Columns. |
| `--rows` | int | `2` | Rows. |
| `--gap` | int | `16` | Gap in px. |

```bash
figma-kit card bento --cols 4 --rows 3 --gap 12
```

---

### `card neumorphic`

**Usage:** `figma-kit card neumorphic`

**Description:** Soft UI (neumorphism) card with inset/outset shadow pairs.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Card` | Card title |
| `--desc` | string | `""` | Description text |
| `--depth` | string | `medium` | `shallow`, `medium`, `deep` |
| `--inset` | bool | `false` | Inset (pressed) shadow style |

```bash
figma-kit card neumorphic --title "Settings" --depth deep
figma-kit card neumorphic --title "Volume" --inset
```

---

### `card clay`

**Usage:** `figma-kit card clay`

**Description:** Claymorphism / puffy 3D card with soft shadows and rounded corners.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Card` | Card title |
| `--desc` | string | `""` | Description text |
| `--color` | string | `""` | Accent color hex (overrides theme primary) |

```bash
figma-kit card clay --title "Welcome" --color "#A78BFA"
```

---

### `card outline`

**Usage:** `figma-kit card outline`

**Description:** Ghost / outline card with optional glow border effect.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Card` | Card title |
| `--desc` | string | `""` | Description text |
| `--glow-color` | string | `""` | Glow border color hex |
| `--glow-spread` | float | `4` | Glow spread radius |

```bash
figma-kit card outline --title "API Docs" --glow-color "#3B82F6" --glow-spread 8
```

---

### `ui button`

**Usage:** `figma-kit ui button`

**Description:** Themed auto-layout button (uses active theme tokens).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--variant` | string | `primary` | `primary`, `secondary`, `ghost`, `destructive`, `outline`. |
| `--label` | string | `Button` | Label text. |
| `--size` | string | `md` | `sm`, `md`, `lg`. |

```bash
figma-kit ui button --variant outline --label "Continue" --size sm -t noir
```

---

### `ui input`

**Usage:** `figma-kit ui input`

**Description:** Labeled field + placeholder + type hint row.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--label` | string | `Email` | Label. |
| `--placeholder` | string | `you@example.com` | Placeholder. |
| `--type` | string | `text` | `text`, `email`, `password`. |

```bash
figma-kit ui input --label "Password" --placeholder "••••••••" --type password
```

---

### `ui badge`

**Usage:** `figma-kit ui badge`

**Description:** Pill badge.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | `New` | Badge copy (single-badge mode). |
| `--color` | string | `blue` | `blue`, `green`, `red`, `yellow`, `gray` (single-badge mode). |
| `--items` | string (JSON) | `""` | Batch mode: `[{"text":"…","color":"blue"},…]` builds a horizontal row; omit `--text`/`--color` when set. |
| `--parent` | string | `""` | Parent node id or `_results[N]` in compose. |

```bash
figma-kit ui badge --text "Beta" --color yellow
figma-kit ui badge --items '[{"text":"v2.1","color":"blue"},{"text":"New","color":"green"}]'
```

---

### `ui avatar`

**Usage:** `figma-kit ui avatar`

**Description:** Circle avatar with initials.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--initials` | string | `AK` | 1–3 characters. |
| `--size` | int | `40` | Diameter (min 16). |

```bash
figma-kit ui avatar --initials "MJ" --size 56
```

---

### `ui divider`

**Usage:** `figma-kit ui divider`

**Description:** 1px-thick horizontal or vertical divider.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | `H` | `H` or `V`. |
| `--length` | int | `240` | Length in px. |
| `--color` | string | `""` | Empty = stroke token; `muted` uses muted fill. |

```bash
figma-kit ui divider --dir V --length 320 --color muted
```

---

### `ui icon`

**Usage:** `figma-kit ui icon`

**Description:** Icon placeholder (card + inner glyph rect).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--shape` | string | `square` | `circle` or `square`. |
| `--size` | int | `32` | Outer size (min 8). |
| `--color` | string | `#94A3B8` | Glyph hex. |

```bash
figma-kit ui icon --shape circle --size 24 --color "#38BDF8"
```

---

### `ui progress`

**Usage:** `figma-kit ui progress`

**Description:** Progress bar.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--value` | int | `42` | 0–100. |
| `--width` | int | `200` | Track width (min 40). |

```bash
figma-kit ui progress --value 75 --width 280
```

---

### `ui toggle`

**Usage:** `figma-kit ui toggle`

**Description:** Switch track + thumb.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--state` | string | `off` | `on` or `off`. |
| `--size` | string | `md` | `sm`, `md`, `lg`. |

```bash
figma-kit ui toggle --state on --size lg
```

---

### `ui tooltip`

**Usage:** `figma-kit ui tooltip`

**Description:** Tooltip with bubble + caret.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | `Copied!` | Tooltip copy. |
| `--position` | string | `top` | `top`, `bottom`, `left`, `right`. |

```bash
figma-kit ui tooltip --text "Saved" --position bottom
```

---

### `ui stat`

**Usage:** `figma-kit ui stat`

**Description:** Value + label + optional trend glyph.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--value` | string | `4.2x` | Main value (single-stat mode). |
| `--label` | string | `Faster` | Caption (single-stat mode). |
| `--trend` | string | `""` | `up`, `down`, or `neutral` (single-stat mode only). |
| `--items` | string (JSON) | `""` | Batch mode: `[{"value":"…","label":"…"},…]` as a horizontal stats row (no per-item `trend`; use single-stat mode for `--trend`). |
| `--parent` | string | `""` | Parent node id or `_results[N]` in compose. |

```bash
figma-kit ui stat --value "12%" --label "MoM" --trend up
figma-kit ui stat --items '[{"value":"150+","label":"Commands"},{"value":"4.2x","label":"Faster"}]'
```

---

### `ui table`

**Usage:** `figma-kit ui table`

**Description:** Table from JSON rows file (read at CLI time); emits JS with embedded data.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--data` | string | `./data.json` | Path to JSON array of objects. |
| `--cols` | string | `Name,Role,Status` | Comma-separated keys matching JSON. |

```bash
figma-kit ui table --data ./team.json --cols "Name,Role,Status" -t default
```

---

### `ui nav`

**Usage:** `figma-kit ui nav`

**Description:** Top bar or sidebar link stack.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string | `Home,Products,Pricing` | Comma labels. |
| `--style` | string | `topbar` | `topbar` or `sidebar`. |

```bash
figma-kit ui nav --items "Docs,Blog,Contact" --style sidebar
```

---

### `ui footer`

**Usage:** `figma-kit ui footer`

**Description:** Footer columns with placeholder links.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--cols` | int | `3` | 1–6 columns. |
| `--copyright` | string | `""` | Optional copyright line. |

```bash
figma-kit ui footer --cols 4 --copyright "© 2026 Acme"
```

---

### `ui checkbox`

**Usage:** `figma-kit ui checkbox`

**Description:** Themed checkbox with label.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--label` | string | `Option` | Checkbox label |
| `--checked` | bool | `false` | Pre-checked state |

```bash
figma-kit ui checkbox --label "Accept terms" --checked
```

---

### `ui radio`

**Usage:** `figma-kit ui radio`

**Description:** Themed radio button with label.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--label` | string | `Option` | Radio label |
| `--selected` | bool | `false` | Pre-selected state |

```bash
figma-kit ui radio --label "Monthly" --selected
```

---

### `ui tabs`

**Usage:** `figma-kit ui tabs`

**Description:** Tab bar with selectable items.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string | `Tab 1,Tab 2,Tab 3` | Comma-separated tab labels |
| `--active` | int | `0` | Index of active tab |

```bash
figma-kit ui tabs --items "Overview,Features,Pricing" --active 1
```

---

### `ui dropdown`

**Usage:** `figma-kit ui dropdown`

**Description:** Dropdown / select component.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--options` | string | `Option 1,Option 2,Option 3` | Comma-separated options |
| `--selected` | string | `""` | Pre-selected value |

```bash
figma-kit ui dropdown --options "Draft,Published,Archived" --selected "Draft"
```

---

### `ui breadcrumb`

**Usage:** `figma-kit ui breadcrumb`

**Description:** Breadcrumb navigation trail.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string | `Home,Products,Detail` | Comma-separated breadcrumb items |

```bash
figma-kit ui breadcrumb --items "Home,Dashboard,Settings,Profile"
```

---

### `ui skeleton`

**Usage:** `figma-kit ui skeleton`

**Description:** Skeleton loading placeholder.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--rows` | int | `3` | Number of skeleton rows |

```bash
figma-kit ui skeleton --rows 5
```

---

### `ui chip`

**Usage:** `figma-kit ui chip`

**Description:** Tag / filter chip component.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--label` | string | `Tag` | Chip label |
| `--variant` | string | `filled` | `filled`, `outlined`, `tonal` |
| `--dismissible` | bool | `false` | Show dismiss icon |

```bash
figma-kit ui chip --label "React" --variant outlined
figma-kit ui chip --label "Filter" --dismissible
```

---

### `ui toast`

**Usage:** `figma-kit ui toast`

**Description:** Notification / toast popup.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--message` | string | `Notification` | Toast message |
| `--type` | string | `info` | `info`, `success`, `warning`, `error` |

```bash
figma-kit ui toast --message "Changes saved" --type success
figma-kit ui toast --message "Connection lost" --type error
```

---

### `ui modal`

**Usage:** `figma-kit ui modal`

**Description:** Modal / dialog with overlay.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Modal` | Modal title |
| `--body` | string | `""` | Body content text |
| `--primary-action` | string | `Confirm` | Primary button label |
| `--secondary-action` | string | `Cancel` | Secondary button label |

```bash
figma-kit ui modal --title "Delete item?" --body "This action cannot be undone." --primary-action "Delete"
```

---

### `ui card-list`

**Usage:** `figma-kit ui card-list`

**Description:** Vertical list of data cards.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string (JSON) | `[]` | Array of `{title, desc}` objects |

```bash
figma-kit ui card-list --items '[{"title":"Item 1","desc":"Description"},{"title":"Item 2","desc":"Details"}]'
```

---

### `ui sidebar`

**Usage:** `figma-kit ui sidebar`

**Description:** Sidebar navigation panel.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string | `Dashboard,Settings,Profile` | Comma-separated nav items |
| `--active` | int | `0` | Active item index |

```bash
figma-kit ui sidebar --items "Home,Projects,Team,Settings" --active 1
```

---

### `ui avatar-group`

**Usage:** `figma-kit ui avatar-group`

**Description:** Overlapping avatar stack.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--count` | int | `4` | Number of avatars |
| `--size` | int | `40` | Avatar diameter in px |
| `--overlap` | int | `12` | Overlap in px |

```bash
figma-kit ui avatar-group --count 5 --size 48 --overlap 16
```

---

### `ui rating`

**Usage:** `figma-kit ui rating`

**Description:** Star rating display.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--value` | float | `3.5` | Rating value (0–5) |
| `--max` | int | `5` | Maximum stars |

```bash
figma-kit ui rating --value 4.5
```

---

### `ui search`

**Usage:** `figma-kit ui search`

**Description:** Search input with icon.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--placeholder` | string | `Search...` | Placeholder text |

```bash
figma-kit ui search --placeholder "Search components..."
```

---

### `ui section`

**Usage:** `figma-kit ui section`

**Description:** Centered vertical section wrapper (label, heading, subtitle) with auto-layout — **recommended** for page blocks inside `compose`. Children from later steps attach with `--parent _results[N]` where `N` is this step’s index.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Section Title` | Main heading text. |
| `--label` | string | `""` | Monospace eyebrow above heading (stored uppercased). |
| `--subtitle` | string | `""` | Subtitle below heading. |
| `--label-color` | string | `""` | Label fill: hex or theme token; default uses theme `BL`. |
| `--width` | int | `1440` | Section frame width. |
| `--padding` | int | `80` | Uniform padding. |
| `--spacing` | int | `24` | Vertical spacing between stacked items. |
| `--divider` | bool | `false` | Draw a top divider line inside the section. |
| `--parent` | string | `""` | Parent node id or `_results[N]` in compose. |

```bash
figma-kit ui section -t noir --title "Features" --label "PRODUCT" --subtitle "Everything you need"
figma-kit compose "ui section --title Features" "card glass --parent _results[0] --title Card1"
```

---

### `ui pagination`

**Usage:** `figma-kit ui pagination`

**Description:** Page number navigation bar.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--pages` | int | `5` | Total page count |
| `--current` | int | `1` | Current active page |

```bash
figma-kit ui pagination --pages 10 --current 3
```

---

### `ui color-picker`

**Usage:** `figma-kit ui color-picker`

**Description:** Color swatch grid.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--colors` | string | `""` | Comma-separated hex colors (defaults to theme palette) |
| `--cols` | int | `6` | Grid columns |

```bash
figma-kit ui color-picker --colors "#FF0000,#00FF00,#0000FF,#FFFF00" --cols 4
```

---

### `ui hero`

**Usage:** `figma-kit ui hero`

**Description:** Complete hero section with headline, subtitle, CTA button, and optional badge.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | `Headline` | Hero headline |
| `--subtitle` | string | `""` | Subtitle text |
| `--cta` | string | `Get Started` | CTA button label |
| `--badge` | string | `""` | Optional top badge text |

```bash
figma-kit ui hero -t noir --title "Ship Faster" --subtitle "Build with AI" --cta "Start Free" --badge "New"
```

---

### `ui pricing`

**Usage:** `figma-kit ui pricing`

**Description:** Pricing table with tier cards. Highlighted tier gets accent border.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--tiers` | string (JSON) | *(3 defaults)* | Array of `{name, price, features[], highlighted?}` |

```bash
figma-kit ui pricing -t noir --tiers '[{"name":"Free","price":"$0","features":["5 projects"]},{"name":"Pro","price":"$29","highlighted":true,"features":["Unlimited","Priority support"]}]'
```

---

### `ui feature-grid`

**Usage:** `figma-kit ui feature-grid`

**Description:** Grid of feature cards with icon placeholders.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--cols` | int | `3` | Number of columns |
| `--features` | string (JSON) | *(3 defaults)* | Array of `{title, desc}` |

```bash
figma-kit ui feature-grid -t noir --cols 3
```

---

### `ui testimonial`

**Usage:** `figma-kit ui testimonial`

**Description:** Quote / testimonial card with avatar, name, and star rating.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | `Jane Doe` | Author name |
| `--role` | string | `""` | Author role/title |
| `--quote` | string | `""` | Quote text |
| `--rating` | int | `5` | Star rating (0–5) |

```bash
figma-kit ui testimonial -t noir --name "Jane" --quote "Changed everything" --rating 5
```

---

### `ui timeline`

**Usage:** `figma-kit ui timeline`

**Description:** Vertical timeline with date markers and event descriptions.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--entries` | string (JSON) | *(3 defaults)* | Array of `{date, title, desc?}` |

```bash
figma-kit ui timeline -t noir --entries '[{"date":"Jan 2024","title":"Launch"},{"date":"Mar 2024","title":"1K users"}]'
```

---

### `ui stepper`

**Usage:** `figma-kit ui stepper`

**Description:** Step progress indicator (wizard / form steps).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--steps` | string | `Step 1,Step 2,Step 3` | Comma-separated step labels |
| `--current` | int | `1` | Current active step (1-based) |

```bash
figma-kit ui stepper --steps "Account,Profile,Billing,Confirm" --current 2
```

---

### `ui accordion`

**Usage:** `figma-kit ui accordion`

**Description:** Expandable FAQ / accordion section.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--items` | string (JSON) | *(3 defaults)* | Array of `{question, answer}` |

```bash
figma-kit ui accordion -t noir --items '[{"question":"Is it free?","answer":"Yes, free for individuals under BSL 1.1."}]'
```

---

### `fx` — compose chaining (`--last`)

Every **`fx`** subcommand accepts **`--last`**. In **`figma-kit compose`**, pass **`--last`** instead of a positional `<nodeId>` (or `<parentId>` where applicable) to target **`_results[_results.length - 1]`** — the main node from the immediately previous step. Standalone CLI use still requires the explicit id argument.

---

### `fx glow`

**Usage:** `figma-kit fx glow <nodeId>` (or `figma-kit fx glow --last` in compose)

**Description:** Layered radial gradient glows; optional hex tint overrides inner color.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--position` | string | `subtle` | `topRight`, `center`, `subtle`, `cta`. |
| `--intensity` | float | `1` | Multiplier for glow alphas. |
| `--color` | string | `""` | Optional `#RRGGBB` tint. |

```bash
figma-kit fx glow "123:456" --position cta --intensity 1.2 --color "#6366F1"
```

---

### `fx mesh`

**Usage:** `figma-kit fx mesh <nodeId>` (or `--last` in compose)

**Description:** Mesh-like stacked radial fills.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--points` | int | `5` | Count (min 2). |
| `--palette` | string | `#2563eb,#14b8a6,#8b5cf6` | Comma hex or names (`blue`, `teal`, …). |

```bash
figma-kit fx mesh "123:456" --points 7 --palette "purple,pink,orange"
```

---

### `fx noise`

**Usage:** `figma-kit fx noise <nodeId>` (or `--last` in compose)

**Description:** Overlay rectangle with dithered gradient (child of target).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--opacity` | float | `0.35` | Strength 0–1. |

```bash
figma-kit fx noise "123:456" --opacity 0.25
```

---

### `fx vignette`

**Usage:** `figma-kit fx vignette <nodeId>` (or `--last` in compose)

**Description:** Radial darkening overlay child.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--strength` | float | `0.6` | Vignette strength. |

```bash
figma-kit fx vignette "123:456" --strength 0.8
```

---

### `fx grain`

**Usage:** `figma-kit fx grain <nodeId>` (or `--last` in compose)

**Description:** Film grain overlay.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--amount` | string | `medium` | `light`, `medium`, or `heavy`. |

```bash
figma-kit fx grain "123:456" --amount light
```

---

### `fx blur-bg`

**Usage:** `figma-kit fx blur-bg <nodeId>` (or `--last` in compose)

**Description:** Frosted rectangle + background blur child.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--radius` | int | `24` | Blur radius. |
| `--tint` | string | `""` | Optional `rgba(r,g,b,a)` (0–1 or 0–255 channels). |

```bash
figma-kit fx blur-bg "123:456" --radius 32 --tint "rgba(255,255,255,0.25)"
```

---

### `fx accent-bar`

**Usage:** `figma-kit fx accent-bar <parentId>` (or `--last` in compose)

**Description:** Gradient bar child inside parent.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: use previous step’s result as parent instead of `<parentId>`. |
| `--from` | string | `#3366FF` | Start hex. |
| `--to` | string | `#14B8A6` | End hex. |
| `--w` | int | `240` | Bar width. |
| `--h` | int | `4` | Bar height. |
| `--x` | int | `0` | X in parent. |
| `--y` | int | `0` | Y in parent. |

```bash
figma-kit fx accent-bar "123:456" --w 400 --h 6 --from "#F97316" --to "#EF4444"
```

---

### `fx shadow`

**Usage:** `figma-kit fx shadow <nodeId>` (or `--last` in compose)

**Description:** Replace `effects` with one shadow preset.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--preset` | string | `md` | `sm`, `md`, `lg`, `xl`, `glow`, `inner`. |

```bash
figma-kit fx shadow "123:456" --preset xl
```

---

### `fx parallax-layer`

**Usage:** `figma-kit fx parallax-layer <parentId>` (or `--last` in compose)

**Description:** Stack of inset frames inside parent.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: use previous step’s result as parent instead of `<parentId>`. |
| `--layers` | int | `3` | Count (min 2). |

```bash
figma-kit fx parallax-layer "123:456" --layers 5
```

---

### `fx aurora`

**Usage:** `figma-kit fx aurora <nodeId>` (or `--last` in compose)

**Description:** Northern lights gradient overlay effect with layered ellipses.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--palette` | string | `default` | Color palette: `default`, `sunset`, `ocean`, `forest` |

```bash
figma-kit fx aurora "123:456" --palette sunset
```

---

### `fx morph`

**Usage:** `figma-kit fx morph <nodeId>` (or `--last` in compose)

**Description:** Organic blob / morphism shapes inside a parent frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--count` | int | `3` | Number of blobs |
| `--spread` | float | `0.8` | Spread factor (0–1) |

```bash
figma-kit fx morph "123:456" --count 5 --spread 0.6
```

---

### `fx gradient-border`

**Usage:** `figma-kit fx gradient-border <nodeId>` (or `--last` in compose)

**Description:** Simulated gradient stroke by layering a slightly larger gradient-filled frame behind the target.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--from` | string | theme primary | Start color hex |
| `--to` | string | theme accent | End color hex |
| `--width` | int | `2` | Border thickness in px |

```bash
figma-kit fx gradient-border "123:456" --from "#3B82F6" --to "#8B5CF6"
```

---

### `fx spotlight`

**Usage:** `figma-kit fx spotlight <nodeId>` (or `--last` in compose)

**Description:** Circular radial highlight effect.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--x` | float | `0.5` | Horizontal position (0–1) |
| `--y` | float | `0.3` | Vertical position (0–1) |
| `--intensity` | float | `0.4` | Opacity of the highlight |

```bash
figma-kit fx spotlight "123:456" --x 0.7 --y 0.2 --intensity 0.6
```

---

### `fx pattern`

**Usage:** `figma-kit fx pattern <nodeId>` (or `--last` in compose)

**Description:** Repeating geometric background patterns.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--last` | bool | `false` | Compose: target previous step’s result instead of `<nodeId>`. |
| `--type` | string | `dots` | Pattern type: `dots`, `lines`, `crosses`, `diagonal`, `grid` |
| `--scale` | float | `1.0` | Scale factor |
| `--opacity` | float | `0.1` | Pattern opacity |

```bash
figma-kit fx pattern "123:456" --type crosses --scale 1.5
figma-kit fx pattern "123:456" --type diagonal --opacity 0.15
```

---

## Layer 3 — Deliverables (`make`)

All `make` subcommands use the active theme (`-t`) and page (`-p`) in the generated preamble where applicable.

### `make carousel`

**Usage:** `figma-kit make carousel`

**Description:** LinkedIn-style slides from YAML (`slides` array) + embedded slide template.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--content` | string | *(required)* | YAML path. |
| `--slides` | int | `0` | Max slides; `0` = all. |

```bash
figma-kit make carousel --content ./deck.yml --slides 5 -t default
```

---

### `make instagram-post`

**Usage:** `figma-kit make instagram-post`

**Description:** 1080×1080 feed post frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `image` | `quote` or `image`. |
| `--content` | string | *(required)* | Headline / caption. |

```bash
figma-kit make instagram-post --type quote --content "Ship fast"
```

---

### `make instagram-story`

**Usage:** `figma-kit make instagram-story`

**Description:** 1080×1920 story from YAML (`title`, `subtitle`, …).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--content` | string | *(required)* | YAML path. |

```bash
figma-kit make instagram-story --content ./story.yml
```

---

### `make twitter-card`

**Usage:** `figma-kit make twitter-card`

**Description:** 1200×675 large card.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--headline` | string | *(required)* | Headline. |
| `--image` | string | `hero` | `hero` or `minimal`. |

```bash
figma-kit make twitter-card --headline "We raised our Series A" --image minimal
```

---

### `make facebook-cover`

**Usage:** `figma-kit make facebook-cover`

**Description:** 820×312 cover.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scheme` | string | `dark` | `dark` or `light` (cover palette; not the global theme flag). |

```bash
figma-kit make facebook-cover --scheme light
```

---

### `make youtube-thumb`

**Usage:** `figma-kit make youtube-thumb`

**Description:** 1280×720 thumbnail with face ellipse placeholder.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | *(required)* | Title on thumb. |
| `--face` | string | `center` | `left`, `right`, or `center`. |

```bash
figma-kit make youtube-thumb --title "Figma automation deep dive" --face left
```

---

### `make og-image`

**Usage:** `figma-kit make og-image`

**Description:** 1200×630 Open Graph frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | *(required)* | Title. |
| `--description` | string | *(required)* | Subtitle. |

```bash
figma-kit make og-image --title "Acme" --description "Design systems at scale"
```

---

### `make banner`

**Usage:** `figma-kit make banner`

**Description:** IAB banner frames from size keys + YAML (file must exist).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sizes` | string | `leaderboard,mrec` | Comma keys: `leaderboard`, `mrec`, `skyscraper`, `billboard`, `mobile`. |
| `--content` | string | *(required)* | YAML path. |

```bash
figma-kit make banner --sizes "leaderboard,mobile" --content ./ad.yml
```

---

### `make email-header`

**Usage:** `figma-kit make email-header`

**Description:** Email header strip (default height 120).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--width` | int | `600` | Width in px. |

```bash
figma-kit make email-header --width 640
```

---

### `make ad-set`

**Usage:** `figma-kit make ad-set`

**Description:** One frame per platform from campaign YAML.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--content` | string | *(required)* | `campaign.yml`. |
| `--platforms` | string | `linkedin,instagram,twitter` | Comma platform keys (`linkedin`, `instagram`, `twitter`, `facebook`). |

```bash
figma-kit make ad-set --content ./campaign.yml --platforms "linkedin,facebook"
```

---

### `make one-pager`

**Usage:** `figma-kit make one-pager`

**Description:** B2B one-pager via embedded `one-pager-print` template.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `letter` | `letter` or `A4` (metadata). |
| `--mode` | string | `print` | `print` or `digital`. |
| `--content` | string | *(required)* | `one-pager.yml`. |

```bash
figma-kit make one-pager --content ./one-pager.yml --format A4 --mode digital
```

---

### `make pitch-deck`

**Usage:** `figma-kit make pitch-deck`

**Description:** 1920×1080 slide frames in a row.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--slides` | int | `10` | Slide count. |
| `--template` | string | `saas` | `saas`, `agency`, or `startup`. |

```bash
figma-kit make pitch-deck --slides 12 --template agency
```

---

### `make case-study`

**Usage:** `figma-kit make case-study`

**Description:** Vertical section frames.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sections` | string | `overview,challenge,solution,results` | Comma section names. |

```bash
figma-kit make case-study --sections "problem,solution,metrics"
```

---

### `make proposal`

**Usage:** `figma-kit make proposal`

**Description:** Proposal cover + client + scope summary from YAML.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--client` | string | *(required)* | Client name. |
| `--scope` | string | *(required)* | `scope.yml` path. |

```bash
figma-kit make proposal --client "Globex" --scope ./scope.yml
```

---

### `make invoice`

**Usage:** `figma-kit make invoice`

**Description:** Invoice layout frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--template` | string | `modern` | `modern` or `minimal`. |

```bash
figma-kit make invoice --template minimal
```

---

### `make business-card`

**Usage:** `figma-kit make business-card`

**Description:** 1050×600 card.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | *(required)* | Name. |
| `--title` | string | *(required)* | Role / title. |

```bash
figma-kit make business-card --name "Alex Kim" --title "Design Lead"
```

---

### `make letterhead`

**Usage:** `figma-kit make letterhead`

**Description:** Letterhead document frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `letter` | `letter` or `A4`. |

```bash
figma-kit make letterhead --format A4
```

---

### `make contract`

**Usage:** `figma-kit make contract`

**Description:** Contract title page.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--title` | string | *(required)* | Contract title. |

```bash
figma-kit make contract --title "Master Services Agreement"
```

---

### `make storyboard`

**Usage:** `figma-kit make storyboard`

**Description:** Panels from YAML `scenes` + embedded storyboard template.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--content` | string | *(required)* | `storyboard.yml`. |
| `--scenes` | int | `0` | Max scenes; `0` = all. |

```bash
figma-kit make storyboard --content ./boards.yml --scenes 8
```

---

### `make styleframe`

**Usage:** `figma-kit make styleframe`

**Description:** Single 1920×1080 motion styleframe.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--mood` | string | `neutral` | `warm`, `cool`, or `neutral`. |
| `--scene` | string | *(required)* | Scene description. |

```bash
figma-kit make styleframe --mood cool --scene "City at dusk — neon rain"
```

---

### `make animatic`

**Usage:** `figma-kit make animatic`

**Description:** Wide timeline overview + strip of placeholder cells.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fps` | int | `24` | FPS label. |
| `--duration` | string | `30s` | Duration label. |

```bash
figma-kit make animatic --fps 30 --duration "0:45"
```

---

### `make transition-spec`

**Usage:** `figma-kit make transition-spec`

**Description:** UI transition spec card.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `fade` | `page-enter`, `fade`, or `slide`. |
| `--easing` | string | `ease-out` | Easing label text. |

```bash
figma-kit make transition-spec --type slide --easing "cubic-bezier(0.4, 0, 0.2, 1)"
```

---

### `make wireframe`

**Usage:** `figma-kit make wireframe`

**Description:** Low-fi shell rectangles.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `landing` | `landing`, `dashboard`, or `form`. |
| `--breakpoint` | string | `desktop` | `desktop`, `tablet`, or `mobile`. |

```bash
figma-kit make wireframe --type dashboard --breakpoint tablet
```

---

### `make screen`

**Usage:** `figma-kit make screen`

**Description:** Marketing screen with alternating section strips.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `landing` | `landing`, `pricing`, or `features`. |
| `--sections` | string | `hero,features,pricing,cta` | Comma section names. |

```bash
figma-kit make screen --type features --sections "hero,grid,testimonials,cta"
```

---

### `make dashboard`

**Usage:** `figma-kit make dashboard`

**Description:** Dashboard shell + widget grid.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--widgets` | string | `stat,chart,table,list` | Comma widget ids. |
| `--cols` | int | `2` | Column count. |

```bash
figma-kit make dashboard --widgets "stat,stat,chart" --cols 3
```

---

### `make form`

**Usage:** `figma-kit make form`

**Description:** Vertical form from schema file with `fields` array.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fields` | string | *(required)* | JSON or YAML path. |

```bash
figma-kit make form --fields ./form-schema.yml
```

---

### `make modal`

**Usage:** `figma-kit make modal`

**Description:** Dimmed overlay + centered modal.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--size` | string | `md` | `sm`, `md`, `lg`. |
| `--type` | string | `confirmation` | `confirmation`, `form`, or `alert`. |

```bash
figma-kit make modal --size lg --type alert
```

---

### `make empty-state`

**Usage:** `figma-kit make empty-state`

**Description:** Empty state block.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--message` | string | `Nothing here yet.` | Body copy. |

```bash
figma-kit make empty-state --message "No projects yet — create one to get started."
```

---

### `make error-page`

**Usage:** `figma-kit make error-page`

**Description:** Full-page error frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `404` | `404`, `500`, or `offline`. |

```bash
figma-kit make error-page --type offline
```

---

### `make onboarding`

**Usage:** `figma-kit make onboarding`

**Description:** Horizontal row of mobile-sized step frames.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--steps` | int | `3` | Number of screens. |

```bash
figma-kit make onboarding --steps 4
```

---

### `make settings`

**Usage:** `figma-kit make settings`

**Description:** Settings layout: nav column + detail panel.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sections` | string | `profile,notifications,billing,security` | Comma nav ids. |

```bash
figma-kit make settings --sections "account,privacy,billing"
```

---

### `make poster`

**Usage:** `figma-kit make poster`

**Description:** Large print poster frame (pixel dimensions vary by size preset).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--size` | string | `A2` | `A2`, `A3`, `24x36`, or `custom`. |
| `--bleed` | string | `3mm` | Bleed label in header. |

```bash
figma-kit make poster --size 24x36 --bleed "0.125in"
```

---

### `make brochure`

**Usage:** `figma-kit make brochure`

**Description:** Multi-panel brochure frames side by side.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--fold` | string | `trifold` | `trifold` or `bifold`. |
| `--format` | string | `letter` | `letter` or `A4`. |

```bash
figma-kit make brochure --fold bifold --format A4
```

---

### `make packaging`

**Usage:** `figma-kit make packaging`

**Description:** Flat die-line style layout (labeled rects).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--type` | string | `box` | `box` (flat layout). |
| `--w` | float | `200` | Front width. |
| `--h` | float | `280` | Front height. |
| `--d` | float | `100` | Depth. |

```bash
figma-kit make packaging --w 240 --h 300 --d 80
```

---

### `make signage`

**Usage:** `figma-kit make signage`

**Description:** Large board; dimensions in comment + labels (canvas uses fixed 1920×960).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--w` | string | `48in` | Width label. |
| `--h` | string | `24in` | Height label. |

```bash
figma-kit make signage --w "36in" --h "96in"
```

---

### `make menu`

**Usage:** `figma-kit make menu`

**Description:** Restaurant-style menu columns on a page frame.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--sections` | int | `2` | Column count. |
| `--format` | string | `letter` | `letter` or `A4`. |

```bash
figma-kit make menu --sections 3 --format letter
```

---

### `make changelog`

**Usage:** `figma-kit make changelog`

**Description:** Styled changelog / release notes page with version entries, date stamps, and type badges (added, changed, fixed, removed).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--entries` | string (JSON) | *(defaults)* | Array of `{version, date, changes: [{type, text}]}` |

```bash
figma-kit make changelog -t noir --entries '[{"version":"1.0.0","date":"2026-04-01","changes":[{"type":"added","text":"Direct MCP execution"},{"type":"added","text":"30 new design commands"}]}]'
```

---

## Layer 4 — Design system (`ds`)

### `ds create`

**Usage:** `figma-kit ds create`

**Description:** Create “Design System” page with swatches, type specimens, spacing bars.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Theme drives colors and type (`-t`). |

```bash
figma-kit ds create -t light -p 0
```

---

### `ds colors`

**Usage:** `figma-kit ds colors`

**Description:** Tints and shades from a primary hex.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--primary` | string | `#3B82F6` | Base hex. |

```bash
figma-kit ds colors --primary "#7C3AED"
```

---

### `ds type-scale`

**Usage:** `figma-kit ds type-scale`

**Description:** Type scale specimen frame from theme.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Uses `-t`. |

```bash
figma-kit ds type-scale -t default
```

---

### `ds spacing`

**Usage:** `figma-kit ds spacing`

**Description:** Theme spacing presets as labeled bars.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Uses `-t`. |

```bash
figma-kit ds spacing -t noir
```

---

### `ds elevation`

**Usage:** `figma-kit ds elevation`

**Description:** Cards showing theme shadow presets.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Uses `-t`. |

```bash
figma-kit ds elevation
```

---

### `ds radius`

**Usage:** `figma-kit ds radius`

**Description:** Corner radius reference chips.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Uses `-t` (preamble only). |

```bash
figma-kit ds radius
```

---

### `ds icons`

**Usage:** `figma-kit ds icons`

**Description:** 6×6 placeholder icon grid.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | — |

```bash
figma-kit ds icons
```

---

### `ds component`

**Usage:** `figma-kit ds component`

**Description:** Starter component set with two button variants.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | — |

```bash
figma-kit ds component
```

---

### `ds variables`

**Usage:** `figma-kit ds variables`

**Description:** JS that lists local variable collections via `figma.variables`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | — |

```bash
figma-kit ds variables
```

---

### `ds variables-create`

**Usage:** `figma-kit ds variables-create`

**Description:** Create a Figma variable collection with COLOR variables from all theme tokens.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | `Theme Tokens` | Variable collection name. |
| *(global)* | — | — | `--theme` / `-t` required. |

```bash
figma-kit ds variables-create -t noir --name "Noir Tokens"
```

---

### `ds search`

**Usage:** `figma-kit ds search`

**Description:** Prints MCP guidance for **`search_design_system`** (no JS).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit ds search
```

---

### `ds import`

**Usage:** `figma-kit ds import`

**Description:** Stub JS/comments for mapping tokens to variables.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | — |

```bash
figma-kit ds import
```

---

### `ds sync-tokens`

**Usage:** `figma-kit ds sync-tokens`

**Description:** Emit tokens from Go to stdout (**not** plugin JS): CSS, Tailwind snippet, or full theme JSON.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `css` | `css`, `tailwind`, or `json`. |

```bash
figma-kit ds sync-tokens --format json > theme.json
figma-kit ds sync-tokens --format tailwind
```

---

### `ds audit`

**Usage:** `figma-kit ds audit`

**Description:** JS walk of current page; flags solid fills not near theme palette.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | Uses `-t` for palette. |

```bash
figma-kit ds audit -t default
```

---

### `ds tokens` *(extra)*

**Usage:** `figma-kit ds tokens`

**Description:** Same as `ds sync-tokens --format json` (theme JSON to stdout).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | No flags. |

```bash
figma-kit ds tokens -t light > theme-light.json
```

---

### `ds component-sheet`

**Usage:** `figma-kit ds component-sheet`

**Description:** Generate a component inventory page showing all theme-aware components in a labeled grid. Useful for design system documentation.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--cols` | int | `4` | Grid columns |

```bash
figma-kit ds component-sheet -t noir --cols 3
```

---

## Layer 5 — Inspect and QA

### `inspect`

**Usage:** `figma-kit inspect <nodeId>`

**Description:** Return JSON-ish snapshot of node properties.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--deep` | bool | `false` | Include geometry, effects, auto-layout fields. |

```bash
figma-kit inspect "123:456" --deep
```

---

### `screenshot`

**Usage:** `figma-kit screenshot`

**Description:** Instructions for MCP **`get_screenshot`** (no JS).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit screenshot
```

---

### `tree`

**Usage:** `figma-kit tree [nodeId]`

**Description:** JS prints hierarchical tree text; omit id to use current page.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--max-depth` | int | `12` | Max depth. |

```bash
figma-kit tree --max-depth 6
figma-kit tree "123:456" --max-depth 4
```

---

### `find`

**Usage:** `figma-kit find <pattern>`

**Description:** Case-insensitive substring match on layer names in current page.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit find "Button"
```

---

### `measure`

**Usage:** `figma-kit measure <nodeIdA> <nodeIdB>`

**Description:** Axis-aligned gap between bounding boxes.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit measure "1:2" "1:3"
```

---

### `diff`

**Usage:** `figma-kit diff <nodeIdA> <nodeIdB>`

**Description:** Compare size and first solid fill between nodes.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit diff "1:2" "1:3"
```

---

### `qa contrast`

**Usage:** `figma-kit qa contrast`

**Description:** Heuristic text contrast vs assumed white page background.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| *(global)* | — | — | — |

```bash
figma-kit qa contrast
```

---

### `qa touch-targets`

**Usage:** `figma-kit qa touch-targets`

**Description:** Frames/components smaller than minimum size.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--min` | int | `44` | Min width/height (px). |

```bash
figma-kit qa touch-targets --min 48
```

---

### `qa orphans`

**Usage:** `figma-kit qa orphans`

**Description:** Top-level visible children of the page (excluding names starting with `.`).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa orphans
```

---

### `qa fonts`

**Usage:** `figma-kit qa fonts`

**Description:** Aggregate font family/style/size usage counts.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa fonts
```

---

### `qa colors`

**Usage:** `figma-kit qa colors`

**Description:** Histogram of solid fill RGB keys.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa colors
```

---

### `qa spacing`

**Usage:** `figma-kit qa spacing`

**Description:** Auto-layout frames with `itemSpacing` less than 8.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa spacing
```

---

### `qa naming`

**Usage:** `figma-kit qa naming`

**Description:** Empty or default layer names (e.g. `Frame 12`).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa naming
```

---

### `qa responsive`

**Usage:** `figma-kit qa responsive`

**Description:** Count constraints categories (`SCALE`, `STRETCH`, etc.).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa responsive
```

---

### `qa checklist`

**Usage:** `figma-kit qa checklist`

**Description:** Bundled pass: empty names, tiny targets, very tight spacing.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit qa checklist
```

---

## Layer 6 — Export and handoff

### `export png`

**Usage:** `figma-kit export png <nodeId>`

**Description:** JS: `exportAsync` PNG bytes (meta returned; bytes handled by host).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scale` | float | `2` | Scale factor. |

```bash
figma-kit export png "123:456" --scale 3
```

---

### `export svg`

**Usage:** `figma-kit export svg <nodeId>`

**Description:** JS: export SVG.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit export svg "123:456"
```

---

### `export pdf`

**Usage:** `figma-kit export pdf <nodeId>`

**Description:** JS: export PDF bytes.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit export pdf "123:456"
```

---

### `export page`

**Usage:** `figma-kit export page`

**Description:** JS: export each top-level frame/section/component as PNG @2x.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit export page -p 0
```

---

### `export sprites`

**Usage:** `figma-kit export sprites <frameId>`

**Description:** JS: export each direct child as PNG @1x.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit export sprites "123:456"
```

---

### `export tokens`

**Usage:** `figma-kit export tokens`

**Description:** Go stdout: theme as JSON or CSS variables (`--fk-*` color keys).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `json` | `json` or `css`. |

```bash
figma-kit export tokens --format css -t noir
```

---

### `handoff spec`

**Usage:** `figma-kit handoff spec <nodeId>`

**Description:** JS returns Markdown spec string for the node.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit handoff spec "123:456"
```

---

### `handoff redline`

**Usage:** `figma-kit handoff redline <nodeId>`

**Description:** JS creates measurement overlay above target.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit handoff redline "123:456"
```

---

### `handoff css`

**Usage:** `figma-kit handoff css <nodeId>`

**Description:** JS builds a simple CSS class from dimensions, radius, fills, strokes.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit handoff css "123:456"
```

---

### `handoff react`

**Usage:** `figma-kit handoff react <nodeId>`

**Description:** Prints instructions to call MCP **`get_design_context`** for that id (no JS).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | Positional `nodeId` only. |

```bash
figma-kit handoff react "123:456"
```

---

### `handoff assets`

**Usage:** `figma-kit handoff assets <nodeId>`

**Description:** JS walks subtree listing `exportSettings` per node.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | — |

```bash
figma-kit handoff assets "123:456"
```

---

## Layer 7 — Orchestration

### `batch`

**Usage:** `figma-kit batch <recipe.yaml>`

**Description:** Read YAML with `title` and `steps[]` (`title`, `js` per step); prints concatenated JS blocks with comments (local Go, not Figma JS).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| — | — | — | Positional YAML path. |

**YAML shape:**

```yaml
title: My recipe
steps:
  - title: Create frame
    js: |
      // use_figma block body...
  - title: Style
    js: |
      ...
```

```bash
figma-kit batch ./recipe.yaml > combined.js
```

---

### `compose`

**Usage:** `figma-kit compose [flags] "cmd1 args..." "cmd2 args..." ...`

**Description:** Batch N figma-kit commands into a single JavaScript payload for one `use_figma` call. The shared preamble (page setup, **theme-aware** font loading, theme colors, type scale) is emitted once; **only helper functions referenced by the merged steps** are included (**tree-shaking** via `detectNeededHelpers`). Each command body is scope-isolated in `{ }`. Compose emits **`const _results = [];`** and pushes each step’s primary node so later steps can use **`--parent _results[N]`** (or **`fx … --last`**) for chaining. See [STANDARDS.md](STANDARDS.md) compose contracts.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--recipe` | string | `""` | Path to a compose recipe YAML file |
| *(global)* | — | — | `--theme` / `-t` and `--page` / `-p` apply to all steps |

**Recipe YAML shape:**

```yaml
theme: noir
page: 0
steps:
  - "ui section --title 'Features' --label 'PRODUCT'"
  - "card glass --parent _results[0] --title 'Feature 1'"
  - "card glass --parent _results[0] --title 'Feature 2'"
  - "fx glow --last --position subtle"
  - "ui pricing --tiers '[{\"name\":\"Pro\",\"price\":\"$29\"}]'"
```

```bash
# Inline commands
figma-kit compose -t noir \
  "ui section --title 'Features' --label 'PRODUCT'" \
  "card glass --parent _results[0] --title 'Feature 1'" \
  "fx glow --last --position subtle"

# From a recipe YAML
figma-kit compose --recipe landing.yml

# Direct execution (generate + send to Figma in one shot)
figma-kit exec compose -t noir --recipe landing.yml
```

Only commands annotated as `composable` are allowed as steps (most Layer 1–2 commands that emit Plugin API JS). Composable **`ui`** and **`card`** commands accept **`--parent`** with a node id or **`_results[N]`** (emitted as a JS expression). Non-composable commands (e.g. `init`, `config`, `auth`) are rejected with an error.

---

### `image`

**Usage:** `figma-kit image <subcommand>`

**Description:** Place local images or URLs into Figma. Local files are base64-encoded and embedded directly in the generated JavaScript — no server, no public URL needed. Files up to ~33 KB work inline; for larger files, use a URL or `image serve`.

---

### `image place`

**Usage:** `figma-kit image place <path-or-url>`

**Description:** Create a new image frame from a local file or URL.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--width` | int | `400` | Frame width in pixels |
| `--height` | int | `300` | Frame height in pixels |
| `--name` | string | *(filename)* | Frame name |
| `--scale-mode` | string | `FILL` | `FILL`, `FIT`, `CROP`, or `TILE` |

```bash
figma-kit image place ./logo.png --name "Brand Logo" --width 200 --height 60
figma-kit image place https://example.com/hero.jpg --width 1440 --height 900
```

---

### `image fill`

**Usage:** `figma-kit image fill <path-or-url>`

**Description:** Replace an existing node's fill with an image.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--node` | string | *(required)* | Target node ID (e.g. `"2:5"`) |
| `--scale-mode` | string | `FILL` | `FILL`, `FIT`, `CROP`, or `TILE` |

```bash
figma-kit image fill ./hero.jpg --node "2:5"
figma-kit image fill https://example.com/bg.png --node "12:34" --scale-mode FIT
```

---

### `image serve`

**Usage:** `figma-kit image serve [directory]`

**Description:** Start a local HTTP server for image files too large for base64 embedding. Prints URLs you can use with `image place` or `card image`.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--port` | int | *(random)* | Port to listen on |

```bash
figma-kit image serve ./assets
# → Serving on http://127.0.0.1:8741
figma-kit image place http://127.0.0.1:8741/hero.jpg --width 1440 --height 900
```

---

## Theme Management

### `theme init`

**Usage:** `figma-kit theme init [flags]`

**Description:** Generate a complete theme JSON from hex colors. Derives a full palette (14 color tokens), typography, effects, spacing, and gradients from 3 seed colors. Use `--from` to extend an existing theme. With no flags, prints a starter template.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | `"My Theme"` | Theme name |
| `--desc` | string | *(auto)* | Theme description |
| `--bg` | string | `#0D0F17` | Background hex color |
| `--primary` | string | `#3366FF` | Primary accent hex color |
| `--accent` | string | `#14B8A6` | Secondary accent hex color |
| `--font-heading` | string | `Inter` | Heading font family |
| `--font-body` | string | `Inter` | Body font family |
| `--font-mono` | string | `Geist Mono` | Monospace font family |
| `--warn` | string | *(derived)* | Warning color hex |
| `--error` | string | *(derived)* | Error color hex |
| `--success` | string | *(derived)* | Success color hex |
| `--spacing` | string | *(standard)* | Spacing preset: `compact`, `spacious` |
| `--from` | string | — | Base theme file to extend (override specific flags) |
| `--output` / `-o` | string | *(stdout)* | Output file path |

```bash
# Basic: 3 colors
figma-kit theme init --name "Ocean" --bg "#0A1628" --primary "#2196F3" --accent "#00BCD4" -o themes/ocean.json

# Full: custom fonts, status colors, compact spacing
figma-kit theme init --name "Brand" --bg "#1a1a2e" --primary "#e94560" --accent "#0f3460" \
  --font-heading "Poppins" --font-body "DM Sans" --spacing compact -o brand.json

# Extend an existing theme
figma-kit theme init --from themes/brand.json --name "Brand Light" --bg "#F8F9FA" -o brand-light.json
```

---

### `theme preview`

**Usage:** `figma-kit theme preview [flags]`

**Description:** Outputs `use_figma` JS that creates a theme preview page in Figma. Includes: color swatches for all tokens, type scale using the theme's actual fonts, sample UI components (button, badge, status chips), gradient swatches from actual theme data, and brand info if present.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--theme` / `-t` | string | *(resolved)* | Theme to preview |
| `--page` / `-p` | int | *(resolved)* | Target page index |

```bash
figma-kit theme preview -t noir
figma-kit theme preview -t themes/brand.json
```

---

### `export tokens`

**Usage:** `figma-kit export tokens [flags]`

**Description:** Output theme tokens in JSON or CSS. Runs locally (no Figma plugin needed).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `json` | Output format: `json` or `css` |
| `--theme` / `-t` | string | *(resolved)* | Theme to export |

The `css` format emits CSS custom properties for colors, fonts, typography scale, and spacing:

```bash
figma-kit export tokens --format css -t noir
# → :root { --fk-BG: #0D0F17; --fk-font-heading: 'Inter'; --fk-h1-size: 72px; ... }
```

---

### `validate theme`

**Usage:** `figma-kit validate theme <path-or-name>`

**Description:** Parse and validate a theme. Reports color count, type scale, fonts, brand info. Warns (without failing) about missing conventional tokens (BG, WT, BL, CARD, STK), empty type/fonts/effects/spacing sections.

```bash
figma-kit validate theme themes/my-theme.json
figma-kit validate theme noir
```

---

## See also

- Root help: `figma-kit --help`
- Command help: `figma-kit <command> --help`
- Additional top-level commands (not in layers 0–7 above): `preamble`, `helpers`, `template`, `themes`, `scaffold`, `info`, `theme`, `cookbook`, `examples`, `docs`, `completion`.
