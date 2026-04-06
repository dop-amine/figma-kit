Here's a comprehensive command architecture for `figma-kit`, organized from lowest to highest abstraction. Every command generates `use_figma`-compatible JS code or calls the MCP tools directly. I've grounded each in what's actually possible through the 17 MCP tools and the Figma Plugin API.

---

## Complete `figma-kit` CLI Command Recommendations

### LAYER 0 -- File & Session Management

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit init <name>` | Create a new Figma file via `create_new_file` and store the fileKey in a local `.figmarc.json` project config | `figma-kit init "Client X - Brand Refresh"` |
| `figma-kit config` | Manage `.figmarc.json` -- set default fileKey, theme, page, export dir, Figma token | `figma-kit config set fileKey abc123` |
| `figma-kit whoami` | Show authenticated user, plans, teams (wraps `whoami` MCP tool) | `figma-kit whoami` |
| `figma-kit open` | Open the current file in Figma browser | `figma-kit open --page 2` |
| `figma-kit status` | Show file structure -- pages, top-level frames, node counts (wraps `get_metadata` on root) | `figma-kit status` |

---

### LAYER 1 -- Low-Level Node Operations

These generate atomic `use_figma` code blocks for direct node manipulation.

**Node CRUD:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit node create <type>` | Create any node: `frame`, `text`, `rect`, `ellipse`, `line`, `polygon`, `star`, `vector`, `group`, `component`, `component-set` | `figma-kit node create frame --name "Hero" --w 1440 --h 900 --page 0` |
| `figma-kit node clone <nodeId>` | Duplicate a node (with optional offset) | `figma-kit node clone 123:456 --dx 100 --dy 0` |
| `figma-kit node delete <nodeId>` | Remove a node | `figma-kit node delete 123:456` |
| `figma-kit node move <nodeId>` | Reposition a node | `figma-kit node move 123:456 --x 200 --y 300` |
| `figma-kit node resize <nodeId>` | Resize a node | `figma-kit node resize 123:456 --w 800 --h 600` |
| `figma-kit node rename <nodeId>` | Rename a node | `figma-kit node rename 123:456 --name "CTA Button"` |
| `figma-kit node reparent <nodeId> <parentId>` | Move a node to a different parent | `figma-kit node reparent 123:456 789:012` |
| `figma-kit node lock <nodeId>` | Lock/unlock a node | `figma-kit node lock 123:456 --unlock` |
| `figma-kit node visible <nodeId>` | Toggle visibility | `figma-kit node visible 123:456 --hide` |

**Styling:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit style fill <nodeId>` | Set fills -- solid, linear gradient, radial gradient, image | `figma-kit style fill 123:456 --solid "#3366FF" --opacity 0.8` |
| `figma-kit style stroke <nodeId>` | Set strokes -- color, weight, alignment, dashes | `figma-kit style stroke 123:456 --color "#FFF" --weight 2 --align inside` |
| `figma-kit style effect <nodeId>` | Apply effects -- drop shadow, inner shadow, layer blur, background blur | `figma-kit style effect 123:456 --shadow "0 4 12 rgba(0,0,0,0.15)"` |
| `figma-kit style corner <nodeId>` | Set corner radius (uniform or per-corner) | `figma-kit style corner 123:456 --radius 16` or `--tl 16 --tr 0 --br 16 --bl 0` |
| `figma-kit style blend <nodeId>` | Set blend mode and opacity | `figma-kit style blend 123:456 --mode OVERLAY --opacity 0.6` |
| `figma-kit style gradient <nodeId>` | Interactive gradient builder -- type, stops, angle/transform | `figma-kit style gradient 123:456 --type linear --angle 135 --stops "0:#3366FF,1:#0EB8A5"` |
| `figma-kit style clip <nodeId>` | Toggle clip content on frames | `figma-kit style clip 123:456 --on` |

**Text:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit text create` | Create a text node with full typography control | `figma-kit text create --content "Hello" --font Inter --weight Bold --size 48 --color "#FFF" --parent 123:456` |
| `figma-kit text edit <nodeId>` | Modify existing text content | `figma-kit text edit 123:456 --content "Updated copy"` |
| `figma-kit text style <nodeId>` | Change typography -- font, size, weight, line height, letter spacing, alignment, decoration | `figma-kit text style 123:456 --size 24 --lh 32 --ls 2 --align CENTER` |
| `figma-kit text range <nodeId>` | Apply mixed styles to text ranges (bold a word, color a span) | `figma-kit text range 123:456 --start 0 --end 5 --weight Bold --color "#3366FF"` |
| `figma-kit text list-fonts` | List all available fonts in the file (via Plugin API `figma.listAvailableFontsAsync()`) | `figma-kit text list-fonts` |
| `figma-kit text load-fonts` | Generate font loading code for specified families | `figma-kit text load-fonts --families "Inter,Geist Mono,SF Pro"` |

**Layout:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit layout auto <nodeId>` | Set auto layout -- direction, spacing, padding, alignment, wrapping | `figma-kit layout auto 123:456 --dir VERTICAL --gap 16 --pad 24 --align CENTER` |
| `figma-kit layout grid <nodeId>` | Add layout grids -- columns, rows, stretch/fixed | `figma-kit layout grid 123:456 --columns 12 --gutter 24 --margin 80` |
| `figma-kit layout constraints <nodeId>` | Set constraints (pin to edges, scale, stretch) | `figma-kit layout constraints 123:456 --h STRETCH --v TOP` |
| `figma-kit layout sizing <nodeId>` | Set sizing behavior -- fixed, hug, fill | `figma-kit layout sizing 123:456 --w FILL --h HUG` |
| `figma-kit layout align <nodeId>` | Align children within an auto-layout frame | `figma-kit layout align 123:456 --primary CENTER --counter BASELINE` |
| `figma-kit layout distribute <nodeIds>` | Distribute nodes evenly (horizontal or vertical) | `figma-kit layout distribute 1:1,1:2,1:3 --axis H --gap 24` |

---

### LAYER 2 -- Mid-Level Design Patterns

These compose multiple low-level operations into reusable design primitives.

**Cards & Containers:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit card glass` | Glassmorphism card with configurable preset (subtle/default/strong/pill) | `figma-kit card glass --preset strong --w 400 --h 300 --parent 123:456` |
| `figma-kit card solid` | Solid card with optional border, shadow, radius | `figma-kit card solid --bg "#1A1B2E" --border "#2A2B3E" --shadow md --radius 16` |
| `figma-kit card gradient` | Card with gradient fill | `figma-kit card gradient --from "#3366FF" --to "#0EB8A5" --angle 135` |
| `figma-kit card image` | Card with image fill, overlay, and content slot | `figma-kit card image --url "https://..." --overlay dark --title "Feature"` |
| `figma-kit card bento` | Auto-arranged bento grid of cards | `figma-kit card bento --cols 3 --rows 2 --gap 16` |

**UI Primitives:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit ui button` | Button with variants (primary, secondary, ghost, destructive, outline) and sizes (sm, md, lg) | `figma-kit ui button --variant primary --label "Get Started" --size lg` |
| `figma-kit ui input` | Text input field with label, placeholder, states | `figma-kit ui input --label "Email" --placeholder "you@company.com"` |
| `figma-kit ui badge` | Badge/chip/tag with color and text | `figma-kit ui badge --text "NEW" --color blue` |
| `figma-kit ui avatar` | Avatar circle with initials fallback | `figma-kit ui avatar --initials "AK" --size 40` |
| `figma-kit ui divider` | Horizontal/vertical divider line | `figma-kit ui divider --dir H --w 400 --color muted` |
| `figma-kit ui icon` | Placeholder icon shape (circle, square, or from a glyph set) | `figma-kit ui icon --shape circle --size 24 --color blue` |
| `figma-kit ui progress` | Progress bar with percentage | `figma-kit ui progress --value 75 --w 300` |
| `figma-kit ui toggle` | Toggle switch | `figma-kit ui toggle --state on --size md` |
| `figma-kit ui tooltip` | Tooltip callout shape with text | `figma-kit ui tooltip --text "Hover info" --position top` |
| `figma-kit ui stat` | Stat/metric display (big number + label) | `figma-kit ui stat --value "4.2x" --label "Faster decisions" --trend up` |
| `figma-kit ui table` | Data table from JSON/CSV | `figma-kit ui table --data ./data.json --cols "Name,Role,Status"` |
| `figma-kit ui nav` | Navigation bar (horizontal or vertical) | `figma-kit ui nav --items "Home,Products,Pricing,Contact" --style topbar` |
| `figma-kit ui footer` | Footer section with columns, links, copyright | `figma-kit ui footer --cols 4 --copyright "2026 Acme Inc."` |

**Visual Effects:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit fx glow` | Add radial gradient glow to a frame background | `figma-kit fx glow --position topRight --intensity 0.14 --color blue` |
| `figma-kit fx mesh` | Multi-point mesh gradient (layered radials simulating mesh) | `figma-kit fx mesh --points 3 --palette "blue,teal,purple"` |
| `figma-kit fx noise` | Subtle noise texture overlay (via gradient dithering pattern) | `figma-kit fx noise --opacity 0.03` |
| `figma-kit fx vignette` | Edge vignette darkening | `figma-kit fx vignette --strength 0.4` |
| `figma-kit fx grain` | Film grain effect | `figma-kit fx grain --amount light` |
| `figma-kit fx blur-bg` | Background blur layer (frosted glass backdrop) | `figma-kit fx blur-bg --radius 40 --tint "rgba(0,0,0,0.3)"` |
| `figma-kit fx accent-bar` | Gradient accent bar/line | `figma-kit fx accent-bar --from "#3366FF" --to "#0EB8A5" --w 200 --h 4` |
| `figma-kit fx shadow` | Shadow preset (sm, md, lg, xl, glow, inner) | `figma-kit fx shadow 123:456 --preset lg` |
| `figma-kit fx parallax-layer` | Create layered depth composition (background, midground, foreground with opacity/blur stacking) | `figma-kit fx parallax-layer --layers 3` |

---

### LAYER 3 -- High-Level Deliverables & Templates

Complete, production-ready deliverable generators.

**Marketing & Social:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit make carousel` | LinkedIn carousel (1080x1350, 4-10 slides) from YAML/JSON content spec | `figma-kit make carousel --content ./carousel.yml --theme arkham --slides 7` |
| `figma-kit make instagram-post` | Instagram post (1080x1080) | `figma-kit make instagram-post --type quote --content "..."` |
| `figma-kit make instagram-story` | Instagram story (1080x1920) | `figma-kit make instagram-story --content ./story.yml` |
| `figma-kit make twitter-card` | Twitter/X card (1200x675) | `figma-kit make twitter-card --headline "..." --image hero` |
| `figma-kit make facebook-cover` | Facebook cover photo (820x312) | `figma-kit make facebook-cover --theme dark` |
| `figma-kit make youtube-thumb` | YouTube thumbnail (1280x720) | `figma-kit make youtube-thumb --title "..." --face left` |
| `figma-kit make og-image` | Open Graph image for link previews (1200x630) | `figma-kit make og-image --title "..." --description "..."` |
| `figma-kit make banner` | IAB standard banner ads (300x250, 728x90, 160x600, 320x50, etc.) | `figma-kit make banner --sizes "leaderboard,mrec" --content ./ad.yml` |
| `figma-kit make email-header` | Email header/hero (600px wide) | `figma-kit make email-header --w 600` |
| `figma-kit make ad-set` | Generate all social sizes at once from a single content spec | `figma-kit make ad-set --content ./campaign.yml --platforms "linkedin,instagram,twitter,facebook"` |

**Sales & Business:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit make one-pager` | B2B sales one-pager (letter/A4, light/dark) | `figma-kit make one-pager --format letter --mode print --content ./one-pager.yml` |
| `figma-kit make pitch-deck` | Investor/sales pitch deck (16:9 slides) | `figma-kit make pitch-deck --slides 12 --template saas` |
| `figma-kit make case-study` | Portfolio case study layout (multi-section scrolling page) | `figma-kit make case-study --sections "overview,challenge,solution,results"` |
| `figma-kit make proposal` | Client project proposal | `figma-kit make proposal --client "Acme" --scope ./scope.yml` |
| `figma-kit make invoice` | Invoice template (print-ready) | `figma-kit make invoice --template modern` |
| `figma-kit make business-card` | Business card (3.5"x2", front and back) | `figma-kit make business-card --name "..." --title "..."` |
| `figma-kit make letterhead` | Branded letterhead (A4/letter) | `figma-kit make letterhead --format A4` |
| `figma-kit make contract` | Contract/NDA cover page with branding | `figma-kit make contract --title "..."` |

**Motion & Storyboard:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit make storyboard` | Full storyboard with styleframes + panoramic strip | `figma-kit make storyboard --scenes 5 --duration 12s --content ./storyboard.yml` |
| `figma-kit make styleframe` | Single high-fidelity styleframe (1920x1080) | `figma-kit make styleframe --mood warm --scene "The problem"` |
| `figma-kit make animatic` | Sequential frame series with timing annotations | `figma-kit make animatic --fps 24 --duration 30s` |
| `figma-kit make transition-spec` | Motion design specification (easing, duration, properties) | `figma-kit make transition-spec --type page-enter --easing ease-out` |

**UI/UX Design:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit make wireframe` | Low-fidelity wireframe for a screen type (landing, dashboard, form, list, detail, auth) | `figma-kit make wireframe --type dashboard --breakpoint desktop` |
| `figma-kit make screen` | High-fidelity screen from wireframe + theme | `figma-kit make screen --type landing --theme arkham --sections "hero,features,pricing,cta"` |
| `figma-kit make dashboard` | Data dashboard layout (charts, stats, tables, sidebar nav) | `figma-kit make dashboard --widgets "stat,chart,table,list" --cols 3` |
| `figma-kit make form` | Form layout with field types from schema | `figma-kit make form --fields ./form-schema.json` |
| `figma-kit make modal` | Modal/dialog in various sizes | `figma-kit make modal --size md --type confirmation` |
| `figma-kit make empty-state` | Empty state illustration placeholder + copy | `figma-kit make empty-state --message "No results found"` |
| `figma-kit make error-page` | Error page (404, 500, offline) | `figma-kit make error-page --type 404` |
| `figma-kit make onboarding` | Multi-step onboarding flow | `figma-kit make onboarding --steps 4` |
| `figma-kit make settings` | Settings/preferences page layout | `figma-kit make settings --sections "profile,notifications,billing,security"` |

**Print & Physical:**

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit make poster` | Poster layout (A3, A2, 24x36, custom) | `figma-kit make poster --size A2 --bleed 3mm` |
| `figma-kit make brochure` | Tri-fold or bi-fold brochure | `figma-kit make brochure --fold trifold --format letter` |
| `figma-kit make packaging` | Product packaging dieline with design zones | `figma-kit make packaging --type box --w 100 --h 150 --d 50` |
| `figma-kit make signage` | Signage/wayfinding layout at real-world dimensions | `figma-kit make signage --w 48in --h 24in` |
| `figma-kit make menu` | Restaurant/service menu layout | `figma-kit make menu --sections 4 --format letter` |

---

### LAYER 4 -- Design System Management

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit ds create` | Initialize a design system page with color, type, spacing, and effect tokens | `figma-kit ds create --theme arkham` |
| `figma-kit ds colors` | Generate color palette page -- primary, secondary, neutrals, semantic, with tints/shades | `figma-kit ds colors --primary "#3366FF" --secondary "#0EB8A5" --tints 9` |
| `figma-kit ds type-scale` | Generate type scale specimen (all heading/body/caption sizes) | `figma-kit ds type-scale --base 16 --ratio 1.25 --font Inter` |
| `figma-kit ds spacing` | Generate spacing scale visualization (4, 8, 12, 16, 24, 32, 48, 64, 96) | `figma-kit ds spacing --base 4 --steps 10` |
| `figma-kit ds elevation` | Generate shadow/elevation specimen | `figma-kit ds elevation --levels 5` |
| `figma-kit ds radius` | Generate border radius specimen | `figma-kit ds radius --values "0,4,8,12,16,24,full"` |
| `figma-kit ds icons` | Generate icon grid placeholder page | `figma-kit ds icons --size 24 --grid 8x6` |
| `figma-kit ds component <name>` | Create a componentized Figma component with variants and properties | `figma-kit ds component button --variants "primary,secondary,ghost" --sizes "sm,md,lg" --states "default,hover,active,disabled"` |
| `figma-kit ds variables` | Read all design variables from a file (wraps `get_variable_defs`) | `figma-kit ds variables --node 0:1` |
| `figma-kit ds search` | Search design system library for components/variables/styles (wraps `search_design_system`) | `figma-kit ds search "button" --include components,styles` |
| `figma-kit ds import <componentKey>` | Generate `importComponentByKeyAsync` code for a library component | `figma-kit ds import abc123def` |
| `figma-kit ds sync-tokens` | Export theme JSON as CSS custom properties, Tailwind config, or Swift/Kotlin tokens | `figma-kit ds sync-tokens --format tailwind --output ./tokens.js` |
| `figma-kit ds audit` | Scan a page for style inconsistencies (off-palette colors, non-standard fonts, irregular spacing) | `figma-kit ds audit --page 0` |

---

### LAYER 5 -- Read, Inspect & QA

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit inspect <nodeId>` | Dump full node properties -- fills, strokes, effects, layout, constraints, text styles | `figma-kit inspect 123:456` |
| `figma-kit screenshot <nodeId>` | Take a screenshot of a node (wraps `get_screenshot`) | `figma-kit screenshot 123:456 --save ./screenshots/` |
| `figma-kit screenshot all` | Screenshot every top-level frame on a page | `figma-kit screenshot all --page 0 --save ./screenshots/` |
| `figma-kit tree` | Print the full node tree for a page (wraps `get_metadata`, renders as indented tree) | `figma-kit tree --page 0 --depth 3` |
| `figma-kit find` | Search nodes by name, type, property, or regex | `figma-kit find --name "Button*" --type COMPONENT` |
| `figma-kit measure <nodeA> <nodeB>` | Measure distance between two nodes | `figma-kit measure 1:1 1:2` |
| `figma-kit diff <nodeA> <nodeB>` | Compare two nodes -- diff their properties | `figma-kit diff 1:1 1:2` |
| `figma-kit qa contrast` | Check all text nodes against WCAG 2.1 contrast ratios (AA/AAA) | `figma-kit qa contrast --page 0 --level AA` |
| `figma-kit qa touch-targets` | Flag interactive elements smaller than 44x44 px | `figma-kit qa touch-targets --page 0` |
| `figma-kit qa orphans` | Find detached instances, hidden nodes, empty frames, unnamed layers | `figma-kit qa orphans --page 0` |
| `figma-kit qa fonts` | List all fonts used, flag missing/unlicensed fonts | `figma-kit qa fonts` |
| `figma-kit qa colors` | Extract all unique colors used, flag off-palette ones | `figma-kit qa colors --page 0 --palette ./theme.json` |
| `figma-kit qa spacing` | Check spacing consistency against a baseline grid | `figma-kit qa spacing --page 0 --base 8` |
| `figma-kit qa naming` | Check layer naming against conventions | `figma-kit qa naming --page 0 --convention kebab` |
| `figma-kit qa responsive` | Check if frames use auto-layout and fill/hug vs. fixed sizing | `figma-kit qa responsive --page 0` |
| `figma-kit qa checklist` | Run all QA checks and output a scored report | `figma-kit qa checklist --page 0 --output ./qa-report.md` |

---

### LAYER 6 -- Export & Handoff

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit export png <nodeId>` | Export node as PNG at 1x, 2x, 3x | `figma-kit export png 123:456 --scale "1,2,3" --output ./assets/` |
| `figma-kit export svg <nodeId>` | Export as SVG with optional optimization | `figma-kit export svg 123:456 --optimize` |
| `figma-kit export pdf <nodeId>` | Export as PDF (for print deliverables) | `figma-kit export pdf 123:456 --output ./print/` |
| `figma-kit export page <index>` | Export all frames on a page as individual files | `figma-kit export page 0 --format png --scale 2` |
| `figma-kit export sprites` | Export all icons/small assets as a sprite sheet | `figma-kit export sprites --page 0 --format svg` |
| `figma-kit export tokens` | Export design tokens as JSON, CSS vars, SCSS vars, Tailwind config, or Swift/Kotlin | `figma-kit export tokens --format css --output ./tokens.css` |
| `figma-kit handoff spec <nodeId>` | Generate developer handoff spec -- dimensions, spacing, colors, fonts, CSS snippets | `figma-kit handoff spec 123:456 --format markdown --output ./specs/` |
| `figma-kit handoff redline <nodeId>` | Generate annotated redline frame overlaying the design with measurements | `figma-kit handoff redline 123:456` |
| `figma-kit handoff css <nodeId>` | Generate CSS/Tailwind code from a node's visual properties | `figma-kit handoff css 123:456 --framework tailwind` |
| `figma-kit handoff react <nodeId>` | Generate React component from design context (wraps `get_design_context`) | `figma-kit handoff react 123:456 --framework next` |
| `figma-kit handoff assets` | Package all exportable assets with manifest | `figma-kit handoff assets --page 0 --output ./handoff/` |

---

### LAYER 7 -- Orchestration & Batch Operations

| Command | Description | Example |
|---------|-------------|---------|
| `figma-kit batch <recipe.yml>` | Execute a YAML/
