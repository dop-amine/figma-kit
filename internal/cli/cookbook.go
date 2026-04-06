package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/assets"
)

func newCookbookCmd() *cobra.Command {
	var listFlag bool

	cmd := &cobra.Command{
		Use:   "cookbook [example]",
		Short: "Browse prompt examples for AI-driven Figma design",
		Long: `Browse real-world prompt examples showing how to use figma-kit with AI.

Each example shows the prompt you'd type, the commands the AI runs, and the
result you get in Figma. Copy these prompts directly into Cursor, Claude Code,
or any MCP-compatible AI agent.

Run with no arguments to see all examples, or specify a section name.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content := assets.CookbookMD

			if listFlag {
				printCookbookIndex(content)
				return nil
			}

			if len(args) > 0 {
				section := findSection(content, args[0])
				if section == "" {
					fmt.Fprintf(os.Stderr, "No example matching %q. Run 'figma-kit cookbook --list' to see available examples.\n", args[0])
					return nil
				}
				fmt.Print(section)
				return nil
			}

			fmt.Print(content)
			return nil
		},
	}

	cmd.Flags().BoolVar(&listFlag, "list", false, "List available examples")

	return cmd
}

var tipsHeadings = map[string]bool{
	"Tips": true, "Quick Wins": true,
	"Always start with a theme": true, "Use content YAML for data-driven designs": true,
	"Verify with screenshots": true, "Compose primitives for custom layouts": true,
	"Export for developers": true,
}

func printCookbookIndex(content string) {
	fmt.Println("Available cookbook examples:")
	fmt.Println()

	inTips := false
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "## ") {
			heading := strings.TrimPrefix(line, "## ")
			if tipsHeadings[heading] {
				inTips = true
				continue
			}
			inTips = false
			fmt.Printf("\n  %s\n", heading)
		} else if strings.HasPrefix(line, "### ") && !inTips {
			name := strings.TrimPrefix(line, "### ")
			if tipsHeadings[name] {
				continue
			}
			slug := slugify(name)
			fmt.Printf("  %-35s figma-kit cookbook %s\n", name, slug)
		}
	}
	fmt.Println()
	fmt.Println("Run 'figma-kit cookbook <name>' to see a specific example.")
	fmt.Println("Run 'figma-kit cookbook' to see everything.")
}

func findSection(content, query string) string {
	query = strings.ToLower(query)
	lines := strings.Split(content, "\n")

	type section struct {
		title string
		body  string
		start int
		end   int
	}

	var sections []section
	var cur *section

	for i, line := range lines {
		if strings.HasPrefix(line, "### ") {
			if cur != nil {
				cur.end = i
				sections = append(sections, *cur)
			}
			title := strings.TrimPrefix(line, "### ")
			cur = &section{title: title, start: i}
		} else if strings.HasPrefix(line, "## ") {
			if cur != nil {
				cur.end = i
				sections = append(sections, *cur)
				cur = nil
			}
		}
	}
	if cur != nil {
		cur.end = len(lines)
		sections = append(sections, *cur)
	}

	for i := range sections {
		bodyLines := lines[sections[i].start:sections[i].end]
		sections[i].body = strings.Join(bodyLines, "\n")
	}

	// First: exact slug/title match
	for _, s := range sections {
		if matchesQuery(s.title, query) {
			return s.body + "\n"
		}
	}
	// Second: body content match
	for _, s := range sections {
		if strings.Contains(strings.ToLower(s.body), query) {
			return s.body + "\n"
		}
	}
	return ""
}

func matchesQuery(title, query string) bool {
	slug := slugify(title)
	lower := strings.ToLower(title)
	return strings.Contains(slug, query) || strings.Contains(lower, query)
}

func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, s)
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	return strings.Trim(s, "-")
}

func newExamplesCmd() *cobra.Command {
	var dumpFlag bool

	cmd := &cobra.Command{
		Use:   "examples [name]",
		Short: "List or dump starter YAML content files",
		Long: `List and access ready-to-use YAML content files for figma-kit commands.

These files work with 'figma-kit make carousel --content <file>' and similar
content-driven commands. They provide real starter content you can customize.

Run with no arguments to list available examples.
Use --dump to write all files to ./examples/ in the current directory.
Specify a name to print a specific file to stdout (pipeable).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := assets.ExamplesFS.ReadDir("examples")
			if err != nil {
				return fmt.Errorf("reading embedded examples: %w", err)
			}

			if dumpFlag {
				return dumpExamples(entries)
			}

			if len(args) > 0 {
				return printExample(entries, args[0])
			}

			listExamples(entries)
			return nil
		},
	}

	cmd.Flags().BoolVar(&dumpFlag, "dump", false, "Write all example files to ./examples/")

	return cmd
}

func listExamples(entries []os.DirEntry) {
	fmt.Println("Available example content files:")
	fmt.Println()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		desc := describeExample(name)
		fmt.Printf("  %-28s %s\n", name, desc)
	}
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  figma-kit examples <name>    Print to stdout")
	fmt.Println("  figma-kit examples --dump    Write all to ./examples/")
	fmt.Println()
	fmt.Println("Then use with: figma-kit make carousel --content examples/<name>")
}

func describeExample(name string) string {
	switch {
	case strings.Contains(name, "saas"):
		return "SaaS landing page carousel (5 slides)"
	case strings.Contains(name, "pitch"):
		return "Startup pitch deck (7 slides)"
	case strings.Contains(name, "social"):
		return "Product launch social media campaign (5 slides)"
	case strings.Contains(name, "portfolio"):
		return "Photographer portfolio showcase (4 slides)"
	case strings.Contains(name, "dashboard"):
		return "Dashboard widget layout reference"
	default:
		return ""
	}
}

func printExample(entries []os.DirEntry, query string) error {
	query = strings.ToLower(query)
	if !strings.HasSuffix(query, ".yml") && !strings.HasSuffix(query, ".yaml") {
		query += ".yml"
	}
	for _, e := range entries {
		if strings.ToLower(e.Name()) == query || strings.Contains(strings.ToLower(e.Name()), strings.TrimSuffix(query, ".yml")) {
			data, err := assets.ExamplesFS.ReadFile("examples/" + e.Name())
			if err != nil {
				return err
			}
			fmt.Print(string(data))
			return nil
		}
	}
	return fmt.Errorf("no example matching %q — run 'figma-kit examples' to see available files", query)
}

func dumpExamples(entries []os.DirEntry) error {
	dir := "examples"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating examples directory: %w", err)
	}

	count := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := assets.ExamplesFS.ReadFile("examples/" + e.Name())
		if err != nil {
			return err
		}
		dest := filepath.Join(dir, e.Name())
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return err
		}
		fmt.Printf("  wrote %s\n", dest)
		count++
	}
	fmt.Printf("\n%d example files written to ./%s/\n", count, dir)
	fmt.Println("Use with: figma-kit make carousel --content examples/saas-landing.yml -t noir")
	return nil
}

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs [topic]",
		Short: "Open documentation in the browser",
		Long: `Open figma-kit documentation in your browser.

Topics:
  (none)     Main site: dop-amine.github.io/figma-kit
  cookbook    Prompt cookbook with real-world examples
  commands   Full CLI command reference
  themes     Theme system documentation
  recipes    YAML batch recipe guide
  examples   Example content files`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			base := "https://github.com/dop-amine/figma-kit/blob/main"
			url := "https://dop-amine.github.io/figma-kit/"

			if len(args) > 0 {
				switch strings.ToLower(args[0]) {
				case "cookbook":
					url = base + "/docs/COOKBOOK.md"
				case "commands":
					url = base + "/docs/COMMANDS.md"
				case "themes":
					url = base + "/docs/THEMES.md"
				case "recipes":
					url = base + "/docs/RECIPES.md"
				case "architecture", "arch":
					url = base + "/docs/ARCHITECTURE.md"
				case "examples":
					url = base + "/assets/examples"
				case "contributing", "contribute":
					url = base + "/CONTRIBUTING-THEMES.md"
				default:
					fmt.Fprintf(os.Stderr, "Unknown topic %q. Available: cookbook, commands, themes, recipes, architecture, examples, contributing\n", args[0])
					return nil
				}
			}

			fmt.Printf("Opening %s\n", url)
			return openBrowser(url)
		},
	}

	return cmd
}
