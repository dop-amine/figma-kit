package codegen

import (
	"fmt"
	"strings"
)

// Builder composes use_figma-compatible JavaScript output.
// Use the fluent API to add sections, then call String() for the final script.
type Builder struct {
	buf      strings.Builder
	bodyOnly bool
}

// New creates a new Builder.
func New() *Builder {
	return &Builder{}
}

// SetBodyOnly enables body-only mode. When true, preamble methods (PageSetup,
// FontLoading, ThemeColors*) and return methods (ReturnIDs, ReturnDone,
// ReturnExpr) become no-ops. Used by the compose engine to avoid duplicate
// preambles when batching multiple commands.
func (b *Builder) SetBodyOnly(v bool) { b.bodyOnly = v }

// IsBodyOnly returns true if the builder is in body-only mode.
func (b *Builder) IsBodyOnly() bool { return b.bodyOnly }

// Comment appends a JS comment line.
func (b *Builder) Comment(text string) *Builder {
	b.buf.WriteString("// ")
	b.buf.WriteString(text)
	b.buf.WriteByte('\n')
	return b
}

// Line appends a raw line of JS.
func (b *Builder) Line(line string) *Builder {
	b.buf.WriteString(line)
	b.buf.WriteByte('\n')
	return b
}

// Linef appends a formatted line of JS.
func (b *Builder) Linef(format string, args ...any) *Builder {
	_, _ = fmt.Fprintf(&b.buf, format, args...)
	b.buf.WriteByte('\n')
	return b
}

// Blank appends an empty line.
func (b *Builder) Blank() *Builder {
	b.buf.WriteByte('\n')
	return b
}

// Raw appends a raw block of JS (multi-line) without modification.
func (b *Builder) Raw(js string) *Builder {
	b.buf.WriteString(js)
	if !strings.HasSuffix(js, "\n") {
		b.buf.WriteByte('\n')
	}
	return b
}

// PageSetup emits the standard page selection and activation boilerplate.
// No-op when bodyOnly is true.
func (b *Builder) PageSetup(pageIndex int) *Builder {
	if b.bodyOnly {
		return b
	}
	b.Linef("const pg = figma.root.children[%d];", pageIndex)
	b.Line("await figma.setCurrentPageAsync(pg);")
	b.Blank()
	return b
}

// FontLoading emits the standard Inter + Geist Mono font loading block.
// No-op when bodyOnly is true.
func (b *Builder) FontLoading() *Builder {
	if b.bodyOnly {
		return b
	}
	b.Line("const fonts = [")
	b.Line("  {family:'Inter',style:'Bold'},{family:'Inter',style:'Semi Bold'},")
	b.Line("  {family:'Inter',style:'Medium'},{family:'Inter',style:'Regular'},")
	b.Line("  {family:'Inter',style:'Light'},")
	b.Line("  {family:'Geist Mono',style:'Regular'},{family:'Geist Mono',style:'Medium'}")
	b.Line("];")
	b.Line("for (const fn of fonts) await figma.loadFontAsync(fn);")
	b.Blank()
	return b
}

// FontLoadingFromTheme generates font loading from theme font families and weights.
// Falls back to standard Inter + Geist Mono when theme fonts are empty.
// No-op when bodyOnly is true.
func (b *Builder) FontLoadingFromTheme(heading, body, mono string, weights []string) *Builder {
	if b.bodyOnly {
		return b
	}
	if heading == "" {
		heading = "Inter"
	}
	if body == "" {
		body = "Inter"
	}
	if mono == "" {
		mono = "Geist Mono"
	}
	if len(weights) == 0 {
		weights = []string{"Bold", "Semi Bold", "Medium", "Regular", "Light"}
	}

	families := []string{heading}
	if body != heading {
		families = append(families, body)
	}

	b.Line("const fonts = [")
	for _, fam := range families {
		for _, w := range weights {
			b.Linef("  {family:%q,style:%q},", fam, w)
		}
	}
	b.Linef("  {family:%q,style:'Regular'},{family:%q,style:'Medium'}", mono, mono)
	b.Line("];")
	b.Line("for (const fn of fonts) await figma.loadFontAsync(fn);")
	b.Blank()
	return b
}

// ThemeColors emits theme color constants from a map of name -> RGB.
// No-op when bodyOnly is true.
func (b *Builder) ThemeColors(colors map[string]struct{ R, G, B float64 }) *Builder {
	if b.bodyOnly {
		return b
	}
	for name, c := range colors {
		b.Linef("const %s={r:%s,g:%s,b:%s};", name, fmtFloat(c.R), fmtFloat(c.G), fmtFloat(c.B))
	}
	b.Blank()
	return b
}

// ThemeColorsOrdered emits theme color constants in a specified order.
// No-op when bodyOnly is true.
func (b *Builder) ThemeColorsOrdered(names []string, colors map[string][3]float64) *Builder {
	if b.bodyOnly {
		return b
	}
	for _, name := range names {
		c := colors[name]
		b.Linef("const %s={r:%s,g:%s,b:%s};", name, fmtFloat(c[0]), fmtFloat(c[1]), fmtFloat(c[2]))
	}
	b.Blank()
	return b
}

// ReturnIDs emits a return statement with the given variable names as created node IDs.
// No-op when bodyOnly is true.
func (b *Builder) ReturnIDs(varNames ...string) *Builder {
	if b.bodyOnly {
		return b
	}
	if len(varNames) == 1 {
		b.Linef("return { createdNodeIds: [%s] };", varNames[0])
	} else {
		b.Linef("return { createdNodeIds: [%s] };", strings.Join(varNames, ", "))
	}
	return b
}

// ReturnDone emits a simple return 'Done' statement.
// No-op when bodyOnly is true.
func (b *Builder) ReturnDone() *Builder {
	if b.bodyOnly {
		return b
	}
	b.Line("return 'Done';")
	return b
}

// ReturnExpr emits a return statement with an arbitrary expression.
// No-op when bodyOnly is true.
func (b *Builder) ReturnExpr(expr string) *Builder {
	if b.bodyOnly {
		return b
	}
	b.Linef("return %s;", expr)
	return b
}

// ---------------------------------------------------------------------------
// Library import helpers — emit Plugin API calls for published assets
// ---------------------------------------------------------------------------

// ImportComponent emits JS to import a published component by key and create an instance.
// varName is used as the JS variable prefix (e.g. "hero" → heroComp, hero).
func (b *Builder) ImportComponent(key, varName string) *Builder {
	b.Linef("const %sComp = await figma.importComponentByKeyAsync(%q);", varName, key)
	b.Linef("const %s = %sComp.createInstance();", varName, varName)
	return b
}

// ImportComponentSet emits JS to import a component set by key, find the variant
// matching the given property string ("Size=Large,State=Default"), and create an instance.
func (b *Builder) ImportComponentSet(key, variantProps, varName string) *Builder {
	b.Linef("const %sSet = await figma.importComponentSetByKeyAsync(%q);", varName, key)
	b.Linef("const %sVariant = %sSet.children.find(c => {", varName, varName)
	b.Linef("  const p = c.variantProperties || {};")
	b.Linef("  return %s;", buildVariantMatch(variantProps, "p"))
	b.Linef("}) || %sSet.defaultVariant || %sSet.children[0];", varName, varName)
	b.Linef("const %s = %sVariant.createInstance();", varName, varName)
	return b
}

// ImportStyle emits JS to import a published style by key.
func (b *Builder) ImportStyle(key, varName string) *Builder {
	b.Linef("const %s = await figma.importStyleByKeyAsync(%q);", varName, key)
	return b
}

// buildVariantMatch converts "Size=Large,State=Default" into a JS boolean expression
// that checks each property: p["Size"]==="Large" && p["State"]==="Default".
func buildVariantMatch(props, jsVar string) string {
	if props == "" {
		return "true"
	}
	parts := strings.Split(props, ",")
	checks := make([]string, 0, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) == 2 {
			checks = append(checks, fmt.Sprintf(`%s[%q]===%q`, jsVar, strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])))
		}
	}
	if len(checks) == 0 {
		return "true"
	}
	return strings.Join(checks, " && ")
}

// String returns the final composed JavaScript string.
func (b *Builder) String() string {
	return b.buf.String()
}

// Len returns the current byte length of the buffer.
func (b *Builder) Len() int {
	return b.buf.Len()
}
