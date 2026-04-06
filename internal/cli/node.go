package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Low-level Figma node operations (create, clone, delete, move, ...)",
		Example: `  # "Create a hero frame for my landing page"
  figma-kit node create frame --name "Hero" -w 1440 --height 800

  # "Clone that card and offset it to the right"
  figma-kit node clone <nodeId> --dx 340 --dy 0

  # "Move the CTA below the features section"
  figma-kit node move <nodeId> --x 0 --y 1200`,
	}

	cmd.AddCommand(newNodeCreateCmd())
	cmd.AddCommand(newNodeCloneCmd())
	cmd.AddCommand(newNodeDeleteCmd())
	cmd.AddCommand(newNodeMoveCmd())
	cmd.AddCommand(newNodeResizeCmd())
	cmd.AddCommand(newNodeRenameCmd())
	cmd.AddCommand(newNodeReparentCmd())
	cmd.AddCommand(newNodeLockCmd())
	cmd.AddCommand(newNodeVisibleCmd())
	cmd.AddCommand(newNodeOrderCmd())
	cmd.AddCommand(newNodeGroupCmd())
	cmd.AddCommand(newNodeUngroupCmd())
	cmd.AddCommand(newNodeComponentCmd())
	cmd.AddCommand(newNodeFlattenCmd())
	return cmd
}

func newNodeCreateCmd() *cobra.Command {
	var (
		name   string
		width  int
		height int
		x, y   int
	)
	cmd := &cobra.Command{
		Use:   "create <type>",
		Short: "Create a Figma node (frame, rect, text, ellipse, line, polygon, star, vector, component)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			nodeType := args[0]
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)

			figmaMethod := nodeTypeToMethod(nodeType)
			if figmaMethod == "" {
				return fmt.Errorf("unknown node type %q (valid: frame, rect, text, ellipse, line, polygon, star, vector, component)", nodeType)
			}

			if nodeType == "text" {
				b.Line("await figma.loadFontAsync({family:'Inter',style:'Regular'});")
			}

			b.Linef("const node = figma.%s();", figmaMethod)
			b.Linef("node.name = %q;", name)
			if nodeType != "text" && nodeType != "line" {
				b.Linef("node.resize(%d, %d);", width, height)
			}
			b.Linef("node.x = %d;", x)
			b.Linef("node.y = %d;", y)
			if nodeType == "text" {
				b.Line("node.characters = 'Text';")
			}
			b.Line("figma.currentPage.appendChild(node);")
			b.ReturnIDs("node.id")

			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "Untitled", "Node name")
	cmd.Flags().IntVarP(&width, "width", "w", 400, "Width")
	cmd.Flags().IntVar(&height, "height", 300, "Height")
	cmd.Flags().IntVar(&x, "x", 0, "X position")
	cmd.Flags().IntVar(&y, "y", 0, "Y position")
	return cmd
}

func newNodeCloneCmd() *cobra.Command {
	var dx, dy int
	cmd := &cobra.Command{
		Use:   "clone <nodeId>",
		Short: "Duplicate a node with optional offset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const src = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!src) throw new Error('Node not found: ' + " + fmt.Sprintf("%q", args[0]) + ");")
			b.Line("const clone = src.clone();")
			b.Linef("clone.x = src.x + %d;", dx)
			b.Linef("clone.y = src.y + %d;", dy)
			b.Line("src.parent.appendChild(clone);")
			b.ReturnIDs("clone.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&dx, "dx", 100, "X offset from original")
	cmd.Flags().IntVar(&dy, "dy", 0, "Y offset from original")
	return cmd
}

func newNodeDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <nodeId>",
		Short: "Remove a node from the canvas",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("node.remove();")
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
}

func newNodeMoveCmd() *cobra.Command {
	var x, y int
	cmd := &cobra.Command{
		Use:   "move <nodeId>",
		Short: "Reposition a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.x = %d;", x)
			b.Linef("node.y = %d;", y)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&x, "x", 0, "X position")
	cmd.Flags().IntVar(&y, "y", 0, "Y position")
	return cmd
}

func newNodeResizeCmd() *cobra.Command {
	var w, h int
	cmd := &cobra.Command{
		Use:   "resize <nodeId>",
		Short: "Resize a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.resize(%d, %d);", w, h)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVarP(&w, "width", "w", 400, "Width")
	cmd.Flags().IntVar(&h, "height", 300, "Height")
	return cmd
}

func newNodeRenameCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "rename <nodeId>",
		Short: "Rename a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.name = %q;", name)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "New name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newNodeReparentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reparent <nodeId> <parentId>",
		Short: "Move a node to a different parent",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Linef("const parent = await figma.getNodeByIdAsync(%q);", args[1])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!parent) throw new Error('Parent not found');")
			b.Line("parent.appendChild(node);")
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
}

func newNodeLockCmd() *cobra.Command {
	var unlock bool
	cmd := &cobra.Command{
		Use:   "lock <nodeId>",
		Short: "Lock or unlock a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.locked = %t;", !unlock)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().BoolVar(&unlock, "unlock", false, "Unlock instead of lock")
	return cmd
}

func newNodeVisibleCmd() *cobra.Command {
	var hide bool
	cmd := &cobra.Command{
		Use:   "visible <nodeId>",
		Short: "Toggle node visibility",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.visible = %t;", !hide)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().BoolVar(&hide, "hide", false, "Hide instead of show")
	return cmd
}

func newNodeOrderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "order <direction> <nodeId>",
		Short: "Change layer order: front, back, forward, backward",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := strings.ToLower(args[0])
			id := args[1]
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", id)
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const parent = node.parent;")
			b.Line("if (!parent || !('children' in parent)) throw new Error('Node has no parent with children');")
			b.Line("const idx = parent.children.indexOf(node);")
			switch dir {
			case "front":
				b.Line("parent.insertChild(parent.children.length - 1, node);")
			case "back":
				b.Line("parent.insertChild(0, node);")
			case "forward":
				b.Line("if (idx < parent.children.length - 1) parent.insertChild(idx + 1, node);")
			case "backward":
				b.Line("if (idx > 0) parent.insertChild(idx - 1, node);")
			default:
				return fmt.Errorf("direction must be front|back|forward|backward")
			}
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
}

func newNodeGroupCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "group <nodeId> [nodeId...]",
		Short: "Group two or more nodes together",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			nodeVars := make([]string, len(args))
			for i, id := range args {
				v := fmt.Sprintf("n%d", i)
				nodeVars[i] = v
				b.Linef("const %s = await figma.getNodeByIdAsync(%q);", v, id)
				b.Linef("if (!%s) throw new Error('Node not found: %s');", v, id)
			}
			b.Linef("const grp = figma.group([%s], n0.parent);", strings.Join(nodeVars, ", "))
			b.Linef("grp.name = %q;", name)
			b.ReturnIDs("grp.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "Group", "Group name")
	return cmd
}

func newNodeUngroupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ungroup <nodeId>",
		Short: "Ungroup a group node, returning children to the parent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const grp = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!grp) throw new Error('Node not found');")
			b.Line("if (grp.type !== 'GROUP') throw new Error('Node is not a group');")
			b.Line("const ids = grp.children.map(c => c.id);")
			b.Line("figma.ungroup(grp);")
			b.Line("return { done: true, ungroupedIds: ids };")
			output(b.String())
			return nil
		},
	}
}

func newNodeComponentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "component <nodeId>",
		Short: "Convert an existing node into a Figma component",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const comp = figma.createComponentFromNode(node);")
			b.ReturnIDs("comp.id")
			output(b.String())
			return nil
		},
	}
}

func newNodeFlattenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "flatten <nodeId>",
		Short: "Flatten a node subtree into a single vector",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const flat = figma.flatten([node]);")
			b.ReturnIDs("flat.id")
			output(b.String())
			return nil
		},
	}
}

func nodeTypeToMethod(t string) string {
	m := map[string]string{
		"frame":         "createFrame",
		"rect":          "createRectangle",
		"rectangle":     "createRectangle",
		"text":          "createText",
		"ellipse":       "createEllipse",
		"line":          "createLine",
		"polygon":       "createPolygon",
		"star":          "createStar",
		"vector":        "createVector",
		"component":     "createComponent",
		"component-set": "createComponentSet",
	}
	return m[strings.ToLower(t)]
}
