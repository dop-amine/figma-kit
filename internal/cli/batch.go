package cli

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

// batchRecipe is the YAML shape for figma-kit batch files.
type batchRecipe struct {
	Title string `yaml:"title"`
	Steps []struct {
		Title string `yaml:"title"`
		JS    string `yaml:"js"`
	} `yaml:"steps"`
}

func newBatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "batch <recipe.yaml>",
		Short: "Read a YAML recipe and print numbered use_figma JS blocks",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read recipe: %w", err)
			}
			var rec batchRecipe
			if err := yaml.Unmarshal(data, &rec); err != nil {
				return fmt.Errorf("parse yaml: %w", err)
			}
			if len(rec.Steps) == 0 {
				return fmt.Errorf("recipe has no steps")
			}
			var out strings.Builder
			if rec.Title != "" {
				fmt.Fprintf(&out, "// Recipe: %s\n\n", rec.Title)
			}
			for i, step := range rec.Steps {
				label := step.Title
				if label == "" {
					label = fmt.Sprintf("Step %d", i+1)
				}
				js := strings.TrimSpace(step.JS)
				if js == "" {
					return fmt.Errorf("step %d (%q) has empty js", i+1, label)
				}
				fmt.Fprintf(&out, "// --- Block %d: %s ---\n", i+1, label)
				out.WriteString(js)
				if !strings.HasSuffix(js, "\n") {
					out.WriteByte('\n')
				}
				out.WriteByte('\n')
			}
			output(out.String())
			return nil
		},
	}
}
