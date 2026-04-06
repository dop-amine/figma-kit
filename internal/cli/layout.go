package cli

import (
	"strings"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/spf13/cobra"
)

func newLayoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "layout",
		Short: "Layout operations (auto-layout, grid, constraints, sizing, align, distribute)",
	}
	cmd.AddCommand(newLayoutAutoCmd())
	cmd.AddCommand(newLayoutGridCmd())
	cmd.AddCommand(newLayoutConstraintsCmd())
	cmd.AddCommand(newLayoutSizingCmd())
	cmd.AddCommand(newLayoutAlignCmd())
	cmd.AddCommand(newLayoutDistributeCmd())
	return cmd
}

func newLayoutAutoCmd() *cobra.Command {
	var (
		dir   string
		gap   int
		pad   int
		align string
		wrap  bool
	)
	cmd := &cobra.Command{
		Use:   "auto <nodeId>",
		Short: "Set auto-layout on a frame",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.layoutMode = %q;", strings.ToUpper(dir))
			b.Linef("node.itemSpacing = %d;", gap)
			b.Linef("node.paddingLeft = %d; node.paddingRight = %d;", pad, pad)
			b.Linef("node.paddingTop = %d; node.paddingBottom = %d;", pad, pad)
			if align != "" {
				b.Linef("node.counterAxisAlignItems = %q;", strings.ToUpper(align))
			}
			if wrap {
				b.Line("node.layoutWrap = 'WRAP';")
			}
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "VERTICAL", "Direction (HORIZONTAL, VERTICAL)")
	cmd.Flags().IntVar(&gap, "gap", 16, "Item spacing")
	cmd.Flags().IntVar(&pad, "pad", 0, "Padding (uniform)")
	cmd.Flags().StringVar(&align, "align", "", "Counter-axis alignment (MIN, CENTER, MAX, BASELINE)")
	cmd.Flags().BoolVar(&wrap, "wrap", false, "Enable wrapping")
	return cmd
}

func newLayoutGridCmd() *cobra.Command {
	var (
		columns int
		gutter  int
		margin  int
	)
	cmd := &cobra.Command{
		Use:   "grid <nodeId>",
		Short: "Add layout grid to a frame",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.layoutGrids = [{pattern:'COLUMNS', count:%d, gutterSize:%d, offset:%d, alignment:'STRETCH'}];",
				columns, gutter, margin)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&columns, "columns", 12, "Number of columns")
	cmd.Flags().IntVar(&gutter, "gutter", 24, "Gutter size")
	cmd.Flags().IntVar(&margin, "margin", 80, "Margin offset")
	return cmd
}

func newLayoutConstraintsCmd() *cobra.Command {
	var h, v string
	cmd := &cobra.Command{
		Use:   "constraints <nodeId>",
		Short: "Set constraints on a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.constraints = {horizontal:%q, vertical:%q};", strings.ToUpper(h), strings.ToUpper(v))
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&h, "h", "MIN", "Horizontal constraint (MIN, CENTER, MAX, STRETCH, SCALE)")
	cmd.Flags().StringVar(&v, "v", "MIN", "Vertical constraint (MIN, CENTER, MAX, STRETCH, SCALE)")
	return cmd
}

func newLayoutSizingCmd() *cobra.Command {
	var w, h string
	cmd := &cobra.Command{
		Use:   "sizing <nodeId>",
		Short: "Set sizing behavior (FIXED, HUG, FILL)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			if w != "" {
				b.Linef("node.layoutSizingHorizontal = %q;", strings.ToUpper(w))
			}
			if h != "" {
				b.Linef("node.layoutSizingVertical = %q;", strings.ToUpper(h))
			}
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&w, "width", "w", "", "Horizontal sizing (FIXED, HUG, FILL)")
	cmd.Flags().StringVar(&h, "height", "", "Vertical sizing (FIXED, HUG, FILL)")
	return cmd
}

func newLayoutAlignCmd() *cobra.Command {
	var primary, counter string
	cmd := &cobra.Command{
		Use:   "align <nodeId>",
		Short: "Set alignment for auto-layout children",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			if primary != "" {
				b.Linef("node.primaryAxisAlignItems = %q;", strings.ToUpper(primary))
			}
			if counter != "" {
				b.Linef("node.counterAxisAlignItems = %q;", strings.ToUpper(counter))
			}
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&primary, "primary", "", "Primary axis alignment (MIN, CENTER, MAX, SPACE_BETWEEN)")
	cmd.Flags().StringVar(&counter, "counter", "", "Counter axis alignment (MIN, CENTER, MAX, BASELINE)")
	return cmd
}

func newLayoutDistributeCmd() *cobra.Command {
	var (
		axis string
		gap  int
	)
	cmd := &cobra.Command{
		Use:   "distribute <nodeIds>",
		Short: "Distribute nodes evenly along an axis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := strings.Split(args[0], ",")
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Line("const nodes = [];")
			for _, id := range ids {
				id = strings.TrimSpace(id)
				b.Linef("nodes.push(await figma.getNodeByIdAsync(%q));", id)
			}
			b.Line("nodes.sort((a,b) => a.x - b.x);")
			if strings.ToUpper(axis) == "V" || strings.ToUpper(axis) == "VERTICAL" {
				b.Line("nodes.sort((a,b) => a.y - b.y);")
				b.Line("for (let i = 1; i < nodes.length; i++) {")
				b.Linef("  nodes[i].y = nodes[i-1].y + nodes[i-1].height + %d;", gap)
				b.Line("}")
			} else {
				b.Line("for (let i = 1; i < nodes.length; i++) {")
				b.Linef("  nodes[i].x = nodes[i-1].x + nodes[i-1].width + %d;", gap)
				b.Line("}")
			}
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&axis, "axis", "H", "Distribution axis (H, V)")
	cmd.Flags().IntVar(&gap, "gap", 24, "Gap between nodes")
	return cmd
}
