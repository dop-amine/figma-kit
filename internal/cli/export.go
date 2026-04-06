package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export assets and token payloads from Figma",
	}
	cmd.AddCommand(newExportPNGCmd())
	cmd.AddCommand(newExportSVGCmd())
	cmd.AddCommand(newExportPDFCmd())
	cmd.AddCommand(newExportPageCmd())
	cmd.AddCommand(newExportSpritesCmd())
	cmd.AddCommand(newExportTokensCmd())
	return cmd
}

func newExportPNGCmd() *cobra.Command {
	var scale float64
	cmd := &cobra.Command{
		Use:   "png <nodeId>",
		Short: "Generate JS that exports a node as PNG via exportAsync",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("const bytes = await node.exportAsync({ format: 'PNG', constraint: { type: 'SCALE', value: %s } });", codegen.FmtFloat(scale))
			b.Line("// bytes is Uint8Array — handle in host / MCP pipeline")
			b.Line("return { format: 'PNG', byteLength: bytes.length };")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().Float64Var(&scale, "scale", 2, "Export scale factor")
	return cmd
}

func newExportSVGCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "svg <nodeId>",
		Short: "Generate JS that exports a node as SVG",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const svg = await node.exportAsync({ format: 'SVG' });")
			b.Line("return { format: 'SVG', length: typeof svg === 'string' ? svg.length : svg.byteLength };")
			output(b.String())
			return nil
		},
	}
	return cmd
}

func newExportPDFCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pdf <nodeId>",
		Short: "Generate JS that exports a frame or page subtree as PDF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const bytes = await node.exportAsync({ format: 'PDF' });")
			b.Line("return { format: 'PDF', byteLength: bytes.length };")
			output(b.String())
			return nil
		},
	}
	return cmd
}

func newExportPageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "page",
		Short: "Generate JS that exports the current page as PNG slices per top-level frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const results = [];")
			b.Line("for (const child of figma.currentPage.children) {")
			b.Line("  if (child.type === 'FRAME' || child.type === 'SECTION' || child.type === 'COMPONENT') {")
			b.Line("    const bytes = await child.exportAsync({ format: 'PNG', constraint: { type: 'SCALE', value: 2 } });")
			b.Line("    results.push({ id: child.id, name: child.name, byteLength: bytes.length }); } }")
			b.Line("return { exports: results };")
			output(b.String())
			return nil
		},
	}
	return cmd
}

func newExportSpritesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sprites <frameId>",
		Short: "Generate JS that exports each direct child of a frame as PNG (sprite sheet workflow)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const frame = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!frame || !('children' in frame)) throw new Error('Frame not found');")
			b.Line("const sprites = [];")
			b.Line("for (const c of frame.children) {")
			b.Line("  const bytes = await c.exportAsync({ format: 'PNG', constraint: { type: 'SCALE', value: 1 } });")
			b.Line("  sprites.push({ id: c.id, name: c.name, byteLength: bytes.length }); }")
			b.Line("return { sprites };")
			output(b.String())
			return nil
		},
	}
	return cmd
}

func newExportTokensCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Output theme tokens (colors, type, spacing) in JSON or CSS variables (Go, no plugin)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			switch format {
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(t)
			case "css":
				_, _ = fmt.Fprintf(os.Stdout, "/* export tokens — %s */\n:root {\n", t.Name)
				names := t.ColorNames()
				sort.Strings(names)
				for _, name := range names {
					c := t.Colors[name]
					h := codegen.RGBToHex(codegen.RGB{R: c.R, G: c.G, B: c.B})
					key := name
					_, _ = fmt.Fprintf(os.Stdout, "  --fk-%s: %s;\n", key, h)
				}
				_, _ = fmt.Fprintln(os.Stdout, "}")
				return nil
			default:
				return fmt.Errorf("unknown format %q (json or css)", format)
			}
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "json or css")
	return cmd
}
