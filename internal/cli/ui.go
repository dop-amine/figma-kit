package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/theme"
)

func newUICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ui",
		Short: "UI primitive components (button, input, badge, ...)",
		Example: `  # "Add a primary CTA button"
  figma-kit ui button --variant primary -t noir

  # "Create an email input field"
  figma-kit ui input -t noir

  # "Add a navigation bar and footer"
  figma-kit ui nav -t noir
  figma-kit ui footer -t noir`,
	}
	cmd.AddCommand(newUIButtonCmd())
	cmd.AddCommand(newUIInputCmd())
	cmd.AddCommand(newUIBadgeCmd())
	cmd.AddCommand(newUIAvatarCmd())
	cmd.AddCommand(newUIDividerCmd())
	cmd.AddCommand(newUIIconCmd())
	cmd.AddCommand(newUIProgressCmd())
	cmd.AddCommand(newUIToggleCmd())
	cmd.AddCommand(newUITooltipCmd())
	cmd.AddCommand(newUIStatCmd())
	cmd.AddCommand(newUITableCmd())
	cmd.AddCommand(newUINavCmd())
	cmd.AddCommand(newUIFooterCmd())
	cmd.AddCommand(newUICheckboxCmd())
	cmd.AddCommand(newUIRadioCmd())
	cmd.AddCommand(newUITabsCmd())
	cmd.AddCommand(newUIDropdownCmd())
	cmd.AddCommand(newUIBreadcrumbCmd())
	cmd.AddCommand(newUISkeletonCmd())
	cmd.AddCommand(newUIHeroCmd())
	cmd.AddCommand(newUIPricingCmd())
	cmd.AddCommand(newUIFeatureGridCmd())
	cmd.AddCommand(newUITestimonialCmd())
	cmd.AddCommand(newUITimelineCmd())
	cmd.AddCommand(newUIStepperCmd())
	cmd.AddCommand(newUIAccordionCmd())
	cmd.AddCommand(newUIChipCmd())
	cmd.AddCommand(newUIToastCmd())
	cmd.AddCommand(newUIModalCmd())
	cmd.AddCommand(newUICardListCmd())
	cmd.AddCommand(newUISidebarCmd())
	cmd.AddCommand(newUIAvatarGroupCmd())
	cmd.AddCommand(newUIRatingCmd())
	cmd.AddCommand(newUISearchCmd())
	cmd.AddCommand(newUIPaginationCmd())
	cmd.AddCommand(newUIColorPickerCmd())
	return cmd
}

// emitUIThemeTokens writes sorted theme color constants and common type-scale tokens as JS consts.
func emitUIThemeTokens(b *codegen.Builder, t *theme.Theme) {
	names := make([]string, 0, len(t.Colors))
	for k := range t.Colors {
		names = append(names, k)
	}
	sort.Strings(names)
	b.Comment("Theme colors")
	for _, name := range names {
		c := t.Colors[name]
		b.Linef("const %s=%s;", name, codegen.FormatRGB(codegen.RGB{R: c.R, G: c.G, B: c.B}))
	}
	b.Blank()
	b.Comment("Type scale")
	emitType := func(key, suffix string, defSize int, defStyle string) {
		fs, st := defSize, defStyle
		if spec, ok := t.Type[key]; ok {
			fs = spec.FontSize
			if spec.Style != "" {
				st = spec.Style
			}
		}
		b.Linef("const TY_%s=%d; const ST_%s=%q;", suffix, fs, suffix, st)
	}
	emitType("body", "BODY", 16, "Regular")
	emitType("small", "SMALL", 13, "Regular")
	emitType("label", "LABEL", 11, "Medium")
	emitType("h4", "H4", 22, "Semi Bold")
	b.Blank()
}

func uiPreamble(b *codegen.Builder, t *theme.Theme, page int) {
	b.PageSetup(page)
	b.FontLoading()
	emitUIThemeTokens(b, t)
}

func newUIButtonCmd() *cobra.Command {
	var variant, label, size string
	cmd := &cobra.Command{
		Use:   "button",
		Short: "Themed button frame (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			v := strings.ToLower(variant)
			if v != "primary" && v != "secondary" && v != "ghost" && v != "destructive" && v != "outline" {
				return fmt.Errorf("variant must be primary|secondary|ghost|destructive|outline")
			}
			sz := strings.ToLower(size)
			if sz != "sm" && sz != "md" && sz != "lg" {
				return fmt.Errorf("size must be sm|md|lg")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Button")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Button';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.primaryAxisAlignItems = 'CENTER';")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			switch sz {
			case "sm":
				b.Line("root.paddingLeft = 10; root.paddingRight = 10; root.paddingTop = 6; root.paddingBottom = 6; root.itemSpacing = 6; root.cornerRadius = 6;")
				b.Line("const fs = TY_SMALL; const fw = ST_SMALL;")
			case "lg":
				b.Line("root.paddingLeft = 20; root.paddingRight = 20; root.paddingTop = 12; root.paddingBottom = 12; root.itemSpacing = 8; root.cornerRadius = 10;")
				b.Line("const fs = TY_H4; const fw = ST_H4;")
			default:
				b.Line("root.paddingLeft = 14; root.paddingRight = 14; root.paddingTop = 9; root.paddingBottom = 9; root.itemSpacing = 8; root.cornerRadius = 8;")
				b.Line("const fs = TY_BODY; const fw = ST_BODY;")
			}
			b.Line("root.strokes = [];")
			switch v {
			case "primary":
				b.Line("root.fills = [{type:'SOLID', color: BL}];")
				b.Linef("const tc = WT;")
			case "secondary":
				b.Line("root.fills = [{type:'SOLID', color: CARD}];")
				b.Line("root.strokes = [{type:'SOLID', color: STK}]; root.strokeWeight = 1;")
				b.Line("const tc = WT;")
			case "ghost":
				b.Line("root.fills = [];")
				b.Line("const tc = WT;")
			case "destructive":
				b.Line("root.fills = [{type:'SOLID', color: ERR}];")
				b.Line("const tc = WT;")
			case "outline":
				b.Line("root.fills = [];")
				b.Line("root.strokes = [{type:'SOLID', color: BL}]; root.strokeWeight = 1;")
				b.Line("const tc = BL;")
			}
			b.Line("const txt = figma.createText();")
			b.Line("txt.fontName = { family: 'Inter', style: fw };")
			b.Linef("txt.characters = %q;", label)
			b.Line("txt.fontSize = fs;")
			b.Line("txt.fills = [{ type: 'SOLID', color: tc }];")
			b.Line("txt.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(txt);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&variant, "variant", "primary", "primary|secondary|ghost|destructive|outline")
	cmd.Flags().StringVar(&label, "label", "Button", "Button label")
	cmd.Flags().StringVar(&size, "size", "md", "sm|md|lg")
	return cmd
}

func newUIInputCmd() *cobra.Command {
	var label, placeholder, typ string
	cmd := &cobra.Command{
		Use:   "input",
		Short: "Labeled input field (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			tp := strings.ToLower(typ)
			if tp != "text" && tp != "email" && tp != "password" {
				return fmt.Errorf("type must be text|email|password")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Input")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Input';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisAlignItems = 'MIN';")
			b.Line("root.counterAxisAlignItems = 'STRETCH';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.itemSpacing = 6;")
			b.Line("root.fills = [];")
			b.Line("const lab = figma.createText();")
			b.Line("lab.fontName = { family: 'Inter', style: ST_LABEL };")
			b.Linef("lab.characters = %q;", label)
			b.Line("lab.fontSize = TY_LABEL;")
			b.Line("lab.fills = [{ type: 'SOLID', color: MT }];")
			b.Line("lab.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(lab);")
			b.Line("const field = figma.createFrame();")
			b.Line("field.name = 'Field';")
			b.Line("field.layoutMode = 'HORIZONTAL';")
			b.Line("field.primaryAxisAlignItems = 'MIN';")
			b.Line("field.counterAxisAlignItems = 'CENTER';")
			b.Line("field.paddingLeft = 12; field.paddingRight = 12; field.paddingTop = 10; field.paddingBottom = 10;")
			b.Line("field.itemSpacing = 8;")
			b.Line("field.cornerRadius = 8;")
			b.Line("field.layoutAlign = 'STRETCH';")
			b.Line("field.primaryAxisSizingMode = 'FIXED';")
			b.Line("field.counterAxisSizingMode = 'FIXED';")
			b.Line("field.resize(280, 44);")
			b.Line("field.fills = [{type:'SOLID', color: CARD}];")
			b.Line("field.strokes = [{type:'SOLID', color: STK}]; field.strokeWeight = 1;")
			b.Line("const ph = figma.createText();")
			b.Line("ph.fontName = { family: 'Inter', style: ST_SMALL };")
			b.Linef("ph.characters = %q;", placeholder)
			b.Line("ph.fontSize = TY_SMALL;")
			b.Line("ph.fills = [{ type: 'SOLID', color: MT }];")
			b.Line("ph.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("field.appendChild(ph);")
			b.Linef("const meta = figma.createText(); meta.fontName = { family: 'Geist Mono', style: 'Regular' }; meta.characters = '[%s]'; meta.fontSize = 10; meta.fills = [{ type: 'SOLID', color: MT }]; meta.textAutoResize = 'WIDTH_AND_HEIGHT'; meta.letterSpacing = { value: 8, unit: 'PERCENT' }; field.appendChild(meta);", tp)
			b.Line("root.appendChild(field);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&label, "label", "Email", "Field label")
	cmd.Flags().StringVar(&placeholder, "placeholder", "you@example.com", "Placeholder text")
	cmd.Flags().StringVar(&typ, "type", "text", "text|email|password")
	return cmd
}

func badgeColorConst(name string) (string, error) {
	switch strings.ToLower(name) {
	case "blue":
		return "BL", nil
	case "green":
		return "SUCCESS", nil
	case "red":
		return "ERR", nil
	case "yellow":
		return "WARN", nil
	case "gray":
		return "MT", nil
	default:
		return "", fmt.Errorf("color must be blue|green|red|yellow|gray")
	}
}

func newUIBadgeCmd() *cobra.Command {
	var text, color string
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Pill badge (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			cc, err := badgeColorConst(color)
			if err != nil {
				return err
			}
			ink := "WT"
			if strings.EqualFold(color, "yellow") {
				ink = "BG"
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Badge")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Badge';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.primaryAxisAlignItems = 'CENTER';")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.paddingLeft = 10; root.paddingRight = 10; root.paddingTop = 4; root.paddingBottom = 4;")
			b.Line("root.itemSpacing = 4;")
			b.Line("root.cornerRadius = 999;")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Linef("root.fills = [{type:'SOLID', color: %s}];", cc)
			b.Line("const txt = figma.createText();")
			b.Line("txt.fontName = { family: 'Inter', style: ST_SMALL };")
			b.Linef("txt.characters = %q;", text)
			b.Line("txt.fontSize = TY_SMALL;")
			b.Linef("txt.fills = [{ type: 'SOLID', color: %s }];", ink)
			b.Line("txt.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(txt);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&text, "text", "New", "Badge text")
	cmd.Flags().StringVar(&color, "color", "blue", "blue|green|red|yellow|gray")
	return cmd
}

func newUIAvatarCmd() *cobra.Command {
	var initials string
	var size int
	cmd := &cobra.Command{
		Use:   "avatar",
		Short: "Circular avatar with initials",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			if size < 16 {
				return fmt.Errorf("size must be at least 16")
			}
			initials = strings.TrimSpace(strings.ToUpper(initials))
			if initials == "" {
				initials = "??"
			}
			if len(initials) > 3 {
				initials = initials[:3]
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Avatar")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Avatar';")
			b.Line("root.layoutMode = 'NONE';")
			b.Line("root.clipsContent = true;")
			b.Linef("root.resize(%d, %d);", size, size)
			b.Line("root.cornerRadius = " + strconv.Itoa(size) + " / 2;")
			b.Line("root.fills = [{type:'SOLID', color: STK}];")
			b.Line("root.strokes = [{type:'SOLID', color: BL, opacity: 0.35}]; root.strokeWeight = 1;")
			b.Line("const txt = figma.createText();")
			b.Line("txt.fontName = { family: 'Inter', style: 'Semi Bold' };")
			b.Linef("txt.characters = %q;", initials)
			b.Linef("txt.fontSize = Math.max(10, Math.floor(%d * 0.36));", size)
			b.Line("txt.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("txt.textAlignHorizontal = 'CENTER';")
			b.Line("txt.textAlignVertical = 'CENTER';")
			b.Linef("txt.resize(%d, %d);", size, size)
			b.Line("txt.x = 0; txt.y = 0;")
			b.Line("root.appendChild(txt);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&initials, "initials", "AK", "Initials (1–3 chars)")
	cmd.Flags().IntVar(&size, "size", 40, "Diameter in px")
	return cmd
}

func newUIDividerCmd() *cobra.Command {
	var dir, colorMode string
	var length int
	cmd := &cobra.Command{
		Use:   "divider",
		Short: "Horizontal or vertical divider line",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			d := strings.ToUpper(strings.TrimSpace(dir))
			if d != "H" && d != "V" {
				return fmt.Errorf("dir must be H or V")
			}
			if length < 1 {
				return fmt.Errorf("length must be positive")
			}
			useMuted := strings.EqualFold(colorMode, "muted")
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Divider")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Divider';")
			b.Line("root.layoutMode = 'NONE';")
			b.Line("root.fills = [];")
			var w, h int
			if d == "H" {
				w, h = length, 1
			} else {
				w, h = 1, length
			}
			b.Linef("root.resize(%d, %d);", w, h)
			b.Line("const bar = figma.createRectangle();")
			if useMuted {
				b.Line("bar.fills = [{type:'SOLID', color: MT, opacity: 0.45}];")
			} else {
				b.Line("bar.fills = [{type:'SOLID', color: STK}];")
			}
			b.Linef("bar.resize(%d, %d); bar.x = 0; bar.y = 0;", w, h)
			b.Line("root.appendChild(bar);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "H", "H (horizontal) or V (vertical)")
	cmd.Flags().IntVar(&length, "length", 240, "Line length in px")
	cmd.Flags().StringVar(&colorMode, "color", "", "omit for STK, or \"muted\" for MT @ 45% opacity")
	return cmd
}

func newUIIconCmd() *cobra.Command {
	var shape, hex string
	var size int
	cmd := &cobra.Command{
		Use:   "icon",
		Short: "Icon placeholder frame (square or circle)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			sh := strings.ToLower(shape)
			if sh != "circle" && sh != "square" {
				return fmt.Errorf("shape must be circle|square")
			}
			if size < 8 {
				return fmt.Errorf("size must be at least 8")
			}
			c, err := codegen.HexToRGB(hex)
			if err != nil {
				return err
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Icon")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Icon';")
			b.Line("root.layoutMode = 'NONE';")
			b.Line("root.clipsContent = " + func() string {
				if sh == "circle" {
					return "true"
				}
				return "false"
			}() + ";")
			b.Linef("root.resize(%d, %d);", size, size)
			if sh == "circle" {
				b.Linef("root.cornerRadius = %d / 2;", size)
			} else {
				b.Line("root.cornerRadius = 6;")
			}
			b.Line("root.fills = [{type:'SOLID', color: CARD}];")
			b.Line("root.strokes = [{type:'SOLID', color: STK}]; root.strokeWeight = 1;")
			b.Line("const glyph = figma.createRectangle();")
			b.Linef("glyph.fills = [{type:'SOLID', color: %s}];", codegen.FormatRGB(c))
			inner := size * 5 / 16
			if inner < 4 {
				inner = 4
			}
			off := (size - inner) / 2
			b.Linef("glyph.resize(%d, %d); glyph.cornerRadius = %d;", inner, inner, inner/4)
			b.Linef("glyph.x = %d; glyph.y = %d;", off, off)
			b.Line("root.appendChild(glyph);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&shape, "shape", "square", "circle|square")
	cmd.Flags().IntVar(&size, "size", 32, "Outer size in px")
	cmd.Flags().StringVar(&hex, "color", "#94A3B8", "Fill hex color")
	return cmd
}

func newUIProgressCmd() *cobra.Command {
	var value, width int
	cmd := &cobra.Command{
		Use:   "progress",
		Short: "Progress bar (themed track + fill)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			if width < 40 {
				return fmt.Errorf("width must be at least 40")
			}
			if value < 0 {
				value = 0
			}
			if value > 100 {
				value = 100
			}
			h := 8
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Progress")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Progress';")
			b.Line("root.layoutMode = 'NONE';")
			b.Linef("root.resize(%d, %d);", width, h+4)
			b.Line("root.fills = [];")
			b.Line("const track = figma.createFrame();")
			b.Line("track.name = 'Track';")
			b.Line("track.layoutMode = 'NONE';")
			b.Linef("track.resize(%d, %d); track.x = 0; track.y = 2;", width, h)
			b.Line("track.cornerRadius = " + strconv.Itoa(h/2) + ";")
			b.Line("track.fills = [{type:'SOLID', color: STK}];")
			b.Line("const fill = figma.createRectangle();")
			b.Line("fill.fills = [{type:'SOLID', color: BL}];")
			b.Linef("const fw = Math.max(4, Math.round(%d * (%d / 100)));", width-4, value)
			b.Linef("fill.resize(fw, %d); fill.cornerRadius = %d; fill.x = 2; fill.y = 2;", h-4, (h-4)/2)
			b.Line("track.appendChild(fill);")
			b.Line("root.appendChild(track);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&value, "value", 42, "Progress 0–100")
	cmd.Flags().IntVar(&width, "width", 200, "Track width in px")
	return cmd
}

func newUIToggleCmd() *cobra.Command {
	var state, size string
	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "On/off switch (auto-layout track + thumb)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			st := strings.ToLower(state)
			if st != "on" && st != "off" {
				return fmt.Errorf("state must be on|off")
			}
			sz := strings.ToLower(size)
			if sz != "sm" && sz != "md" && sz != "lg" {
				return fmt.Errorf("size must be sm|md|lg")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Toggle")
			var tw, th, pad, dia int
			switch sz {
			case "sm":
				tw, th, pad, dia = 36, 20, 2, 16
			case "lg":
				tw, th, pad, dia = 52, 28, 3, 22
			default:
				tw, th, pad, dia = 44, 24, 2, 20
			}
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Toggle';")
			b.Line("root.layoutMode = 'NONE';")
			b.Linef("root.resize(%d, %d);", tw, th)
			b.Line("root.fills = [];")
			b.Line("const track = figma.createFrame();")
			b.Line("track.layoutMode = 'NONE';")
			b.Linef("track.resize(%d, %d); track.x = 0; track.y = 0;", tw, th)
			b.Linef("track.cornerRadius = %d;", th/2)
			if st == "on" {
				b.Line("track.fills = [{type:'SOLID', color: TL, opacity: 0.35}];")
				b.Line("track.strokes = [{type:'SOLID', color: TL}]; track.strokeWeight = 1;")
			} else {
				b.Line("track.fills = [{type:'SOLID', color: STK}];")
				b.Line("track.strokes = [];")
			}
			b.Line("const thumb = figma.createEllipse();")
			b.Linef("thumb.resize(%d, %d);", dia, dia)
			offOn := tw - pad - dia
			offOff := pad
			if st == "on" {
				b.Linef("thumb.x = %d; thumb.y = %d;", offOn, (th-dia)/2)
				b.Line("thumb.fills = [{type:'SOLID', color: TL}];")
			} else {
				b.Linef("thumb.x = %d; thumb.y = %d;", offOff, (th-dia)/2)
				b.Line("thumb.fills = [{type:'SOLID', color: MT}];")
			}
			b.Line("track.appendChild(thumb);")
			b.Line("root.appendChild(track);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&state, "state", "off", "on|off")
	cmd.Flags().StringVar(&size, "size", "md", "sm|md|lg")
	return cmd
}

func newUITooltipCmd() *cobra.Command {
	var text, position string
	cmd := &cobra.Command{
		Use:   "tooltip",
		Short: "Tooltip bubble (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			pos := strings.ToLower(position)
			if pos != "top" && pos != "bottom" && pos != "left" && pos != "right" {
				return fmt.Errorf("position must be top|bottom|left|right")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Tooltip")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Tooltip';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisAlignItems = 'CENTER';")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.itemSpacing = 0;")
			b.Line("root.fills = [];")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			// Outer offset simulates caret side
			switch pos {
			case "top":
				b.Line("root.paddingBottom = 6;")
			case "bottom":
				b.Line("root.paddingTop = 6;")
			case "left":
				b.Line("root.paddingRight = 6;")
			case "right":
				b.Line("root.paddingLeft = 6;")
			}
			b.Line("const bubble = figma.createFrame();")
			b.Line("bubble.name = 'Bubble';")
			b.Line("bubble.layoutMode = 'HORIZONTAL';")
			b.Line("bubble.primaryAxisAlignItems = 'CENTER';")
			b.Line("bubble.counterAxisAlignItems = 'CENTER';")
			b.Line("bubble.paddingLeft = 12; bubble.paddingRight = 12; bubble.paddingTop = 8; bubble.paddingBottom = 8;")
			b.Line("bubble.cornerRadius = 6;")
			b.Line("bubble.fills = [{type:'SOLID', color: CARD}];")
			b.Line("bubble.strokes = [{type:'SOLID', color: STK}]; bubble.strokeWeight = 1;")
			b.Line("bubble.effects = [{ type: 'DROP_SHADOW', color: { r: 0, g: 0, b: 0, a: 0.2 }, offset: { x: 0, y: 2 }, radius: 8, spread: 0, visible: true, blendMode: 'NORMAL' }];")
			b.Line("const txt = figma.createText();")
			b.Line("txt.fontName = { family: 'Inter', style: ST_SMALL };")
			b.Linef("txt.characters = %q;", text)
			b.Line("txt.fontSize = TY_SMALL;")
			b.Line("txt.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("txt.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("bubble.appendChild(txt);")
			b.Line("const caretHost = figma.createFrame();")
			b.Line("caretHost.name = 'Caret';")
			b.Line("caretHost.layoutMode = 'HORIZONTAL';")
			b.Line("caretHost.primaryAxisAlignItems = 'CENTER';")
			b.Line("caretHost.counterAxisAlignItems = 'CENTER';")
			b.Line("caretHost.primaryAxisSizingMode = 'FIXED';")
			b.Line("caretHost.counterAxisSizingMode = 'FIXED';")
			b.Line("caretHost.fills = [];")
			b.Line("caretHost.itemSpacing = 0;")
			switch pos {
			case "left", "right":
				b.Line("caretHost.resize(8, 28); caretHost.layoutAlign = 'STRETCH';")
			default:
				b.Line("caretHost.resize(32, 8); caretHost.layoutAlign = 'STRETCH';")
			}
			b.Line("const caret = figma.createPolygon();")
			b.Line("caret.pointCount = 3;")
			b.Line("caret.fills = [{type:'SOLID', color: CARD}];")
			b.Line("caret.strokes = [{type:'SOLID', color: STK}]; caret.strokeWeight = 1;")
			b.Line("caret.resize(10, 6);")
			switch pos {
			case "top":
				b.Line("caret.rotation = Math.PI;")
			case "bottom":
				b.Line("caret.rotation = 0;")
			case "left":
				b.Line("caret.rotation = Math.PI / 2;")
			case "right":
				b.Line("caret.rotation = -Math.PI / 2;")
			}
			b.Line("caretHost.appendChild(caret);")
			switch pos {
			case "top":
				b.Line("root.appendChild(caretHost);")
				b.Line("root.appendChild(bubble);")
			case "bottom":
				b.Line("root.appendChild(bubble);")
				b.Line("root.appendChild(caretHost);")
			case "left":
				b.Line("root.layoutMode = 'HORIZONTAL';")
				b.Line("root.primaryAxisAlignItems = 'CENTER';")
				b.Line("root.counterAxisAlignItems = 'CENTER';")
				b.Line("root.appendChild(caretHost);")
				b.Line("root.appendChild(bubble);")
			case "right":
				b.Line("root.layoutMode = 'HORIZONTAL';")
				b.Line("root.primaryAxisAlignItems = 'CENTER';")
				b.Line("root.counterAxisAlignItems = 'CENTER';")
				b.Line("root.appendChild(bubble);")
				b.Line("root.appendChild(caretHost);")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&text, "text", "Copied!", "Tooltip text")
	cmd.Flags().StringVar(&position, "position", "top", "top|bottom|left|right")
	return cmd
}

func newUIStatCmd() *cobra.Command {
	var value, label, trend string
	cmd := &cobra.Command{
		Use:   "stat",
		Short: "Stat callout (value + label + optional trend)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			tr := strings.ToLower(strings.TrimSpace(trend))
			if trend != "" && tr != "up" && tr != "down" && tr != "neutral" {
				return fmt.Errorf("trend must be up|down|neutral when set")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Stat")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Stat';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisAlignItems = 'MIN';")
			b.Line("root.counterAxisAlignItems = 'MIN';")
			b.Line("root.itemSpacing = 4;")
			b.Line("root.fills = [];")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("const row = figma.createFrame();")
			b.Line("row.layoutMode = 'HORIZONTAL';")
			b.Line("row.primaryAxisAlignItems = 'CENTER';")
			b.Line("row.counterAxisAlignItems = 'CENTER';")
			b.Line("row.itemSpacing = 8;")
			b.Line("row.fills = [];")
			b.Line("row.primaryAxisSizingMode = 'AUTO';")
			b.Line("row.counterAxisSizingMode = 'AUTO';")
			b.Line("const val = figma.createText();")
			b.Line("val.fontName = { family: 'Inter', style: ST_H4 };")
			b.Linef("val.characters = %q;", value)
			b.Line("val.fontSize = TY_H4;")
			b.Line("val.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("val.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("row.appendChild(val);")
			if tr != "" {
				b.Line("const trendMark = figma.createText();")
				b.Line("trendMark.fontName = { family: 'Inter', style: 'Medium' };")
				switch tr {
				case "up":
					b.Line("trendMark.characters = '▲'; trendMark.fontSize = TY_SMALL; trendMark.fills = [{ type: 'SOLID', color: SUCCESS }];")
				case "down":
					b.Line("trendMark.characters = '▼'; trendMark.fontSize = TY_SMALL; trendMark.fills = [{ type: 'SOLID', color: ERR }];")
				default:
					b.Line("trendMark.characters = '—'; trendMark.fontSize = TY_SMALL; trendMark.fills = [{ type: 'SOLID', color: MT }];")
				}
				b.Line("trendMark.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("row.appendChild(trendMark);")
			}
			b.Line("root.appendChild(row);")
			b.Line("const cap = figma.createText();")
			b.Line("cap.fontName = { family: 'Inter', style: ST_LABEL };")
			b.Linef("cap.characters = %q;", label)
			b.Line("cap.fontSize = TY_LABEL;")
			b.Line("cap.fills = [{ type: 'SOLID', color: MT }];")
			b.Line("cap.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(cap);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&value, "value", "4.2x", "Primary stat value")
	cmd.Flags().StringVar(&label, "label", "Faster", "Caption label")
	cmd.Flags().StringVar(&trend, "trend", "", "optional: up|down|neutral")
	return cmd
}

func newUITableCmd() *cobra.Command {
	var dataPath, cols string
	cmd := &cobra.Command{
		Use:   "table",
		Short: "Table from JSON rows + column headers",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			raw, err := os.ReadFile(dataPath)
			if err != nil {
				return fmt.Errorf("read data file: %w", err)
			}
			var rows []map[string]any
			if err := json.Unmarshal(raw, &rows); err != nil {
				return fmt.Errorf("data.json must be a JSON array of objects: %w", err)
			}
			headers := splitCSVLine(cols)
			if len(headers) == 0 {
				return fmt.Errorf("cols must list at least one column name")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Table")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Table';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisAlignItems = 'STRETCH';")
			b.Line("root.counterAxisAlignItems = 'MIN';")
			b.Line("root.itemSpacing = 0;")
			b.Line("root.fills = [{type:'SOLID', color: CARD}];")
			b.Line("root.strokes = [{type:'SOLID', color: STK}]; root.strokeWeight = 1;")
			b.Line("root.cornerRadius = 8;")
			b.Line("root.clipsContent = true;")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'FIXED';")
			b.Line("root.resize(640, 40);")
			// Header
			b.Line("const head = figma.createFrame();")
			b.Line("head.name = 'Header';")
			b.Line("head.layoutMode = 'HORIZONTAL';")
			b.Line("head.primaryAxisAlignItems = 'MIN';")
			b.Line("head.counterAxisAlignItems = 'CENTER';")
			b.Line("head.layoutAlign = 'STRETCH';")
			b.Line("head.paddingLeft = 12; head.paddingRight = 12; head.paddingTop = 10; head.paddingBottom = 10;")
			b.Line("head.itemSpacing = 0;")
			b.Line("head.primaryAxisSizingMode = 'FIXED';")
			b.Line("head.counterAxisSizingMode = 'AUTO';")
			b.Line("head.resize(640, 1);")
			b.Line("head.fills = [{type:'SOLID', color: STK, opacity: 0.5}];")
			colW := 640 / len(headers)
			if colW < 80 {
				colW = 80
			}
			for _, h := range headers {
				b.Line("{")
				b.Line("const c = figma.createFrame();")
				b.Line("c.layoutMode = 'HORIZONTAL';")
				b.Line("c.primaryAxisAlignItems = 'MIN';")
				b.Line("c.counterAxisAlignItems = 'CENTER';")
				b.Line("c.paddingLeft = 4; c.paddingRight = 4;")
				b.Linef("c.resize(%d, 1); c.layoutGrow = 1;", colW)
				b.Line("const tx = figma.createText();")
				b.Line("tx.fontName = { family: 'Inter', style: 'Semi Bold' };")
				b.Linef("tx.characters = %q;", strings.TrimSpace(h))
				b.Line("tx.fontSize = TY_SMALL;")
				b.Line("tx.fills = [{ type: 'SOLID', color: BD }];")
				b.Line("tx.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("c.appendChild(tx);")
				b.Line("head.appendChild(c);")
				b.Line("}")
			}
			b.Line("root.appendChild(head);")
			b.Line("const sep = figma.createFrame();")
			b.Line("sep.layoutMode = 'NONE';")
			b.Line("sep.resize(640, 1);")
			b.Line("sep.fills = [{type:'SOLID', color: STK}];")
			b.Line("sep.layoutAlign = 'STRETCH';")
			b.Line("root.appendChild(sep);")
			// Body rows from Go-generated data
			for _, row := range rows {
				b.Line("{")
				b.Line("const row = figma.createFrame();")
				b.Line("row.layoutMode = 'HORIZONTAL';")
				b.Line("row.primaryAxisAlignItems = 'MIN';")
				b.Line("row.counterAxisAlignItems = 'CENTER';")
				b.Line("row.layoutAlign = 'STRETCH';")
				b.Line("row.paddingLeft = 12; row.paddingRight = 12; row.paddingTop = 8; row.paddingBottom = 8;")
				b.Line("row.itemSpacing = 0;")
				b.Line("row.primaryAxisSizingMode = 'FIXED';")
				b.Line("row.counterAxisSizingMode = 'AUTO';")
				b.Line("row.resize(640, 1);")
				b.Line("row.fills = [];")
				for _, h := range headers {
					key := strings.TrimSpace(h)
					cell := ""
					if row != nil {
						if v, ok := row[key]; ok && v != nil {
							cell = fmt.Sprint(v)
						}
					}
					b.Line("{")
					b.Line("const c = figma.createFrame();")
					b.Line("c.layoutMode = 'HORIZONTAL';")
					b.Line("c.primaryAxisAlignItems = 'MIN';")
					b.Line("c.counterAxisAlignItems = 'CENTER';")
					b.Line("c.paddingLeft = 4; c.paddingRight = 4;")
					b.Linef("c.resize(%d, 1); c.layoutGrow = 1;", colW)
					b.Line("const tx = figma.createText();")
					b.Line("tx.fontName = { family: 'Inter', style: ST_BODY };")
					b.Linef("tx.characters = %q;", cell)
					b.Line("tx.fontSize = TY_SMALL;")
					b.Line("tx.fills = [{ type: 'SOLID', color: WT }];")
					b.Line("tx.textAutoResize = 'WIDTH_AND_HEIGHT';")
					b.Line("c.appendChild(tx);")
					b.Line("row.appendChild(c);")
					b.Line("}")
				}
				b.Line("root.appendChild(row);")
				b.Line("}")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&dataPath, "data", "./data.json", "Path to JSON array of row objects")
	cmd.Flags().StringVar(&cols, "cols", "Name,Role,Status", "Comma-separated column keys (match JSON keys)")
	return cmd
}

func splitCSVLine(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func newUINavCmd() *cobra.Command {
	var items, style string
	cmd := &cobra.Command{
		Use:   "nav",
		Short: "Navigation links (topbar or sidebar)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			st := strings.ToLower(style)
			if st != "topbar" && st != "sidebar" {
				return fmt.Errorf("style must be topbar|sidebar")
			}
			labels := splitCSVLine(items)
			if len(labels) == 0 {
				return fmt.Errorf("items must list at least one label")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Nav")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Nav';")
			if st == "topbar" {
				b.Line("root.layoutMode = 'HORIZONTAL';")
				b.Line("root.primaryAxisAlignItems = 'MIN';")
				b.Line("root.counterAxisAlignItems = 'CENTER';")
				b.Line("root.itemSpacing = 24;")
				b.Line("root.paddingLeft = 16; root.paddingRight = 16; root.paddingTop = 12; root.paddingBottom = 12;")
			} else {
				b.Line("root.layoutMode = 'VERTICAL';")
				b.Line("root.primaryAxisAlignItems = 'MIN';")
				b.Line("root.counterAxisAlignItems = 'STRETCH';")
				b.Line("root.itemSpacing = 4;")
				b.Line("root.paddingLeft = 12; root.paddingRight = 12; root.paddingTop = 12; root.paddingBottom = 12;")
			}
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [{type:'SOLID', color: BG}];")
			b.Line("root.strokes = [{type:'SOLID', color: STK}]; root.strokeWeight = 1;")
			b.Line("root.cornerRadius = 8;")
			for _, lab := range labels {
				b.Line("{")
				b.Line("const link = figma.createText();")
				b.Line("link.fontName = { family: 'Inter', style: 'Medium' };")
				b.Linef("link.characters = %q;", lab)
				b.Line("link.fontSize = TY_BODY;")
				b.Line("link.fills = [{ type: 'SOLID', color: BD }];")
				b.Line("link.textAutoResize = 'WIDTH_AND_HEIGHT';")
				if st == "sidebar" {
					b.Line("link.layoutAlign = 'STRETCH';")
				}
				b.Line("root.appendChild(link);")
				b.Line("}")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&items, "items", "Home,Products,Pricing", "Comma-separated nav labels")
	cmd.Flags().StringVar(&style, "style", "topbar", "topbar|sidebar")
	return cmd
}

func newUIFooterCmd() *cobra.Command {
	var copyright string
	var cols int
	cmd := &cobra.Command{
		Use:   "footer",
		Short: "Footer with N link columns + optional copyright",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			if cols < 1 || cols > 6 {
				return fmt.Errorf("cols must be 1–6")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Footer")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Footer';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisAlignItems = 'STRETCH';")
			b.Line("root.counterAxisAlignItems = 'MIN';")
			b.Line("root.itemSpacing = 20;")
			b.Line("root.paddingLeft = 32; root.paddingRight = 32; root.paddingTop = 28; root.paddingBottom = 28;")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'FIXED';")
			b.Line("root.resize(960, 1);")
			b.Line("root.fills = [{type:'SOLID', color: CARD}];")
			b.Line("root.strokes = [{type:'SOLID', color: STK}]; root.strokeWeight = 1;")
			b.Line("const grid = figma.createFrame();")
			b.Line("grid.name = 'Columns';")
			b.Line("grid.layoutMode = 'HORIZONTAL';")
			b.Line("grid.primaryAxisAlignItems = 'MIN';")
			b.Line("grid.counterAxisAlignItems = 'MIN';")
			b.Line("grid.itemSpacing = 32;")
			b.Line("grid.layoutAlign = 'STRETCH';")
			b.Line("grid.primaryAxisSizingMode = 'FIXED';")
			b.Line("grid.counterAxisSizingMode = 'AUTO';")
			b.Line("grid.resize(896, 1);")
			b.Line("grid.fills = [];")
			for i := 0; i < cols; i++ {
				b.Line("{")
				b.Line("const col = figma.createFrame();")
				b.Line("col.layoutMode = 'VERTICAL';")
				b.Line("col.primaryAxisAlignItems = 'MIN';")
				b.Line("col.counterAxisAlignItems = 'MIN';")
				b.Line("col.itemSpacing = 8;")
				b.Line("col.layoutGrow = 1;")
				b.Line("col.primaryAxisSizingMode = 'AUTO';")
				b.Line("col.counterAxisSizingMode = 'AUTO';")
				b.Line("col.fills = [];")
				b.Line("const h = figma.createText();")
				b.Line("h.fontName = { family: 'Inter', style: 'Semi Bold' };")
				b.Linef("h.characters = 'Section %d';", i+1)
				b.Line("h.fontSize = TY_SMALL;")
				b.Line("h.fills = [{ type: 'SOLID', color: WT }];")
				b.Line("h.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("col.appendChild(h);")
				for j := 1; j <= 3; j++ {
					b.Linef("const l%d = figma.createText();", j)
					b.Linef("l%d.fontName = { family: 'Inter', style: 'Regular' };", j)
					b.Linef("l%d.characters = 'Link %d';", j, j)
					b.Linef("l%d.fontSize = TY_SMALL;", j)
					b.Linef("l%d.fills = [{ type: 'SOLID', color: MT }];", j)
					b.Linef("l%d.textAutoResize = 'WIDTH_AND_HEIGHT';", j)
					b.Linef("col.appendChild(l%d);", j)
				}
				b.Line("grid.appendChild(col);")
				b.Line("}")
			}
			b.Line("root.appendChild(grid);")
			if strings.TrimSpace(copyright) != "" {
				b.Line("const cr = figma.createText();")
				b.Line("cr.fontName = { family: 'Geist Mono', style: 'Regular' };")
				b.Linef("cr.characters = %q;", strings.TrimSpace(copyright))
				b.Line("cr.fontSize = 10;")
				b.Line("cr.fills = [{ type: 'SOLID', color: MT }];")
				b.Line("cr.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("cr.letterSpacing = { value: 6, unit: 'PERCENT' };")
				b.Line("cr.layoutAlign = 'MIN';")
				b.Line("root.appendChild(cr);")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&cols, "cols", 3, "Number of columns (1–6)")
	cmd.Flags().StringVar(&copyright, "copyright", "", "Optional copyright line")
	return cmd
}

func newUICheckboxCmd() *cobra.Command {
	var label string
	var checked bool
	cmd := &cobra.Command{
		Use:   "checkbox",
		Short: "Checkbox control with label (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Checkbox")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Checkbox';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.itemSpacing = 8;")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			b.Line("const box = figma.createFrame();")
			b.Line("box.name = 'Box';")
			b.Line("box.resize(16, 16);")
			b.Line("box.cornerRadius = 4;")
			if checked {
				b.Line("box.fills = [{ type: 'SOLID', color: BL }];")
				b.Line("box.strokes = [];")
				b.Line("const ck = figma.createText();")
				b.Line("ck.fontName = { family: 'Inter', style: 'Bold' };")
				b.Line("ck.characters = '✓';")
				b.Line("ck.fontSize = 11;")
				b.Line("ck.fills = [{ type: 'SOLID', color: WT }];")
				b.Line("ck.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("ck.x = 2; ck.y = 1;")
				b.Line("box.appendChild(ck);")
			} else {
				b.Line("box.fills = [];")
				b.Line("box.strokes = [{ type: 'SOLID', color: STK }]; box.strokeWeight = 1.5;")
			}
			b.Line("root.appendChild(box);")
			b.Line("const lbl = figma.createText();")
			b.Line("lbl.fontName = { family: 'Inter', style: ST_BODY };")
			b.Linef("lbl.characters = %q;", label)
			b.Line("lbl.fontSize = TY_BODY;")
			b.Line("lbl.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("lbl.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(lbl);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "Remember me", "Checkbox label")
	cmd.Flags().BoolVar(&checked, "checked", false, "Render in checked state")
	return cmd
}

func newUIRadioCmd() *cobra.Command {
	var label string
	var selected bool
	cmd := &cobra.Command{
		Use:   "radio",
		Short: "Radio button with label (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Radio")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Radio';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.itemSpacing = 8;")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			b.Line("const circle = figma.createEllipse();")
			b.Line("circle.name = 'Dot';")
			b.Line("circle.resize(16, 16);")
			if selected {
				b.Line("circle.fills = [{ type: 'SOLID', color: BL }];")
				b.Line("circle.strokes = [];")
				b.Line("const inner = figma.createEllipse();")
				b.Line("inner.resize(6, 6);")
				b.Line("inner.fills = [{ type: 'SOLID', color: WT }];")
				b.Line("inner.x = 5; inner.y = 5;")
				b.Line("figma.currentPage.appendChild(inner);")
			} else {
				b.Line("circle.fills = [];")
				b.Line("circle.strokes = [{ type: 'SOLID', color: STK }]; circle.strokeWeight = 1.5;")
			}
			b.Line("root.appendChild(circle);")
			b.Line("const lbl = figma.createText();")
			b.Line("lbl.fontName = { family: 'Inter', style: ST_BODY };")
			b.Linef("lbl.characters = %q;", label)
			b.Line("lbl.fontSize = TY_BODY;")
			b.Line("lbl.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("lbl.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("root.appendChild(lbl);")
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "Option A", "Radio label")
	cmd.Flags().BoolVar(&selected, "selected", false, "Render in selected state")
	return cmd
}

func newUITabsCmd() *cobra.Command {
	var tabsRaw string
	var active int
	cmd := &cobra.Command{
		Use:   "tabs",
		Short: "Horizontal tab bar (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			tabs := splitTrimmed(tabsRaw)
			if len(tabs) < 2 {
				return fmt.Errorf("provide at least 2 tab labels via --tabs")
			}
			if active < 0 || active >= len(tabs) {
				return fmt.Errorf("--active must be between 0 and %d", len(tabs)-1)
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Tabs")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Tabs';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.itemSpacing = 0;")
			b.Line("root.counterAxisAlignItems = 'MIN';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			b.Line("root.strokes = [{ type: 'SOLID', color: STK }];")
			b.Line("root.strokeWeight = 1;")
			b.Line("root.strokeAlign = 'INSIDE';")
			for i, tab := range tabs {
				isActive := i == active
				b.Linef("{ // Tab %d", i)
				b.Line("const tab = figma.createFrame();")
				b.Line("tab.layoutMode = 'VERTICAL';")
				b.Line("tab.primaryAxisSizingMode = 'AUTO';")
				b.Line("tab.counterAxisSizingMode = 'AUTO';")
				b.Line("tab.paddingLeft = 16; tab.paddingRight = 16; tab.paddingTop = 12; tab.paddingBottom = 12;")
				b.Line("tab.fills = [];")
				if isActive {
					b.Line("tab.strokes = [{ type: 'SOLID', color: BL }];")
					b.Line("tab.strokeWeight = 2;")
					b.Line("tab.strokeAlign = 'OUTSIDE';")
				} else {
					b.Line("tab.strokes = [];")
				}
				b.Line("const txt = figma.createText();")
				b.Line("txt.fontName = { family: 'Inter', style: ST_LABEL };")
				b.Linef("txt.characters = %q;", tab)
				b.Line("txt.fontSize = TY_LABEL;")
				if isActive {
					b.Line("txt.fills = [{ type: 'SOLID', color: BL }];")
				} else {
					b.Line("txt.fills = [{ type: 'SOLID', color: MT }];")
				}
				b.Line("txt.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("tab.appendChild(txt);")
				b.Line("root.appendChild(tab);")
				b.Line("}")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&tabsRaw, "tabs", "Overview,Analytics,Settings", "Comma-separated tab labels")
	cmd.Flags().IntVar(&active, "active", 0, "Index of the active tab (0-based)")
	return cmd
}

func newUIDropdownCmd() *cobra.Command {
	var label, placeholder string
	var open bool
	var optionsRaw string
	cmd := &cobra.Command{
		Use:   "dropdown",
		Short: "Select dropdown control (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Dropdown")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Dropdown';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.itemSpacing = 6;")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			if label != "" {
				b.Line("const lab = figma.createText();")
				b.Line("lab.fontName = { family: 'Inter', style: ST_LABEL };")
				b.Linef("lab.characters = %q;", label)
				b.Line("lab.fontSize = TY_LABEL;")
				b.Line("lab.fills = [{ type: 'SOLID', color: MT }];")
				b.Line("lab.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("root.appendChild(lab);")
			}
			b.Line("const trigger = figma.createFrame();")
			b.Line("trigger.name = 'Trigger';")
			b.Line("trigger.layoutMode = 'HORIZONTAL';")
			b.Line("trigger.primaryAxisAlignItems = 'SPACE_BETWEEN';")
			b.Line("trigger.counterAxisAlignItems = 'CENTER';")
			b.Line("trigger.paddingLeft = 12; trigger.paddingRight = 12; trigger.paddingTop = 10; trigger.paddingBottom = 10;")
			b.Line("trigger.cornerRadius = 8;")
			b.Line("trigger.primaryAxisSizingMode = 'FIXED'; trigger.counterAxisSizingMode = 'AUTO';")
			b.Line("trigger.resize(280, 44);")
			b.Line("trigger.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("trigger.strokes = [{ type: 'SOLID', color: STK }]; trigger.strokeWeight = 1;")
			b.Line("const ph = figma.createText();")
			b.Line("ph.fontName = { family: 'Inter', style: ST_SMALL };")
			b.Linef("ph.characters = %q;", placeholder)
			b.Line("ph.fontSize = TY_SMALL;")
			b.Line("ph.fills = [{ type: 'SOLID', color: MT }];")
			b.Line("ph.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("trigger.appendChild(ph);")
			b.Line("const chevron = figma.createText();")
			b.Line("chevron.fontName = { family: 'Inter', style: 'Regular' };")
			b.Line("chevron.characters = '▾';")
			b.Line("chevron.fontSize = 12;")
			b.Line("chevron.fills = [{ type: 'SOLID', color: MT }];")
			b.Line("chevron.textAutoResize = 'WIDTH_AND_HEIGHT';")
			b.Line("trigger.appendChild(chevron);")
			b.Line("root.appendChild(trigger);")
			if open {
				opts := splitTrimmed(optionsRaw)
				b.Line("const menu = figma.createFrame();")
				b.Line("menu.name = 'Menu';")
				b.Line("menu.layoutMode = 'VERTICAL';")
				b.Line("menu.itemSpacing = 0;")
				b.Line("menu.cornerRadius = 8;")
				b.Line("menu.primaryAxisSizingMode = 'AUTO'; menu.counterAxisSizingMode = 'AUTO';")
				b.Line("menu.fills = [{ type: 'SOLID', color: CARD }];")
				b.Line("menu.strokes = [{ type: 'SOLID', color: STK }]; menu.strokeWeight = 1;")
				for _, opt := range opts {
					b.Line("{ const item = figma.createFrame();")
					b.Line("item.layoutMode = 'HORIZONTAL';")
					b.Line("item.paddingLeft = 12; item.paddingRight = 12; item.paddingTop = 10; item.paddingBottom = 10;")
					b.Line("item.primaryAxisSizingMode = 'AUTO'; item.counterAxisSizingMode = 'AUTO';")
					b.Line("item.fills = [];")
					b.Line("const ot = figma.createText();")
					b.Line("ot.fontName = { family: 'Inter', style: ST_SMALL };")
					b.Linef("ot.characters = %q;", opt)
					b.Line("ot.fontSize = TY_SMALL;")
					b.Line("ot.fills = [{ type: 'SOLID', color: WT }];")
					b.Line("ot.textAutoResize = 'WIDTH_AND_HEIGHT';")
					b.Line("item.appendChild(ot);")
					b.Line("menu.appendChild(item); }")
				}
				b.Line("root.appendChild(menu);")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "Sort by", "Optional field label above the trigger")
	cmd.Flags().StringVar(&placeholder, "placeholder", "Select an option…", "Placeholder text in the trigger")
	cmd.Flags().BoolVar(&open, "open", false, "Render the dropdown in open/expanded state")
	cmd.Flags().StringVar(&optionsRaw, "options", "Option A,Option B,Option C", "Comma-separated options shown when --open")
	return cmd
}

func newUIBreadcrumbCmd() *cobra.Command {
	var itemsRaw string
	cmd := &cobra.Command{
		Use:   "breadcrumb",
		Short: "Breadcrumb navigation trail (auto-layout)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			items := splitTrimmed(itemsRaw)
			if len(items) < 1 {
				return fmt.Errorf("provide at least one item via --items")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Breadcrumb")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Breadcrumb';")
			b.Line("root.layoutMode = 'HORIZONTAL';")
			b.Line("root.itemSpacing = 6;")
			b.Line("root.counterAxisAlignItems = 'CENTER';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			for i, item := range items {
				isLast := i == len(items)-1
				b.Linef("{ // crumb %d", i)
				b.Line("const crumb = figma.createText();")
				b.Line("crumb.fontName = { family: 'Inter', style: ST_SMALL };")
				b.Linef("crumb.characters = %q;", item)
				b.Line("crumb.fontSize = TY_SMALL;")
				if isLast {
					b.Line("crumb.fills = [{ type: 'SOLID', color: WT }];")
				} else {
					b.Line("crumb.fills = [{ type: 'SOLID', color: MT }];")
				}
				b.Line("crumb.textAutoResize = 'WIDTH_AND_HEIGHT';")
				b.Line("root.appendChild(crumb);")
				if !isLast {
					b.Line("const sep = figma.createText();")
					b.Line("sep.fontName = { family: 'Inter', style: 'Regular' };")
					b.Line("sep.characters = '/';")
					b.Line("sep.fontSize = TY_SMALL;")
					b.Line("sep.fills = [{ type: 'SOLID', color: STK }];")
					b.Line("sep.textAutoResize = 'WIDTH_AND_HEIGHT';")
					b.Line("root.appendChild(sep);")
				}
				b.Line("}")
			}
			b.Line("figma.currentPage.appendChild(root);")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&itemsRaw, "items", "Home,Products,Checkout", "Comma-separated breadcrumb labels (last = current)")
	return cmd
}

func newUISkeletonCmd() *cobra.Command {
	var variant string
	var rows int
	cmd := &cobra.Command{
		Use:   "skeleton",
		Short: "Loading skeleton placeholder (text, card, or list variant)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			v := strings.ToLower(variant)
			if v != "text" && v != "card" && v != "list" {
				return fmt.Errorf("variant must be text|card|list")
			}
			b := codegen.New()
			uiPreamble(b, t, resolvePage())
			b.Comment("UI / Skeleton — " + v)
			b.Line("const shimmer = { type: 'SOLID', color: STK };")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'UI/Skeleton';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.itemSpacing = 12;")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'AUTO';")
			b.Line("root.fills = [];")
			b.Line("figma.currentPage.appendChild(root);")
			switch v {
			case "text":
				widths := []int{280, 240, 200}
				for i, w := range widths {
					b.Linef("{ const r%d = figma.createRectangle();", i)
					b.Linef("r%d.resize(%d, 14); r%d.cornerRadius = 6;", i, w, i)
					b.Linef("r%d.fills = [shimmer]; root.appendChild(r%d); }", i, i)
				}
			case "card":
				b.Line("const img = figma.createRectangle();")
				b.Line("img.resize(280, 160); img.cornerRadius = 8;")
				b.Line("img.fills = [shimmer]; root.appendChild(img);")
				for i, w := range []int{220, 180} {
					b.Linef("{ const r%d = figma.createRectangle();", i)
					b.Linef("r%d.resize(%d, 14); r%d.cornerRadius = 6;", i, w, i)
					b.Linef("r%d.fills = [shimmer]; root.appendChild(r%d); }", i, i)
				}
			case "list":
				if rows < 1 {
					rows = 3
				}
				for i := 0; i < rows; i++ {
					b.Linef("{ // row %d", i)
					b.Line("const row = figma.createFrame();")
					b.Line("row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 12;")
					b.Line("row.primaryAxisSizingMode = 'AUTO'; row.counterAxisSizingMode = 'AUTO';")
					b.Line("row.fills = [];")
					b.Line("const av = figma.createEllipse(); av.resize(36, 36); av.fills = [shimmer]; row.appendChild(av);")
					b.Line("const col = figma.createFrame(); col.layoutMode = 'VERTICAL'; col.itemSpacing = 8;")
					b.Line("col.primaryAxisSizingMode = 'AUTO'; col.counterAxisSizingMode = 'AUTO'; col.fills = [];")
					b.Line("const t1 = figma.createRectangle(); t1.resize(160, 12); t1.cornerRadius = 5; t1.fills = [shimmer]; col.appendChild(t1);")
					b.Line("const t2 = figma.createRectangle(); t2.resize(100, 10); t2.cornerRadius = 5; t2.fills = [shimmer]; col.appendChild(t2);")
					b.Line("row.appendChild(col);")
					b.Line("root.appendChild(row); }")
				}
			}
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&variant, "variant", "text", "text|card|list")
	cmd.Flags().IntVar(&rows, "rows", 3, "Number of rows for the list variant")
	return cmd
}

// splitTrimmed splits a comma-separated string and trims whitespace from each part.
func splitTrimmed(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func mustMarshalJSON(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func newUIHeroCmd() *cobra.Command {
	var (
		title    string
		subtitle string
		cta      string
		badge    string
	)
	cmd := &cobra.Command{
		Use:     "hero",
		Short:   "Complete hero section with heading, subtitle, and CTA",
		Example: `  figma-kit ui hero -t noir --title "Build Faster" --subtitle "Ship in days, not months" --cta "Get Started"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if title == "" {
				title = "Build Something Amazing"
			}
			if subtitle == "" {
				subtitle = "The fastest way to ship beautiful products"
			}
			if cta == "" {
				cta = "Get Started"
			}

			b.Line("const hero = figma.createFrame();")
			b.Line("hero.name = 'Hero Section';")
			b.Line("hero.resize(1440, 720);")
			b.Line("hero.layoutMode = 'VERTICAL';")
			b.Line("hero.primaryAxisAlignItems = 'CENTER';")
			b.Line("hero.counterAxisAlignItems = 'CENTER';")
			b.Line("hero.paddingTop = hero.paddingBottom = 120;")
			b.Line("hero.paddingLeft = hero.paddingRight = 80;")
			b.Line("hero.itemSpacing = 24;")
			b.Line("hero.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")

			if badge != "" {
				b.Linef("{ const badge = figma.createFrame(); badge.name = 'Badge'; badge.layoutMode = 'HORIZONTAL'; badge.paddingLeft = badge.paddingRight = 16; badge.paddingTop = badge.paddingBottom = 6; badge.cornerRadius = 999; badge.fills = [{type:'SOLID', color:typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96}, opacity:0.15}]; badge.counterAxisSizingMode = 'AUTO'; badge.primaryAxisSizingMode = 'AUTO';")
				b.Linef("const bt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); bt.fontName = {family:'Inter',style:'Medium'}; bt.fontSize = 13; bt.characters = %q; bt.fills = [{type:'SOLID', color:typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96}}]; badge.appendChild(bt); hero.appendChild(badge); }", badge)
			}

			b.Linef("{ const h1 = figma.createText(); h1.name = 'Heading'; await figma.loadFontAsync({family:'Inter',style:'Bold'}); h1.fontName = {family:'Inter',style:'Bold'}; h1.fontSize = 64; h1.characters = %q; h1.fills = [{type:'SOLID', color:typeof fg !== 'undefined' ? fg : {r:0.95,g:0.95,b:0.97}}]; h1.textAlignHorizontal = 'CENTER'; h1.textAutoResize = 'WIDTH_AND_HEIGHT'; hero.appendChild(h1); }", title)
			b.Linef("{ const sub = figma.createText(); sub.name = 'Subtitle'; await figma.loadFontAsync({family:'Inter',style:'Regular'}); sub.fontName = {family:'Inter',style:'Regular'}; sub.fontSize = 20; sub.characters = %q; sub.fills = [{type:'SOLID', color:typeof muted !== 'undefined' ? muted : {r:0.55,g:0.55,b:0.6}}]; sub.textAlignHorizontal = 'CENTER'; sub.textAutoResize = 'WIDTH_AND_HEIGHT'; hero.appendChild(sub); }", subtitle)
			b.Linef("{ const btn = figma.createFrame(); btn.name = 'CTA Button'; btn.layoutMode = 'HORIZONTAL'; btn.paddingLeft = btn.paddingRight = 32; btn.paddingTop = btn.paddingBottom = 14; btn.cornerRadius = 8; btn.fills = [{type:'SOLID', color:typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96}}]; btn.counterAxisSizingMode = 'AUTO'; btn.primaryAxisSizingMode = 'AUTO';")
			b.Linef("const ct = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); ct.fontName = {family:'Inter',style:'Semi Bold'}; ct.fontSize = 16; ct.characters = %q; ct.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}]; btn.appendChild(ct); hero.appendChild(btn); }", cta)

			b.Line("pg.appendChild(hero);")
			b.ReturnIDs("hero.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Hero heading text")
	cmd.Flags().StringVar(&subtitle, "subtitle", "", "Hero subtitle text")
	cmd.Flags().StringVar(&cta, "cta", "", "Call-to-action button text")
	cmd.Flags().StringVar(&badge, "badge", "", "Optional badge text above heading")
	return cmd
}

func newUIPricingCmd() *cobra.Command {
	var tiersJSON string
	cmd := &cobra.Command{
		Use:     "pricing",
		Short:   "Pricing table with tier cards",
		Example: `  figma-kit ui pricing -t noir --tiers '[{"name":"Free","price":"$0","features":["5 projects","1GB storage"]},{"name":"Pro","price":"$29","features":["Unlimited","100GB"],"highlighted":true}]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if tiersJSON == "" {
				tiersJSON = `[{"name":"Starter","price":"$0","period":"/mo","features":["5 projects","1GB storage","Email support"],"cta":"Start Free"},{"name":"Pro","price":"$29","period":"/mo","features":["Unlimited projects","100GB storage","Priority support","API access"],"cta":"Go Pro","highlighted":true},{"name":"Enterprise","price":"Custom","features":["Everything in Pro","SSO","Dedicated support","SLA"],"cta":"Contact Sales"}]`
			}

			b.Linef("const tiers = JSON.parse(%s);", jsStringLiteral(tiersJSON))
			b.Line("const grid = figma.createFrame(); grid.name = 'Pricing'; grid.layoutMode = 'HORIZONTAL'; grid.itemSpacing = 24; grid.paddingLeft = grid.paddingRight = 48; grid.paddingTop = grid.paddingBottom = 48; grid.counterAxisAlignItems = 'CENTER';")
			b.Line("grid.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")
			b.Line("const accColor = typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96};")
			b.Line("for (const tier of tiers) {")
			b.Line("  const card = figma.createFrame(); card.name = tier.name; card.layoutMode = 'VERTICAL'; card.itemSpacing = 16; card.paddingLeft = card.paddingRight = 28; card.paddingTop = card.paddingBottom = 32; card.resize(280, 400); card.cornerRadius = 16;")
			b.Line("  if (tier.highlighted) { card.fills = [{type:'SOLID', color:{r:0.08,g:0.08,b:0.12}}]; card.strokes = [{type:'SOLID', color:accColor}]; card.strokeWeight = 2; card.effects = [{type:'DROP_SHADOW', color:{...accColor, a:0.2}, offset:{x:0,y:0}, radius:20, spread:0, visible:true, blendMode:'NORMAL'}]; }")
			b.Line("  else { card.fills = [{type:'SOLID', color:{r:0.06,g:0.06,b:0.09}}]; card.strokes = [{type:'SOLID', color:{r:0.15,g:0.15,b:0.2}}]; card.strokeWeight = 1; }")
			b.Line("  const nm = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); nm.fontName = {family:'Inter',style:'Semi Bold'}; nm.fontSize = 18; nm.characters = tier.name; nm.fills = [{type:'SOLID', color:{r:0.8,g:0.8,b:0.85}}]; card.appendChild(nm);")
			b.Line("  const pr = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Bold'}); pr.fontName = {family:'Inter',style:'Bold'}; pr.fontSize = 36; pr.characters = tier.price + (tier.period || ''); pr.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.97}}]; card.appendChild(pr);")
			b.Line("  if (tier.features) for (const f of tier.features) { const ft = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); ft.fontName = {family:'Inter',style:'Regular'}; ft.fontSize = 14; ft.characters = '✓ ' + f; ft.fills = [{type:'SOLID', color:{r:0.6,g:0.6,b:0.65}}]; card.appendChild(ft); }")
			b.Line("  if (tier.cta) { const btn = figma.createFrame(); btn.name = 'CTA'; btn.layoutMode = 'HORIZONTAL'; btn.primaryAxisAlignItems = 'CENTER'; btn.counterAxisAlignItems = 'CENTER'; btn.paddingLeft = btn.paddingRight = 24; btn.paddingTop = btn.paddingBottom = 12; btn.cornerRadius = 8; btn.counterAxisSizingMode = 'AUTO'; btn.primaryAxisSizingMode = 'AUTO';")
			b.Line("    if (tier.highlighted) btn.fills = [{type:'SOLID', color:accColor}]; else btn.fills = [{type:'SOLID', color:{r:0.12,g:0.12,b:0.16}}];")
			b.Line("    const ct = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); ct.fontName = {family:'Inter',style:'Semi Bold'}; ct.fontSize = 14; ct.characters = tier.cta; ct.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}]; btn.appendChild(ct); card.appendChild(btn); }")
			b.Line("  grid.appendChild(card);")
			b.Line("}")
			b.Line("grid.resize(grid.children.length * 304 + 96, 500);")
			b.Line("pg.appendChild(grid);")
			b.ReturnIDs("grid.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&tiersJSON, "tiers", "", "Tiers as JSON array")
	return cmd
}

func newUIFeatureGridCmd() *cobra.Command {
	var (
		featuresJSON string
		cols         int
	)
	cmd := &cobra.Command{
		Use:     "feature-grid",
		Short:   "Grid of feature cards with icon, title, and description",
		Example: `  figma-kit ui feature-grid -t noir --cols 3 --features '[{"title":"Fast","desc":"Sub-ms reads"},{"title":"Secure","desc":"E2E encrypted"},{"title":"Global","desc":"Edge network"}]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if featuresJSON == "" {
				featuresJSON = `[{"title":"Lightning Fast","desc":"Sub-millisecond response times"},{"title":"Fully Secure","desc":"End-to-end encryption built in"},{"title":"Global Scale","desc":"Deployed on the edge worldwide"},{"title":"Developer First","desc":"APIs you'll love working with"},{"title":"Real-time","desc":"Live collaboration built in"},{"title":"Open Source","desc":"Community driven development"}]`
			}

			b.Linef("const features = JSON.parse(%s);", jsStringLiteral(featuresJSON))
			b.Linef("const cols = %d;", cols)
			b.Line("const grid = figma.createFrame(); grid.name = 'Feature Grid'; grid.layoutMode = 'VERTICAL'; grid.itemSpacing = 24; grid.paddingLeft = grid.paddingRight = 48; grid.paddingTop = grid.paddingBottom = 48;")
			b.Line("grid.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")
			b.Line("let row;")
			b.Line("for (let i = 0; i < features.length; i++) {")
			b.Line("  if (i % cols === 0) { row = figma.createFrame(); row.name = 'Row'; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 24; row.fills = []; row.counterAxisSizingMode = 'AUTO'; row.primaryAxisSizingMode = 'AUTO'; grid.appendChild(row); }")
			b.Line("  const f = features[i];")
			b.Line("  const card = figma.createFrame(); card.name = f.title; card.layoutMode = 'VERTICAL'; card.itemSpacing = 12; card.paddingLeft = card.paddingRight = 24; card.paddingTop = card.paddingBottom = 24; card.resize(280, 180); card.cornerRadius = 12; card.fills = [{type:'SOLID', color:{r:0.06,g:0.06,b:0.09}}]; card.strokes = [{type:'SOLID', color:{r:0.12,g:0.12,b:0.16}}]; card.strokeWeight = 1;")
			b.Line("  const icon = figma.createFrame(); icon.name = 'icon'; icon.resize(40, 40); icon.cornerRadius = 10; icon.fills = [{type:'SOLID', color:typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96}, opacity:0.15}]; card.appendChild(icon);")
			b.Line("  const tt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tt.fontName = {family:'Inter',style:'Semi Bold'}; tt.fontSize = 16; tt.characters = f.title; tt.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.97}}]; card.appendChild(tt);")
			b.Line("  const dd = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); dd.fontName = {family:'Inter',style:'Regular'}; dd.fontSize = 14; dd.characters = f.desc || ''; dd.fills = [{type:'SOLID', color:{r:0.55,g:0.55,b:0.6}}]; card.appendChild(dd);")
			b.Line("  row.appendChild(card);")
			b.Line("}")
			b.Line("grid.resize(cols * 304 + 96, Math.ceil(features.length / cols) * 204 + 96);")
			b.Line("pg.appendChild(grid);")
			b.ReturnIDs("grid.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&featuresJSON, "features", "", "Features as JSON array [{title,desc}]")
	cmd.Flags().IntVar(&cols, "cols", 3, "Number of columns (2, 3, or 4)")
	return cmd
}

func newUITestimonialCmd() *cobra.Command {
	var (
		name    string
		role    string
		quote   string
		rating  int
		variant string
	)
	cmd := &cobra.Command{
		Use:     "testimonial",
		Short:   "Testimonial/quote card with avatar, name, and rating",
		Example: `  figma-kit ui testimonial -t noir --name "Jane Doe" --role "CEO at Acme" --quote "This changed everything" --rating 5`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if name == "" {
				name = "Alex Johnson"
			}
			if quote == "" {
				quote = "This tool completely transformed our design workflow. We ship 3x faster now."
			}
			if role == "" {
				role = "Head of Design"
			}

			width := 480
			if variant == "large" {
				width = 640
			}

			b.Line("const card = figma.createFrame(); card.name = 'Testimonial';")
			b.Linef("card.resize(%d, 260); card.layoutMode = 'VERTICAL'; card.itemSpacing = 16; card.paddingLeft = card.paddingRight = 32; card.paddingTop = card.paddingBottom = 32; card.cornerRadius = 16;", width)
			b.Line("card.fills = [{type:'SOLID', color:{r:0.06,g:0.06,b:0.09}}]; card.strokes = [{type:'SOLID', color:{r:0.12,g:0.12,b:0.16}}]; card.strokeWeight = 1;")

			if rating > 0 {
				stars := ""
				for i := 0; i < rating && i < 5; i++ {
					stars += "★"
				}
				for i := rating; i < 5; i++ {
					stars += "☆"
				}
				b.Linef("{ const st = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); st.fontName = {family:'Inter',style:'Regular'}; st.fontSize = 18; st.characters = %q; st.fills = [{type:'SOLID', color:{r:0.96,g:0.76,b:0.05}}]; card.appendChild(st); }", stars)
			}

			b.Linef("{ const qt = figma.createText(); qt.name = 'Quote'; await figma.loadFontAsync({family:'Inter',style:'Regular'}); qt.fontName = {family:'Inter',style:'Regular'}; qt.fontSize = 16; qt.lineHeight = {unit:'PIXELS',value:24}; qt.characters = '\"' + %q + '\"'; qt.fills = [{type:'SOLID', color:{r:0.8,g:0.8,b:0.85}}]; qt.textAutoResize = 'HEIGHT'; qt.resize(%d - 64, 1); card.appendChild(qt); }", quote, width)

			b.Line("const info = figma.createFrame(); info.name = 'Author'; info.layoutMode = 'HORIZONTAL'; info.itemSpacing = 12; info.counterAxisAlignItems = 'CENTER'; info.fills = []; info.counterAxisSizingMode = 'AUTO'; info.primaryAxisSizingMode = 'AUTO';")
			b.Line("const av = figma.createEllipse(); av.name = 'Avatar'; av.resize(40, 40); av.fills = [{type:'SOLID', color:typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96}}]; info.appendChild(av);")
			b.Line("const nameCol = figma.createFrame(); nameCol.layoutMode = 'VERTICAL'; nameCol.itemSpacing = 2; nameCol.fills = []; nameCol.counterAxisSizingMode = 'AUTO'; nameCol.primaryAxisSizingMode = 'AUTO';")
			b.Linef("const nm = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); nm.fontName = {family:'Inter',style:'Semi Bold'}; nm.fontSize = 14; nm.characters = %q; nm.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.97}}]; nameCol.appendChild(nm);", name)
			b.Linef("const rl = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); rl.fontName = {family:'Inter',style:'Regular'}; rl.fontSize = 13; rl.characters = %q; rl.fills = [{type:'SOLID', color:{r:0.5,g:0.5,b:0.55}}]; nameCol.appendChild(rl);", role)
			b.Line("info.appendChild(nameCol); card.appendChild(info);")

			b.Line("pg.appendChild(card);")
			b.ReturnIDs("card.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Person name")
	cmd.Flags().StringVar(&role, "role", "", "Person role/company")
	cmd.Flags().StringVar(&quote, "quote", "", "Testimonial quote text")
	cmd.Flags().IntVar(&rating, "rating", 5, "Star rating (1-5)")
	cmd.Flags().StringVar(&variant, "variant", "card", "Variant: card, inline, large")
	return cmd
}

func newUITimelineCmd() *cobra.Command {
	var entriesJSON string
	cmd := &cobra.Command{
		Use:     "timeline",
		Short:   "Vertical timeline with dated entries",
		Example: `  figma-kit ui timeline -t noir --entries '[{"date":"2024","title":"Founded","desc":"Started the journey"},{"date":"2025","title":"Series A","desc":"Raised $10M"}]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if entriesJSON == "" {
				entriesJSON = `[{"date":"Jan 2024","title":"Project Started","desc":"Initial concept and research phase"},{"date":"Jun 2024","title":"Beta Launch","desc":"First public beta release"},{"date":"Jan 2025","title":"v1.0 Release","desc":"Production-ready with full feature set"},{"date":"Jun 2025","title":"10K Users","desc":"Growing community adoption"}]`
			}

			b.Linef("const entries = JSON.parse(%s);", jsStringLiteral(entriesJSON))
			b.Line("const timeline = figma.createFrame(); timeline.name = 'Timeline'; timeline.layoutMode = 'VERTICAL'; timeline.itemSpacing = 0; timeline.paddingLeft = 60; timeline.paddingRight = 40; timeline.paddingTop = timeline.paddingBottom = 40;")
			b.Line("timeline.resize(500, entries.length * 120 + 80);")
			b.Line("timeline.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")
			b.Line("const accColor = typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96};")
			b.Line("for (let i = 0; i < entries.length; i++) {")
			b.Line("  const e = entries[i];")
			b.Line("  const row = figma.createFrame(); row.name = e.title; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 20; row.paddingTop = row.paddingBottom = 16; row.fills = []; row.counterAxisSizingMode = 'AUTO'; row.primaryAxisSizingMode = 'AUTO';")
			b.Line("  const dot = figma.createEllipse(); dot.resize(12, 12); dot.fills = [{type:'SOLID', color:accColor}]; row.appendChild(dot);")
			b.Line("  const content = figma.createFrame(); content.layoutMode = 'VERTICAL'; content.itemSpacing = 4; content.fills = []; content.counterAxisSizingMode = 'AUTO'; content.primaryAxisSizingMode = 'AUTO';")
			b.Line("  const dt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); dt.fontName = {family:'Inter',style:'Medium'}; dt.fontSize = 12; dt.characters = e.date; dt.fills = [{type:'SOLID', color:accColor}]; content.appendChild(dt);")
			b.Line("  const tt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tt.fontName = {family:'Inter',style:'Semi Bold'}; tt.fontSize = 16; tt.characters = e.title; tt.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.97}}]; content.appendChild(tt);")
			b.Line("  if (e.desc) { const dd = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); dd.fontName = {family:'Inter',style:'Regular'}; dd.fontSize = 14; dd.characters = e.desc; dd.fills = [{type:'SOLID', color:{r:0.55,g:0.55,b:0.6}}]; content.appendChild(dd); }")
			b.Line("  row.appendChild(content); timeline.appendChild(row);")
			b.Line("}")
			b.Line("pg.appendChild(timeline);")
			b.ReturnIDs("timeline.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&entriesJSON, "entries", "", "Timeline entries as JSON array")
	return cmd
}

func newUIStepperCmd() *cobra.Command {
	var (
		steps     int
		active    int
		labelsCSV string
		direction string
	)
	cmd := &cobra.Command{
		Use:     "stepper",
		Short:   "Step indicator / progress stepper",
		Example: `  figma-kit ui stepper -t noir --steps 4 --active 2 --labels "Account,Profile,Settings,Done"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			labels := strings.Split(labelsCSV, ",")
			for len(labels) < steps {
				labels = append(labels, fmt.Sprintf("Step %d", len(labels)+1))
			}

			isVertical := direction == "vertical"
			mode := "HORIZONTAL"
			if isVertical {
				mode = "VERTICAL"
			}

			b.Linef("const stepper = figma.createFrame(); stepper.name = 'Stepper'; stepper.layoutMode = %q; stepper.itemSpacing = 8; stepper.counterAxisAlignItems = 'CENTER'; stepper.paddingLeft = stepper.paddingRight = 32; stepper.paddingTop = stepper.paddingBottom = 24;", mode)
			if isVertical {
				b.Linef("stepper.resize(280, %d);", steps*80+48)
			} else {
				b.Linef("stepper.resize(%d, 80);", steps*120+64)
			}
			b.Line("stepper.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")
			b.Line("const accColor = typeof accent !== 'undefined' ? accent : {r:0.23,g:0.51,b:0.96};")
			b.Linef("const labels = %s;", mustMarshalJSON(labels[:steps]))
			b.Linef("const activeIdx = %d;", active)
			b.Line("for (let i = 0; i < labels.length; i++) {")
			b.Line("  if (i > 0) { const conn = figma.createRectangle(); conn.name = 'connector';")
			if isVertical {
				b.Line("    conn.resize(2, 24);")
			} else {
				b.Line("    conn.resize(40, 2);")
			}
			b.Line("    conn.fills = [{type:'SOLID', color: i <= activeIdx ? accColor : {r:0.2,g:0.2,b:0.25}}]; stepper.appendChild(conn); }")
			b.Line("  const step = figma.createFrame(); step.name = labels[i]; step.layoutMode = 'VERTICAL'; step.itemSpacing = 6; step.counterAxisAlignItems = 'CENTER'; step.fills = []; step.counterAxisSizingMode = 'AUTO'; step.primaryAxisSizingMode = 'AUTO';")
			b.Line("  const circle = figma.createEllipse(); circle.resize(32, 32);")
			b.Line("  if (i < activeIdx) circle.fills = [{type:'SOLID', color:accColor}];")
			b.Line("  else if (i === activeIdx) { circle.fills = [{type:'SOLID', color:accColor}]; circle.effects = [{type:'DROP_SHADOW', color:{...accColor,a:0.3}, offset:{x:0,y:0}, radius:8, spread:0, visible:true, blendMode:'NORMAL'}]; }")
			b.Line("  else circle.fills = [{type:'SOLID', color:{r:0.15,g:0.15,b:0.2}}];")
			b.Line("  step.appendChild(circle);")
			b.Line("  const num = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); num.fontName = {family:'Inter',style:'Semi Bold'}; num.fontSize = 13; num.characters = String(i + 1); num.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
			b.Line("  const lbl = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); lbl.fontName = {family:'Inter',style:'Regular'}; lbl.fontSize = 12; lbl.characters = labels[i]; lbl.fills = [{type:'SOLID', color: i <= activeIdx ? {r:0.9,g:0.9,b:0.95} : {r:0.45,g:0.45,b:0.5}}];")
			b.Line("  step.appendChild(lbl); stepper.appendChild(step);")
			b.Line("}")
			b.Line("pg.appendChild(stepper);")
			b.ReturnIDs("stepper.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&steps, "steps", 4, "Number of steps")
	cmd.Flags().IntVar(&active, "active", 1, "Current active step (0-indexed)")
	cmd.Flags().StringVar(&labelsCSV, "labels", "", "Comma-separated step labels")
	cmd.Flags().StringVar(&direction, "direction", "horizontal", "Direction: horizontal, vertical")
	return cmd
}

func newUIAccordionCmd() *cobra.Command {
	var (
		itemsJSON string
		openIdx   int
	)
	cmd := &cobra.Command{
		Use:     "accordion",
		Short:   "Expandable sections (FAQ accordion)",
		Example: `  figma-kit ui accordion -t noir --items '[{"question":"What is this?","answer":"A powerful design tool"},{"question":"How much?","answer":"Free and open source"}]'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			page := resolvePage()
			b := codegen.New()
			b.PageSetup(page)
			emitUIThemeTokens(b, t)

			if itemsJSON == "" {
				itemsJSON = `[{"question":"What is figma-kit?","answer":"A CLI tool for programmatic Figma design, powered by AI and the MCP server."},{"question":"Do I need a Figma account?","answer":"Yes, you need at least a free Figma account to use figma-kit."},{"question":"Can I use it without AI?","answer":"Yes! All commands output JavaScript that can be pasted into Figma plugins."},{"question":"Is it open source?","answer":"Yes, figma-kit is MIT licensed and available on GitHub."}]`
			}

			b.Linef("const items = JSON.parse(%s);", jsStringLiteral(itemsJSON))
			b.Linef("const openIdx = %d;", openIdx)
			b.Line("const acc = figma.createFrame(); acc.name = 'Accordion'; acc.layoutMode = 'VERTICAL'; acc.itemSpacing = 0; acc.paddingLeft = acc.paddingRight = 32; acc.paddingTop = acc.paddingBottom = 16;")
			b.Line("acc.resize(600, items.length * (openIdx >= 0 ? 100 : 60) + 80);")
			b.Line("acc.fills = typeof bg !== 'undefined' ? [{type:'SOLID', color:bg}] : [{type:'SOLID', color:{r:0.02,g:0.02,b:0.05}}];")
			b.Line("for (let i = 0; i < items.length; i++) {")
			b.Line("  const item = items[i]; const isOpen = i === openIdx;")
			b.Line("  const section = figma.createFrame(); section.name = item.question; section.layoutMode = 'VERTICAL'; section.itemSpacing = 8; section.paddingTop = section.paddingBottom = 16; section.fills = [];")
			b.Line("  section.counterAxisSizingMode = 'AUTO'; section.primaryAxisSizingMode = 'AUTO';")
			b.Line("  const header = figma.createFrame(); header.layoutMode = 'HORIZONTAL'; header.fills = []; header.counterAxisSizingMode = 'AUTO'; header.primaryAxisSizingMode = 'AUTO'; header.itemSpacing = 12;")
			b.Line("  const chevron = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); chevron.fontName = {family:'Inter',style:'Regular'}; chevron.fontSize = 16; chevron.characters = isOpen ? '▾' : '▸'; chevron.fills = [{type:'SOLID', color:{r:0.5,g:0.5,b:0.55}}]; header.appendChild(chevron);")
			b.Line("  const q = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); q.fontName = {family:'Inter',style:'Semi Bold'}; q.fontSize = 16; q.characters = item.question; q.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.97}}]; header.appendChild(q);")
			b.Line("  section.appendChild(header);")
			b.Line("  if (isOpen && item.answer) { const a = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); a.fontName = {family:'Inter',style:'Regular'}; a.fontSize = 14; a.lineHeight = {unit:'PIXELS',value:22}; a.characters = item.answer; a.fills = [{type:'SOLID', color:{r:0.6,g:0.6,b:0.65}}]; a.textAutoResize = 'HEIGHT'; a.resize(520, 1); section.appendChild(a); }")
			b.Line("  if (i < items.length - 1) { const div = figma.createRectangle(); div.resize(536, 1); div.fills = [{type:'SOLID', color:{r:0.12,g:0.12,b:0.16}}]; section.appendChild(div); }")
			b.Line("  acc.appendChild(section);")
			b.Line("}")
			b.Line("pg.appendChild(acc);")
			b.ReturnIDs("acc.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&itemsJSON, "items", "", "FAQ items as JSON array [{question,answer}]")
	cmd.Flags().IntVar(&openIdx, "open", 0, "Index of initially expanded item (-1 for all closed)")
	return cmd
}
