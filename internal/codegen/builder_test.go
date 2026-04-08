package codegen

import (
	"strings"
	"testing"
)

func TestNew_LenZero(t *testing.T) {
	b := New()
	if b.Len() != 0 {
		t.Errorf("New().Len() = %d, want 0", b.Len())
	}
	if got := b.String(); got != "" {
		t.Errorf("New().String() = %q, want empty", got)
	}
}

func TestBuilder_Comment(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"simple", "hello", "// hello\n"},
		{"empty", "", "// \n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New().Comment(tt.text).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Comment(%q).String() = %q, want to contain %q", tt.text, got, tt.want)
			}
		})
	}
}

func TestBuilder_Line(t *testing.T) {
	got := New().Line("const x = 1;").String()
	if !strings.Contains(got, "const x = 1;\n") {
		t.Errorf("Line().String() = %q, want substring %q", got, "const x = 1;\n")
	}
}

func TestBuilder_Linef(t *testing.T) {
	got := New().Linef("const n = %d;", 42).String()
	if !strings.Contains(got, "const n = 42;\n") {
		t.Errorf("Linef().String() = %q, want substring %q", got, "const n = 42;\n")
	}
}

func TestBuilder_Blank(t *testing.T) {
	got := New().Blank().String()
	if !strings.Contains(got, "\n") || got != "\n" {
		t.Errorf("Blank().String() = %q, want single newline", got)
	}
}

func TestBuilder_Raw(t *testing.T) {
	t.Run("already has newline", func(t *testing.T) {
		in := "foo\n"
		got := New().Raw(in).String()
		if got != "foo\n" {
			t.Errorf("Raw(%q).String() = %q, want %q", in, got, in)
		}
	})
	t.Run("adds trailing newline", func(t *testing.T) {
		in := "bar"
		got := New().Raw(in).String()
		if got != "bar\n" {
			t.Errorf("Raw(%q).String() = %q, want %q", in, got, "bar\n")
		}
	})
}

func TestBuilder_PageSetup(t *testing.T) {
	got := New().PageSetup(2).String()
	if !strings.Contains(got, "const pg = figma.root.children[2];") {
		t.Errorf("missing page const: %q", got)
	}
	if !strings.Contains(got, "await figma.setCurrentPageAsync(pg);") {
		t.Errorf("missing setCurrentPageAsync: %q", got)
	}
}

func TestBuilder_FontLoading(t *testing.T) {
	got := New().FontLoading().String()
	must := []string{
		"const fonts = [",
		"{family:'Inter',style:'Bold'},{family:'Inter',style:'Semi Bold'},",
		"{family:'Inter',style:'Medium'},{family:'Inter',style:'Regular'},",
		"{family:'Inter',style:'Light'},",
		"{family:'Geist Mono',style:'Regular'},{family:'Geist Mono',style:'Medium'}",
		"];",
		"for (const fn of fonts) await figma.loadFontAsync(fn);",
	}
	for _, s := range must {
		if !strings.Contains(got, s) {
			t.Errorf("FontLoading() output missing %q in:\n%s", s, got)
		}
	}
}

func TestBuilder_ReturnIDs(t *testing.T) {
	t.Run("single arg", func(t *testing.T) {
		got := New().ReturnIDs("frame").String()
		want := "return { createdNodeIds: [frame] };"
		if !strings.Contains(got, want) {
			t.Errorf("ReturnIDs(frame) = %q, want to contain %q", got, want)
		}
	})
	t.Run("multiple args", func(t *testing.T) {
		got := New().ReturnIDs("a", "b", "c").String()
		want := "return { createdNodeIds: [a, b, c] };"
		if !strings.Contains(got, want) {
			t.Errorf("ReturnIDs(a,b,c) = %q, want to contain %q", got, want)
		}
	})
}

func TestBuilder_ReturnDone(t *testing.T) {
	got := New().ReturnDone().String()
	if !strings.Contains(got, "return 'Done';\n") {
		t.Errorf("ReturnDone() = %q, want to contain %q", got, "return 'Done';\n")
	}
}

func TestBuilder_ChainOrder(t *testing.T) {
	got := New().
		Comment("start").
		Line("x++;").
		Linef("y = %q;", "ok").
		Blank().
		Raw("z();").
		String()

	want := "// start\nx++;\ny = \"ok\";\n\nz();\n"
	if got != want {
		t.Errorf("chained output mismatch\ngot:  %q\nwant: %q", got, want)
	}
	if !strings.Contains(got, "// start") || !strings.Contains(got, "x++;") ||
		!strings.Contains(got, `y = "ok";`) || !strings.Contains(got, "z();") {
		t.Errorf("expected substrings missing in %q", got)
	}
}

func TestBuilder_ImportComponent(t *testing.T) {
	got := New().ImportComponent("abc123", "hero").String()
	must := []string{
		`const heroComp = await figma.importComponentByKeyAsync("abc123");`,
		`const hero = heroComp.createInstance();`,
	}
	for _, s := range must {
		if !strings.Contains(got, s) {
			t.Errorf("ImportComponent output missing %q in:\n%s", s, got)
		}
	}
}

func TestBuilder_ImportComponentSet(t *testing.T) {
	got := New().ImportComponentSet("setKey", "Size=Large,State=Default", "btn").String()
	must := []string{
		`const btnSet = await figma.importComponentSetByKeyAsync("setKey");`,
		`p["Size"]==="Large"`,
		`p["State"]==="Default"`,
		`const btn = btnVariant.createInstance();`,
	}
	for _, s := range must {
		if !strings.Contains(got, s) {
			t.Errorf("ImportComponentSet output missing %q in:\n%s", s, got)
		}
	}
}

func TestBuilder_ImportComponentSet_Empty(t *testing.T) {
	got := New().ImportComponentSet("k", "", "x").String()
	if !strings.Contains(got, "true") {
		t.Errorf("empty variant should produce 'true' match, got:\n%s", got)
	}
}

func TestBuilder_ImportStyle(t *testing.T) {
	got := New().ImportStyle("s1key", "myStyle").String()
	want := `const myStyle = await figma.importStyleByKeyAsync("s1key");`
	if !strings.Contains(got, want) {
		t.Errorf("ImportStyle output missing %q in:\n%s", want, got)
	}
}
