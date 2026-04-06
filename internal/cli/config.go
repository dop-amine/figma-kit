package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/amine/figma-kit/internal/codegen"
	"github.com/amine/figma-kit/internal/config"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new figma-kit project with .figmarc.json",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "Untitled"
			if len(args) > 0 {
				name = args[0]
			}
			if err := config.Init(name); err != nil {
				return err
			}
			fmt.Printf("Created .figmarc.json for %q\n", name)
			fmt.Println("Tip: use 'figma-kit config set fileKey <key>' to link a Figma file.")
			return nil
		},
	}
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage .figmarc.json project configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a config value (fileKey, theme, page, exportDir)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Set(args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Set %s = %s\n", args[0], args[1])
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := config.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all config values",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			fmt.Printf("fileKey:   %s\n", c.FileKey)
			fmt.Printf("theme:     %s\n", c.Theme)
			fmt.Printf("page:      %d\n", c.Page)
			fmt.Printf("exportDir: %s\n", c.ExportDir)
			return nil
		},
	})

	return cmd
}

func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show Figma authentication status (wraps MCP whoami tool)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("To check your Figma identity, ask the AI to call the 'whoami' MCP tool.")
			fmt.Println("This returns your email, plans, and seat type.")
		},
	}
}

func newOpenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open",
		Short: "Open the current Figma file in a browser",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			if c.FileKey == "" {
				return fmt.Errorf("no fileKey configured — run 'figma-kit config set fileKey <key>'")
			}
			url := fmt.Sprintf("https://www.figma.com/file/%s", c.FileKey)
			fmt.Printf("Opening %s\n", url)
			return openBrowser(url)
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Generate JS to inspect file structure (pages, frames, node counts)",
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.Comment("figma-kit status: inspect file structure")
			b.Line("const pages = figma.root.children;")
			b.Line("const result = [];")
			b.Line("for (const page of pages) {")
			b.Line("  await figma.setCurrentPageAsync(page);")
			b.Line("  const frames = page.children.filter(n => n.type === 'FRAME');")
			b.Line("  result.push({")
			b.Line("    name: page.name,")
			b.Line("    id: page.id,")
			b.Line("    frameCount: frames.length,")
			b.Line("    frames: frames.map(f => ({ name: f.name, id: f.id, w: f.width, h: f.height }))")
			b.Line("  });")
			b.Line("}")
			b.ReturnExpr("result")
			output(b.String())
			return nil
		},
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported OS for browser open: %s", runtime.GOOS)
	}
	return cmd.Start()
}
