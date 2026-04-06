package codegen

import (
	"math"
	"strings"
	"testing"
)

func approxRGB(t *testing.T, got, want RGB) {
	t.Helper()
	const eps = 1e-9
	for _, tc := range []struct {
		name string
		a, b float64
	}{
		{"R", got.R, want.R},
		{"G", got.G, want.G},
		{"B", got.B, want.B},
	} {
		if math.Abs(tc.a-tc.b) > eps {
			t.Errorf("%s: got %g, want %g (full RGB %+v vs %+v)", tc.name, tc.a, tc.b, got, want)
		}
	}
}

func TestHexToRGB_valid(t *testing.T) {
	tests := []struct {
		name string
		hex  string
		want RGB
	}{
		{"6-digit with hash", "#3366FF", RGB{R: 0.2, G: 0.4, B: 1}},
		{"6-digit no hash", "3366FF", RGB{R: 0.2, G: 0.4, B: 1}},
		{"3-digit with hash", "#FFF", RGB{R: 1, G: 1, B: 1}},
		{"3-digit no hash", "FFF", RGB{R: 1, G: 1, B: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToRGB(tt.hex)
			if err != nil {
				t.Fatalf("HexToRGB(%q): %v", tt.hex, err)
			}
			approxRGB(t, got, tt.want)
		})
	}
}

func TestHexToRGB_errors(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantSub string
	}{
		{"empty after trim", "", "must be 3 or 6 hex digits"},
		{"too short", "#12", "must be 3 or 6 hex digits"},
		{"too long", "#1122334", "must be 3 or 6 hex digits"},
		{"invalid chars", "#GGGGGG", "invalid hex color"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := HexToRGB(tt.hex)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

func TestRGBToHex_roundTrip(t *testing.T) {
	// HexToRGB rounds to 2 decimal places; only some inputs round-trip exactly.
	hexes := []string{"#3366FF", "#000000", "#FFFFFF", "#FF0000", "#00FF00", "#0000FF"}
	for _, h := range hexes {
		t.Run(h, func(t *testing.T) {
			rgb, err := HexToRGB(h)
			if err != nil {
				t.Fatal(err)
			}
			back := RGBToHex(rgb)
			if !strings.EqualFold(back, h) {
				t.Errorf("round-trip: %s -> %+v -> %s, want %s", h, rgb, back, h)
			}
		})
	}
}

func TestFormatRGB(t *testing.T) {
	tests := []struct {
		name string
		c    RGB
		want string
	}{
		{"example", RGB{R: 0.2, G: 0.4, B: 1}, "{r:0.2,g:0.4,b:1}"},
		{"ints", RGB{R: 1, G: 0, B: 0}, "{r:1,g:0,b:0}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRGB(tt.c)
			if got != tt.want {
				t.Errorf("FormatRGB(%+v) = %q, want %q", tt.c, got, tt.want)
			}
		})
	}
}

func TestFormatRGBA(t *testing.T) {
	got := FormatRGBA(RGB{R: 0.2, G: 0.4, B: 1}, 0.5)
	want := "{r:0.2,g:0.4,b:1,a:0.5}"
	if got != want {
		t.Errorf("FormatRGBA(...) = %q, want %q", got, want)
	}
}

func TestFmtFloat(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		want string
	}{
		{"one", 1.0, "1"},
		{"half", 0.5, "0.5"},
		{"decimal", 0.123, "0.123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FmtFloat(tt.f)
			if got != tt.want {
				t.Errorf("FmtFloat(%v) = %q, want %q", tt.f, got, tt.want)
			}
		})
	}
}
