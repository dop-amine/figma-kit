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
	var token string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the Figma MCP server",
		Long: `Authenticate with Figma for direct command execution.

Without flags, runs the full OAuth flow (requires a Figma PAT for first-time setup).
With --token, saves a pre-existing OAuth access token directly.`,
		Example: `  figma-kit auth login                    # Full OAuth flow
  figma-kit auth login --token <token>    # Direct token injection`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if token != "" {
				td := &mcpclient.TokenData{AccessToken: token}
				if err := mcpclient.SaveToken(td); err != nil {
					return fmt.Errorf("saving token: %w", err)
				}
				fmt.Println("Token saved. Verifying...")

				ctx := context.Background()
				session, err := mcpclient.Connect(ctx)
				if err != nil {
					fmt.Println("Token saved but verification failed:", err)
					fmt.Println("The token may still work — try 'figma-kit auth status'.")
					return nil
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
			}

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
	cmd.Flags().StringVar(&token, "token", "", "Save a pre-existing OAuth access token directly")
	return cmd
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
