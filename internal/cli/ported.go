package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/assets"
	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/theme"
)

var embeddedTemplates = map[string]string{
	"slide":            assets.TemplateSlide,
	"one-pager-print":  assets.TemplateOnePager,
	"storyboard-panel": assets.TemplateStoryboard,
}

func newPreambleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "preamble",
		Short: "Generate the use_figma preamble (theme colors + font loading)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.Preamble(b, t)
			output(b.String())
			return nil
		},
	}
}

func newHelpersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "helpers",
		Short: "Output the full helpers.js code for injection into use_figma",
		Run: func(cmd *cobra.Command, args []string) {
			output(codegen.AllHelpers())
		},
	}
}

func newTemplateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "template <name>",
		Short: "Output a template (slide, one-pager-print, storyboard-panel)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			src, ok := embeddedTemplates[name]
			if !ok {
				names := make([]string, 0, len(embeddedTemplates))
				for k := range embeddedTemplates {
					names = append(names, k)
				}
				sort.Strings(names)
				return fmt.Errorf("template %q not found (available: %s)", name, strings.Join(names, ", "))
			}
			output(src)
			return nil
		},
	}
}

func newThemesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "themes",
		Short: "List available themes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("\nAvailable themes:\n\n")
			for _, info := range theme.List() {
				fmt.Printf("  %-12s %s\n", info.Key, info.Description)
			}
			fmt.Println()
		},
	}
}

func newScaffoldCmd() *cobra.Command {
	var templateName string

	cmd := &cobra.Command{
		Use:   "scaffold",
		Short: "Generate a full use_figma code block ready for execution",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()

			b := codegen.New()
			codegen.PreambleWithPage(b, t, page)
			b.Comment("--- Helpers ---")
			b.Raw(codegen.AllHelpers())

			if templateName != "" {
				src, ok := embeddedTemplates[templateName]
				if !ok {
					return fmt.Errorf("template %q not found", templateName)
				}
				b.Blank()
				b.Comment("--- Template ---")
				b.Raw(src)
			}

			b.Blank()
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}

	cmd.Flags().StringVar(&templateName, "template", "", "Template to include")
	return cmd
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show the figma-kit directory structure and capabilities",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			fmt.Println("figma-kit — CLI for programmatic Figma design")
			fmt.Println()
			fmt.Println("  Layers:")
			fmt.Println("    0  File & Session     init, config, whoami, open, status")
			fmt.Println("    1  Primitives         node, style, text, layout")
			fmt.Println("    2  Patterns           card, ui, fx")
			fmt.Println("    3  Deliverables       make carousel, make one-pager, ...")
			fmt.Println("    4  Design System      ds create, ds colors, ds component, ...")
			fmt.Println("    5  Inspect & QA       inspect, tree, qa contrast, qa checklist, ...")
			fmt.Println("    6  Export & Handoff    export png, handoff css, ...")
			fmt.Println("    7  Orchestration      batch <recipe.yml>")
			fmt.Println()
			fmt.Println("  Themes:")
			for _, info := range theme.List() {
				fmt.Printf("    %-12s %s\n", info.Key, info.Description)
			}
			fmt.Println()

			helpers := codegen.AvailableHelpers()
			sort.Strings(helpers)
			fmt.Printf("  Helpers:     %s\n", strings.Join(helpers, ", "))
			fmt.Println()

			templates := make([]string, 0, len(embeddedTemplates))
			for k := range embeddedTemplates {
				templates = append(templates, k)
			}
			sort.Strings(templates)
			fmt.Printf("  Templates:   %s\n", strings.Join(templates, ", "))
			fmt.Println()

			exe, _ := os.Executable()
			fmt.Printf("  Binary:      %s\n", filepath.Base(exe))
			fmt.Printf("  Version:     %s\n", Version)
			fmt.Println()
		},
	}
}
