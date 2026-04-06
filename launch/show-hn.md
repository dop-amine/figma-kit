# Show HN: figma-kit — a CLI that generates production Figma designs from the terminal

**Title:** Show HN: figma-kit – a CLI that generates production Figma designs from the terminal

**URL:** https://github.com/dop-amine/figma-kit

---

**First comment (post immediately after submission):**

Hey HN! I built figma-kit because I wanted a faster way to create Figma designs without clicking through menus.

**What it does:** figma-kit is a single Go binary with 120+ commands that generate JavaScript for Figma's official MCP (Model Context Protocol) server. You run a command like `figma-kit make carousel --content slides.yml -t arkham`, it outputs Figma Plugin API JavaScript, and the MCP server executes it in your Figma file.

**How it's different from existing tools:**
- Uses the official Figma MCP server (no CDP/browser hacking)
- Goes beyond CRUD — it generates complete deliverables (carousels, pitch decks, OG images, wireframes)
- Built-in theme system with design tokens
- YAML-driven content specs for data-driven designs
- 10 automated QA checks (contrast, touch targets, orphaned layers, etc.)

**Architecture:** Pure Go binary, no runtime dependencies. JavaScript files are embedded via `go:embed` — they're the output format for Figma's Plugin API, not implementation code. The binary is essentially a sophisticated code generator.

**The 8 layers:**
0. Session management (init, config, open)
1. Primitives (node create/delete/move, style, text, layout)
2. Patterns (cards, UI components, visual effects)
3. Deliverables (36 templates: carousel, pitch-deck, wireframe, storyboard, ...)
4. Design system (create token specimens, audit palette usage)
5. QA (contrast checker, touch target validator, naming conventions)
6. Export (PNG/SVG/PDF, CSS/React handoff)
7. Batch orchestration (YAML recipes for multi-step workflows)

**Install:** `brew install dop-amine/tap/figma-kit` or `go install github.com/dop-amine/figma-kit@latest`

It's MIT licensed, free, and open source. Built this as a portfolio project to explore what's possible with Figma's new MCP integration.

Happy to answer any questions about the architecture, the Figma MCP protocol, or the Go codegen approach.
