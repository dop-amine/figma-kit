package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/mcpclient"
)

func newInspectCmd() *cobra.Command {
	var deep bool
	cmd := &cobra.Command{
		Use:   "inspect <nodeId>",
		Short: "Generate JS that dumps node properties (fills, strokes, layout, geometry)",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const base = { id: node.id, name: node.name, type: node.type, visible: node.visible, locked: node.locked };")
			if deep {
				b.Line("const box = 'absoluteBoundingBox' in node ? node.absoluteBoundingBox : null;")
				b.Line("const out = { ...base, x: 'x' in node ? node.x : null, y: 'y' in node ? node.y : null,")
				b.Line("  width: 'width' in node ? node.width : null, height: 'height' in node ? node.height : null,")
				b.Line("  fills: 'fills' in node ? node.fills : undefined, strokes: 'strokes' in node ? node.strokes : undefined,")
				b.Line("  effects: 'effects' in node ? node.effects : undefined, layoutMode: 'layoutMode' in node ? node.layoutMode : undefined,")
				b.Line("  absoluteBoundingBox: box };")
			} else {
				b.Line("const out = { ...base, fills: 'fills' in node ? node.fills : undefined, strokes: 'strokes' in node ? node.strokes : undefined };")
			}
			b.Line("return out;")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().BoolVar(&deep, "deep", false, "Include geometry, effects, and auto-layout fields")
	return cmd
}

func newScreenshotCmd() *cobra.Command {
	var nodeID string
	cmd := &cobra.Command{
		Use:   "screenshot",
		Short: "Capture a screenshot of a Figma node via MCP (falls back to instructions)",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fk := resolveFileKey()
			if fk == "" || nodeID == "" {
				msg := strings.TrimSpace(`
Usage: figma-kit screenshot --node <nodeId>

Requires authentication (figma-kit auth login) and a file key.
Set the file key via .figmarc.json or FIGMA_FILE_KEY env var.
`)
				_, _ = fmt.Fprint(os.Stdout, msg, "\n")
				return nil
			}
			ctx := cmd.Context()
			session, err := mcpclient.Connect(ctx)
			if err != nil {
				fmt.Println("Not authenticated. Run 'figma-kit auth login' first.")
				return nil
			}
			defer session.Close()
			result, sErr := session.CallScreenshot(ctx, fk, nodeID)
			if sErr != nil {
				return sErr
			}
			fmt.Println(result)
			return nil
		},
	}
	cmd.Flags().StringVar(&nodeID, "node", "", "Node ID to screenshot")
	return cmd
}

func newTreeCmd() *cobra.Command {
	var maxDepth int
	cmd := &cobra.Command{
		Use:   "tree [nodeId]",
		Short: "Generate JS that prints a hierarchical node tree (default: current page)",
		Args:  cobra.MaximumNArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const maxDepth = %d;", maxDepth)
			b.Line("function describe(n, depth) {")
			b.Line("  const pad = '  '.repeat(depth);")
			b.Line("  let line = pad + n.type + ' ' + JSON.stringify(n.name) + ' [' + n.id + ']';")
			b.Line("  if ('width' in n && 'height' in n) line += ' ' + Math.round(n.width) + 'x' + Math.round(n.height);")
			b.Line("  const lines = [line];")
			b.Line("  if (depth < maxDepth && 'children' in n) for (const c of n.children) lines.push(...describe(c, depth + 1));")
			b.Line("  return lines;")
			b.Line("}")
			if len(args) == 1 {
				b.Linef("const root = await figma.getNodeByIdAsync(%q);", args[0])
				b.Line("if (!root) throw new Error('Node not found');")
				b.Line("const text = describe(root, 0).join('\\n');")
			} else {
				b.Line("const text = describe(figma.currentPage, 0).join('\\n');")
			}
			b.Line("return { tree: text };")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&maxDepth, "max-depth", 12, "Maximum tree depth")
	return cmd
}

func newFindCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find <pattern>",
		Short: "Generate JS that finds nodes by substring match on name",
		Args:  cobra.ExactArgs(1),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			pat := strings.ToLower(args[0])
			b.Linef("const needle = %q;", pat)
			b.Line("const hits = [];")
			b.Line("function walk(n){ if (n.name && n.name.toLowerCase().includes(needle)) hits.push({id:n.id,name:n.name,type:n.type});")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { matches: hits };")
			output(b.String())
			return nil
		},
	}
	return cmd
}

func newMeasureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "measure <nodeIdA> <nodeIdB>",
		Short: "Generate JS that measures axis-aligned distance between two nodes' bounding boxes",
		Args:  cobra.ExactArgs(2),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const a = await figma.getNodeByIdAsync(%q);", args[0])
			b.Linef("const b = await figma.getNodeByIdAsync(%q);", args[1])
			b.Line("if (!a || !b) throw new Error('Node not found');")
			b.Line("const ba = a.absoluteBoundingBox; const bb = b.absoluteBoundingBox;")
			b.Line("if (!ba || !bb) throw new Error('Bounding boxes unavailable');")
			b.Line("const dx = Math.max(0, Math.max(bb.x - (ba.x + ba.width), ba.x - (bb.x + bb.width)));")
			b.Line("const dy = Math.max(0, Math.max(bb.y - (ba.y + ba.height), ba.y - (bb.y + bb.height)));")
			b.Line("const gap = Math.hypot(dx, dy);")
			b.Line("return { gapPx: gap, deltaX: dx, deltaY: dy, a: ba, b: bb };")
			output(b.String())
			return nil
		},
	}
}

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <nodeIdA> <nodeIdB>",
		Short: "Generate JS that compares dimensions and fill colors between two nodes",
		Args:  cobra.ExactArgs(2),
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const a = await figma.getNodeByIdAsync(%q);", args[0])
			b.Linef("const b = await figma.getNodeByIdAsync(%q);", args[1])
			b.Line("if (!a || !b) throw new Error('Node not found');")
			b.Line("function snap(n){ return {")
			b.Line("  name: n.name, type: n.type,")
			b.Line("  w: 'width' in n ? n.width : null, h: 'height' in n ? n.height : null,")
			b.Line("  fill0: ('fills' in n && n.fills && n.fills[0] && n.fills[0].type === 'SOLID') ? n.fills[0].color : null }; }")
			b.Line("return { a: snap(a), b: snap(b) };")
			output(b.String())
			return nil
		},
	}
}

func newQACmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qa",
		Short: "Design quality checks (contrast, targets, typography, ...)",
		Example: `  # "Run a full QA check on my design"
  figma-kit qa checklist --page 0

  # "Check text contrast for accessibility"
  figma-kit qa contrast --page 0

  # "Find buttons that are too small for touch"
  figma-kit qa touch-targets --page 0`,
	}
	cmd.AddCommand(newQAContrastCmd())
	cmd.AddCommand(newQATouchTargetsCmd())
	cmd.AddCommand(newQAOrphansCmd())
	cmd.AddCommand(newQAFontsCmd())
	cmd.AddCommand(newQAColorsCmd())
	cmd.AddCommand(newQASpacingCmd())
	cmd.AddCommand(newQANamingCmd())
	cmd.AddCommand(newQAResponsiveCmd())
	cmd.AddCommand(newQAChecklistCmd())
	return cmd
}

func newQAContrastCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contrast",
		Short: "Traverse text nodes and flag potential contrast issues (heuristic)",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("function lum(c){ const r = c.r <= 0.03928 ? c.r/12.92 : Math.pow((c.r+0.055)/1.055,2.4);")
			b.Line("  const g = c.g <= 0.03928 ? c.g/12.92 : Math.pow((c.g+0.055)/1.055,2.4);")
			b.Line("  const bl = c.b <= 0.03928 ? c.b/12.92 : Math.pow((c.b+0.055)/1.055,2.4); return 0.2126*r + 0.7152*g + 0.0722*bl; }")
			b.Line("function ratio(fg,bg){ const L1 = lum(fg) + 0.05; const L2 = lum(bg) + 0.05; return L1 > L2 ? L1/L2 : L2/L1; }")
			b.Line("const issues = [];")
			b.Line("function walk(n){")
			b.Line("  if (n.type === 'TEXT' && n.fills && n.fills[0] && n.fills[0].type === 'SOLID') {")
			b.Line("    const fg = n.fills[0].color; const bg = {r:1,g:1,b:1}; const r = ratio(fg, bg);")
			b.Line("    if (r < 4.5) issues.push({ id: n.id, name: n.name, approxRatio: r, note: 'Assumes white page bg' }); }")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { lowContrastText: issues };")
			output(b.String())
			return nil
		},
	}
}

func newQATouchTargetsCmd() *cobra.Command {
	var min int
	cmd := &cobra.Command{
		Use:   "touch-targets",
		Short: "Flag frames/components smaller than minimum touch target size",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Linef("const minSize = %d;", min)
			b.Line("const bad = [];")
			b.Line("function walk(n){")
			b.Line("  if ((n.type === 'FRAME' || n.type === 'COMPONENT' || n.type === 'INSTANCE') && 'width' in n) {")
			b.Line("    if (n.width < minSize || n.height < minSize) bad.push({id:n.id,name:n.name,w:n.width,h:n.height}); }")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { undersized: bad };")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&min, "min", 44, "Minimum width/height in px")
	return cmd
}

func newQAOrphansCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "orphans",
		Short: "List top-level nodes on the page that look like stray layers",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const orphans = [];")
			b.Line("for (const c of figma.currentPage.children) {")
			b.Line("  if (c.visible && c.name && !c.name.startsWith('.')) orphans.push({ id: c.id, name: c.name, type: c.type }); }")
			b.Line("return { topLevel: orphans };")
			output(b.String())
			return nil
		},
	}
}

func newQAFontsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "fonts",
		Short: "Collect font combinations used by text nodes",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const map = {};")
			b.Line("function walk(n){")
			b.Line("  if (n.type === 'TEXT') { const k = n.fontName.family + ' / ' + n.fontName.style + ' @' + n.fontSize; map[k] = (map[k]||0)+1; }")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { fontUsage: map };")
			output(b.String())
			return nil
		},
	}
}

func newQAColorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "colors",
		Short: "Aggregate solid fill colors used in the subtree",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const map = {};")
			b.Line("function key(c){ return c.r.toFixed(3)+','+c.g.toFixed(3)+','+c.b.toFixed(3); }")
			b.Line("function walk(n){")
			b.Line("  if ('fills' in n && n.fills) for (const f of n.fills) {")
			b.Line("    if (f.type === 'SOLID' && f.color) { const k = key(f.color); map[k] = (map[k]||0)+1; } }")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { colorHistogram: map };")
			output(b.String())
			return nil
		},
	}
}

func newQASpacingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spacing",
		Short: "Flag auto-layout frames with itemSpacing below a threshold",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const tight = [];")
			b.Line("function walk(n){")
			b.Line("  if (n.type === 'FRAME' && n.layoutMode && n.layoutMode !== 'NONE' && n.itemSpacing < 8)")
			b.Line("    tight.push({ id: n.id, name: n.name, itemSpacing: n.itemSpacing });")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { tightAutoLayoutSpacing: tight };")
			output(b.String())
			return nil
		},
	}
}

func newQANamingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "naming",
		Short: "Flag default or empty layer names",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const bad = []; const defaults = /^(Frame|Group|Rectangle|Ellipse|Vector|Line|Text|Star|Polygon) \\d+$/i;")
			b.Line("function walk(n){")
			b.Line("  if (!n.name || !n.name.trim() || defaults.test(n.name)) bad.push({ id: n.id, type: n.type, name: n.name });")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { namingIssues: bad };")
			output(b.String())
			return nil
		},
	}
}

func newQAResponsiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "responsive",
		Short: "Summarize horizontal/vertical constraints usage",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const summary = { scaleH: 0, scaleV: 0, fixedH: 0, fixedV: 0, other: 0 };")
			b.Line("function walk(n){")
			b.Line("  if ('constraints' in n) {")
			b.Line("    const c = n.constraints;")
			b.Line("    if (c.horizontal === 'SCALE') summary.scaleH++; else if (c.horizontal === 'STRETCH' || c.horizontal === 'LEFT_RIGHT') summary.other++; else summary.fixedH++;")
			b.Line("    if (c.vertical === 'SCALE') summary.scaleV++; else if (c.vertical === 'STRETCH' || c.vertical === 'TOP_BOTTOM') summary.other++; else summary.fixedV++; }")
			b.Line("  if ('children' in n) for (const ch of n.children) walk(ch); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { constraints: summary };")
			output(b.String())
			return nil
		},
	}
}

func newQAChecklistCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checklist",
		Short: "Run a lightweight bundled QA pass and return combined findings",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := newBuilder()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const report = { emptyNames: [], tinyTargets: [], tightSpacing: [] };")
			b.Line("function walk(n){")
			b.Line("  if (!n.name || !String(n.name).trim()) report.emptyNames.push({id:n.id,type:n.type});")
			b.Line("  if ((n.type==='FRAME'||n.type==='COMPONENT') && 'width' in n && (n.width<40 || n.height<40))")
			b.Line("    report.tinyTargets.push({id:n.id,name:n.name,w:n.width,h:n.height});")
			b.Line("  if (n.type==='FRAME' && n.layoutMode && n.layoutMode!=='NONE' && n.itemSpacing<4)")
			b.Line("    report.tightSpacing.push({id:n.id,name:n.name,itemSpacing:n.itemSpacing});")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c); }")
			b.Line("walk(figma.currentPage);")
			b.Line("return { checklist: report };")
			output(b.String())
			return nil
		},
	}
}
