# Launch Checklist

## Pre-Launch (1 week before)

- [ ] Record 60-second demo: `figma-kit make carousel` producing a real carousel in Figma
  - Tool: `asciinema` or `vhs` (Charm) for terminal recording, screen record Figma side
  - Convert to GIF for GitHub/social, keep MP4 for website
- [ ] Generate OG image using `figma-kit make og-image --title "figma-kit" --description "Design from the command line"`
- [ ] Get 2-3 beta testers (designer friends) to try it, collect quotes
- [ ] Tag `v0.1.0` and push — triggers goreleaser, creates GitHub Release
- [ ] Verify Homebrew tap works: `brew install amine/tap/figma-kit`
- [ ] Verify `install.sh` works on clean machine
- [ ] Deploy website to GitHub Pages (or Cloudflare Pages)
- [ ] Final README review — ensure demo GIF is embedded

## Launch Day

### Hacker News (7-8am ET)
- [ ] Title: "Show HN: figma-kit – a CLI that generates production Figma designs from the terminal"
- [ ] URL: link to GitHub repo
- [ ] First comment: paste from `show-hn.md`
- [ ] Monitor and respond to all comments for 24 hours

### Product Hunt
- [ ] Create maker profile if needed
- [ ] Launch with demo GIF, tagline, 3 key features, install instructions
- [ ] Tagline: "Design from the command line — 120+ commands for Figma"

### Reddit
- [ ] r/golang — "I built a Go CLI that generates Figma designs from YAML specs"
- [ ] r/figma — "Open source CLI for programmatic Figma design via MCP"
- [ ] r/commandline — "figma-kit: generate production designs from your terminal"
- [ ] r/webdev — "CLI tool that creates Figma designs from templates and YAML content specs"

### Twitter/X
- [ ] Thread: build story + demo GIF + link
- [ ] Tag @figma, @golang, relevant design/dev accounts

### Figma Community Forum
- [ ] Post in Developer section with MCP integration angle

## Post-Launch (week 1)

- [ ] Respond to all GitHub issues within 24 hours
- [ ] Write Dev.to follow-up based on feedback
- [ ] Add community-requested themes/templates
- [ ] Submit to Figma's MCP community resource list
