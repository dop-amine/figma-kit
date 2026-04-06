package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/mcpclient"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Figma MCP authentication",
		Long: `Authenticate with the Figma MCP server for direct command execution.

After logging in, figma-kit can execute commands directly in Figma
without needing an AI agent as a middleman.`,
		Example: `  figma-kit auth login    # authenticate with Figma
  figma-kit auth status   # check current auth state
  figma-kit auth logout   # clear cached credentials`,
	}
	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthLogoutCmd())
	cmd.AddCommand(newAuthStatusCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the Figma MCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Connecting to Figma MCP server...")
			ctx := context.Background()
			session, err := mcpclient.Connect(ctx)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}
			defer session.Close()

			result, err := session.CallWhoami(ctx)
			if err != nil {
				fmt.Println("Authenticated (could not fetch user info).")
				return nil
			}

			fmt.Println("Authenticated successfully!")
			fmt.Println(result.Raw)
			return nil
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear cached Figma credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := mcpclient.ClearToken(); err != nil {
				fmt.Println("No cached credentials found.")
				return nil
			}
			fmt.Println("Logged out. Cached credentials cleared.")
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check current Figma authentication state",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := mcpclient.LoadToken()
			if err != nil || !token.IsValid() {
				fmt.Println("Not authenticated.")
				fmt.Println("Run 'figma-kit auth login' to connect to Figma.")
				return
			}
			fmt.Println("Authenticated.")
			if !token.Expiry.IsZero() {
				fmt.Printf("Token expires: %s\n", token.Expiry.Format("2006-01-02 15:04:05"))
			}
		},
	}
}
