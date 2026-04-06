package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RGB holds a Figma-compatible color in 0-1 range.
type RGB struct {
	R float64 `json:"r"`
	G float64 `json:"g"`
	B float64 `json:"b"`
}

// TypeSpec defines a typography preset.
type TypeSpec struct {
	FontSize   int     `json:"fontSize"`
	Style      string  `json:"style"`
	LineHeight *int    `json:"lineHeight"`
	Family     string  `json:"family,omitempty"`
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
	Type      string  `json:"type"`
	Color     RGBA    `json:"color"`
	Offset    XY      `json:"offset"`
	Radius    int     `json:"radius"`
	Spread    int     `json:"spread"`
	Visible   bool    `json:"visible"`
	BlendMode string  `json:"blendMode"`
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
	Type              string       `json:"type"`
	GradientTransform [2][3]float64 `json:"gradientTransform"`
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
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Website     string                 `json:"website,omitempty"`
	Colors      map[string]RGB         `json:"colors"`
	Type        map[string]TypeSpec    `json:"type"`
	Fonts       FontSpec               `json:"fonts"`
	Effects     EffectsSpec            `json:"effects"`
	Spacing     SpacingSpec            `json:"spacing"`
	Gradients   map[string]GradientSpec `json:"gradients"`
	Brand       *BrandSpec             `json:"brand,omitempty"`
}

// Load reads a theme by name, searching embedded themes first, then user directories.
func Load(name string) (*Theme, error) {
	// Check embedded themes
	data, ok := embeddedThemes[name]
	if ok {
		return parseTheme(data)
	}

	// Check user config dir
	configDir, err := os.UserConfigDir()
	if err == nil {
		p := filepath.Join(configDir, "figma-kit", "themes", name+".json")
		if d, err := os.ReadFile(p); err == nil {
			return parseTheme(d)
		}
	}

	// Check local ./themes directory
	p := filepath.Join("themes", name+".json")
	if d, err := os.ReadFile(p); err == nil {
		return parseTheme(d)
	}

	return nil, fmt.Errorf("theme %q not found (searched: embedded, ~/.config/figma-kit/themes/, ./themes/)", name)
}

// LoadFile reads a theme from an explicit file path.
func LoadFile(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading theme file: %w", err)
	}
	return parseTheme(data)
}

// List returns the names and descriptions of all available themes (embedded only).
func List() []ThemeInfo {
	var infos []ThemeInfo
	for name, data := range embeddedThemes {
		var partial struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := json.Unmarshal(data, &partial); err == nil {
			infos = append(infos, ThemeInfo{
				Key:         name,
				Name:        partial.Name,
				Description: partial.Description,
			})
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Key < infos[j].Key
	})
	return infos
}

// ThemeInfo holds summary data for listing themes.
type ThemeInfo struct {
	Key         string
	Name        string
	Description string
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
