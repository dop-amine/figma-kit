package theme

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RGB holds a Figma-compatible color in 0-1 range.
// JSON accepts either {"r":0.2,"g":0.4,"b":1.0} or "#3366FF".
type RGB struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

// UnmarshalJSON accepts hex strings ("#RRGGBB") or standard {r,g,b} objects.
func (c *RGB) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		var hex string
		if err := json.Unmarshal(data, &hex); err != nil {
			return err
		}
		hex = strings.TrimPrefix(hex, "#")
		if len(hex) != 6 {
			return fmt.Errorf("invalid hex color %q: must be 6 hex digits", hex)
		}
		var r, g, b uint8
		if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
			return fmt.Errorf("invalid hex color %q: %w", hex, err)
		}
		c.R = math.Round(float64(r)/255*100) / 100
		c.G = math.Round(float64(g)/255*100) / 100
		c.B = math.Round(float64(b)/255*100) / 100
		return nil
	}
	type plain RGB
	return json.Unmarshal(data, (*plain)(c))
}

// TypeSpec defines a typography preset.
type TypeSpec struct {
	FontSize   int    `json:"fontSize"`
	Style      string `json:"style"`
	LineHeight *int   `json:"lineHeight"`
	Family     string `json:"family,omitempty"`
}

// FontSpec defines the font families and weights used by a theme.
type FontSpec struct {
	Heading string   `json:"heading"`
	Body    string   `json:"body"`
	Mono    string   `json:"mono"`
	Weights []string `json:"weights"`
}

// GlassPreset defines glassmorphism effect parameters.
type GlassPreset struct {
	R  int     `json:"r"`
	F  float64 `json:"f"`
	S  float64 `json:"s"`
	GA float64 `json:"ga"`
	BL int     `json:"bl"`
}

// ShadowPreset defines a drop shadow effect.
type ShadowPreset struct {
	Type      string `json:"type"`
	Color     RGBA   `json:"color"`
	Offset    XY     `json:"offset"`
	Radius    int    `json:"radius"`
	Spread    int    `json:"spread"`
	Visible   bool   `json:"visible"`
	BlendMode string `json:"blendMode"`
}

// RGBA is an RGB color with alpha.
type RGBA struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
	A float64 `json:"a"`
}

// XY is a 2D offset.
type XY struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// EffectsSpec defines the effects section of a theme.
type EffectsSpec struct {
	Glass  map[string]GlassPreset  `json:"glass"`
	Shadow map[string]ShadowPreset `json:"shadow"`
}

// SpacingPreset defines spacing values for a context.
type SpacingPreset struct {
	Width   int `json:"width,omitempty"`
	Height  int `json:"height,omitempty"`
	Padding int `json:"padding,omitempty"`
	Margin  int `json:"margin,omitempty"`
	Gap     int `json:"gap,omitempty"`
}

// SpacingSpec defines the spacing section of a theme.
type SpacingSpec struct {
	Page    SpacingPreset `json:"page"`
	Card    SpacingPreset `json:"card"`
	Slide   SpacingPreset `json:"slide"`
	Frame16 SpacingPreset `json:"frame16"`
	Letter  SpacingPreset `json:"letter"`
}

// GradientStop defines a gradient color stop.
type GradientStop struct {
	Position float64 `json:"position"`
	Color    RGBA    `json:"color"`
}

// GradientSpec defines a gradient paint.
type GradientSpec struct {
	Type              string         `json:"type"`
	GradientTransform [2][3]float64  `json:"gradientTransform"`
	GradientStops     []GradientStop `json:"gradientStops"`
}

// BrandSpec defines brand-specific metadata (optional).
type BrandSpec struct {
	Primary  string   `json:"primary,omitempty"`
	Logo     string   `json:"logo,omitempty"`
	LogoFull string   `json:"logoFull,omitempty"`
	Tagline  string   `json:"tagline,omitempty"`
	URL      string   `json:"url,omitempty"`
	Clients  []string `json:"clients,omitempty"`
	Product  string   `json:"product,omitempty"`
	Features []string `json:"features,omitempty"`
}

// Theme is the complete theme configuration.
type Theme struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Website     string                  `json:"website,omitempty"`
	Colors      map[string]RGB          `json:"colors"`
	Type        map[string]TypeSpec     `json:"type"`
	Fonts       FontSpec                `json:"fonts"`
	Effects     EffectsSpec             `json:"effects"`
	Spacing     SpacingSpec             `json:"spacing"`
	Gradients   map[string]GradientSpec `json:"gradients"`
	Brand       *BrandSpec              `json:"brand,omitempty"`
}

// Load reads a theme by name, searching embedded, community, user config, then local directories.
func Load(name string) (*Theme, error) {
	if data, ok := embeddedThemes[name]; ok {
		return parseTheme(data)
	}

	if data, ok := communityThemes[name]; ok {
		return parseTheme(data)
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		p := filepath.Join(configDir, "figma-kit", "themes", name+".json")
		if d, err := os.ReadFile(p); err == nil {
			return parseTheme(d)
		}
	}

	p := filepath.Join("themes", name+".json")
	if d, err := os.ReadFile(p); err == nil {
		return parseTheme(d)
	}

	return nil, fmt.Errorf("theme %q not found (searched: embedded, community, ~/.config/figma-kit/themes/, ./themes/)", name)
}

// LoadFile reads a theme from an explicit file path.
func LoadFile(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading theme file: %w", err)
	}
	return parseTheme(data)
}

// ThemeSource indicates where a theme was discovered.
type ThemeSource string

const (
	SourceBuiltIn   ThemeSource = "built-in"
	SourceCommunity ThemeSource = "community"
	SourceUser      ThemeSource = "user"
	SourceLocal     ThemeSource = "local"
)

// ThemeInfo holds summary data for listing themes.
type ThemeInfo struct {
	Key         string
	Name        string
	Description string
	Source      ThemeSource
}

// List returns all discoverable themes grouped by source.
func List() []ThemeInfo {
	seen := map[string]bool{}
	var infos []ThemeInfo

	for name, data := range embeddedThemes {
		if info, ok := infoFromJSON(name, data, SourceBuiltIn); ok {
			infos = append(infos, info)
			seen[name] = true
		}
	}

	for name, data := range communityThemes {
		if seen[name] {
			continue
		}
		if info, ok := infoFromJSON(name, data, SourceCommunity); ok {
			infos = append(infos, info)
			seen[name] = true
		}
	}

	configDir, err := os.UserConfigDir()
	if err == nil {
		scanDir(filepath.Join(configDir, "figma-kit", "themes"), SourceUser, seen, &infos)
	}

	scanDir("themes", SourceLocal, seen, &infos)

	sort.Slice(infos, func(i, j int) bool {
		if infos[i].Source != infos[j].Source {
			return sourceOrder(infos[i].Source) < sourceOrder(infos[j].Source)
		}
		return infos[i].Key < infos[j].Key
	})
	return infos
}

func sourceOrder(s ThemeSource) int {
	switch s {
	case SourceBuiltIn:
		return 0
	case SourceCommunity:
		return 1
	case SourceUser:
		return 2
	case SourceLocal:
		return 3
	}
	return 4
}

func infoFromJSON(key string, data []byte, src ThemeSource) (ThemeInfo, bool) {
	var partial struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(data, &partial); err != nil {
		return ThemeInfo{}, false
	}
	return ThemeInfo{Key: key, Name: partial.Name, Description: partial.Description, Source: src}, true
}

func scanDir(dir string, src ThemeSource, seen map[string]bool, infos *[]ThemeInfo) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		key := strings.TrimSuffix(e.Name(), ".json")
		if seen[key] {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if info, ok := infoFromJSON(key, data, src); ok {
			*infos = append(*infos, info)
			seen[key] = true
		}
	}
}

// ColorNames returns the color key names in sorted order.
func (t *Theme) ColorNames() []string {
	names := make([]string, 0, len(t.Colors))
	for name := range t.Colors {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func parseTheme(data []byte) (*Theme, error) {
	var t Theme
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing theme: %w", err)
	}
	if t.Name == "" {
		return nil, fmt.Errorf("theme missing required field 'name'")
	}
	if len(t.Colors) == 0 {
		return nil, fmt.Errorf("theme %q has no colors defined", t.Name)
	}
	// Normalize empty description
	if t.Description == "" {
		t.Description = strings.TrimSpace(t.Name)
	}
	return &t, nil
}
