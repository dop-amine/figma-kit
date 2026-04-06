package theme

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestLoad_embeddedDefault(t *testing.T) {
	th, err := Load("default")
	if err != nil {
		t.Fatal(err)
	}
	if th.Name == "" {
		t.Error("expected non-empty Name")
	}
	if len(th.Colors) == 0 {
		t.Error("expected Colors")
	}
	if len(th.Type) == 0 {
		t.Error("expected Type")
	}
	if th.Fonts.Heading == "" || th.Fonts.Body == "" || th.Fonts.Mono == "" {
		t.Errorf("expected Fonts (heading/body/mono), got %+v", th.Fonts)
	}
}

func TestLoad_embeddedLight(t *testing.T) {
	_, err := Load("light")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad_embeddedNoirHasBrand(t *testing.T) {
	th, err := Load("noir")
	if err != nil {
		t.Fatal(err)
	}
	if th.Brand == nil {
		t.Fatal("expected Brand on noir theme")
	}
	if th.Brand.Primary == "" {
		t.Error("expected non-empty Brand.Primary")
	}
}

func TestLoad_nonexistent(t *testing.T) {
	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestList_threeEmbeddedThemes(t *testing.T) {
	infos := List()
	if len(infos) != 3 {
		t.Fatalf("expected 3 themes, got %d: %+v", len(infos), infos)
	}
	keys := make([]string, len(infos))
	for i, info := range infos {
		keys[i] = info.Key
	}
	want := []string{"default", "light", "noir"}
	if !slices.Equal(keys, want) {
		t.Errorf("keys = %v, want %v", keys, want)
	}
}

func TestTheme_ColorNames_sorted(t *testing.T) {
	th, err := Load("default")
	if err != nil {
		t.Fatal(err)
	}
	names := th.ColorNames()
	if len(names) < 2 {
		t.Fatal("expected multiple color names")
	}
	if !slices.IsSorted(names) {
		t.Errorf("ColorNames not sorted: %v", names)
	}
}

func TestLoadFile_parseValidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "t.json")
	valid := `{
  "name": "Parse Test",
  "colors": {
    "accent": { "r": 0.2, "g": 0.4, "b": 1.0 }
  }
}`
	if err := os.WriteFile(p, []byte(valid), 0644); err != nil {
		t.Fatal(err)
	}
	th, err := LoadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if th.Name != "Parse Test" {
		t.Errorf("Name = %q", th.Name)
	}
	c, ok := th.Colors["accent"]
	if !ok {
		t.Fatal("missing accent color")
	}
	if c.R != 0.2 || c.G != 0.4 || c.B != 1.0 {
		t.Errorf("RGB = %+v", c)
	}
}

func TestLoadFile_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(p, []byte(`{not json`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_missingName(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "noname.json")
	js := `{"colors":{"x":{"r":0,"g":0,"b":0}}}`
	if err := os.WriteFile(p, []byte(js), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_missingColors(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "nocolors.json")
	js := `{"name":"No Colors"}`
	if err := os.WriteFile(p, []byte(js), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadFile_emptyColorsMap(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "emptycolors.json")
	js := `{"name":"Empty","colors":{}}`
	if err := os.WriteFile(p, []byte(js), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadFile(p)
	if err == nil {
		t.Fatal("expected error for empty colors")
	}
}
