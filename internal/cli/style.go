package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newStyleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "style",
		Short: "Node styling operations (fill, stroke, effect, corner, blend, gradient, clip)",
	}

	cmd.AddCommand(newStyleFillCmd())
	cmd.AddCommand(newStyleStrokeCmd())
	cmd.AddCommand(newStyleEffectCmd())
	cmd.AddCommand(newStyleCornerCmd())
	cmd.AddCommand(newStyleBlendCmd())
	cmd.AddCommand(newStyleGradientCmd())
	cmd.AddCommand(newStyleClipCmd())
	return cmd
}

func newStyleFillCmd() *cobra.Command {
	var (
		solid   string
		opacity float64
	)
	cmd := &cobra.Command{
		Use:   "fill <nodeId>",
		Short: "Set solid fill on a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := codegen.HexToRGB(solid)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.fills = [{type:'SOLID', color:%s, opacity:%s}];",
				codegen.FormatRGB(c), codegen.FmtFloat(opacity))
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&solid, "solid", "#FFFFFF", "Hex color")
	cmd.Flags().Float64Var(&opacity, "opacity", 1.0, "Fill opacity (0-1)")
	return cmd
}

func newStyleStrokeCmd() *cobra.Command {
	var (
		color  string
		weight int
		align  string
	)
	cmd := &cobra.Command{
		Use:   "stroke <nodeId>",
		Short: "Set stroke on a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := codegen.HexToRGB(color)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.strokes = [{type:'SOLID', color:%s}];", codegen.FormatRGB(c))
			b.Linef("node.strokeWeight = %d;", weight)
			b.Linef("node.strokeAlign = %q;", mapStrokeAlign(align))
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&color, "color", "#FFFFFF", "Stroke color (hex)")
	cmd.Flags().IntVar(&weight, "weight", 1, "Stroke weight")
	cmd.Flags().StringVar(&align, "align", "inside", "Stroke alignment (inside, outside, center)")
	return cmd
}

func newStyleEffectCmd() *cobra.Command {
	var (
		shadow   string
		blur     int
		blurType string
	)
	cmd := &cobra.Command{
		Use:   "effect <nodeId>",
		Short: "Apply effects (shadow, blur) to a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")

			if blur > 0 {
				bt := "LAYER_BLUR"
				if blurType == "background" {
					bt = "BACKGROUND_BLUR"
				}
				b.Linef("node.effects = [{type:'%s', radius:%d, visible:true}];", bt, blur)
			} else if shadow != "" {
				b.Linef("node.effects = [{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.15}, offset:{x:0,y:4}, radius:12, spread:0, visible:true, blendMode:'NORMAL'}];")
			}

			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&shadow, "shadow", "", "Shadow specification")
	cmd.Flags().IntVar(&blur, "blur", 0, "Blur radius")
	cmd.Flags().StringVar(&blurType, "blur-type", "layer", "Blur type (layer, background)")
	return cmd
}

func newStyleCornerCmd() *cobra.Command {
	var (
		radius         int
		tl, tr, br, bl int
	)
	cmd := &cobra.Command{
		Use:   "corner <nodeId>",
		Short: "Set corner radius on a node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")

			if cmd.Flags().Changed("tl") || cmd.Flags().Changed("tr") || cmd.Flags().Changed("br") || cmd.Flags().Changed("bl") {
				b.Linef("node.topLeftRadius = %d;", tl)
				b.Linef("node.topRightRadius = %d;", tr)
				b.Linef("node.bottomRightRadius = %d;", br)
				b.Linef("node.bottomLeftRadius = %d;", bl)
			} else {
				b.Linef("node.cornerRadius = %d;", radius)
			}

			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&radius, "radius", 0, "Uniform corner radius")
	cmd.Flags().IntVar(&tl, "tl", 0, "Top-left radius")
	cmd.Flags().IntVar(&tr, "tr", 0, "Top-right radius")
	cmd.Flags().IntVar(&br, "br", 0, "Bottom-right radius")
	cmd.Flags().IntVar(&bl, "bl", 0, "Bottom-left radius")
	return cmd
}

func newStyleBlendCmd() *cobra.Command {
	var (
		mode    string
		opacity float64
	)
	cmd := &cobra.Command{
		Use:   "blend <nodeId>",
		Short: "Set blend mode and opacity",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.blendMode = %q;", mode)
			b.Linef("node.opacity = %s;", codegen.FmtFloat(opacity))
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "NORMAL", "Blend mode (NORMAL, MULTIPLY, SCREEN, OVERLAY, ...)")
	cmd.Flags().Float64Var(&opacity, "opacity", 1.0, "Opacity (0-1)")
	return cmd
}

func newStyleGradientCmd() *cobra.Command {
	var (
		gradType string
		angle    int
		stops    string
	)
	cmd := &cobra.Command{
		Use:   "gradient <nodeId>",
		Short: "Apply a gradient fill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")

			figType := "GRADIENT_LINEAR"
			if gradType == "radial" {
				figType = "GRADIENT_RADIAL"
			}

			b.Linef("node.fills = [{type:'%s', gradientTransform:[[1,0,0],[0,1,0]], gradientStops:[", figType)
			parsedStops := parseGradientStops(stops)
			for i, s := range parsedStops {
				comma := ","
				if i == len(parsedStops)-1 {
					comma = ""
				}
				b.Linef("  {position:%s, color:%s}%s", codegen.FmtFloat(s.pos), codegen.FormatRGBA(s.color, 1), comma)
			}
			b.Line("]}];")

			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&gradType, "type", "linear", "Gradient type (linear, radial)")
	cmd.Flags().IntVar(&angle, "angle", 0, "Gradient angle in degrees")
	cmd.Flags().StringVar(&stops, "stops", "0:#000000,1:#FFFFFF", "Gradient stops (pos:color,...)")
	return cmd
}

func newStyleClipCmd() *cobra.Command {
	var off bool
	cmd := &cobra.Command{
		Use:   "clip <nodeId>",
		Short: "Toggle clip content on a frame",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Linef("node.clipsContent = %t;", !off)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().BoolVar(&off, "off", false, "Disable clipping")
	return cmd
}

func mapStrokeAlign(a string) string {
	switch a {
	case "outside":
		return "OUTSIDE"
	case "center":
		return "CENTER"
	default:
		return "INSIDE"
	}
}

type gradientStop struct {
	pos   float64
	color codegen.RGB
}

func parseGradientStops(s string) []gradientStop {
	var stops []gradientStop
	for _, part := range splitCommaIgnoreColon(s) {
		var pos float64
		var hex string
		if _, err := fmt.Sscanf(part, "%f:%s", &pos, &hex); err == nil {
			if c, err := codegen.HexToRGB(hex); err == nil {
				stops = append(stops, gradientStop{pos: pos, color: c})
			}
		}
	}
	if len(stops) == 0 {
		stops = []gradientStop{
			{pos: 0, color: codegen.RGB{R: 0, G: 0, B: 0}},
			{pos: 1, color: codegen.RGB{R: 1, G: 1, B: 1}},
		}
	}
	return stops
}

func splitCommaIgnoreColon(s string) []string {
	var parts []string
	current := ""
	for _, ch := range s {
		if ch == ',' {
			if len(current) > 0 && current[len(current)-1] != ':' {
				parts = append(parts, current)
				current = ""
				continue
			}
		}
		current += string(ch)
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
