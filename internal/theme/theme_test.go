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

func TestList_includesBuiltInAndCommunity(t *testing.T) {
	infos := List()
	if len(infos) < 4 {
		t.Fatalf("expected at least 4 themes (3 built-in + community), got %d: %+v", len(infos), infos)
	}
	keys := make(map[string]bool)
	for _, info := range infos {
		keys[info.Key] = true
	}
	for _, want := range []string{"default", "light", "noir", "ocean"} {
		if !keys[want] {
			t.Errorf("missing expected theme %q in list", want)
		}
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

func TestLoadFile_hexColors(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "hex.json")
	js := `{
  "name": "Hex Test",
  "colors": {
    "BG": "#0D0F17",
    "BL": "#3366FF",
    "mixed": { "r": 0.5, "g": 0.5, "b": 0.5 }
  }
}`
	if err := os.WriteFile(p, []byte(js), 0o644); err != nil {
		t.Fatal(err)
	}
	th, err := LoadFile(p)
	if err != nil {
		t.Fatalf("LoadFile hex theme: %v", err)
	}
	bg := th.Colors["BG"]
	if bg.R != 0.05 || bg.G != 0.06 || bg.B != 0.09 {
		t.Errorf("BG hex parse: got %+v", bg)
	}
	bl := th.Colors["BL"]
	if bl.R != 0.2 || bl.G != 0.4 || bl.B != 1.0 {
		t.Errorf("BL hex parse: got %+v", bl)
	}
	mixed := th.Colors["mixed"]
	if mixed.R != 0.5 {
		t.Errorf("mixed obj parse: got %+v", mixed)
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
