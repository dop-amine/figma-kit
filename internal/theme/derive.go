package theme

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// HexToRGB converts a hex color string (#RRGGBB or RRGGBB) to Figma 0-1 RGB.
func HexToRGB(hex string) (RGB, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color %q: must be 6 hex digits", hex)
	}
	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}
	return RGB{
		R: math.Round(float64(r)/255*100) / 100,
		G: math.Round(float64(g)/255*100) / 100,
		B: math.Round(float64(b)/255*100) / 100,
	}, nil
}

func rgbToHSL(c RGB) (h, s, l float64) {
	max := math.Max(c.R, math.Max(c.G, c.B))
	min := math.Min(c.R, math.Min(c.G, c.B))
	l = (max + min) / 2

	if max == min {
		return 0, 0, l
	}

	d := max - min
	if l > 0.5 {
		s = d / (2 - max - min)
	} else {
		s = d / (max + min)
	}

	switch max {
	case c.R:
		h = (c.G - c.B) / d
		if c.G < c.B {
			h += 6
		}
	case c.G:
		h = (c.B-c.R)/d + 2
	case c.B:
		h = (c.R-c.G)/d + 4
	}
	h /= 6
	return
}

func hslToRGB(h, s, l float64) RGB {
	if s == 0 {
		return RGB{R: l, G: l, B: l}
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	return RGB{
		R: round2(hue2rgb(p, q, h+1.0/3.0)),
		G: round2(hue2rgb(p, q, h)),
		B: round2(hue2rgb(p, q, h-1.0/3.0)),
	}
}

func hue2rgb(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6.0:
		return p + (q-p)*6*t
	case t < 1.0/2.0:
		return q
	case t < 2.0/3.0:
		return p + (q-p)*(2.0/3.0-t)*6
	default:
		return p
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func lighten(c RGB, amount float64) RGB {
	h, s, l := rgbToHSL(c)
	return hslToRGB(h, s, clamp01(l+amount))
}

func darken(c RGB, amount float64) RGB {
	h, s, l := rgbToHSL(c)
	return hslToRGB(h, s, clamp01(l-amount))
}


func isDark(c RGB) bool {
	luminance := 0.299*c.R + 0.587*c.G + 0.114*c.B
	return luminance < 0.5
}

// DerivePalette generates a full color palette from background, primary, and accent colors.
// Optional overrides in opts are applied after derivation.
func DerivePalette(bg, primary, accent RGB, opts *ThemeOptions) map[string]RGB {
	dark := isDark(bg)

	colors := map[string]RGB{
		"BG": bg,
		"BL": primary,
		"TL": accent,
	}

	if dark {
		colors["CARD"] = lighten(bg, 0.04)
		colors["CARD2"] = lighten(bg, 0.07)
		colors["STK"] = lighten(bg, 0.08)
		colors["WT"] = RGB{R: 0.96, G: 0.97, B: 0.98}
		colors["BD"] = RGB{R: 0.78, G: 0.81, B: 0.85}
		colors["MT"] = RGB{R: 0.45, G: 0.48, B: 0.55}
		colors["HOVER"] = lighten(bg, 0.05)
	} else {
		colors["CARD"] = darken(bg, 0.03)
		colors["CARD2"] = darken(bg, 0.05)
		colors["STK"] = darken(bg, 0.08)
		colors["WT"] = RGB{R: 0.08, G: 0.09, B: 0.12}
		colors["BD"] = RGB{R: 0.25, G: 0.27, B: 0.32}
		colors["MT"] = RGB{R: 0.50, G: 0.52, B: 0.56}
		colors["HOVER"] = darken(bg, 0.04)
	}

	colors["LINK"] = primary
	colors["WARN"] = RGB{R: 1.00, G: 0.60, B: 0.20}
	colors["ERR"] = RGB{R: 1.00, G: 0.35, B: 0.35}
	colors["SUCCESS"] = RGB{R: 0.20, G: 0.80, B: 0.50}

	if opts != nil {
		if (opts.Warn != RGB{}) {
			colors["WARN"] = opts.Warn
		}
		if (opts.Error != RGB{}) {
			colors["ERR"] = opts.Error
		}
		if (opts.Success != RGB{}) {
			colors["SUCCESS"] = opts.Success
		}
	}

	return colors
}

// ThemeOptions controls theme generation beyond the 3 seed colors.
type ThemeOptions struct {
	FontHeading string
	FontBody    string
	FontMono    string
	Warn        RGB
	Error       RGB
	Success     RGB
	Spacing     SpacingMode
	FromPath    string
}

// SpacingMode selects a spacing preset.
type SpacingMode string

const (
	SpacingDefault  SpacingMode = ""
	SpacingCompact  SpacingMode = "compact"
	SpacingSpacious SpacingMode = "spacious"
)

// SpacingForMode returns the spacing preset for the given mode.
func SpacingForMode(mode SpacingMode) SpacingSpec {
	switch mode {
	case SpacingCompact:
		return SpacingSpec{
			Page:    SpacingPreset{Padding: 48, Gap: 12},
			Card:    SpacingPreset{Padding: 16, Gap: 8},
			Slide:   SpacingPreset{Width: 1080, Height: 1350, Padding: 48},
			Frame16: SpacingPreset{Width: 1920, Height: 1080},
			Letter:  SpacingPreset{Width: 1224, Height: 1584, Margin: 40},
		}
	case SpacingSpacious:
		return SpacingSpec{
			Page:    SpacingPreset{Padding: 120, Gap: 24},
			Card:    SpacingPreset{Padding: 36, Gap: 16},
			Slide:   SpacingPreset{Width: 1080, Height: 1350, Padding: 120},
			Frame16: SpacingPreset{Width: 1920, Height: 1080},
			Letter:  SpacingPreset{Width: 1224, Height: 1584, Margin: 80},
		}
	default:
		return SpacingSpec{
			Page:    SpacingPreset{Padding: 80, Gap: 16},
			Card:    SpacingPreset{Padding: 24, Gap: 12},
			Slide:   SpacingPreset{Width: 1080, Height: 1350, Padding: 80},
			Frame16: SpacingPreset{Width: 1920, Height: 1080},
			Letter:  SpacingPreset{Width: 1224, Height: 1584, Margin: 60},
		}
	}
}

// GenerateThemeJSON creates a complete, valid theme JSON string.
// If opts is nil, defaults are used for fonts, status colors, and spacing.
func GenerateThemeJSON(name, description string, bg, primary, accent RGB, opts *ThemeOptions) (string, error) {
	colors := DerivePalette(bg, primary, accent, opts)

	fontHeading := "Inter"
	fontBody := "Inter"
	fontMono := "Geist Mono"
	spacing := SpacingForMode(SpacingDefault)

	if opts != nil {
		if opts.FontHeading != "" {
			fontHeading = opts.FontHeading
		}
		if opts.FontBody != "" {
			fontBody = opts.FontBody
		}
		if opts.FontMono != "" {
			fontMono = opts.FontMono
		}
		spacing = SpacingForMode(opts.Spacing)
	}

	t := Theme{
		Name:        name,
		Description: description,
		Colors:      colors,
		Type: map[string]TypeSpec{
			"h1":    {FontSize: 72, Style: "Bold", LineHeight: intPtr(86)},
			"h2":    {FontSize: 48, Style: "Bold", LineHeight: intPtr(58)},
			"h3":    {FontSize: 32, Style: "Semi Bold", LineHeight: intPtr(40)},
			"h4":    {FontSize: 22, Style: "Semi Bold", LineHeight: intPtr(30)},
			"body":  {FontSize: 16, Style: "Regular", LineHeight: intPtr(26)},
			"small": {FontSize: 13, Style: "Regular", LineHeight: intPtr(20)},
			"label": {FontSize: 11, Style: "Medium"},
			"mono":  {FontSize: 11, Style: "Medium", Family: fontMono},
		},
		Fonts: FontSpec{
			Heading: fontHeading,
			Body:    fontBody,
			Mono:    fontMono,
			Weights: []string{"Bold", "Semi Bold", "Medium", "Regular"},
		},
		Effects: EffectsSpec{
			Glass: map[string]GlassPreset{
				"subtle":  {R: 16, F: 0.03, S: 0.06, GA: 0.04, BL: 20},
				"default": {R: 20, F: 0.04, S: 0.08, GA: 0.06, BL: 24},
				"strong":  {R: 24, F: 0.06, S: 0.12, GA: 0.12, BL: 24},
			},
			Shadow: map[string]ShadowPreset{
				"card": {Type: "DROP_SHADOW", Color: RGBA{R: primary.R * 0.4, G: primary.G * 0.3, B: primary.B * 0.15, A: 0.10}, Offset: XY{X: 0, Y: 2}, Radius: 16, Visible: true, BlendMode: "NORMAL"},
				"glow": {Type: "DROP_SHADOW", Color: RGBA{R: primary.R, G: primary.G, B: primary.B, A: 0.20}, Offset: XY{X: 0, Y: 4}, Radius: 32, Visible: true, BlendMode: "NORMAL"},
			},
		},
		Spacing: spacing,
		Gradients: map[string]GradientSpec{
			"accent": {
				Type:              "GRADIENT_LINEAR",
				GradientTransform: [2][3]float64{{1, 0, 0}, {0, 1, 0}},
				GradientStops: []GradientStop{
					{Position: 0, Color: RGBA{R: primary.R, G: primary.G, B: primary.B, A: 1}},
					{Position: 1, Color: RGBA{R: accent.R, G: accent.G, B: accent.B, A: 1}},
				},
			},
		},
	}

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling theme: %w", err)
	}
	return string(data), nil
}

func intPtr(v int) *int {
	return &v
}
