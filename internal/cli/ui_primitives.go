package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newUIChipCmd() *cobra.Command {
	var (
		label     string
		removable bool
		variant   string
		parent    string
	)
	cmd := &cobra.Command{
		Use:         "chip",
		Short:       "Tag/filter chip with optional close icon",
		Example:     `  figma-kit ui chip -t noir --label "React" --removable`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if label == "" {
				label = "Tag"
			}
			b.Line("const chip = figma.createFrame(); chip.name = 'Chip';")
			b.Line("chip.layoutMode = 'HORIZONTAL'; chip.itemSpacing = 6; chip.counterAxisAlignItems = 'CENTER';")
			b.Line("chip.paddingLeft = 12; chip.paddingRight = " + fmt.Sprintf("%d", func() int {
				if removable {
					return 8
				}
				return 12
			}()) + "; chip.paddingTop = chip.paddingBottom = 6;")
			b.Line("chip.cornerRadius = 999;")
			b.Line("chip.counterAxisSizingMode = 'AUTO'; chip.primaryAxisSizingMode = 'AUTO';")
			if variant == "outline" {
				b.Line("chip.fills = []; chip.strokes = [{type:'SOLID', color:STK}]; chip.strokeWeight = 1;")
			} else {
				b.Line("chip.fills = [{type:'SOLID', color:STK}];")
			}
			b.Linef("const tx = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); tx.fontName = {family:'Inter',style:'Medium'}; tx.fontSize = 13; tx.characters = %q; tx.fills = [{type:'SOLID', color:BD}]; chip.appendChild(tx);", label)
			if removable {
				b.Line("const x = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); x.fontName = {family:'Inter',style:'Regular'}; x.fontSize = 14; x.characters = '✕'; x.fills = [{type:'SOLID', color:{r:0.5,g:0.5,b:0.55}}]; chip.appendChild(x);")
			}
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(chip); else pg.appendChild(chip);")
			} else {
				b.Line("pg.appendChild(chip);")
			}
			b.ReturnIDs("chip.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "Chip label text")
	cmd.Flags().BoolVar(&removable, "removable", false, "Show close/remove icon")
	cmd.Flags().StringVar(&variant, "variant", "filled", "Variant: filled, outline")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIToastCmd() *cobra.Command {
	var (
		message     string
		toastType   string
		dismissible bool
		parent      string
	)
	cmd := &cobra.Command{
		Use:         "toast",
		Short:       "Notification toast popup",
		Example:     `  figma-kit ui toast -t noir --message "File saved!" --type success`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if message == "" {
				message = "Operation completed successfully"
			}

			var iconChar, bgColor string
			switch toastType {
			case "error":
				iconChar = "✕"
				bgColor = "{r:0.6,g:0.15,b:0.15}"
			case "warning":
				iconChar = "⚠"
				bgColor = "{r:0.55,g:0.45,b:0.1}"
			case "info":
				iconChar = "ℹ"
				bgColor = "{r:0.15,g:0.3,b:0.55}"
			default:
				iconChar = "✓"
				bgColor = "{r:0.1,g:0.45,b:0.25}"
			}

			b.Line("const toast = figma.createFrame(); toast.name = 'Toast';")
			b.Line("toast.layoutMode = 'HORIZONTAL'; toast.itemSpacing = 10; toast.counterAxisAlignItems = 'CENTER';")
			b.Line("toast.paddingLeft = 16; toast.paddingRight = " + fmt.Sprintf("%d", func() int {
				if dismissible {
					return 12
				}
				return 16
			}()) + "; toast.paddingTop = toast.paddingBottom = 12;")
			b.Line("toast.cornerRadius = 8;")
			b.Line("toast.counterAxisSizingMode = 'AUTO'; toast.primaryAxisSizingMode = 'AUTO';")
			b.Linef("toast.fills = [{type:'SOLID', color:%s}];", bgColor)
			b.Linef("const icon = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Bold'}); icon.fontName = {family:'Inter',style:'Bold'}; icon.fontSize = 16; icon.characters = %q; icon.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}]; toast.appendChild(icon);", iconChar)
			b.Linef("const msg = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); msg.fontName = {family:'Inter',style:'Medium'}; msg.fontSize = 14; msg.characters = %q; msg.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}]; toast.appendChild(msg);", message)
			if dismissible {
				b.Line("const close = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); close.fontName = {family:'Inter',style:'Regular'}; close.fontSize = 16; close.characters = '✕'; close.fills = [{type:'SOLID', color:{r:1,g:1,b:1,a:0.6}}]; toast.appendChild(close);")
			}
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(toast); else pg.appendChild(toast);")
			} else {
				b.Line("pg.appendChild(toast);")
			}
			b.ReturnIDs("toast.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&message, "message", "", "Toast message text")
	cmd.Flags().StringVar(&toastType, "type", "success", "Type: success, error, warning, info")
	cmd.Flags().BoolVar(&dismissible, "dismissible", true, "Show dismiss button")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIModalCmd() *cobra.Command {
	var (
		title   string
		body    string
		confirm string
		cancel  string
		width   int
		parent  string
	)
	cmd := &cobra.Command{
		Use:         "modal",
		Short:       "Modal dialog with title, body, and action buttons",
		Example:     `  figma-kit ui modal -t noir --title "Delete item?" --body "This cannot be undone." --confirm "Delete" --cancel "Cancel"`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if title == "" {
				title = "Confirm Action"
			}
			if body == "" {
				body = "Are you sure you want to proceed?"
			}
			if confirm == "" {
				confirm = "Confirm"
			}
			if cancel == "" {
				cancel = "Cancel"
			}

			b.Line("const overlay = figma.createFrame(); overlay.name = 'Modal Overlay';")
			b.Line("overlay.resize(1440, 900); overlay.fills = [{type:'SOLID', color:{r:0,g:0,b:0}, opacity:0.5}];")
			b.Line("overlay.layoutMode = 'VERTICAL'; overlay.primaryAxisAlignItems = 'CENTER'; overlay.counterAxisAlignItems = 'CENTER';")
			b.Linef("const modal = figma.createFrame(); modal.name = 'Modal'; modal.resize(%d, 240); modal.cornerRadius = 16;", width)
			b.Line("modal.fills = [{type:'SOLID', color:CARD}]; modal.strokes = [{type:'SOLID', color:STK}]; modal.strokeWeight = 1;")
			b.Line("modal.layoutMode = 'VERTICAL'; modal.itemSpacing = 16; modal.paddingLeft = modal.paddingRight = 28; modal.paddingTop = modal.paddingBottom = 24;")
			b.Linef("const tt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tt.fontName = {family:'Inter',style:'Semi Bold'}; tt.fontSize = 18; tt.characters = %q; tt.fills = [{type:'SOLID', color:WT}]; modal.appendChild(tt);", title)
			b.Linef("const bd = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); bd.fontName = {family:'Inter',style:'Regular'}; bd.fontSize = 14; bd.lineHeight = {unit:'PIXELS',value:22}; bd.characters = %q; bd.fills = [{type:'SOLID', color:MT}]; bd.textAutoResize = 'HEIGHT'; bd.resize(%d - 56, 1); modal.appendChild(bd);", body, width)
			b.Line("const actions = figma.createFrame(); actions.name = 'Actions'; actions.layoutMode = 'HORIZONTAL'; actions.itemSpacing = 12; actions.fills = []; actions.counterAxisSizingMode = 'AUTO'; actions.primaryAxisSizingMode = 'AUTO'; actions.primaryAxisAlignItems = 'MAX';")
			b.Linef("const cancelBtn = figma.createFrame(); cancelBtn.name = 'Cancel'; cancelBtn.layoutMode = 'HORIZONTAL'; cancelBtn.paddingLeft = cancelBtn.paddingRight = 20; cancelBtn.paddingTop = cancelBtn.paddingBottom = 10; cancelBtn.cornerRadius = 8; cancelBtn.fills = [{type:'SOLID', color:STK}]; cancelBtn.counterAxisSizingMode = 'AUTO'; cancelBtn.primaryAxisSizingMode = 'AUTO';")
			b.Linef("const ct1 = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); ct1.fontName = {family:'Inter',style:'Medium'}; ct1.fontSize = 14; ct1.characters = %q; ct1.fills = [{type:'SOLID', color:BD}]; cancelBtn.appendChild(ct1); actions.appendChild(cancelBtn);", cancel)
			b.Linef("const confirmBtn = figma.createFrame(); confirmBtn.name = 'Confirm'; confirmBtn.layoutMode = 'HORIZONTAL'; confirmBtn.paddingLeft = confirmBtn.paddingRight = 20; confirmBtn.paddingTop = confirmBtn.paddingBottom = 10; confirmBtn.cornerRadius = 8; confirmBtn.fills = [{type:'SOLID', color:BL}]; confirmBtn.counterAxisSizingMode = 'AUTO'; confirmBtn.primaryAxisSizingMode = 'AUTO';")
			b.Linef("const ct2 = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); ct2.fontName = {family:'Inter',style:'Semi Bold'}; ct2.fontSize = 14; ct2.characters = %q; ct2.fills = [{type:'SOLID', color:{r:1,g:1,b:1}}]; confirmBtn.appendChild(ct2); actions.appendChild(confirmBtn);", confirm)
			b.Line("modal.appendChild(actions); overlay.appendChild(modal);")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(overlay); else pg.appendChild(overlay);")
			} else {
				b.Line("pg.appendChild(overlay);")
			}
			b.ReturnIDs("overlay.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Modal title")
	cmd.Flags().StringVar(&body, "body", "", "Modal body text")
	cmd.Flags().StringVar(&confirm, "confirm", "", "Confirm button text")
	cmd.Flags().StringVar(&cancel, "cancel", "", "Cancel button text")
	cmd.Flags().IntVar(&width, "width", 420, "Modal width")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUICardListCmd() *cobra.Command {
	var (
		itemsJSON string
		cardType  string
		parent    string
	)
	cmd := &cobra.Command{
		Use:         "card-list",
		Short:       "Generate a list of cards from data",
		Example:     `  figma-kit ui card-list -t noir --items '[{"title":"Item 1","desc":"First"},{"title":"Item 2","desc":"Second"}]'`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if itemsJSON == "" {
				itemsJSON = `[{"title":"Dashboard","desc":"View analytics"},{"title":"Settings","desc":"Configure preferences"},{"title":"Users","desc":"Manage team members"}]`
			}
			b.Linef("const items = JSON.parse(%s);", jsStringLiteral(itemsJSON))
			b.Line("const list = figma.createFrame(); list.name = 'Card List'; list.layoutMode = 'VERTICAL'; list.itemSpacing = 16; list.paddingLeft = list.paddingRight = 24; list.paddingTop = list.paddingBottom = 24; list.fills = []; list.counterAxisSizingMode = 'AUTO'; list.primaryAxisSizingMode = 'AUTO';")
			b.Line("for (const item of items) {")
			b.Line("  const card = figma.createFrame(); card.name = item.title; card.layoutMode = 'VERTICAL'; card.itemSpacing = 8; card.paddingLeft = card.paddingRight = 20; card.paddingTop = card.paddingBottom = 16; card.resize(320, 100); card.cornerRadius = 12;")
			switch cardType {
			case "glass":
				b.Line("  card.fills = [{type:'SOLID', color:{r:1,g:1,b:1}, opacity:0.08}]; card.effects = [{type:'BACKGROUND_BLUR', radius:20, visible:true}];")
			case "outline":
				b.Line("  card.fills = []; card.strokes = [{type:'SOLID', color:STK}]; card.strokeWeight = 1;")
			default:
				b.Line("  card.fills = [{type:'SOLID', color:CARD}];")
			}
			b.Line("  const tt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Semi Bold'}); tt.fontName = {family:'Inter',style:'Semi Bold'}; tt.fontSize = 16; tt.characters = item.title; tt.fills = [{type:'SOLID', color:WT}]; card.appendChild(tt);")
			b.Line("  if (item.desc) { const dd = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); dd.fontName = {family:'Inter',style:'Regular'}; dd.fontSize = 13; dd.characters = item.desc; dd.fills = [{type:'SOLID', color:MT}]; card.appendChild(dd); }")
			b.Line("  list.appendChild(card);")
			b.Line("}")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(list); else pg.appendChild(list);")
			} else {
				b.Line("pg.appendChild(list);")
			}
			b.ReturnIDs("list.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&itemsJSON, "items", "", "Items as JSON array [{title,desc}]")
	cmd.Flags().StringVar(&cardType, "card-type", "solid", "Card type: solid, glass, outline")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUISidebarCmd() *cobra.Command {
	var (
		sectionsJSON string
		parent       string
	)
	cmd := &cobra.Command{
		Use:         "sidebar",
		Short:       "Sidebar navigation with sections and items",
		Example:     `  figma-kit ui sidebar -t noir --sections '[{"title":"Main","items":[{"label":"Dashboard","active":true},{"label":"Analytics"}]},{"title":"Settings","items":[{"label":"Profile"},{"label":"Team"}]}]'`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if sectionsJSON == "" {
				sectionsJSON = `[{"title":"Main","items":[{"label":"Dashboard","active":true},{"label":"Analytics"},{"label":"Projects"}]},{"title":"Settings","items":[{"label":"Profile"},{"label":"Team"},{"label":"Billing"}]}]`
			}
			b.Linef("const sections = JSON.parse(%s);", jsStringLiteral(sectionsJSON))
			b.Line("const sidebar = figma.createFrame(); sidebar.name = 'Sidebar'; sidebar.layoutMode = 'VERTICAL'; sidebar.itemSpacing = 24; sidebar.paddingLeft = sidebar.paddingRight = 16; sidebar.paddingTop = sidebar.paddingBottom = 24;")
			b.Line("sidebar.resize(240, 600); sidebar.fills = [{type:'SOLID', color:BG}]; sidebar.strokes = [{type:'SOLID', color:STK}]; sidebar.strokeWeight = 1;")
			b.Line("for (const section of sections) {")
			b.Line("  const sec = figma.createFrame(); sec.name = section.title; sec.layoutMode = 'VERTICAL'; sec.itemSpacing = 2; sec.fills = []; sec.counterAxisSizingMode = 'AUTO'; sec.primaryAxisSizingMode = 'AUTO';")
			b.Line("  const hdr = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); hdr.fontName = {family:'Inter',style:'Medium'}; hdr.fontSize = 11; hdr.characters = section.title.toUpperCase(); hdr.fills = [{type:'SOLID', color:MT}]; hdr.letterSpacing = {value:1, unit:'PIXELS'}; sec.appendChild(hdr);")
			b.Line("  for (const item of section.items) {")
			b.Line("    const row = figma.createFrame(); row.name = item.label; row.layoutMode = 'HORIZONTAL'; row.itemSpacing = 10; row.paddingLeft = 12; row.paddingRight = 12; row.paddingTop = row.paddingBottom = 8; row.cornerRadius = 6; row.counterAxisSizingMode = 'AUTO'; row.primaryAxisSizingMode = 'AUTO';")
			b.Line("    if (item.active) row.fills = [{type:'SOLID', color:STK}]; else row.fills = [];")
			b.Line("    const lbl = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); lbl.fontName = {family:'Inter',style:'Regular'}; lbl.fontSize = 14;")
			b.Line("    lbl.characters = item.label; lbl.fills = [{type:'SOLID', color: item.active ? WT : MT}];")
			b.Line("    row.appendChild(lbl); sec.appendChild(row);")
			b.Line("  }")
			b.Line("  sidebar.appendChild(sec);")
			b.Line("}")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(sidebar); else pg.appendChild(sidebar);")
			} else {
				b.Line("pg.appendChild(sidebar);")
			}
			b.ReturnIDs("sidebar.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&sectionsJSON, "sections", "", "Sections as JSON array")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIAvatarGroupCmd() *cobra.Command {
	var (
		count   int
		maxShow int
		size    int
		parent  string
	)
	cmd := &cobra.Command{
		Use:         "avatar-group",
		Short:       "Overlapping avatar stack with +N overflow",
		Example:     `  figma-kit ui avatar-group -t noir --count 8 --max 4 --size 36`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			show := maxShow
			if show > count {
				show = count
			}
			b.Linef("const count = %d; const show = %d; const sz = %d;", count, show, size)
			b.Line("const group = figma.createFrame(); group.name = 'Avatar Group'; group.layoutMode = 'HORIZONTAL'; group.itemSpacing = " + fmt.Sprintf("%d", -size/4) + "; group.fills = []; group.counterAxisSizingMode = 'AUTO'; group.primaryAxisSizingMode = 'AUTO';")
			b.Line("const colors = [{r:0.23,g:0.51,b:0.96},{r:0.55,g:0.36,b:0.96},{r:0.1,g:0.7,b:0.5},{r:0.9,g:0.4,b:0.3},{r:0.85,g:0.6,b:0.1}];")
			b.Line("for (let i = 0; i < show; i++) {")
			b.Line("  const av = figma.createEllipse(); av.resize(sz, sz);")
			b.Line("  av.fills = [{type:'SOLID', color:colors[i % colors.length]}];")
			b.Line("  av.strokes = [{type:'SOLID', color:BG}]; av.strokeWeight = 2;")
			b.Line("  av.name = 'Avatar ' + (i + 1); group.appendChild(av);")
			b.Line("}")
			b.Line("if (count > show) {")
			b.Line("  const overflow = figma.createEllipse(); overflow.resize(sz, sz);")
			b.Line("  overflow.fills = [{type:'SOLID', color:STK}];")
			b.Line("  overflow.strokes = [{type:'SOLID', color:BG}]; overflow.strokeWeight = 2;")
			b.Line("  overflow.name = '+' + (count - show); group.appendChild(overflow);")
			b.Line("}")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(group); else pg.appendChild(group);")
			} else {
				b.Line("pg.appendChild(group);")
			}
			b.ReturnIDs("group.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&count, "count", 6, "Total number of avatars")
	cmd.Flags().IntVar(&maxShow, "max", 4, "Maximum visible avatars")
	cmd.Flags().IntVar(&size, "size", 36, "Avatar size in pixels")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIRatingCmd() *cobra.Command {
	var (
		value  float64
		size   int
		color  string
		parent string
	)
	cmd := &cobra.Command{
		Use:         "rating",
		Short:       "Star rating display",
		Example:     `  figma-kit ui rating --value 4.5 --size 24`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			b.Line("const rating = figma.createFrame(); rating.name = 'Rating'; rating.layoutMode = 'HORIZONTAL'; rating.itemSpacing = 4; rating.fills = []; rating.counterAxisSizingMode = 'AUTO'; rating.primaryAxisSizingMode = 'AUTO';")
			filled := int(value)
			hasHalf := value-float64(filled) >= 0.5
			for i := 0; i < 5; i++ {
				if i < filled {
					b.Linef("{ const s = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); s.fontName = {family:'Inter',style:'Regular'}; s.fontSize = %d; s.characters = '★'; s.fills = [{type:'SOLID', color:{r:0.96,g:0.76,b:0.05}}]; rating.appendChild(s); }", size)
				} else if i == filled && hasHalf {
					b.Linef("{ const s = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); s.fontName = {family:'Inter',style:'Regular'}; s.fontSize = %d; s.characters = '★'; s.fills = [{type:'SOLID', color:{r:0.96,g:0.76,b:0.05}, opacity:0.5}]; rating.appendChild(s); }", size)
				} else {
					b.Linef("{ const s = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); s.fontName = {family:'Inter',style:'Regular'}; s.fontSize = %d; s.characters = '☆'; s.fills = [{type:'SOLID', color:MT}]; rating.appendChild(s); }", size)
				}
			}
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(rating); else pg.appendChild(rating);")
			} else {
				b.Line("pg.appendChild(rating);")
			}
			b.ReturnIDs("rating.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().Float64Var(&value, "value", 4, "Rating value (1-5, supports halves)")
	cmd.Flags().IntVar(&size, "size", 20, "Star size")
	cmd.Flags().StringVar(&color, "color", "#F5C205", "Star color hex")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUISearchCmd() *cobra.Command {
	var (
		placeholder string
		withClear   bool
		size        string
		parent      string
	)
	cmd := &cobra.Command{
		Use:         "search",
		Short:       "Search input with magnifying glass icon",
		Example:     `  figma-kit ui search -t noir --placeholder "Search components..."`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if placeholder == "" {
				placeholder = "Search..."
			}
			h := 40
			if size == "lg" {
				h = 48
			}
			if size == "sm" {
				h = 32
			}
			b.Line("const search = figma.createFrame(); search.name = 'Search';")
			b.Linef("search.resize(320, %d); search.layoutMode = 'HORIZONTAL'; search.itemSpacing = 8; search.counterAxisAlignItems = 'CENTER';", h)
			b.Linef("search.paddingLeft = 12; search.paddingRight = %d; search.paddingTop = search.paddingBottom = 8;", func() int {
				if withClear {
					return 8
				}
				return 12
			}())
			b.Line("search.cornerRadius = 8; search.fills = [{type:'SOLID', color:CARD}]; search.strokes = [{type:'SOLID', color:STK}]; search.strokeWeight = 1;")
			b.Line("const icon = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); icon.fontName = {family:'Inter',style:'Regular'}; icon.fontSize = 16; icon.characters = '🔍'; icon.fills = [{type:'SOLID', color:MT}]; search.appendChild(icon);")
			b.Linef("const input = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); input.fontName = {family:'Inter',style:'Regular'}; input.fontSize = 14; input.characters = %q; input.fills = [{type:'SOLID', color:MT}]; search.appendChild(input);", placeholder)
			if withClear {
				b.Line("const clear = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); clear.fontName = {family:'Inter',style:'Regular'}; clear.fontSize = 14; clear.characters = '✕'; clear.fills = [{type:'SOLID', color:MT}]; search.appendChild(clear);")
			}
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(search); else pg.appendChild(search);")
			} else {
				b.Line("pg.appendChild(search);")
			}
			b.ReturnIDs("search.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&placeholder, "placeholder", "", "Placeholder text")
	cmd.Flags().BoolVar(&withClear, "with-clear", false, "Show clear button")
	cmd.Flags().StringVar(&size, "size", "md", "Size: sm, md, lg")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIPaginationCmd() *cobra.Command {
	var (
		total      int
		current    int
		maxVisible int
		parent     string
	)
	cmd := &cobra.Command{
		Use:         "pagination",
		Short:       "Page number pagination bar",
		Example:     `  figma-kit ui pagination --total 10 --current 3 --max-visible 5`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			b.Linef("const total = %d; const current = %d; const maxVis = %d;", total, current, maxVisible)
			b.Line("const bar = figma.createFrame(); bar.name = 'Pagination'; bar.layoutMode = 'HORIZONTAL'; bar.itemSpacing = 4; bar.fills = []; bar.counterAxisAlignItems = 'CENTER'; bar.counterAxisSizingMode = 'AUTO'; bar.primaryAxisSizingMode = 'AUTO';")
			b.Line("const prev = figma.createFrame(); prev.name = 'Prev'; prev.resize(36, 36); prev.cornerRadius = 8; prev.fills = [{type:'SOLID', color:STK}]; prev.layoutMode = 'HORIZONTAL'; prev.primaryAxisAlignItems = 'CENTER'; prev.counterAxisAlignItems = 'CENTER';")
			b.Line("const pt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); pt.fontName = {family:'Inter',style:'Regular'}; pt.fontSize = 14; pt.characters = '‹'; pt.fills = [{type:'SOLID', color:MT}]; prev.appendChild(pt); bar.appendChild(prev);")
			b.Line("const start = Math.max(1, current - Math.floor(maxVis / 2));")
			b.Line("const end = Math.min(total, start + maxVis - 1);")
			b.Line("for (let i = start; i <= end; i++) {")
			b.Line("  const pg = figma.createFrame(); pg.name = 'Page ' + i; pg.resize(36, 36); pg.cornerRadius = 8; pg.layoutMode = 'HORIZONTAL'; pg.primaryAxisAlignItems = 'CENTER'; pg.counterAxisAlignItems = 'CENTER';")
			b.Line("  if (i === current) pg.fills = [{type:'SOLID', color:BL}]; else pg.fills = [{type:'SOLID', color:STK}];")
			b.Line("  const t = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Medium'}); t.fontName = {family:'Inter',style:'Medium'}; t.fontSize = 14; t.characters = String(i);")
			b.Line("  t.fills = [{type:'SOLID', color: i === current ? {r:1,g:1,b:1} : MT}]; pg.appendChild(t); bar.appendChild(pg);")
			b.Line("}")
			b.Line("const next = figma.createFrame(); next.name = 'Next'; next.resize(36, 36); next.cornerRadius = 8; next.fills = [{type:'SOLID', color:STK}]; next.layoutMode = 'HORIZONTAL'; next.primaryAxisAlignItems = 'CENTER'; next.counterAxisAlignItems = 'CENTER';")
			b.Line("const nt = figma.createText(); await figma.loadFontAsync({family:'Inter',style:'Regular'}); nt.fontName = {family:'Inter',style:'Regular'}; nt.fontSize = 14; nt.characters = '›'; nt.fills = [{type:'SOLID', color:MT}]; next.appendChild(nt); bar.appendChild(next);")
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Line("if (_par && 'appendChild' in _par) _par.appendChild(bar); else figma.currentPage.appendChild(bar);")
			} else {
				b.Line("figma.currentPage.appendChild(bar);")
			}
			b.ReturnIDs("bar.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&total, "total", 10, "Total pages")
	cmd.Flags().IntVar(&current, "current", 1, "Current page")
	cmd.Flags().IntVar(&maxVisible, "max-visible", 5, "Max visible page numbers")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}

func newUIColorPickerCmd() *cobra.Command {
	var (
		colorsCSV string
		selected  int
		parent    string
	)
	cmd := &cobra.Command{
		Use:         "color-picker",
		Short:       "Color swatch grid for palette selection",
		Example:     `  figma-kit ui color-picker --colors "#EF4444,#F59E0B,#10B981,#3B82F6,#8B5CF6,#EC4899"`,
		Annotations: map[string]string{"composable": "true"},
		RunE: func(cmd *cobra.Command, args []string) error {
			page := resolvePage()
			b := newBuilder()
			b.PageSetup(page)
			t, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			emitUIThemeTokens(b, t)
			if colorsCSV == "" {
				colorsCSV = "#EF4444,#F59E0B,#10B981,#3B82F6,#8B5CF6,#EC4899,#06B6D4,#84CC16,#F97316,#6366F1"
			}
			colors := strings.Split(colorsCSV, ",")
			b.Line("const grid = figma.createFrame(); grid.name = 'Color Picker'; grid.layoutMode = 'HORIZONTAL'; grid.itemSpacing = 8; grid.paddingLeft = grid.paddingRight = 12; grid.paddingTop = grid.paddingBottom = 12; grid.cornerRadius = 12; grid.fills = [{type:'SOLID', color:CARD}]; grid.counterAxisSizingMode = 'AUTO'; grid.primaryAxisSizingMode = 'AUTO';")
			for i, hex := range colors {
				hex = strings.TrimSpace(hex)
				rgb, err := codegen.HexToRGB(hex)
				if err != nil {
					continue
				}
				b.Linef("{ const sw = figma.createEllipse(); sw.name = %q; sw.resize(28, 28); sw.fills = [{type:'SOLID', color:{r:%.3f,g:%.3f,b:%.3f}}];", hex, rgb.R, rgb.G, rgb.B)
				if i == selected {
					b.Line("  sw.strokes = [{type:'SOLID', color:{r:1,g:1,b:1}}]; sw.strokeWeight = 2; sw.effects = [{type:'DROP_SHADOW', color:{r:0,g:0,b:0,a:0.3}, offset:{x:0,y:0}, radius:4, spread:0, visible:true, blendMode:'NORMAL'}];")
				}
				b.Line("  grid.appendChild(sw); }")
			}
			if parent != "" {
				if strings.HasPrefix(parent, "_results[") {
					b.Linef("const _par = %s;", parent)
				} else {
					b.Linef("const _par = await figma.getNodeByIdAsync(%q);", parent)
				}
				b.Linef("if (_par && 'appendChild' in _par) _par.appendChild(grid); else pg.appendChild(grid);")
			} else {
				b.Line("pg.appendChild(grid);")
			}
			b.ReturnIDs("grid.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&colorsCSV, "colors", "", "Comma-separated hex colors")
	cmd.Flags().IntVar(&selected, "selected", -1, "Index of selected color (-1 for none)")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent node ID (or _results[N] in compose)")
	return cmd
}
