package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

const (
	defaultCardW = 320
	defaultCardH = 200
	imageCardW   = 400
	imageCardH   = 240
	bentoCellW   = 160
	bentoCellH   = 120
)

func newCardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "card",
		Short: "Card primitives (glass, solid, gradient, image, bento grid)",
		Example: `  # "Create a glassmorphism feature card about AI integration"
  figma-kit card glass -t noir --title "AI-Native" --desc "Prompt in natural language"

  # "Add a solid card with pricing info"
  figma-kit card solid -t noir --title "Pro Plan" --desc "$29/mo"

  # "Create a gradient card for a case study"
  figma-kit card gradient -t noir --title "Nike React" --desc "E-commerce redesign"`,
	}
	cmd.AddCommand(newCardGlassCmd())
	cmd.AddCommand(newCardSolidCmd())
	cmd.AddCommand(newCardGradientCmd())
	cmd.AddCommand(newCardImageCmd())
	cmd.AddCommand(newCardBentoCmd())
	cmd.AddCommand(newCardNeumorphicCmd())
	cmd.AddCommand(newCardClayCmd())
	cmd.AddCommand(newCardOutlineCmd())
	return cmd
}

func newCardGlassCmd() *cobra.Command {
	var (
		preset string
		w, h   int
		parent string
		title  string
		desc   string
	)
	cmd := &cobra.Command{
		Use:   "glass",
		Short: "Glassmorphism card via G() helper (presets: subtle, default, strong, pill)",
		Example: `  # "Create three glass feature cards for my landing page"
  figma-kit card glass -t noir --title "Speed" --desc "Sub-millisecond reads"
  figma-kit card glass -t noir --title "Scale" --desc "From prototype to planet-scale"
  figma-kit card glass -t noir --title "Sync" --desc "Real-time everywhere" --preset strong`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			p := strings.ToLower(strings.TrimSpace(preset))
			optsJS, err := glassPresetOptsJS(p, w, h)
			if err != nil {
				return err
			}
			needFonts := title != "" || desc != ""
			page := resolvePage()
			b := newBuilder()
			if needFonts {
				t, err := resolveTheme(cmd)
				if err != nil {
					return err
				}
				codegen.PreambleWithPage(b, t, page)
			} else {
				b.PageSetup(page)
			}
			if !b.IsBodyOnly() {
				b.Raw(codegen.AllHelpers())
			}
			b.Line("let par = pg;")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _p = %s;", parent)
				} else {
					b.Linef("const _p = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Line("if (!_p || typeof _p.appendChild !== 'function') throw new Error('Invalid parent');")
				b.Line("par = _p;")
			}
			b.Linef("const card = G(par, 0, 0, %d, %d, %s);", w, h, optsJS)
			b.Line(`card.name = 'Glass card';`)
			if title != "" {
				b.Linef("T(card, %q, 24, 24, %d, 20, 'Semi Bold', WT);", title, w-48)
			}
			if desc != "" {
				y := 56
				if title == "" {
					y = 24
				}
				b.Linef("T(card, %q, 24, %d, %d, 14, 'Regular', MT, 22);", desc, y, w-48)
			}
			b.ReturnIDs("card.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&preset, "preset", "default", "Glass preset (subtle, default, strong, pill)")
	cmd.Flags().IntVarP(&w, "width", "w", 320, "Width")
	cmd.Flags().IntVar(&h, "height", 200, "Height")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID or _results[N] in compose (defaults to page)")
	cmd.Flags().StringVar(&title, "title", "", "Optional title text inside the card")
	cmd.Flags().StringVar(&desc, "desc", "", "Optional description text inside the card")
	return cmd
}

func glassPresetOptsJS(preset string, w, h int) (string, error) {
	switch preset {
	case "subtle":
		return "{ r: 20, f: 0.03, s: 0.06, ga: 0.04, bl: 16, glow: true }", nil
	case "strong":
		return "{ r: 24, f: 0.1, s: 0.14, ga: 0.12, bl: 40, glow: true }", nil
	case "pill":
		return fmt.Sprintf("{ r: Math.min(%d, %d) / 2, f: 0.04, s: 0.08, ga: 0.06, bl: 24, glow: true }", w, h), nil
	case "default":
		return "{ r: 20, f: 0.04, s: 0.08, ga: 0.06, bl: 24, glow: true }", nil
	default:
		return "", fmt.Errorf("invalid --preset %q (use subtle, default, strong, or pill)", preset)
	}
}

func newCardSolidCmd() *cobra.Command {
	var (
		bg     string
		border string
		shadow string
		radius int
		w, h   int
		title  string
		desc   string
	)
	cmd := &cobra.Command{
		Use:         "solid",
		Short:       "Solid fill card with optional border, shadow, and corner radius",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := codegen.HexToRGB(bg)
			if err != nil {
				return err
			}
			needFonts := title != "" || desc != ""
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			if needFonts {
				b.FontLoading()
			}
			b.Line("const card = figma.createFrame();")
			b.Line(`card.name = 'Solid card';`)
			b.Linef("card.resize(%d, %d);", w, h)
			b.Linef("card.cornerRadius = %d;", radius)
			b.Linef("card.fills = [{type:'SOLID', color:%s}];", codegen.FormatRGB(c))
			b.Line("card.clipsContent = true;")

			if border != "" {
				bc, err := codegen.HexToRGB(border)
				if err != nil {
					return err
				}
				b.Linef("card.strokes = [{type:'SOLID', color:%s}];", codegen.FormatRGB(bc))
				b.Line("card.strokeWeight = 1;")
			}

			if shadow != "" {
				fx, err := solidShadowEffectJS(shadow)
				if err != nil {
					return err
				}
				b.Linef("card.effects = [%s];", fx)
			}

			if title != "" {
				b.Linef("await figma.loadFontAsync({family:'Inter',style:'Semi Bold'});")
				b.Line("const t = figma.createText();")
				b.Line("t.fontName = {family:'Inter',style:'Semi Bold'};")
				b.Linef("t.characters = %q;", title)
				b.Line("t.fontSize = 20;")
				b.Line("t.fills = [{type:'SOLID', color:{r:0.96,g:0.97,b:0.98}}];")
				b.Line("t.x = 24; t.y = 24;")
				b.Line("t.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("card.appendChild(t);")
			}

			if desc != "" {
				b.Linef("await figma.loadFontAsync({family:'Inter',style:'Regular'});")
				b.Line("const d = figma.createText();")
				b.Line("d.fontName = {family:'Inter',style:'Regular'};")
				b.Linef("d.characters = %q;", desc)
				b.Line("d.fontSize = 14;")
				b.Line("d.lineHeight = {value:22,unit:'PIXELS'};")
				b.Line("d.fills = [{type:'SOLID', color:{r:0.45,g:0.48,b:0.55}}];")
				y := 56
				if title == "" {
					y = 24
				}
				b.Linef("d.x = 24; d.y = %d;", y)
				b.Linef("d.resize(%d, d.height); d.textAutoResize = 'HEIGHT';", w-48)
				b.Line("card.appendChild(d);")
			}

			b.Line("pg.appendChild(card);")
			b.ReturnIDs("card.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&bg, "bg", "#1A1D24", "Background color (hex)")
	cmd.Flags().StringVar(&border, "border", "", "Optional border color (hex)")
	cmd.Flags().StringVar(&shadow, "shadow", "", "Drop shadow size (sm, md, lg)")
	cmd.Flags().IntVar(&radius, "radius", 16, "Corner radius")
	cmd.Flags().IntVarP(&w, "width", "w", defaultCardW, "Card width")
	cmd.Flags().IntVar(&h, "height", defaultCardH, "Card height")
	cmd.Flags().StringVar(&title, "title", "", "Optional title text inside the card")
	cmd.Flags().StringVar(&desc, "desc", "", "Optional description text inside the card")
	return cmd
}

func solidShadowEffectJS(size string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(size)) {
	case "sm":
		return "{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.12}, offset:{x:0,y:2}, radius:6, spread:0, visible:true, blendMode:'NORMAL'}", nil
	case "md":
		return "{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.15}, offset:{x:0,y:4}, radius:12, spread:0, visible:true, blendMode:'NORMAL'}", nil
	case "lg":
		return "{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.22}, offset:{x:0,y:12}, radius:28, spread:0, visible:true, blendMode:'NORMAL'}", nil
	default:
		return "", fmt.Errorf("invalid --shadow %q (use sm, md, or lg)", size)
	}
}

func newCardGradientCmd() *cobra.Command {
	var (
		from  string
		to    string
		angle float64
	)
	cmd := &cobra.Command{
		Use:         "gradient",
		Short:       "Linear gradient fill card",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c0, err := codegen.HexToRGB(from)
			if err != nil {
				return err
			}
			c1, err := codegen.HexToRGB(to)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			b.Linef("const ang = (%s * Math.PI) / 180;", codegen.FmtFloat(angle))
			b.Line("const cos = Math.cos(ang), sin = Math.sin(ang);")
			b.Line("const gt = [[cos, sin, 0], [-sin, cos, 0]];")
			b.Line("const card = figma.createFrame();")
			b.Line(`card.name = 'Gradient card';`)
			b.Linef("card.resize(%d, %d);", defaultCardW, defaultCardH)
			b.Linef("card.cornerRadius = 16;")
			b.Line("card.fills = [{")
			b.Line("  type: 'GRADIENT_LINEAR',")
			b.Line("  gradientTransform: gt,")
			b.Linef("  gradientStops: [")
			b.Linef("    { position: 0, color: %s },", codegen.FormatRGBA(c0, 1))
			b.Linef("    { position: 1, color: %s },", codegen.FormatRGBA(c1, 1))
			b.Line("  ],")
			b.Line("}];")
			b.Line("pg.appendChild(card);")
			b.ReturnIDs("card.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "#3B5BFF", "Gradient start color (hex)")
	cmd.Flags().StringVar(&to, "to", "#14B8A6", "Gradient end color (hex)")
	cmd.Flags().Float64Var(&angle, "angle", 135, "Gradient angle in degrees")
	return cmd
}

func newCardImageCmd() *cobra.Command {
	var (
		url     string
		overlay string
		title   string
	)
	cmd := &cobra.Command{
		Use:         "image",
		Short:       "Image fill card with optional overlay and title",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ov := strings.ToLower(strings.TrimSpace(overlay))
			switch ov {
			case "", "dark", "light":
			default:
				return fmt.Errorf("invalid --overlay %q (use dark or light)", overlay)
			}

			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)

			if title != "" {
				b.FontLoading()
			}

			b.Linef("const res = await fetch(%q);", url)
			b.Line("if (!res.ok) throw new Error('Image fetch failed: ' + res.status);")
			b.Line("const buf = new Uint8Array(await res.arrayBuffer());")
			b.Line("const img = figma.createImage(buf);")
			b.Line("const card = figma.createFrame();")
			b.Line(`card.name = 'Image card';`)
			b.Linef("card.resize(%d, %d);", imageCardW, imageCardH)
			b.Line("card.cornerRadius = 12;")
			b.Line("card.clipsContent = true;")
			b.Line("card.fills = [{ type: 'IMAGE', imageHash: img.hash, scaleMode: 'FILL' }];")

			switch ov {
			case "dark":
				b.Line("const ov = figma.createRectangle();")
				b.Linef("ov.resize(%d, %d);", imageCardW, imageCardH)
				b.Line("ov.fills = [{ type: 'SOLID', color: { r: 0, g: 0, b: 0 }, opacity: 0.45 }];")
				b.Line("card.appendChild(ov);")
			case "light":
				b.Line("const ov = figma.createRectangle();")
				b.Linef("ov.resize(%d, %d);", imageCardW, imageCardH)
				b.Line("ov.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.35 }];")
				b.Line("card.appendChild(ov);")
			}

			if title != "" {
				b.Line("const t = figma.createText();")
				b.Line("t.fontName = { family: 'Inter', style: 'Semi Bold' };")
				b.Linef("t.characters = %q;", title)
				b.Line("t.fontSize = 20;")
				b.Line("t.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 } }];")
				b.Line("t.x = 16;")
				b.Line("t.y = 16;")
				b.Line("card.appendChild(t);")
			}

			b.Line("pg.appendChild(card);")
			b.ReturnIDs("card.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "Image URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringVar(&overlay, "overlay", "", "Overlay tint (dark, light)")
	cmd.Flags().StringVar(&title, "title", "", "Optional title text")
	return cmd
}

func newCardBentoCmd() *cobra.Command {
	var cols, rows, gap int
	cmd := &cobra.Command{
		Use:         "bento",
		Short:       "Grid of card frames (bento layout)",
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if cols < 1 || rows < 1 {
				return fmt.Errorf("cols and rows must be at least 1")
			}
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			b.Line("const ids = [];")
			b.Linef("const cols = %d, rows = %d, gap = %d;", cols, rows, gap)
			b.Linef("const cw = %d, ch = %d;", bentoCellW, bentoCellH)
			b.Line("for (let r = 0; r < rows; r++) {")
			b.Line("  for (let c = 0; c < cols; c++) {")
			b.Line("    const cell = figma.createFrame();")
			b.Line("    cell.name = 'Bento ' + r + '-' + c;")
			b.Line("    cell.resize(cw, ch);")
			b.Line("    cell.x = c * (cw + gap);")
			b.Line("    cell.y = r * (ch + gap);")
			b.Line("    cell.cornerRadius = 12;")
			b.Line("    cell.fills = [{ type: 'SOLID', color: { r: 0.14, g: 0.16, b: 0.22 } }];")
			b.Line("    cell.strokes = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.06 }];")
			b.Line("    cell.strokeWeight = 1;")
			b.Line("    pg.appendChild(cell);")
			b.Line("    ids.push(cell.id);")
			b.Line("  }")
			b.Line("}")
			b.ReturnExpr("{ createdNodeIds: ids }")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&cols, "cols", 3, "Column count")
	cmd.Flags().IntVar(&rows, "rows", 2, "Row count")
	cmd.Flags().IntVar(&gap, "gap", 16, "Gap between cells (px)")
	return cmd
}

func cond(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func newCardNeumorphicCmd() *cobra.Command {
	var (
		w, h  int
		title string
		desc  string
		depth string
		inset bool
	)
	cmd := &cobra.Command{
		Use:   "neumorphic",
		Short: "Neumorphic (soft UI) card with inset/outset shadow pair",
		Example: `  # "Create a soft UI card for settings"
  figma-kit card neumorphic -t light --title "Settings" --depth deep`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)

			var offset, blur, spread int
			switch depth {
			case "subtle":
				offset, blur, spread = 3, 6, 0
			case "deep":
				offset, blur, spread = 8, 20, -2
			default:
				offset, blur, spread = 5, 12, 0
			}

			_ = t
			b.Line("const frame = figma.createFrame();")
			b.Linef("frame.name = %q;", cond(title != "", title, "Neumorphic Card"))
			b.Linef("frame.resize(%d, %d);", w, h)
			b.Line("frame.cornerRadius = 20;")
			b.Line("frame.fills = [{type:'SOLID', color:{r:0.93,g:0.93,b:0.95}}];")
			b.Line("frame.layoutMode = 'VERTICAL';")
			b.Line("frame.paddingLeft = frame.paddingRight = 24;")
			b.Line("frame.paddingTop = frame.paddingBottom = 24;")
			b.Line("frame.itemSpacing = 12;")
			b.Line("frame.primaryAxisSizingMode = 'FIXED';")
			b.Line("frame.counterAxisSizingMode = 'FIXED';")

			if inset {
				b.Linef("frame.effects = [")
				b.Linef("  {type:'INNER_SHADOW', color:{r:1,g:1,b:1,a:0.7}, offset:{x:-%d,y:-%d}, radius:%d, spread:%d, visible:true, blendMode:'NORMAL'},", offset, offset, blur, spread)
				b.Linef("  {type:'INNER_SHADOW', color:{r:0,g:0,b:0,a:0.15}, offset:{x:%d,y:%d}, radius:%d, spread:%d, visible:true, blendMode:'NORMAL'}", offset, offset, blur, spread)
				b.Line("];")
			} else {
				b.Linef("frame.effects = [")
				b.Linef("  {type:'DROP_SHADOW', color:{r:1,g:1,b:1,a:0.7}, offset:{x:-%d,y:-%d}, radius:%d, spread:%d, visible:true, blendMode:'NORMAL'},", offset, offset, blur, spread)
				b.Linef("  {type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.15}, offset:{x:%d,y:%d}, radius:%d, spread:%d, visible:true, blendMode:'NORMAL'}", offset, offset, blur, spread)
				b.Line("];")
			}

			if title != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tx.fontName = {family:'Inter',style:'Semi Bold'}; tx.fontSize = 18; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:0.2,g:0.2,b:0.25}}]; frame.appendChild(tx); }", title)
			}
			if desc != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); tx.fontName = {family:'Inter',style:'Regular'}; tx.fontSize = 14; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:0.45,g:0.45,b:0.5}}]; frame.appendChild(tx); }", desc)
			}

			b.Line("pg.appendChild(frame);")
			b.ReturnIDs("frame.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&w, "width", defaultCardW, "Card width")
	cmd.Flags().IntVar(&h, "height", defaultCardH, "Card height")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&desc, "desc", "", "Card description")
	cmd.Flags().StringVar(&depth, "depth", "default", "Shadow depth: subtle, default, deep")
	cmd.Flags().BoolVar(&inset, "inset", false, "Use inset shadows (pressed state)")
	return cmd
}

func newCardClayCmd() *cobra.Command {
	var (
		w, h  int
		title string
		desc  string
		color string
	)
	cmd := &cobra.Command{
		Use:   "clay",
		Short: "Claymorphism card with puffy 3D appearance",
		Example: `  # "Create a playful clay card for onboarding"
  figma-kit card clay --title "Welcome" --color "#A78BFA"`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)

			rgb, cErr := codegen.HexToRGB(color)
			if cErr != nil {
				rgb = codegen.RGB{R: 0.65, G: 0.55, B: 0.98}
			}

			b.Line("const frame = figma.createFrame();")
			b.Linef("frame.name = %q;", cond(title != "", title, "Clay Card"))
			b.Linef("frame.resize(%d, %d);", w, h)
			b.Line("frame.cornerRadius = 24;")
			b.Linef("frame.fills = [{type:'GRADIENT_LINEAR', gradientStops:[")
			b.Linef("  {position:0,color:{r:%.3f,g:%.3f,b:%.3f,a:1}},", minF(rgb.R+0.12, 1), minF(rgb.G+0.12, 1), minF(rgb.B+0.12, 1))
			b.Linef("  {position:1,color:{r:%.3f,g:%.3f,b:%.3f,a:1}}", rgb.R, rgb.G, rgb.B)
			b.Line("], gradientTransform:[[1,0,0],[0,1,0]]}];")
			b.Line("frame.effects = [")
			b.Linef("  {type:'INNER_SHADOW', color:{r:1,g:1,b:1,a:0.3}, offset:{x:-2,y:-2}, radius:4, spread:0, visible:true, blendMode:'NORMAL'},")
			b.Linef("  {type:'DROP_SHADOW', color:{r:%.3f,g:%.3f,b:%.3f,a:0.4}, offset:{x:0,y:8}, radius:20, spread:-4, visible:true, blendMode:'NORMAL'}", maxF(rgb.R-0.2, 0), maxF(rgb.G-0.2, 0), maxF(rgb.B-0.2, 0))
			b.Line("];")
			b.Line("frame.layoutMode = 'VERTICAL';")
			b.Line("frame.paddingLeft = frame.paddingRight = 28;")
			b.Line("frame.paddingTop = frame.paddingBottom = 28;")
			b.Line("frame.itemSpacing = 12;")
			b.Line("frame.primaryAxisSizingMode = 'FIXED';")
			b.Line("frame.counterAxisSizingMode = 'FIXED';")

			if title != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Bold'}); tx.fontName = {family:'Inter',style:'Bold'}; tx.fontSize = 20; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:1,g:1,b:1}}]; frame.appendChild(tx); }", title)
			}
			if desc != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); tx.fontName = {family:'Inter',style:'Regular'}; tx.fontSize = 14; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:1,g:1,b:1,a:0.8}}]; frame.appendChild(tx); }", desc)
			}

			b.Line("pg.appendChild(frame);")
			b.ReturnIDs("frame.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&w, "width", defaultCardW, "Card width")
	cmd.Flags().IntVar(&h, "height", defaultCardH, "Card height")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&desc, "desc", "", "Card description")
	cmd.Flags().StringVar(&color, "color", "#A78BFA", "Base pastel color (hex)")
	return cmd
}

func newCardOutlineCmd() *cobra.Command {
	var (
		w, h       int
		title      string
		desc       string
		glowColor  string
		glowSpread int
	)
	cmd := &cobra.Command{
		Use:   "outline",
		Short: "Ghost/outline card with glow border (dark-mode SaaS staple)",
		Example: `  # "Create an outline feature card for dark mode"
  figma-kit card outline -t noir --title "API Access" --glow-color "#3B82F6"`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)

			glow, cErr := codegen.HexToRGB(glowColor)
			if cErr != nil {
				glow = codegen.RGB{R: 0.23, G: 0.51, B: 0.96}
			}

			b.Line("const frame = figma.createFrame();")
			b.Linef("frame.name = %q;", cond(title != "", title, "Outline Card"))
			b.Linef("frame.resize(%d, %d);", w, h)
			b.Line("frame.cornerRadius = 12;")
			b.Line("frame.fills = [{type:'SOLID', color:{r:0,g:0,b:0}, opacity:0}];")
			b.Linef("frame.strokes = [{type:'SOLID', color:{r:%.3f,g:%.3f,b:%.3f}}];", glow.R, glow.G, glow.B)
			b.Line("frame.strokeWeight = 1;")
			b.Linef("frame.effects = [{type:'DROP_SHADOW', color:{r:%.3f,g:%.3f,b:%.3f,a:0.3}, offset:{x:0,y:0}, radius:%d, spread:0, visible:true, blendMode:'NORMAL'}];", glow.R, glow.G, glow.B, glowSpread)
			b.Line("frame.layoutMode = 'VERTICAL';")
			b.Line("frame.paddingLeft = frame.paddingRight = 24;")
			b.Line("frame.paddingTop = frame.paddingBottom = 24;")
			b.Line("frame.itemSpacing = 12;")
			b.Line("frame.primaryAxisSizingMode = 'FIXED';")
			b.Line("frame.counterAxisSizingMode = 'FIXED';")

			if title != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tx.fontName = {family:'Inter',style:'Semi Bold'}; tx.fontSize = 18; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:0.95,g:0.95,b:0.97}}]; frame.appendChild(tx); }", title)
			}
			if desc != "" {
				b.Linef("{ const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); tx.fontName = {family:'Inter',style:'Regular'}; tx.fontSize = 14; tx.characters = %q; tx.fills = [{type:'SOLID',color:{r:0.65,g:0.65,b:0.7}}]; frame.appendChild(tx); }", desc)
			}

			b.Line("pg.appendChild(frame);")
			b.ReturnIDs("frame.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&w, "width", defaultCardW, "Card width")
	cmd.Flags().IntVar(&h, "height", defaultCardH, "Card height")
	cmd.Flags().StringVar(&title, "title", "", "Card title")
	cmd.Flags().StringVar(&desc, "desc", "", "Card description")
	cmd.Flags().StringVar(&glowColor, "glow-color", "#3B82F6", "Glow/border color (hex)")
	cmd.Flags().IntVar(&glowSpread, "glow-spread", 12, "Glow spread radius")
	return cmd
}
