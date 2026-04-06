package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func executeCmd(args ...string) (string, error) {
	cmd := newRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	var err error
	out := captureOutput(func() {
		err = cmd.Execute()
	})
	return out, err
}

func TestCLICommandsValidJSOutput(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "themes",
			args:     []string{"themes"},
			expected: []string{"default", "light", "noir"},
		},
		{
			name:     "info",
			args:     []string{"info"},
			expected: []string{"figma-kit", "Layers:", "Themes:"},
		},
		{
			name:     "preamble_default_theme",
			args:     []string{"preamble", "-t", "default"},
			expected: []string{"figma.loadFontAsync", "Theme Colors"},
		},
		{
			name:     "helpers",
			args:     []string{"helpers"},
			expected: []string{"function T(", "function G("},
		},
		{
			name:     "template_slide",
			args:     []string{"template", "slide"},
			expected: []string{"createSlide"},
		},
		{
			name:     "node_create_frame",
			args:     []string{"node", "create", "frame", "--name", "Hero", "-w", "1440"},
			expected: []string{"createFrame", "Hero", "1440"},
		},
		{
			name:     "node_create_text",
			args:     []string{"node", "create", "text", "--name", "Title"},
			expected: []string{"createText"},
		},
		{
			name:     "style_fill_solid",
			args:     []string{"style", "fill", "abc123", "--solid", "#FF0000"},
			expected: []string{"fills", "SOLID"},
		},
		{
			name:     "text_create",
			args:     []string{"text", "create", "--content", "Hello", "--font", "Inter", "--size", "24"},
			expected: []string{"characters", "Hello", "fontSize"},
		},
		{
			name:     "layout_auto",
			args:     []string{"layout", "auto", "abc", "--dir", "VERTICAL", "--gap", "16", "--pad", "24"},
			expected: []string{"layoutMode", "VERTICAL", "itemSpacing"},
		},
		{
			name:     "card_glass",
			args:     []string{"card", "glass", "-w", "400", "--height", "300"},
			expected: []string{"function G("},
		},
		{
			name:     "make_og_image",
			args:     []string{"make", "og-image", "--title", "Test", "--description", "Desc"},
			expected: []string{"1200", "630"},
		},
		{
			name:     "status",
			args:     []string{"status"},
			expected: []string{"figma.root.children"},
		},
		{
			name:     "export_tokens_css",
			args:     []string{"export", "tokens", "-t", "default", "--format", "css"},
			expected: []string{":root", "--fk-"},
		},
		{
			name:     "scaffold_noir",
			args:     []string{"scaffold", "-t", "noir"},
			expected: []string{"Noir", "function T("},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := executeCmd(tt.args...)
			if err != nil {
				t.Fatalf("executeCmd(%v): %v", tt.args, err)
			}
			for _, sub := range tt.expected {
				if !strings.Contains(out, sub) {
					t.Errorf("output missing substring %q\n--- output ---\n%s", sub, out)
				}
			}
		})
	}
}
