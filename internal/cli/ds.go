package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/theme"
)

func newDSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ds",
		Short: "Design system management",
	}
	cmd.AddCommand(newDSCreateCmd())
	cmd.AddCommand(newDSColorsCmd())
	cmd.AddCommand(newDSTypeScaleCmd())
	cmd.AddCommand(newDSSpacingCmd())
	cmd.AddCommand(newDSElevationCmd())
	cmd.AddCommand(newDSRadiusCmd())
	cmd.AddCommand(newDSIconsCmd())
	cmd.AddCommand(newDSComponentCmd())
	cmd.AddCommand(newDSVariablesCmd())
	cmd.AddCommand(newDSSearchCmd())
	cmd.AddCommand(newDSImportCmd())
	cmd.AddCommand(newDSSyncTokensCmd())
	cmd.AddCommand(newDSAuditCmd())
	cmd.AddCommand(newDSTokensCmd())
	return cmd
}

func newDSCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create a design system page with swatches, type specimens, and spacing scale",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Comment("--- Design system page ---")
			b.Line("const page = figma.currentPage;")
			b.Line("page.name = 'Design System';")
			b.Line("const root = figma.createFrame();")
			b.Line("root.name = 'DS / Overview';")
			b.Line("root.layoutMode = 'VERTICAL';")
			b.Line("root.primaryAxisSizingMode = 'AUTO';")
			b.Line("root.counterAxisSizingMode = 'FIXED';")
			b.Line("root.resize(1200, 2000);")
			b.Line("root.itemSpacing = 48;")
			b.Line("root.paddingLeft = root.paddingRight = 64;")
			b.Line("root.paddingTop = root.paddingBottom = 64;")
			b.Line("root.fills = [{type:'SOLID', color:{r:0.98,g:0.98,b:0.99}}];")
			b.Line("page.appendChild(root);")
			b.Blank()
			b.Line("const swRow = figma.createFrame();")
			b.Line("swRow.name = 'Color swatches';")
			b.Line("swRow.layoutMode = 'HORIZONTAL';")
			b.Line("swRow.itemSpacing = 16;")
			b.Line("swRow.primaryAxisSizingMode = 'AUTO';")
			b.Line("swRow.counterAxisSizingMode = 'AUTO';")
			b.Line("swRow.fills = [];")
			b.Line("root.appendChild(swRow);")
			for _, name := range t.ColorNames() {
				c := t.Colors[name]
				b.Linef("{ const cell = figma.createFrame(); cell.name = %q; cell.resize(72, 72); cell.cornerRadius = 8;", name)
				b.Linef("  cell.fills = [{type:'SOLID', color:{r:%s,g:%s,b:%s}}];", codegen.FmtFloat(c.R), codegen.FmtFloat(c.G), codegen.FmtFloat(c.B))
				b.Line("  swRow.appendChild(cell); }")
			}
			b.Blank()
			b.Line("const typeCol = figma.createFrame();")
			b.Line("typeCol.name = 'Type scale';")
			b.Line("typeCol.layoutMode = 'VERTICAL';")
			b.Line("typeCol.itemSpacing = 12;")
			b.Line("typeCol.primaryAxisSizingMode = 'AUTO';")
			b.Line("typeCol.counterAxisSizingMode = 'FIXED';")
			b.Line("typeCol.resize(800, 1);")
			b.Line("typeCol.fills = [];")
			b.Line("root.appendChild(typeCol);")
			var typeNames []string
			for k := range t.Type {
				typeNames = append(typeNames, k)
			}
			sort.Strings(typeNames)
			for _, k := range typeNames {
				spec := t.Type[k]
				fam := spec.Family
				if fam == "" {
					fam = t.Fonts.Body
				}
				b.Linef("{ const tx = figma.createText(); tx.name = %q;", k)
				b.Linef("  await figma.loadFontAsync({family:%q, style:%q});", fam, spec.Style)
				b.Linef("  tx.fontName = {family:%q, style:%q};", fam, spec.Style)
				b.Linef("  tx.fontSize = %d;", spec.FontSize)
				if spec.LineHeight != nil {
					b.Linef("  tx.lineHeight = {unit:'PIXELS', value:%d};", *spec.LineHeight)
				}
				b.Linef("  tx.characters = %q;", k+" — The quick brown fox")
				b.Line("  typeCol.appendChild(tx); }")
			}
			b.Blank()
			b.Line("const sp = figma.createFrame();")
			b.Line("sp.name = 'Spacing scale';")
			b.Line("sp.layoutMode = 'HORIZONTAL';")
			b.Line("sp.itemSpacing = 8;")
			b.Line("sp.fills = [];")
			b.Line("root.appendChild(sp);")
			b.Line("const spaceVals = [4,8,12,16,24,32,48,64,96];")
			b.Line("for (const n of spaceVals) { const bar = figma.createRectangle(); bar.resize(n, 40);")
			b.Line("  bar.name = 'space-' + n; bar.fills = [{type:'SOLID', color:{r:0.2,g:0.4,b:0.9}}]; sp.appendChild(bar); }")
			b.ReturnIDs("root.id")
			output(b.String())
			return nil
		},
	}
}

func newDSColorsCmd() *cobra.Command {
	var primary string
	cmd := &cobra.Command{
		Use:   "colors",
		Short: "Create a palette page with tints and shades from a primary color",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			base, err := codegen.HexToRGB(primary)
			if err != nil {
				return err
			}
			white := codegen.RGB{R: 1, G: 1, B: 1}
			black := codegen.RGB{R: 0, G: 0, B: 0}
			var tints, shades []codegen.RGB
			for i := 1; i <= 5; i++ {
				f := float64(i) / 6
				tints = append(tints, mixRGB(base, white, f))
				shades = append(shades, mixRGB(base, black, f))
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Comment("--- Palette from primary ---")
			b.Line("const page = figma.currentPage;")
			b.Line("const wrap = figma.createFrame();")
			b.Line("wrap.name = 'Palette / Tints & Shades';")
			b.Line("wrap.layoutMode = 'VERTICAL'; wrap.itemSpacing = 32;")
			b.Line("wrap.paddingLeft = wrap.paddingRight = 48; wrap.paddingTop = wrap.paddingBottom = 48;")
			b.Line("wrap.resize(900, 1200);")
			b.Line("wrap.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
			b.Line("page.appendChild(wrap);")
			b.Line("// primary")
			b.Line("{ const row = figma.createFrame(); row.name = 'Primary'; row.layoutMode = 'HORIZONTAL'; row.fills = []; wrap.appendChild(row);")
			b.Linef("  const sq = figma.createFrame(); sq.resize(72, 72); sq.cornerRadius = 8;")
			b.Linef("  sq.fills = [{type:'SOLID', color:{r:%s,g:%s,b:%s}}]; row.appendChild(sq); }",
				codegen.FmtFloat(base.R), codegen.FmtFloat(base.G), codegen.FmtFloat(base.B))
			b.Line("}")
			b.Line("{ const row = figma.createFrame(); row.name = 'Tints'; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 12; row.fills = []; wrap.appendChild(row);")
			for i, c := range tints {
				b.Linef("  { const sq = figma.createFrame(); sq.name = 'tint-%d'; sq.resize(56, 56); sq.cornerRadius = 6;", i+1)
				b.Linef("    sq.fills = [{type:'SOLID', color:{r:%s,g:%s,b:%s}}]; row.appendChild(sq); }",
					codegen.FmtFloat(c.R), codegen.FmtFloat(c.G), codegen.FmtFloat(c.B))
			}
			b.Line("}")
			b.Line("{ const row = figma.createFrame(); row.name = 'Shades'; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 12; row.fills = []; wrap.appendChild(row);")
			for i, c := range shades {
				b.Linef("  { const sq = figma.createFrame(); sq.name = 'shade-%d'; sq.resize(56, 56); sq.cornerRadius = 6;", i+1)
				b.Linef("    sq.fills = [{type:'SOLID', color:{r:%s,g:%s,b:%s}}]; row.appendChild(sq); }",
					codegen.FmtFloat(c.R), codegen.FmtFloat(c.G), codegen.FmtFloat(c.B))
			}
			b.Line("}")
			b.ReturnIDs("wrap.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&primary, "primary", "#3B82F6", "Primary hex color")
	return cmd
}

func mixRGB(a, b codegen.RGB, t float64) codegen.RGB {
	return codegen.RGB{
		R: a.R*(1-t) + b.R*t,
		G: a.G*(1-t) + b.G*t,
		B: a.B*(1-t) + b.B*t,
	}
}

func newDSTypeScaleCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "type-scale",
		Short: "Create a dedicated type-scale specimen frame from the theme",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const col = figma.createFrame(); col.name = 'Type scale';")
			b.Line("col.layoutMode = 'VERTICAL'; col.itemSpacing = 20; col.paddingLeft = col.paddingRight = 48; col.paddingTop = col.paddingBottom = 48;")
			b.Line("col.resize(720, 1600); col.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
			b.Line("figma.currentPage.appendChild(col);")
			var names []string
			for k := range t.Type {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				spec := t.Type[k]
				fam := spec.Family
				if fam == "" {
					fam = t.Fonts.Body
				}
				b.Linef("{ const tx = figma.createText(); tx.name = %q;", k)
				b.Linef("  await figma.loadFontAsync({family:%q, style:%q});", fam, spec.Style)
				b.Linef("  tx.fontName = {family:%q, style:%q}; tx.fontSize = %d;", fam, spec.Style, spec.FontSize)
				b.Linef("  tx.characters = %q;", fmt.Sprintf("%s — %dpx %s", k, spec.FontSize, spec.Style))
				b.Line("  col.appendChild(tx); }")
			}
			b.ReturnIDs("col.id")
			output(b.String())
			return nil
		},
	}
}

func newDSSpacingCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spacing",
		Short: "Visualize theme spacing presets as labeled bars",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const f = figma.createFrame(); f.name = 'Spacing / Theme';")
			b.Line("f.layoutMode = 'VERTICAL'; f.itemSpacing = 16; f.paddingLeft = 40; f.paddingTop = 40;")
			b.Line("f.resize(640, 900); f.fills = [{type:'SOLID', color:{r:0.97,g:0.97,b:0.98}}];")
			b.Line("figma.currentPage.appendChild(f);")
			presets := []struct {
				label string
				p     theme.SpacingPreset
			}{
				{"page", t.Spacing.Page},
				{"card", t.Spacing.Card},
				{"slide", t.Spacing.Slide},
				{"frame16", t.Spacing.Frame16},
				{"letter", t.Spacing.Letter},
			}
			for _, pr := range presets {
				b.Linef("{ const row = figma.createFrame(); row.name = %q; row.layoutMode = 'VERTICAL'; row.itemSpacing = 6; row.fills = []; f.appendChild(row);", pr.label)
				if pr.p.Gap > 0 {
					b.Linef("  { const r = figma.createRectangle(); r.resize(%d, 24); r.name = 'gap'; r.fills = [{type:'SOLID', color:{r:0.2,g:0.5,b:0.9}}]; row.appendChild(r); }", pr.p.Gap)
				}
				if pr.p.Padding > 0 {
					b.Linef("  { const r = figma.createRectangle(); r.resize(%d, 24); r.name = 'padding'; r.fills = [{type:'SOLID', color:{r:0.3,g:0.7,b:0.4}}]; row.appendChild(r); }", pr.p.Padding)
				}
				if pr.p.Margin > 0 {
					b.Linef("  { const r = figma.createRectangle(); r.resize(%d, 24); r.name = 'margin'; r.fills = [{type:'SOLID', color:{r:0.9,g:0.5,b:0.2}}]; row.appendChild(r); }", pr.p.Margin)
				}
				b.Line("}")
			}
			b.ReturnIDs("f.id")
			output(b.String())
			return nil
		},
	}
}

func newDSElevationCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "elevation",
		Short: "Create frames demonstrating theme shadow presets",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const grid = figma.createFrame(); grid.name = 'Elevation';")
			b.Line("grid.layoutMode = 'HORIZONTAL'; grid.itemSpacing = 24; grid.paddingLeft = grid.paddingRight = 40; grid.paddingTop = grid.paddingBottom = 40;")
			b.Line("grid.resize(960, 220); grid.fills = [{type:'SOLID', color:{r:0.95,g:0.95,b:0.96}}];")
			b.Line("figma.currentPage.appendChild(grid);")
			var names []string
			for k := range t.Effects.Shadow {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, k := range names {
				sh := t.Effects.Shadow[k]
				b.Linef("{ const card = figma.createFrame(); card.name = %q; card.resize(160, 120); card.cornerRadius = 12;", k)
				b.Line("  card.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
				b.Linef("  card.effects = [{type:%q, color:{r:%s,g:%s,b:%s,a:%s}, offset:{x:%d,y:%d}, radius:%d, spread:%d, visible:%t, blendMode:%q}];",
					sh.Type,
					codegen.FmtFloat(sh.Color.R), codegen.FmtFloat(sh.Color.G), codegen.FmtFloat(sh.Color.B), codegen.FmtFloat(sh.Color.A),
					sh.Offset.X, sh.Offset.Y, sh.Radius, sh.Spread, sh.Visible, sh.BlendMode)
				b.Line("  grid.appendChild(card); }")
			}
			if len(names) == 0 {
				b.Line("{ const card = figma.createFrame(); card.name = 'shadow-sample'; card.resize(160, 120); card.cornerRadius = 12;")
				b.Line("  card.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
				b.Line("  card.effects = [{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.2}, offset:{x:0,y:8}, radius:24, spread:0, visible:true, blendMode:'NORMAL'}];")
				b.Line("  grid.appendChild(card); }")
			}
			b.ReturnIDs("grid.id")
			output(b.String())
			return nil
		},
	}
}

func newDSRadiusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "radius",
		Short: "Create corner-radius reference chips",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const row = figma.createFrame(); row.name = 'Radius scale'; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 16;")
			b.Line("row.paddingLeft = row.paddingRight = row.paddingTop = row.paddingBottom = 32; row.fills = [{type:'SOLID', color:{r:0.93,g:0.93,b:0.94}}];")
			b.Line("figma.currentPage.appendChild(row);")
			radii := []int{0, 2, 4, 8, 12, 16, 24, 32}
			for _, r := range radii {
				b.Linef("{ const sq = figma.createRectangle(); sq.resize(80, 80); sq.cornerRadius = %d; sq.name = 'r-%d';", r, r)
				b.Line("  sq.fills = [{type:'SOLID', color:{r:0.2,g:0.45,b:0.95}}]; row.appendChild(sq); }")
			}
			b.ReturnIDs("row.id")
			output(b.String())
			return nil
		},
	}
}

func newDSIconsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "icons",
		Short: "Create a placeholder grid for iconography slots",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const grid = figma.createFrame(); grid.name = 'Icons / Grid';")
			b.Line("grid.layoutMode = 'VERTICAL'; grid.itemSpacing = 8; grid.paddingLeft = grid.paddingRight = 24; grid.paddingTop = grid.paddingBottom = 24;")
			b.Line("grid.resize(520, 520); grid.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}];")
			b.Line("figma.currentPage.appendChild(grid);")
			b.Line("for (let r = 0; r < 6; r++) { const row = figma.createFrame(); row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 8; row.fills = [];")
			b.Line("  for (let c = 0; c < 6; c++) { const cell = figma.createFrame(); cell.resize(72, 72); cell.cornerRadius = 8;")
			b.Line("    cell.name = 'icon-' + r + '-' + c; cell.fills = [{type:'SOLID', color:{r:0.94,g:0.94,b:0.95}}];")
			b.Line("    cell.strokes = [{type:'SOLID', color:{r:0.85,g:0.85,b:0.88}}]; cell.strokeWeight = 1; row.appendChild(cell); }")
			b.Line("  grid.appendChild(row); }")
			b.ReturnIDs("grid.id")
			output(b.String())
			return nil
		},
	}
}

func newDSComponentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "component",
		Short: "Create a starter component with variant placeholders",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const set = figma.createComponentSet();")
			b.Line("set.name = 'Button'; set.layoutMode = 'HORIZONTAL'; set.itemSpacing = 16; set.paddingLeft = set.paddingRight = 24; set.paddingTop = set.paddingBottom = 16;")
			b.Line("set.resize(320, 120); set.fills = [{type:'SOLID', color:{r:0.96,g:0.96,b:0.97}}];")
			b.Line("figma.currentPage.appendChild(set);")
			b.Line("const a = figma.createComponent(); a.name = 'Property 1=default'; a.resize(120, 44); a.cornerRadius = 8;")
			b.Line("a.fills = [{type:'SOLID', color:{r:0.2,g:0.45,b:0.95}}]; set.appendChild(a);")
			b.Line("const b = figma.createComponent(); b.name = 'Property 1=hover'; b.resize(120, 44); b.cornerRadius = 8;")
			b.Line("b.fills = [{type:'SOLID', color:{r:0.15,g:0.35,b:0.85}}]; set.appendChild(b);")
			b.ReturnIDs("set.id")
			output(b.String())
			return nil
		},
	}
}

func newDSVariablesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "variables",
		Short: "List local variable collections (figma.variables API)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const cols = await figma.variables.getLocalVariableCollectionsAsync();")
			b.Line("const out = cols.map(c => ({ id: c.id, name: c.name, modes: c.modes.map(m => m.name) }));")
			b.Line("return { variableCollections: out };")
			output(b.String())
			return nil
		},
	}
}

func newDSSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "How to search the design system via MCP",
		RunE: func(cmd *cobra.Command, args []string) error {
			msg := strings.TrimSpace(`
Design system search is performed through the Figma MCP tool **search_design_system**.

1. In your MCP-enabled client, invoke the tool with a natural-language query (components, styles, variables).
2. Scope queries to your team library or file as supported by your MCP server configuration.
3. Use the returned node or style IDs with figma-kit commands (e.g. node clone, style fill) to apply results.

This CLI does not call MCP directly; run the tool from the assistant or IDE integration that exposes Figma MCP.
`)
			_, _ = fmt.Fprint(os.Stdout, msg, "\n")
			return nil
		},
	}
}

func newDSImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "Stub workflow for importing external tokens into variables (extend as needed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Comment("Import tokens: map theme colors to variables — customize collection/mode names")
			b.Line("const collections = await figma.variables.getLocalVariableCollectionsAsync();")
			b.Line("const col = collections[0];")
			b.Line("if (!col) throw new Error('Create a variable collection first, or extend this script.');")
			b.Line("const modeId = col.modes[0].modeId;")
			b.Comment("Example: create COLOR variable from theme constant primary — duplicate per token")
			b.Line("// const v = figma.variables.createVariable('primary', col, 'COLOR');")
			b.Line("// v.setValueForMode(modeId, { type: 'color', r: primary.r, g: primary.g, b: primary.b });")
			b.ReturnDone()
			output(b.String())
			return nil
		},
	}
}

func newDSSyncTokensCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "sync-tokens",
		Short: "Emit theme tokens as CSS variables, Tailwind extend, or JSON (no plugin JS)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			switch strings.ToLower(format) {
			case "css":
				_, _ = fmt.Fprintf(os.Stdout, "/* theme: %s */\n:root {\n", t.Name)
				for _, name := range t.ColorNames() {
					c := t.Colors[name]
					h := codegen.RGBToHex(codegen.RGB{R: c.R, G: c.G, B: c.B})
					key := strings.ReplaceAll(strings.ToLower(name), "_", "-")
					_, _ = fmt.Fprintf(os.Stdout, "  --color-%s: %s;\n", key, h)
				}
				_, _ = fmt.Fprintln(os.Stdout, "}")
			case "tailwind":
				_, _ = fmt.Fprintf(os.Stdout, "// tailwind.config theme.extend — %s\n", t.Name)
				_, _ = fmt.Fprintln(os.Stdout, "colors: {")
				for _, name := range t.ColorNames() {
					c := t.Colors[name]
					h := codegen.RGBToHex(codegen.RGB{R: c.R, G: c.G, B: c.B})
					_, _ = fmt.Fprintf(os.Stdout, "  '%s': '%s',\n", name, h)
				}
				_, _ = fmt.Fprintln(os.Stdout, "},")
			case "json":
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(t); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown format %q (use css, tailwind, or json)", format)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "css", "Output format: css, tailwind, json")
	return cmd
}

func newDSAuditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "audit",
		Short: "Scan nodes for solid fills that do not match theme palette (approximate)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			codegen.PreambleWithPage(b, t, resolvePage())
			b.Line("const palette = [")
			for _, name := range t.ColorNames() {
				c := t.Colors[name]
				b.Linef("  {r:%s,g:%s,b:%s},", codegen.FmtFloat(c.R), codegen.FmtFloat(c.G), codegen.FmtFloat(c.B))
			}
			b.Line("];")
			b.Line("function near(a,b,eps=0.02){ return Math.abs(a.r-b.r)<eps && Math.abs(a.g-b.g)<eps && Math.abs(a.b-b.b)<eps; }")
			b.Line("function inPalette(c){ return palette.some(p => near(p,c)); }")
			b.Line("const issues = [];")
			b.Line("function walk(n){")
			b.Line("  if ('fills' in n && Array.isArray(n.fills)) {")
			b.Line("    for (const f of n.fills) {")
			b.Line("      if (f.type === 'SOLID' && f.color && !inPalette(f.color))")
			b.Line("        issues.push({ id: n.id, name: n.name, type: n.type, color: f.color });")
			b.Line("    }")
			b.Line("  }")
			b.Line("  if ('children' in n) for (const c of n.children) walk(c);")
			b.Line("}")
			b.Line("walk(figma.currentPage);")
			b.Line("return { offPalette: issues };")
			output(b.String())
			return nil
		},
	}
}

func newDSTokensCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tokens",
		Short: "Print theme tokens as JSON (same as sync-tokens --format json)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(t)
		},
	}
}
