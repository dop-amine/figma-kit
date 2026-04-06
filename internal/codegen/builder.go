package codegen

import (
	"fmt"
	"strings"
)

// Builder composes use_figma-compatible JavaScript output.
// Use the fluent API to add sections, then call String() for the final script.
type Builder struct {
	buf strings.Builder
}

// New creates a new Builder.
func New() *Builder {
	return &Builder{}
}

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
	b.buf.WriteString(fmt.Sprintf(format, args...))
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
func (b *Builder) PageSetup(pageIndex int) *Builder {
	b.Linef("const pg = figma.root.children[%d];", pageIndex)
	b.Line("await figma.setCurrentPageAsync(pg);")
	b.Blank()
	return b
}

// FontLoading emits the standard Inter + Geist Mono font loading block.
func (b *Builder) FontLoading() *Builder {
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

// ThemeColors emits theme color constants from a map of name -> RGB.
func (b *Builder) ThemeColors(colors map[string]struct{ R, G, B float64 }) *Builder {
	for name, c := range colors {
		b.Linef("const %s={r:%s,g:%s,b:%s};", name, fmtFloat(c.R), fmtFloat(c.G), fmtFloat(c.B))
	}
	b.Blank()
	return b
}

// ThemeColorsOrdered emits theme color constants in a specified order.
func (b *Builder) ThemeColorsOrdered(names []string, colors map[string][3]float64) *Builder {
	for _, name := range names {
		c := colors[name]
		b.Linef("const %s={r:%s,g:%s,b:%s};", name, fmtFloat(c[0]), fmtFloat(c[1]), fmtFloat(c[2]))
	}
	b.Blank()
	return b
}

// ReturnIDs emits a return statement with the given variable names as created node IDs.
func (b *Builder) ReturnIDs(varNames ...string) *Builder {
	if len(varNames) == 1 {
		b.Linef("return { createdNodeIds: [%s] };", varNames[0])
	} else {
		b.Linef("return { createdNodeIds: [%s] };", strings.Join(varNames, ", "))
	}
	return b
}

// ReturnDone emits a simple return 'Done' statement.
func (b *Builder) ReturnDone() *Builder {
	b.Line("return 'Done';")
	return b
}

// ReturnExpr emits a return statement with an arbitrary expression.
func (b *Builder) ReturnExpr(expr string) *Builder {
	b.Linef("return %s;", expr)
	return b
}

// String returns the final composed JavaScript string.
func (b *Builder) String() string {
	return b.buf.String()
}

// Len returns the current byte length of the buffer.
func (b *Builder) Len() int {
	return b.buf.Len()
}
