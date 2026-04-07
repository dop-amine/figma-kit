# Prompt Cookbook

Real-world prompt sessions for figma-kit. Each example shows what you say, what the AI does, and the figma-kit commands it runs. Copy these prompts directly into Cursor, Claude Code, or any MCP-compatible AI agent.

> **Installed via Brew?** Run `figma-kit cookbook` to browse these examples from your terminal.

---

## From a Reference

### Match a Website

**Prompt:**

```
Create a figma-kit theme that matches stripe.com — dark navy background, 
purple primary, teal accent. Then build a landing page with a hero section, 
feature grid, and pricing table.
```

**What the AI does:**

```bash
# 1. Create a theme from extracted colors
figma-kit theme init \
  --name "Stripe" \
  --bg "#0A2540" \
  --primary "#635BFF" \
  --accent "#00D4AA" \
  --font-heading "Sohne" \
  --font-body "Sohne" \
  --spacing spacious \
  -o themes/stripe.json

# 2. Preview the theme in Figma
figma-kit theme preview -t stripe
# → AI verifies with get_screenshot

# 3. Set up the page with theme colors and fonts
figma-kit preamble -t stripe
# → AI executes via use_figma

# 4. Create the hero section (one command)
figma-kit ui hero -t stripe \
  --title "Payments infrastructure for the internet" \
  --subtitle "Millions of companies use Stripe to accept payments" \
  --cta "Start now" \
  --badge "New"

# 5. Build the feature grid (one command)
figma-kit ui feature-grid -t stripe --cols 3

# 6. Add individual feature cards with glass effect
figma-kit card glass -t stripe --title "Payments" --desc "Accept and optimise payments globally"
figma-kit card glass -t stripe --title "Billing" --desc "Build and scale your recurring business model"
figma-kit card glass -t stripe --title "Connect" --desc "Set up multi-party payments and payouts"

# 7. Add pricing table (one command)
figma-kit ui pricing -t stripe --tiers '[
  {"name":"Starter","price":"Free","features":["Core APIs","Test mode"]},
  {"name":"Pro","price":"$29","highlighted":true,"features":["All APIs","Priority support"]},
  {"name":"Enterprise","price":"Custom","features":["Everything","Dedicated CSM"]}
]'

# 8. QA the full page
figma-kit qa checklist --page 0
```

**Result:** A complete Stripe-inspired landing page with matching colors, fonts, hero section, three glass feature cards, and a pricing table — all in your Figma file.

---

### Match a Screenshot

**Prompt:**

```
Here's a screenshot of a dark dashboard UI [paste image]. 
Create a figma-kit theme from the colors you see, then build 
a dashboard with stat cards, a chart widget, and a data table.
```

**What the AI does:**

The AI analyzes the image, identifies the dominant colors (dark background, blue accents, white text), and runs:

```bash
# 1. Create theme from extracted colors
figma-kit theme init \
  --name "Dashboard Dark" \
  --bg "#0F1117" \
  --primary "#3B82F6" \
  --accent "#10B981" \
  --font-heading "Inter" \
  --font-body "Inter" \
  --font-mono "JetBrains Mono" \
  --spacing compact \
  -o themes/dashboard.json

# 2. Preview and set up
figma-kit theme preview -t dashboard
figma-kit preamble -t dashboard

# 3. Build the dashboard
figma-kit make dashboard --widgets "stat,stat,stat,chart,table,list" --cols 3 -t dashboard

# 4. Add individual stat cards with real data
figma-kit ui stat -t dashboard --parent <widget-id>
```

---

### Match a Brand Guide

**Prompt:**

```
Here's our brand guide [paste PDF/screenshot]. It uses Montserrat for 
headings, DM Sans for body, #1E3A5F as primary, #FF6B35 as accent. 
Create a theme, then generate a full design system page in Figma.
```

**What the AI does:**

```bash
# 1. Create theme with exact brand fonts and colors
figma-kit theme init \
  --name "Acme Brand" \
  --bg "#0C1824" \
  --primary "#1E3A5F" \
  --accent "#FF6B35" \
  --font-heading "Montserrat" \
  --font-body "DM Sans" \
  --font-mono "Fira Code" \
  --warn "#FFB020" \
  --error "#E53E3E" \
  --success "#38A169" \
  -o themes/acme.json

# 2. Generate the design system
figma-kit ds create -t acme
# → Creates: color swatches, type scale specimens, spacing scale,
#    elevation demos, radius tokens, icon grid, component specimens

# 3. Push tokens to Figma variables
figma-kit ds variables-create -t acme

# 4. Export tokens for the dev team
figma-kit export tokens -t acme --format css > acme-tokens.css
```

**Result:** A complete design system page in Figma with all tokens from the brand guide, plus CSS variables ready for developers.

---

## Landing Pages

### SaaS Product Page

**Prompt:**

```
Build a landing page for a developer tool called "Flux" — a real-time 
database. Dark theme with purple accents. I want:
- Hero with headline and terminal mockup
- 3 feature cards (speed, scale, sync)  
- Code example section
- Pricing grid (free, pro, enterprise)
- Footer with links
```

**What the AI does:**

```bash
# 1. Create a custom theme
figma-kit theme init \
  --name "Flux" \
  --bg "#09090B" \
  --primary "#7C3AED" \
  --accent "#06B6D4" \
  --font-heading "Inter" \
  --font-body "Inter" \
  --font-mono "Geist Mono" \
  -o themes/flux.json

# 2. Set up page
figma-kit preamble -t flux

# 3. Hero section
figma-kit node create frame --name "Hero" -w 1440 --height 800
figma-kit fx mesh <hero-id> -t flux
figma-kit text create --content "Real-time data.\nZero latency." \
  --font "Inter" --weight "Extra Bold" --size 80 --color "#FFFFFF" --parent <hero-id>
figma-kit text create --content "Flux syncs your data across every edge in under 10ms." \
  --font "Inter" --size 20 --color "#A1A1AA" --parent <hero-id>
figma-kit ui button --variant primary -t flux --parent <hero-id>

# 4. Feature cards
figma-kit card glass -t flux --title "Blazing Fast" --desc "Sub-millisecond reads from the edge. Your data, everywhere." --preset strong
figma-kit card glass -t flux --title "Infinite Scale" --desc "From prototype to planet-scale without changing a line of code." --preset strong
figma-kit card glass -t flux --title "Real-time Sync" --desc "Every write propagates globally in under 10ms. Always consistent." --preset strong

# 5. Pricing section
figma-kit make screen --type pricing -t flux

# 6. Footer
figma-kit ui footer -t flux

# 7. QA
figma-kit qa checklist --page 0
```

---

### Photographer Portfolio

**Prompt:**

```
Create a minimal portfolio site for a photographer. Light theme, lots of 
whitespace. Use Playfair Display for headings and Lato for body text. 
I want 4 project cards and a simple contact section at the bottom.
```

**What the AI does:**

```bash
# 1. Create a light, editorial theme
figma-kit theme init \
  --name "Portfolio" \
  --bg "#FAFAFA" \
  --primary "#1A1A1A" \
  --accent "#C9A87C" \
  --font-heading "Playfair Display" \
  --font-body "Lato" \
  --spacing spacious \
  -o themes/portfolio.json

# 2. Set up page
figma-kit preamble -t portfolio

# 3. Hero — minimal, just a name and tagline
figma-kit node create frame --name "Hero" -w 1440 --height 600
figma-kit text create --content "Elena Voss" \
  --font "Playfair Display" --weight Bold --size 64 --color "#1A1A1A" --parent <hero-id>
figma-kit text create --content "Documentary & Portrait Photography" \
  --font "Lato" --size 18 --color "#71717A" --parent <hero-id>

# 4. Project cards — image cards with overlay
figma-kit card image -t portfolio --title "Patagonia" --desc "Landscapes of southern Chile" -w 640 --height 480
figma-kit card image -t portfolio --title "Berlin Nights" --desc "Street photography after dark" -w 640 --height 480
figma-kit card image -t portfolio --title "Faces of Tokyo" --desc "Portraits from Shinjuku to Shibuya" -w 640 --height 480
figma-kit card image -t portfolio --title "Desert Light" --desc "Morning in the Sahara" -w 640 --height 480

# 5. Contact section
figma-kit node create frame --name "Contact" -w 1440 --height 400
figma-kit text create --content "Get in Touch" \
  --font "Playfair Display" --weight Bold --size 48 --color "#1A1A1A" --parent <contact-id>
figma-kit text create --content "hello@elenavoss.com" \
  --font "Lato" --size 20 --color "#C9A87C" --parent <contact-id>
```

---

### Agency Homepage

**Prompt:**

```
Build a creative agency landing page. Bold, maximalist typography. 
Gradient cards. Use Space Grotesk for headings. Dark background with 
bright green accents. Include a case study grid with 4 projects.
```

**What the AI does:**

```bash
figma-kit theme init \
  --name "Agency" \
  --bg "#0A0A0A" \
  --primary "#BEFF00" \
  --accent "#FF3366" \
  --font-heading "Space Grotesk" \
  --font-body "Inter" \
  -o themes/agency.json

figma-kit preamble -t agency

# Hero with bold type
figma-kit node create frame --name "Hero" -w 1440 --height 900
figma-kit text create --content "We build\ndigital\nexperiences." \
  --font "Space Grotesk" --weight Bold --size 120 --color "#BEFF00" --parent <hero-id>

# Case study cards
figma-kit card gradient -t agency --title "Nike React" --desc "E-commerce redesign" -w 640 --height 400
figma-kit card gradient -t agency --title "Spotify Wrapped" --desc "Annual campaign" -w 640 --height 400
figma-kit card gradient -t agency --title "Airbnb Luxe" --desc "Premium tier launch" -w 640 --height 400
figma-kit card gradient -t agency --title "Tesla Configurator" --desc "3D product builder" -w 640 --height 400
```

---

## Marketing Materials

### Pitch Deck

**Prompt:**

```
Create a 7-slide pitch deck for a climate-tech startup called "Canopy". 
Dark green theme, clean typography. Slides: cover, problem, solution, 
market size, traction, team, closing CTA with "Book a demo".
```

**What the AI does:**

First, create a content YAML file (`canopy-deck.yml`):

```yaml
theme: canopy
total: 7
slides:
  - name: Cover
    glow: topRight
    headline: "Canopy"
    subtitle: "Carbon intelligence for the built environment"
  - name: Problem
    glow: subtle
    headline: "Buildings account for\n40% of global emissions"
    chips: ["Energy waste", "No visibility", "Compliance risk"]
  - name: Solution
    headline: "AI-powered carbon\nmonitoring"
    subtitle: "Real-time emissions tracking across your entire portfolio"
  - name: Market
    headline: "$340B\nCarbon management"
    chips: ["Growing 24% YoY", "Regulatory tailwind", "ESG mandates"]
  - name: Traction
    headline: "12 enterprise\ncustomers"
    chips: ["$2.4M ARR", "3x YoY growth", "NPS 72"]
  - name: Team
    headline: "Built by climate\nand ML experts"
    subtitle: "Ex-Google, ex-McKinsey, 2 PhD climate scientists"
  - name: CTA
    glow: cta
    headline: "Let's build a\nsustainable future."
    cta:
      text: "Book a demo"
      centered: true
```

Then the commands:

```bash
figma-kit theme init \
  --name "Canopy" \
  --bg "#0A1F14" \
  --primary "#22C55E" \
  --accent "#06B6D4" \
  --font-heading "Inter" \
  -o themes/canopy.json

figma-kit make carousel --content canopy-deck.yml -t canopy
```

**Result:** A 7-slide pitch deck with climate-green theme, glow effects, and proper slide hierarchy — all from one YAML file and two commands.

---

### Social Media Campaign

**Prompt:**

```
Generate an Instagram carousel (5 slides), an OG image, and a Twitter 
card for our product launch. Brand colors: #1a1a2e background, #e94560 
primary. The product is called "Nimbus" — a cloud development platform.
```

**What the AI does:**

```bash
# 1. Create theme
figma-kit theme init \
  --name "Nimbus" \
  --bg "#1A1A2E" \
  --primary "#E94560" \
  --accent "#0F3460" \
  -o themes/nimbus.json

# 2. Instagram carousel
# (create slides.yml with 5 slides about the product)
figma-kit make carousel --content nimbus-slides.yml --slides 5 -t nimbus

# 3. OG image
figma-kit make og-image --title "Nimbus" --description "Cloud development, reimagined" -t nimbus

# 4. Twitter card
figma-kit make twitter-card -t nimbus

# 5. Instagram post (standalone)
figma-kit make instagram-post --type quote --content "Ship to the cloud in seconds, not hours." -t nimbus
```

**Result:** A complete social media kit with consistent branding across Instagram, Twitter, and Open Graph — ready for export.

---

### Email Header

**Prompt:**

```
Design an email header for a Black Friday sale. Dark background, gold 
accents. Title: "Black Friday — 50% Off Everything". Make it 600px wide.
```

**What the AI does:**

```bash
figma-kit theme init \
  --name "BFCM" \
  --bg "#0A0A0A" \
  --primary "#FFD700" \
  --accent "#FF4500" \
  -o themes/bfcm.json

figma-kit make email-header -t bfcm

# Then customize the text
figma-kit text create --content "BLACK FRIDAY" \
  --font "Inter" --weight "Extra Bold" --size 48 --color "#FFD700" --parent <header-id>
figma-kit text create --content "50% Off Everything" \
  --font "Inter" --weight Bold --size 28 --color "#FFFFFF" --parent <header-id>
figma-kit ui button --variant primary -t bfcm --parent <header-id>
```

---

## Design Systems

### From Scratch

**Prompt:**

```
Create a complete design system for our app. Primary color: #5E6AD2, 
accent: #26B5CE. I want the full setup: color swatches, type scale, 
spacing, elevation, components. Use Inter for everything. Push tokens 
to Figma variables when done.
```

**What the AI does:**

```bash
# 1. Create theme
figma-kit theme init \
  --name "AppDS" \
  --bg "#1B1B25" \
  --primary "#5E6AD2" \
  --accent "#26B5CE" \
  --font-heading "Inter" \
  --font-body "Inter" \
  --font-mono "JetBrains Mono" \
  -o themes/app-ds.json

# 2. Generate the full design system page
figma-kit ds create -t app-ds
# → Creates: color palette, type scale, spacing, elevation, radius, icons, components

# 3. Push to Figma variables
figma-kit ds variables-create -t app-ds

# 4. Export CSS tokens
figma-kit export tokens -t app-ds --format css

# 5. Audit the system
figma-kit ds audit -t app-ds
```

**Result:** A design system page in Figma with 14 color tokens, 8 type scale entries, spacing bars, elevation samples, and every component specimen — plus Figma variables and a CSS export.

---

### From an Existing Website

**Prompt:**

```
Look at linear.app and create a figma-kit design system that matches their 
aesthetic. I want the theme to feel like Linear — compact spacing, Inter 
font, that signature purple-blue.
```

**What the AI does:**

The AI visits the site (or analyzes a screenshot), extracts the visual language, and runs:

```bash
# 1. Create Linear-inspired theme
figma-kit theme init \
  --name "Linear Dark" \
  --bg "#1B1B25" \
  --primary "#5E6AD2" \
  --accent "#26B5CE" \
  --font-heading "Inter" \
  --font-body "Inter" \
  --font-mono "Geist Mono" \
  --spacing compact \
  -o themes/linear-dark.json

# 2. Preview
figma-kit theme preview -t linear-dark
# → AI checks with get_screenshot to verify the vibe matches

# 3. Build the DS
figma-kit ds create -t linear-dark
```

---

## Theme-First Workflows

### Theme from a Mood

**Prompt:**

```
Create a warm, editorial theme. Earth tones — think coffee shop meets 
library. Serif headings, generous spacing. Then build a magazine-style 
article layout with it.
```

**What the AI does:**

```bash
# 1. Create the warm, editorial theme
figma-kit theme init \
  --name "Editorial" \
  --bg "#F5F0EB" \
  --primary "#8B4513" \
  --accent "#2F4F4F" \
  --font-heading "Playfair Display" \
  --font-body "Source Serif Pro" \
  --font-mono "IBM Plex Mono" \
  --spacing spacious \
  -o themes/editorial.json

# 2. Preview
figma-kit theme preview -t editorial

# 3. Build an article layout
figma-kit node create frame --name "Article" -w 800 --height 2000
figma-kit text create --content "The Lost Art of\nSlow Reading" \
  --font "Playfair Display" --weight Bold --size 56 --color "#8B4513" --parent <article-id>
figma-kit text create --content "By Margaret Chen  ·  12 min read" \
  --font "Source Serif Pro" --size 14 --color "#8B7355" --parent <article-id>
figma-kit ui divider -t editorial --parent <article-id>
figma-kit text create --content "In an age of infinite scrolling and algorithmic feeds, the simple act of reading a book cover to cover has become almost revolutionary..." \
  --font "Source Serif Pro" --size 18 --color "#3E2723" --parent <article-id> -w 700
```

---

### Iterate on an Existing Theme

**Prompt:**

```
Start from the Ocean community theme, but make it lighter for a SaaS 
app. Switch to Poppins headings, use spacious layout. Preview it, 
then build a settings page with it.
```

**What the AI does:**

```bash
# 1. Extend the Ocean theme with overrides
figma-kit theme init \
  --from themes/ocean.json \
  --name "Ocean Light" \
  --bg "#F0F4F8" \
  --font-heading "Poppins" \
  --spacing spacious \
  -o themes/ocean-light.json

# 2. Preview
figma-kit theme preview -t ocean-light

# 3. Build a settings page
figma-kit make settings -t ocean-light

# 4. Add specific UI components
figma-kit ui toggle -t ocean-light
figma-kit ui input -t ocean-light
figma-kit ui checkbox -t ocean-light
figma-kit ui dropdown -t ocean-light
```

---

### Build Everything from One Theme

**Prompt:**

```
I want to create a complete project in Figma for a fintech app called 
"Ledger". Dark blue theme, green accents. Build me: a landing page, 
a dashboard, a design system, OG image, and pitch deck. All consistent.
```

**What the AI does:**

```bash
# 1. Create the master theme
figma-kit theme init \
  --name "Ledger" \
  --bg "#0B1121" \
  --primary "#2563EB" \
  --accent "#10B981" \
  --font-heading "Inter" \
  --font-body "Inter" \
  --font-mono "Geist Mono" \
  -o themes/ledger.json

# 2. Page 1: Landing page
figma-kit page create "Landing" -p 0
figma-kit make screen --type landing --sections "hero,features,pricing,cta" -t ledger

# 3. Page 2: Dashboard
figma-kit page create "Dashboard"
figma-kit make dashboard --widgets "stat,stat,stat,chart,table" --cols 3 -t ledger

# 4. Page 3: Design system
figma-kit page create "Design System"
figma-kit ds create -t ledger

# 5. OG image
figma-kit make og-image --title "Ledger" --description "Financial intelligence for modern teams" -t ledger

# 6. Pitch deck
figma-kit make pitch-deck --slides 10 --template saas -t ledger

# 7. QA everything
figma-kit qa checklist --page 0
figma-kit qa checklist --page 1

# 8. Export tokens
figma-kit export tokens -t ledger --format css
```

**Result:** Five Figma pages with a complete, consistent project — all from one theme and a series of prompts. Every color, font, and spacing value matches across landing page, dashboard, design system, social assets, and pitch deck.

---

## Quick Wins

Short prompts for common tasks — each is one or two commands:

| Prompt | Command |
|--------|---------|
| "Create an OG image for my blog post" | `figma-kit make og-image --title "My Post" --description "..." -t noir` |
| "Add a primary button" | `figma-kit ui button --variant primary -t noir` |
| "Make a glass card about pricing" | `figma-kit card glass -t noir --title "Pro Plan" --desc "$29/mo"` |
| "Generate the design system" | `figma-kit ds create -t noir` |
| "Show me all color tokens as CSS" | `figma-kit export tokens -t noir --format css` |
| "Run a full QA check" | `figma-kit qa checklist --page 0` |
| "Create a wireframe" | `figma-kit make wireframe --type landing -t default` |
| "Add a mesh gradient background" | `figma-kit fx mesh <frame-id> -t noir` |
| "List all pages and frames" | `figma-kit status` |
| "Preview my custom theme" | `figma-kit theme preview -t my-theme` |

---

## Compose Workflows

### Landing Page with Compose

Instead of running individual commands one at a time, use `compose` to batch everything into a single `use_figma` call:

**Prompt:**

```
Build a SaaS landing page for "Apex" using compose — hero, 3 feature cards,
pricing, and a CTA. Dark theme with blue accents.
```

**What the AI does:**

```bash
# 1. Create the theme
figma-kit theme init \
  --name "Apex" \
  --bg "#0B0F1A" \
  --primary "#3B82F6" \
  --accent "#14B8A6" \
  -o themes/apex.json

# 2. Compose the entire page in one call
figma-kit exec compose -t apex \
  "ui hero --title 'Ship 10x Faster' --subtitle 'AI-powered deployment' --cta 'Get Started' --badge 'New'" \
  "card glass --title 'Lightning Fast' --preset strong" \
  "card glass --title 'Auto-Scale' --preset strong" \
  "card glass --title 'Zero Downtime' --preset strong" \
  "ui pricing --tiers '[{\"name\":\"Free\",\"price\":\"$0\",\"features\":[\"5 deploys\"]},{\"name\":\"Pro\",\"price\":\"$49\",\"highlighted\":true,\"features\":[\"Unlimited\",\"Priority\"]}]'"
```

Or define it as a recipe YAML for reuse:

```yaml
# apex-landing.yml
theme: apex
steps:
  - "ui hero --title 'Ship 10x Faster' --subtitle 'AI-powered deployment' --cta 'Get Started' --badge 'New'"
  - "card glass --title 'Lightning Fast' --preset strong"
  - "card glass --title 'Auto-Scale' --preset strong"
  - "card glass --title 'Zero Downtime' --preset strong"
  - "ui pricing --tiers '[{\"name\":\"Pro\",\"price\":\"$49\",\"highlighted\":true}]'"
```

```bash
figma-kit exec compose --recipe apex-landing.yml
```

**Result:** One round-trip to Figma creates the entire page — shared preamble, all components, all theme tokens. Roughly 5× faster than individual `exec` calls.

---

### Compose with Images

Embed local images alongside UI components in a single compose call:

**Prompt:**

```
Create a product page with our logo from ./brand/logo.png and a hero
background from ./assets/hero.jpg. Add a headline and CTA below.
```

**What the AI does:**

```bash
figma-kit exec compose -t noir \
  "image place ./brand/logo.png --name 'Logo' --width 200 --height 60" \
  "ui hero --title 'Welcome to Apex' --cta 'Learn More'" \
  "card glass --title 'Feature 1'" \
  "card glass --title 'Feature 2'"
```

For images larger than ~33 KB, start a local server first:

```bash
figma-kit image serve ./assets
# → Serving on http://127.0.0.1:8741

figma-kit exec compose -t noir \
  "image place http://127.0.0.1:8741/hero.jpg --width 1440 --height 900" \
  "ui hero --title 'Welcome'"
```

---

### Landing Page Section

**Prompt:**

```
Add a Features section with a label, subtitle, and two glass cards (Fast / Secure) under the noir theme, in one compose run.
```

**What the AI does:**

```bash
figma-kit compose -t noir \
  "ui section --title 'Features' --label 'WHAT WE OFFER' --subtitle 'Everything you need'" \
  "card glass --parent _results[0] --title 'Fast' --desc 'Blazing speed'" \
  "card glass --parent _results[0] --title 'Secure' --desc 'Enterprise ready'"
```

---

### Card with Effects

**Prompt:**

```
Create a glass Dashboard card, then add noise and a centered glow on that card in one shot.
```

**What the AI does:**

```bash
figma-kit compose -t noir \
  "card glass --title 'Dashboard'" \
  "fx noise --last" \
  "fx glow --last --position center"
```

---

### Stats Row

**Prompt:**

```
Show a stats row with commands count, themes count, and payload size using compose.
```

**What the AI does:**

```bash
figma-kit compose -t noir \
  "ui stat --items '[{\"value\":\"150+\",\"label\":\"Commands\"},{\"value\":\"3\",\"label\":\"Themes\"},{\"value\":\"<3KB\",\"label\":\"Payload\"}]'"
```

---

### Text with Typography

**Prompt:**

```
Add a bold centered headline "Hello World" in Inter at 48px with a 58px line height.
```

**What the AI does:**

```bash
figma-kit text create --content "Hello World" --font Inter --weight Bold --size 48 --align CENTER --line-height 58
```

---

## Tips

### Always start with a theme

Every project should begin with `figma-kit theme init`. Even if you only have vague color preferences, giving the AI 3 hex values produces a complete, consistent palette. All subsequent commands inherit those colors.

### Use content YAML for data-driven designs

For carousels, pitch decks, and any multi-slide deliverable, write a YAML content file first. This separates content from design and makes iteration fast — change the text without re-running design commands.

### Verify with screenshots

After every major step, the AI should call `get_screenshot` via the Figma MCP server to verify the result visually. This catches layout issues before they compound.

### Compose primitives for custom layouts

The `make` templates give you a head start, but the real power is composing primitives:
- `node create frame` → `layout auto` → `text create` → `style fill` → `fx glow`

This pipeline lets you build anything — the templates are just pre-composed versions of this.

### Use local images

You can place images from your local machine directly into Figma — no public URL needed:

```bash
# Small files (< 33 KB) — base64 embedded, zero setup
figma-kit image place ./assets/logo.png --name "Logo" --width 200 --height 60

# Fill an existing frame with an image
figma-kit image fill ./hero.jpg --node "2:5"

# Larger files — start a local server, then fetch
figma-kit image serve ./assets
# → Serving on http://127.0.0.1:8741
figma-kit image place http://127.0.0.1:8741/hero.jpg --width 1440 --height 900
```

**Prompt example:**

> "Add our logo from ./brand/logo.png to the hero section, 200x60. Then use the hero-bg.jpg from that same folder as a full-width background image."

The AI runs `figma-kit image place ./brand/logo.png --width 200 --height 60`, and if the background is too large for inline, it starts `figma-kit image serve ./brand` and fetches from the local URL.

### Direct execution (no AI middleman)

If you want to run commands directly without an AI agent:

```bash
# One-time authentication
figma-kit auth login

# Execute any command directly in your Figma file
figma-kit exec make carousel -t noir --content slides.yml
figma-kit exec card glass -t noir --title "Feature" --screenshot
figma-kit exec ui hero -t noir --title "Build Faster"

# Create a new file
figma-kit new-file "My Landing Page"
```

### New design patterns

**Trending card styles:**
```bash
figma-kit card neumorphic --title "Settings" --depth deep     # Soft UI
figma-kit card clay --title "Welcome" --color "#A78BFA"        # Puffy 3D
figma-kit card outline --title "API" --glow-color "#3B82F6"    # Ghost card
```

**Layout compositions:**
```bash
figma-kit ui hero -t noir --title "Ship Faster" --cta "Get Started" --badge "New"
figma-kit ui pricing -t noir --tiers '[{"name":"Pro","price":"$29","highlighted":true}]'
figma-kit ui feature-grid -t noir --cols 3
figma-kit ui testimonial -t noir --name "Jane" --quote "Changed everything" --rating 5
figma-kit ui timeline -t noir --entries '[{"date":"2024","title":"Launch"}]'
figma-kit ui accordion -t noir  # FAQ section
```

**Effects:**
```bash
figma-kit fx aurora <frameId> --palette sunset    # Northern lights
figma-kit fx morph <frameId> --count 5            # Organic blobs
figma-kit fx gradient-border <frameId>            # Gradient stroke
figma-kit fx spotlight <frameId>                  # Radial highlight
figma-kit fx pattern <frameId> --type dots        # Geometric patterns
```

**Shape operations:**
```bash
figma-kit node boolean union "1:2" "1:3"          # Boolean ops
figma-kit node svg "M10 10 L90 90" --fill "#3B82F6"  # SVG paths
```

### Export for developers

Once your design is done:
```bash
figma-kit export tokens -t my-theme --format css   # CSS variables
figma-kit handoff css <frame-id>                    # CSS for a specific frame
figma-kit handoff react <frame-id>                  # React component spec
```
