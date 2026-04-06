package codegen

import (
	"fmt"
	"strconv"
	"strings"
)

// RGB holds a Figma-compatible color in 0-1 range.
type RGB struct {
	R float64
	G float64
	B float64
}

// HexToRGB converts a hex color string (#RGB, #RRGGBB) to 0-1 range RGB.
func HexToRGB(hex string) (RGB, error) {
	hex = strings.TrimPrefix(hex, "#")
	switch len(hex) {
	case 3:
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	case 6:
		// already fine
	default:
		return RGB{}, fmt.Errorf("invalid hex color %q: must be 3 or 6 hex digits", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color %q: %w", hex, err)
	}

	return RGB{
		R: round(float64(r) / 255),
		G: round(float64(g) / 255),
		B: round(float64(b) / 255),
	}, nil
}

// RGBToHex converts 0-1 range RGB to a hex string (#RRGGBB).
func RGBToHex(c RGB) string {
	r := int(c.R*255 + 0.5)
	g := int(c.G*255 + 0.5)
	b := int(c.B*255 + 0.5)
	return fmt.Sprintf("#%02X%02X%02X", clamp8(r), clamp8(g), clamp8(b))
}

// FormatRGB returns a Figma Plugin API color literal: {r:0.2,g:0.4,b:1}.
func FormatRGB(c RGB) string {
	return fmt.Sprintf("{r:%s,g:%s,b:%s}", fmtFloat(c.R), fmtFloat(c.G), fmtFloat(c.B))
}

// FormatRGBA returns a Figma color literal with alpha: {r:0.2,g:0.4,b:1,a:0.5}.
func FormatRGBA(c RGB, a float64) string {
	return fmt.Sprintf("{r:%s,g:%s,b:%s,a:%s}", fmtFloat(c.R), fmtFloat(c.G), fmtFloat(c.B), fmtFloat(a))
}

func round(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

func clamp8(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

// FmtFloat formats a float for JS output, trimming unnecessary trailing zeros.
func FmtFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func fmtFloat(f float64) string {
	return FmtFloat(f)
}
