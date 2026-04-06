package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/config"
	"github.com/dop-amine/figma-kit/internal/mcpclient"
)

func newExecCmd() *cobra.Command {
	var (
		fileKey    string
		screenshot bool
		timeout    int
	)
	cmd := &cobra.Command{
		Use:   "exec <command> [flags...]",
		Short: "Generate JS and execute it directly in Figma via the MCP server",
		Long: `Run any figma-kit command and execute the generated JavaScript directly
in your Figma file. No AI agent middleman needed.

Requires authentication: run 'figma-kit auth login' first.
Requires a file key: set via --file-key, .figmarc.json, or FIGMA_FILE_KEY env var.`,
		Example: `  # Execute a carousel creation directly in Figma
  figma-kit exec make carousel -t noir --content slides.yml

  # Create a glass card and verify with screenshot
  figma-kit exec --screenshot card glass -t noir --title "Feature"

  # Use a specific file key
  figma-kit exec --file-key abc123 make og-image -t noir`,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Re-parse our own flags from the raw args
			fk, ss, to, subArgs := parseExecFlags(args)
			if fk != "" {
				fileKey = fk
			}
			if ss {
				screenshot = true
			}
			if to > 0 {
				timeout = to
			}

			if len(subArgs) == 0 {
				return fmt.Errorf("no command specified — usage: figma-kit exec <command> [flags...]")
			}

			// Resolve file key
			if fileKey == "" {
				fileKey = os.Getenv("FIGMA_FILE_KEY")
			}
			if fileKey == "" {
				c, err := config.Load()
				if err == nil && c.FileKey != "" {
					fileKey = c.FileKey
				}
			}
			if fileKey == "" {
				return fmt.Errorf("no file key — set via --file-key, FIGMA_FILE_KEY, or 'figma-kit config set fileKey <key>'")
			}

			// Capture the sub-command's JS output
			js, err := captureSubCommand(subArgs)
			if err != nil {
				return err
			}
			if strings.TrimSpace(js) == "" {
				return fmt.Errorf("sub-command produced no output")
			}

			// Connect to Figma MCP
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
				defer cancel()
			}

			fmt.Fprintf(os.Stderr, "Connecting to Figma MCP server...\n")
			session, err := mcpclient.Connect(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "MCP connection failed: %v\n", err)
				fmt.Fprintf(os.Stderr, "Falling back to stdout output:\n\n")
				fmt.Print(js)
				return nil
			}
			defer session.Close()

			desc := fmt.Sprintf("figma-kit exec %s", strings.Join(subArgs, " "))
			if len(desc) > 200 {
				desc = desc[:200]
			}

			fmt.Fprintf(os.Stderr, "Executing in file %s...\n", fileKey)
			result, err := session.CallUseFigma(ctx, fileKey, js, desc)
			if err != nil {
				return fmt.Errorf("execution failed: %w", err)
			}

			if result.IsError {
				fmt.Fprintf(os.Stderr, "Error from Figma:\n")
				for _, c := range result.Content {
					fmt.Fprintln(os.Stderr, c)
				}
				return fmt.Errorf("Figma returned an error")
			}

			fmt.Fprintf(os.Stderr, "Done.\n")
			for _, c := range result.Content {
				fmt.Println(c)
			}

			if screenshot {
				fmt.Fprintf(os.Stderr, "Taking screenshot...\n")
				ssResult, ssErr := session.CallScreenshot(ctx, fileKey, "0:1")
				if ssErr != nil {
					fmt.Fprintf(os.Stderr, "Screenshot failed: %v\n", ssErr)
				} else {
					fmt.Println(ssResult)
				}
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&fileKey, "file-key", "", "Figma file key")
	cmd.Flags().BoolVar(&screenshot, "screenshot", false, "Take a screenshot after execution")
	cmd.Flags().IntVar(&timeout, "timeout", 60, "Timeout in seconds")
	return cmd
}

func parseExecFlags(args []string) (fileKey string, screenshot bool, timeout int, remaining []string) {
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == "--file-key" && i+1 < len(args):
			fileKey = args[i+1]
			i++
		case strings.HasPrefix(args[i], "--file-key="):
			fileKey = strings.TrimPrefix(args[i], "--file-key=")
		case args[i] == "--screenshot":
			screenshot = true
		case args[i] == "--timeout" && i+1 < len(args):
			fmt.Sscanf(args[i+1], "%d", &timeout)
			i++
		case strings.HasPrefix(args[i], "--timeout="):
			fmt.Sscanf(strings.TrimPrefix(args[i], "--timeout="), "%d", &timeout)
		default:
			remaining = append(remaining, args[i:]...)
			return
		}
	}
	return
}

// captureSubCommand runs the figma-kit sub-command and captures its stdout.
func captureSubCommand(args []string) (string, error) {
	root := newRootCmd()
	root.SetArgs(args)

	var buf bytes.Buffer
	root.SetOut(&buf)
	// Redirect the output function for this execution
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w

	execErr := root.Execute()

	w.Close()
	os.Stdout = origStdout

	var captured bytes.Buffer
	captured.ReadFrom(r)

	if execErr != nil {
		return "", fmt.Errorf("sub-command failed: %w", execErr)
	}

	return captured.String(), nil
}
