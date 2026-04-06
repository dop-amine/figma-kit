package cli

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/theme"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate theme JSON or recipe YAML without executing",
	}
	cmd.AddCommand(newValidateThemeCmd())
	cmd.AddCommand(newValidateRecipeCmd())
	return cmd
}

func newValidateThemeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "theme <path-or-name>",
		Short: "Parse and type-check a theme file (or built-in name)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			arg := args[0]

			var th *theme.Theme
			var err error

			// Try as a built-in name first, then as a file path.
			if !strings.Contains(arg, "/") && !strings.HasSuffix(arg, ".json") {
				th, err = theme.Load(arg)
				if err != nil {
					return fmt.Errorf("unknown built-in theme %q: %w", arg, err)
				}
			} else {
				th, err = theme.LoadFile(arg)
				if err != nil {
					return fmt.Errorf("theme file %q: %w", arg, err)
				}
			}

			_, _ = fmt.Fprintf(os.Stdout, "✓ theme: %q\n", th.Name)
			_, _ = fmt.Fprintf(os.Stdout, "  colors:   %d tokens\n", len(th.Colors))
			_, _ = fmt.Fprintf(os.Stdout, "  type:     %d scale entries\n", len(th.Type))
			_, _ = fmt.Fprintf(os.Stdout, "  fonts:    heading=%q  body=%q  mono=%q\n", th.Fonts.Heading, th.Fonts.Body, th.Fonts.Mono)
			if th.Brand != nil {
				_, _ = fmt.Fprintf(os.Stdout, "  brand:    primary=%q  url=%q\n", th.Brand.Primary, th.Brand.URL)
			}

			var warnings []string
			for _, key := range []string{"BG", "WT", "BL", "CARD", "STK"} {
				if _, ok := th.Colors[key]; !ok {
					warnings = append(warnings, fmt.Sprintf("missing conventional color token %q (used by most templates)", key))
				}
			}
			if len(th.Type) == 0 {
				warnings = append(warnings, "no type scale defined (needed for text-heavy templates)")
			}
			if th.Fonts.Heading == "" && th.Fonts.Body == "" {
				warnings = append(warnings, "no fonts defined (heading/body will default to Inter)")
			}
			if th.Effects.Glass == nil && th.Effects.Shadow == nil {
				warnings = append(warnings, "no effects defined (glass/shadow presets used by card and make commands)")
			}
			if th.Spacing == (theme.SpacingSpec{}) {
				warnings = append(warnings, "no spacing defined (used by make carousel, one-pager, storyboard)")
			}

			if len(warnings) > 0 {
				_, _ = fmt.Fprintf(os.Stdout, "  warnings: %d\n", len(warnings))
				for _, w := range warnings {
					_, _ = fmt.Fprintf(os.Stdout, "    ⚠ %s\n", w)
				}
			}

			_, _ = fmt.Fprintln(os.Stdout, "  → OK")
			return nil
		},
	}
}

func newValidateRecipeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recipe <path>",
		Short: "Parse and dry-run a batch recipe YAML",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read: %w", err)
			}

			var rec batchRecipe
			if err := yaml.Unmarshal(data, &rec); err != nil {
				return fmt.Errorf("parse YAML: %w", err)
			}

			var errs []string
			if rec.Title == "" {
				errs = append(errs, "missing top-level 'title' field")
			}
			if len(rec.Steps) == 0 {
				errs = append(errs, "recipe has no steps")
			}
			for i, step := range rec.Steps {
				if strings.TrimSpace(step.JS) == "" {
					errs = append(errs, fmt.Sprintf("step %d (%q): empty 'js' field", i+1, step.Title))
				}
			}

			if len(errs) > 0 {
				_, _ = fmt.Fprintf(os.Stdout, "✗ recipe: %q\n", args[0])
				for _, e := range errs {
					_, _ = fmt.Fprintf(os.Stdout, "  error: %s\n", e)
				}
				return fmt.Errorf("%d validation error(s)", len(errs))
			}

			_, _ = fmt.Fprintf(os.Stdout, "✓ recipe: %q\n", rec.Title)
			_, _ = fmt.Fprintf(os.Stdout, "  steps: %d\n", len(rec.Steps))
			for i, step := range rec.Steps {
				lines := len(strings.Split(strings.TrimSpace(step.JS), "\n"))
				_, _ = fmt.Fprintf(os.Stdout, "  [%d] %q  (%d lines of JS)\n", i+1, step.Title, lines)
			}
			_, _ = fmt.Fprintln(os.Stdout, "  → OK")
			return nil
		},
	}
}
