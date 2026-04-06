# figma-kit command reference

`figma-kit` is a CLI that mostly prints **JavaScript** meant to run inside Figma via the **`use_figma`** / MCP workflow (theme preamble, helpers, and plugin-style API calls). A few commands only print **instructions**, **token/CSS/JSON** to stdout, or perform **local** actions (`init`, `open`).

Unless noted, **pipe or paste the output** into your Figma MCP execution path. Commands that resolve a theme use **`--theme` / `-t`** and page index **`--page` / `-p`** when relevant.

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

```bash
figma-kit text create --content "Hello" --font Inter --weight Bold --size 24 --color "#0F172A"
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
| `--text` | string | `New` | Badge copy. |
| `--color` | string | `blue` | `blue`, `green`, `red`, `yellow`, `gray`. |

```bash
figma-kit ui badge --text "Beta" --color yellow
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
| `--value` | string | `4.2x` | Main value. |
| `--label` | string | `Faster` | Caption. |
| `--trend` | string | `""` | `up`, `down`, or `neutral`. |

```bash
figma-kit ui stat --value "12%" --label "MoM" --trend up
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

### `fx glow`

**Usage:** `figma-kit fx glow <nodeId>`

**Description:** Layered radial gradient glows; optional hex tint overrides inner color.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--position` | string | `subtle` | `topRight`, `center`, `subtle`, `cta`. |
| `--intensity` | float | `1` | Multiplier for glow alphas. |
| `--color` | string | `""` | Optional `#RRGGBB` tint. |

```bash
figma-kit fx glow "123:456" --position cta --intensity 1.2 --color "#6366F1"
```

---

### `fx mesh`

**Usage:** `figma-kit fx mesh <nodeId>`

**Description:** Mesh-like stacked radial fills.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--points` | int | `5` | Count (min 2). |
| `--palette` | string | `#2563eb,#14b8a6,#8b5cf6` | Comma hex or names (`blue`, `teal`, …). |

```bash
figma-kit fx mesh "123:456" --points 7 --palette "purple,pink,orange"
```

---

### `fx noise`

**Usage:** `figma-kit fx noise <nodeId>`

**Description:** Overlay rectangle with dithered gradient (child of target).

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--opacity` | float | `0.35` | Strength 0–1. |

```bash
figma-kit fx noise "123:456" --opacity 0.25
```

---

### `fx vignette`

**Usage:** `figma-kit fx vignette <nodeId>`

**Description:** Radial darkening overlay child.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--strength` | float | `0.6` | Vignette strength. |

```bash
figma-kit fx vignette "123:456" --strength 0.8
```

---

### `fx grain`

**Usage:** `figma-kit fx grain <nodeId>`

**Description:** Film grain overlay.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--amount` | string | `medium` | `light`, `medium`, or `heavy`. |

```bash
figma-kit fx grain "123:456" --amount light
```

---

### `fx blur-bg`

**Usage:** `figma-kit fx blur-bg <nodeId>`

**Description:** Frosted rectangle + background blur child.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--radius` | int | `24` | Blur radius. |
| `--tint` | string | `""` | Optional `rgba(r,g,b,a)` (0–1 or 0–255 channels). |

```bash
figma-kit fx blur-bg "123:456" --radius 32 --tint "rgba(255,255,255,0.25)"
```

---

### `fx accent-bar`

**Usage:** `figma-kit fx accent-bar <parentId>`

**Description:** Gradient bar child inside parent.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
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

**Usage:** `figma-kit fx shadow <nodeId>`

**Description:** Replace `effects` with one shadow preset.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--preset` | string | `md` | `sm`, `md`, `lg`, `xl`, `glow`, `inner`. |

```bash
figma-kit fx shadow "123:456" --preset xl
```

---

### `fx parallax-layer`

**Usage:** `figma-kit fx parallax-layer <parentId>`

**Description:** Stack of inset frames inside parent.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--layers` | int | `3` | Count (min 2). |

```bash
figma-kit fx parallax-layer "123:456" --layers 5
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
- Additional top-level commands (not in layers 0–7 above): `preamble`, `helpers`, `template`, `themes`, `scaffold`, `info`, `theme`, `completion`.
