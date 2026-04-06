package config

import (
	"os"
	"path/filepath"
	"testing"
)

func chdirTemp(t *testing.T) (restore func()) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	}
}

func TestInit_createsFigmarc(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Init("myproj"); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(".", configFile)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected %s: %v", configFile, err)
	}
}

func TestLoad_noFile_returnsDefaults(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.Theme != "default" {
		t.Errorf("Theme = %q, want default", c.Theme)
	}
}

func TestSaveLoad_roundTrip(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	want := &Config{
		FileKey:   "abc123",
		Theme:     "light",
		Page:      2,
		ExportDir: "out",
	}
	if err := Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.FileKey != want.FileKey || got.Theme != want.Theme || got.Page != want.Page || got.ExportDir != want.ExportDir {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestGetSet_fileKey(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Set("fileKey", "fk123"); err != nil {
		t.Fatal(err)
	}
	v, err := Get("fileKey")
	if err != nil {
		t.Fatal(err)
	}
	if v != "fk123" {
		t.Errorf("Get fileKey = %q", v)
	}
}

func TestGetSet_theme(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Set("theme", "noir"); err != nil {
		t.Fatal(err)
	}
	v, err := Get("theme")
	if err != nil {
		t.Fatal(err)
	}
	if v != "noir" {
		t.Errorf("Get theme = %q", v)
	}
}

func TestGetSet_page(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Set("page", "42"); err != nil {
		t.Fatal(err)
	}
	v, err := Get("page")
	if err != nil {
		t.Fatal(err)
	}
	if v != "42" {
		t.Errorf("Get page = %q", v)
	}
}

func TestGetSet_exportDir(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Set("exportDir", "exports"); err != nil {
		t.Fatal(err)
	}
	v, err := Get("exportDir")
	if err != nil {
		t.Fatal(err)
	}
	if v != "exports" {
		t.Errorf("Get exportDir = %q", v)
	}
}

func TestGet_invalidKey(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	_, err := Get("notAKey")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSet_invalidKey(t *testing.T) {
	restore := chdirTemp(t)
	defer restore()

	if err := Set("badKey", "x"); err == nil {
		t.Fatal("expected error")
	}
}
