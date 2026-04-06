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
