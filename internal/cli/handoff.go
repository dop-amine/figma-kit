package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newHandoffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handoff",
		Short: "Specs, redlines, and code-oriented handoff helpers",
		Example: `  # "Generate CSS for this component"
  figma-kit handoff css <nodeId>

  # "Create a React spec for the card"
  figma-kit handoff react <nodeId>

  # "Add measurement redlines"
  figma-kit handoff redline <nodeId>`,
	}
	cmd.AddCommand(newHandoffSpecCmd())
	cmd.AddCommand(newHandoffRedlineCmd())
	cmd.AddCommand(newHandoffCSSCmd())
	cmd.AddCommand(newHandoffReactCmd())
	cmd.AddCommand(newHandoffAssetsCmd())
	return cmd
}

func newHandoffSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec <nodeId>",
		Short: "Generate JS that returns a Markdown spec from node properties",
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
			b.Line("const lines = [];")
			b.Line("lines.push('# ' + node.name); lines.push(''); lines.push('- **Type:** ' + node.type); lines.push('- **ID:** `' + node.id + '`');")
			b.Line("if ('width' in node) { lines.push('- **Size:** ' + Math.round(node.width) + ' × ' + Math.round(node.height) + ' px'); }")
			b.Line("if ('layoutMode' in node && node.layoutMode !== 'NONE') {")
			b.Line("  lines.push('- **Auto-layout:** ' + node.layoutMode + ', gap ' + node.itemSpacing + ' px'); }")
			b.Line("if ('fills' in node && node.fills && node.fills[0]) { lines.push('- **Fill:** ' + JSON.stringify(node.fills[0])); }")
			b.Line("if ('cornerRadius' in node && node.cornerRadius) { lines.push('- **Radius:** ' + node.cornerRadius + ' px'); }")
			b.Line("const md = lines.join('\\n');")
			b.Line("return { markdown: md };")
			output(b.String())
			return nil
		},
	}
}

func newHandoffRedlineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "redline <nodeId>",
		Short: "Generate JS that draws measurement overlay lines on a duplicate above the node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const target = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!target || !('absoluteBoundingBox' in target) || !target.absoluteBoundingBox) throw new Error('Node or bounds not found');")
			b.Line("const box = target.absoluteBoundingBox;")
			b.Line("const overlay = figma.createFrame(); overlay.name = 'Redline / ' + target.name;")
			b.Line("overlay.x = box.x; overlay.y = box.y - 24; overlay.resize(box.width, 24);")
			b.Line("overlay.fills = [{type:'SOLID', color:{r:1,g:0.2,b:0.2}, opacity:0.12}];")
			b.Line("if (target.parent && 'appendChild' in target.parent) target.parent.insertChild(target.parent.children.indexOf(target) + 1, overlay);")
			b.Line("else figma.currentPage.appendChild(overlay);")
			b.Line("const rule = figma.createLine(); rule.name = 'width';")
			b.Line("rule.strokes = [{type:'SOLID', color:{r:1,g:0,b:0}}]; rule.strokeWeight = 2;")
			b.Line("rule.x = 0; rule.y = 12; rule.resize(box.width, 0); overlay.appendChild(rule);")
			b.ReturnIDs("overlay.id", "rule.id")
			output(b.String())
			return nil
		},
	}
}

func newHandoffCSSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "css <nodeId>",
		Short: "Generate JS that builds a CSS snippet from solid fills, radius, and dimensions",
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
			b.Line("function toHex(c){ const h = x => Math.round(x*255).toString(16).padStart(2,'0'); return '#' + h(c.r)+h(c.g)+h(c.b); }")
			b.Line("const parts = [];")
			b.Line("if ('width' in node) { parts.push('width: ' + Math.round(node.width) + 'px;'); parts.push('height: ' + Math.round(node.height) + 'px;'); }")
			b.Line("if ('cornerRadius' in node && node.cornerRadius) parts.push('border-radius: ' + node.cornerRadius + 'px;');")
			b.Line("if ('fills' in node && node.fills && node.fills[0] && node.fills[0].type === 'SOLID') {")
			b.Line("  parts.push('background: ' + toHex(node.fills[0].color) + ';'); }")
			b.Line("if ('strokes' in node && node.strokes && node.strokes[0] && node.strokes[0].type === 'SOLID') {")
			b.Line("  parts.push('border: ' + (node.strokeWeight||1) + 'px solid ' + toHex(node.strokes[0].color) + ';'); }")
			b.Line("const css = '.' + node.name.replace(/\\s+/g,'-').toLowerCase() + ' {\\n  ' + parts.join('\\n  ') + '\\n}';")
			b.Line("return { css };")
			output(b.String())
			return nil
		},
	}
}

func newHandoffReactCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "react <nodeId>",
		Short: "Instructions for generating React via the get_design_context MCP tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := strings.TrimSpace(fmt.Sprintf(`
Use the Figma MCP tool **get_design_context** with node ID **%s**.

1. Pass the node ID so the server returns structured layout + style context for that subtree.
2. Ask your assistant to translate the response into React (or TSX) components.
3. Cross-check colors and spacing against your active figma-kit theme: run **figma-kit preamble** or **figma-kit ds sync-tokens** for token alignment.

figma-kit does not call MCP directly; invoke get_design_context from the MCP client.
`, args[0]))
			_, _ = fmt.Fprint(os.Stdout, msg, "\n")
			return nil
		},
	}
}

func newHandoffAssetsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "assets <nodeId>",
		Short: "Generate JS that lists export settings and child assets under a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const root = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!root) throw new Error('Node not found');")
			b.Line("const assets = [];")
			b.Line("function walk(n, depth){")
			b.Line("  const exp = 'exportSettings' in n ? n.exportSettings : [];")
			b.Line("  assets.push({ id: n.id, name: n.name, type: n.type, depth, exportSettings: exp });")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c, depth+1); }")
			b.Line("walk(root, 0);")
			b.Line("return { assets };")
			output(b.String())
			return nil
		},
	}
}
