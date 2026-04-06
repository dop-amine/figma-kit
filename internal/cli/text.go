package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newTextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Text node operations (create, edit, style, range, fonts)",
	}
	cmd.AddCommand(newTextCreateCmd())
	cmd.AddCommand(newTextEditCmd())
	cmd.AddCommand(newTextStyleCmd())
	cmd.AddCommand(newTextRangeCmd())
	cmd.AddCommand(newTextListFontsCmd())
	cmd.AddCommand(newTextLoadFontsCmd())
	return cmd
}

func newTextCreateCmd() *cobra.Command {
	var (
		content string
		font    string
		weight  string
		size    int
		color   string
		parent  string
		x, y    int
		width   int
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a text node",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := codegen.HexToRGB(color)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("await figma.loadFontAsync({family:%q, style:%q});", font, weight)
			b.Line("const t = figma.createText();")
			b.Linef("t.fontName = {family:%q, style:%q};", font, weight)
			b.Linef("t.characters = %q;", content)
			b.Linef("t.fontSize = %d;", size)
			b.Linef("t.fills = [{type:'SOLID', color:%s}];", codegen.FormatRGB(c))
			b.Linef("t.x = %d;", x)
			b.Linef("t.y = %d;", y)
			if width > 0 {
				b.Linef("t.resize(%d, t.height);", width)
				b.Line("t.textAutoResize = 'HEIGHT';")
			}
			if parent != "" {
				b.Linef("const par = await figma.getNodeByIdAsync(%q);", parent)
				b.Line("if (par) par.appendChild(t); else figma.currentPage.appendChild(t);")
			} else {
				b.Line("figma.currentPage.appendChild(t);")
			}
			b.ReturnIDs("t.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&content, "content", "Text", "Text content")
	cmd.Flags().StringVar(&font, "font", "Inter", "Font family")
	cmd.Flags().StringVar(&weight, "weight", "Regular", "Font style/weight")
	cmd.Flags().IntVar(&size, "size", 16, "Font size")
	cmd.Flags().StringVar(&color, "color", "#FFFFFF", "Text color (hex)")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID")
	cmd.Flags().IntVar(&x, "x", 0, "X position")
	cmd.Flags().IntVar(&y, "y", 0, "Y position")
	cmd.Flags().IntVarP(&width, "width", "w", 0, "Text width (0 = auto)")
	return cmd
}

func newTextEditCmd() *cobra.Command {
	var content string
	cmd := &cobra.Command{
		Use:   "edit <nodeId>",
		Short: "Edit text content of a text node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const t = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!t || t.type !== 'TEXT') throw new Error('Text node not found');")
			b.Line("await figma.loadFontAsync(t.fontName);")
			b.Linef("t.characters = %q;", content)
			b.ReturnIDs("t.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "New text content")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newTextStyleCmd() *cobra.Command {
	var (
		size  int
		lh    int
		ls    int
		align string
	)
	cmd := &cobra.Command{
		Use:   "style <nodeId>",
		Short: "Change typography on a text node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const t = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!t || t.type !== 'TEXT') throw new Error('Text node not found');")
			b.Line("await figma.loadFontAsync(t.fontName);")
			if cmd.Flags().Changed("size") {
				b.Linef("t.fontSize = %d;", size)
			}
			if cmd.Flags().Changed("lh") {
				b.Linef("t.lineHeight = {value:%d, unit:'PIXELS'};", lh)
			}
			if cmd.Flags().Changed("ls") {
				b.Linef("t.letterSpacing = {value:%d, unit:'PERCENT'};", ls)
			}
			if cmd.Flags().Changed("align") {
				b.Linef("t.textAlignHorizontal = %q;", strings.ToUpper(align))
			}
			b.ReturnIDs("t.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&size, "size", 16, "Font size")
	cmd.Flags().IntVar(&lh, "lh", 0, "Line height (px)")
	cmd.Flags().IntVar(&ls, "ls", 0, "Letter spacing (%)")
	cmd.Flags().StringVar(&align, "align", "", "Horizontal alignment (LEFT, CENTER, RIGHT, JUSTIFIED)")
	return cmd
}

func newTextRangeCmd() *cobra.Command {
	var (
		start  int
		end    int
		weight string
		color  string
	)
	cmd := &cobra.Command{
		Use:   "range <nodeId>",
		Short: "Apply mixed styles to a text range",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const t = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!t || t.type !== 'TEXT') throw new Error('Text node not found');")
			if weight != "" {
				b.Linef("await figma.loadFontAsync({family:'Inter', style:%q});", weight)
				b.Linef("t.setRangeFontName(%d, %d, {family:'Inter', style:%q});", start, end, weight)
			}
			if color != "" {
				c, err := codegen.HexToRGB(color)
				if err != nil {
					return err
				}
				b.Linef("t.setRangeFills(%d, %d, [{type:'SOLID', color:%s}]);", start, end, codegen.FormatRGB(c))
			}
			b.ReturnIDs("t.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&start, "start", 0, "Range start index")
	cmd.Flags().IntVar(&end, "end", 0, "Range end index")
	cmd.Flags().StringVar(&weight, "weight", "", "Font weight for range")
	cmd.Flags().StringVar(&color, "color", "", "Color for range (hex)")
	return cmd
}

func newTextListFontsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-fonts",
		Short: "Generate JS to list all available fonts in the file",
		Run: func(cmd *cobra.Command, args []string) {
			b := codegen.New()
			b.Line("const fonts = await figma.listAvailableFontsAsync();")
			b.Line("const families = [...new Set(fonts.map(f => f.fontName.family))].sort();")
			b.ReturnExpr("{ count: families.length, families }")
			output(b.String())
		},
	}
}

func newTextLoadFontsCmd() *cobra.Command {
	var families string
	cmd := &cobra.Command{
		Use:   "load-fonts",
		Short: "Generate font loading code for specified families",
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			fams := strings.Split(families, ",")
			for _, f := range fams {
				f = strings.TrimSpace(f)
				if f == "" {
					continue
				}
				b.Linef("const %sFonts = await figma.listAvailableFontsAsync();", sanitizeVarName(f))
				b.Linef("const %sStyles = %sFonts.filter(fn => fn.fontName.family === %q).map(fn => fn.fontName.style);",
					sanitizeVarName(f), sanitizeVarName(f), f)
				b.Linef("for (const st of %sStyles) await figma.loadFontAsync({family:%q, style:st});",
					sanitizeVarName(f), f)
			}
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&families, "families", "Inter,Geist Mono", "Comma-separated font families")
	return cmd
}

func sanitizeVarName(s string) string {
	r := strings.NewReplacer(" ", "", "-", "", ".", "")
	return fmt.Sprintf("_%s", strings.ToLower(r.Replace(s)))
}
