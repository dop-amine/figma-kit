package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/mcpclient"
)

func newNewFileCmd() *cobra.Command {
	var planKey string
	cmd := &cobra.Command{
		Use:   "new-file <name>",
		Short: "Create a new Figma file via MCP",
		Long: `Create a new Figma design file. Requires authentication.

The file is created in your default Figma workspace. Use --plan-key
to specify a team/plan if you have multiple.`,
		Example: `  figma-kit new-file "My Landing Page"
  figma-kit new-file "Portfolio" --plan-key team123`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			session, err := mcpclient.Connect(ctx)
			if err != nil {
				return fmt.Errorf("not authenticated — run 'figma-kit auth login' first: %w", err)
			}
			defer func() { _ = session.Close() }()

			result, err := session.CallCreateFile(ctx, args[0], planKey)
			if err != nil {
				return err
			}
			if result.FileKey != "" {
				fmt.Printf("Created file: %s\n", result.FileKey)
			}
			if result.URL != "" {
				fmt.Printf("URL: %s\n", result.URL)
			}
			if result.FileKey == "" && result.URL == "" {
				fmt.Println(result.Raw)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&planKey, "plan-key", "", "Team/plan key for file creation")
	return cmd
}
