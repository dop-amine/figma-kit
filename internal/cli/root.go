package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/config"
	"github.com/dop-amine/figma-kit/internal/theme"
)

var (
	// Version is set at build time via ldflags.
	Version = "dev"

	flagTheme string
	flagPage  int
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
  Direct exec:   figma-kit exec make carousel -t noir --content slides.yml
  Standalone:    figma-kit card glass -t noir | pipe to use_figma

Commands span 8 layers: node primitives, styles, cards (glass, neumorphic,
clay, outline), 29 UI components (hero, pricing, modal, pagination...),
14 effects (aurora, morph, spotlight...), 37 templates, design systems,
QA audits, export, and batch orchestration.

Run 'figma-kit cookbook' to browse real-world prompt examples.
Run 'figma-kit examples' to get starter content YAML files.
Run 'figma-kit docs' to read the full documentation.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       Version,
	}

	cmd.PersistentFlags().StringVarP(&flagTheme, "theme", "t", "", "Theme name (default, light, noir, or path)")
	cmd.PersistentFlags().IntVarP(&flagPage, "page", "p", -1, "Page index (0-based)")

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
