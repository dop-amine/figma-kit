package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/config"
	"github.com/dop-amine/figma-kit/internal/theme"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"

	flagTheme    string
	flagPage     int
	flagBodyOnly bool
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "figma-kit",
		Short: "CLI for programmatic Figma design via the MCP server",
		Long: `figma-kit — 150+ commands for AI-powered Figma design.

Figma's use_figma MCP tool runs raw Plugin API JavaScript. figma-kit turns
that into named, composable, theme-aware commands — so you (or your AI agent)
don't have to write 40 lines of JS for a glass card.

  AI workflow:   Prompt in Cursor / Claude Code → AI picks commands → Figma renders
  Compose:       figma-kit compose -t noir "ui section --title X" "card glass --parent _results[0]"
                 (batch N commands into 1 use_figma call, tree-shaken helpers)
  Direct exec:   figma-kit exec compose -t noir --recipe landing.yml
  Standalone:    figma-kit card glass -t noir | pipe to use_figma

Compose features: _results[] for cross-step references, --last on fx commands
for chaining, --parent on all composable commands, tree-shaken helpers (~3KB).

Commands span 8 layers: node primitives, styles, cards (glass, neumorphic,
clay, outline), 30 UI components (hero, section, pricing, modal, pagination...),
14 effects (aurora, morph, spotlight...), 37 templates, design systems,
QA audits, export, and compose (batch N commands into 1 call).

Run 'figma-kit cookbook' to browse real-world prompt examples.
Run 'figma-kit examples' to get starter content YAML files.
Run 'figma-kit docs' to read the full documentation.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
	}

	cmd.PersistentFlags().StringVarP(&flagTheme, "theme", "t", "", "Theme name (default, light, noir, or path)")
	cmd.PersistentFlags().IntVarP(&flagPage, "page", "p", -1, "Page index (0-based)")
	cmd.PersistentFlags().BoolVar(&flagBodyOnly, "body-only", false, "Emit only the command body (used internally by compose)")
	_ = cmd.PersistentFlags().MarkHidden("body-only")

	// Phase 1: ported commands
	cmd.AddCommand(newPreambleCmd())
	cmd.AddCommand(newHelpersCmd())
	cmd.AddCommand(newTemplateCmd())
	cmd.AddCommand(newThemesCmd())
	cmd.AddCommand(newScaffoldCmd())
	cmd.AddCommand(newInfoCmd())

	// Phase 2: Layer 0
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newConfigCmd())
	cmd.AddCommand(newWhoamiCmd())
	cmd.AddCommand(newOpenCmd())
	cmd.AddCommand(newStatusCmd())

	// Phase 3: Layer 1
	cmd.AddCommand(newNodeCmd())
	cmd.AddCommand(newStyleCmd())
	cmd.AddCommand(newTextCmd())
	cmd.AddCommand(newLayoutCmd())

	// Phase 4: Layer 2
	cmd.AddCommand(newCardCmd())
	cmd.AddCommand(newUICmd())
	cmd.AddCommand(newFXCmd())
	cmd.AddCommand(newImageCmd())

	// Phase 5: Layer 3
	cmd.AddCommand(newMakeCmd())

	// Phase 6: Layer 4
	cmd.AddCommand(newDSCmd())

	// Phase 7: Layer 5
	cmd.AddCommand(newInspectCmd())
	cmd.AddCommand(newScreenshotCmd())
	cmd.AddCommand(newTreeCmd())
	cmd.AddCommand(newFindCmd())
	cmd.AddCommand(newMeasureCmd())
	cmd.AddCommand(newDiffCmd())
	cmd.AddCommand(newQACmd())

	// Phase 8: Layer 6
	cmd.AddCommand(newExportCmd())
	cmd.AddCommand(newHandoffCmd())

	// Phase 9: Layer 7
	cmd.AddCommand(newBatchCmd())

	// Phase 10: v0.2 additions
	cmd.AddCommand(newPageCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newWatchCmd())

	// Phase 11: Theme management
	cmd.AddCommand(newThemeCmd())

	// Phase 12: Documentation & onboarding
	cmd.AddCommand(newCookbookCmd())
	cmd.AddCommand(newExamplesCmd())
	cmd.AddCommand(newDocsCmd())

	// Phase 13: Direct MCP execution
	cmd.AddCommand(newExecCmd())
	cmd.AddCommand(newAuthCmd())
	cmd.AddCommand(newNewFileCmd())

	// Phase 14: Compose
	cmd.AddCommand(newComposeCmd())

	return cmd
}

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

// resolveTheme loads the theme from flags, config, or falls back to "default".
func resolveTheme(cmd *cobra.Command) (*theme.Theme, error) {
	name := flagTheme
	if name == "" {
		c, err := config.Load()
		if err == nil && c.Theme != "" {
			name = c.Theme
		} else {
			name = "default"
		}
	}
	return theme.Load(name)
}

// resolvePage returns the page index from flags or config.
func resolvePage() int {
	if flagPage >= 0 {
		return flagPage
	}
	c, err := config.Load()
	if err == nil {
		return c.Page
	}
	return 0
}

// resolveFileKey returns the file key from env or config.
func resolveFileKey() string {
	if fk := os.Getenv("FIGMA_FILE_KEY"); fk != "" {
		return fk
	}
	c, err := config.Load()
	if err == nil && c.FileKey != "" {
		return c.FileKey
	}
	return ""
}

// output writes the generated JS to stdout.
func output(js string) {
	_, _ = fmt.Fprint(os.Stdout, js)
}

// newBuilder creates a codegen.Builder that respects the --body-only flag.
func newBuilder() *codegen.Builder {
	b := codegen.New()
	if flagBodyOnly {
		b.SetBodyOnly(true)
	}
	return b
}
