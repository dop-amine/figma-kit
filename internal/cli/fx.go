package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amine/figma-kit/internal/codegen"
	"github.com/spf13/cobra"
)

func newFXCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fx",
		Short: "Visual effects (glow, mesh, noise, vignette, grain, blur, shadow, ...)",
	}
	cmd.AddCommand(newFXGlowCmd())
	cmd.AddCommand(newFXMeshCmd())
	cmd.AddCommand(newFXNoiseCmd())
	cmd.AddCommand(newFXVignetteCmd())
	cmd.AddCommand(newFXGrainCmd())
	cmd.AddCommand(newFXBlurBgCmd())
	cmd.AddCommand(newFXAccentBarCmd())
	cmd.AddCommand(newFXShadowCmd())
	cmd.AddCommand(newFXParallaxLayerCmd())
	return cmd
}

// --- glow ---

type glowRadial struct {
	transform [2][3]float64
	innerA    float64
	innerRGB  codegen.RGB
	midPos    float64 // optional second stop position (0 = use 1.0 end only with two stops)
}

func glowPresetLayers(position string) []glowRadial {
	switch position {
	case "topRight":
		return []glowRadial{
			{transform: [2][3]float64{{1.8, 0, 0.35}, {0, 1.5, -0.1}}, innerA: 0.14, innerRGB: codegen.RGB{R: 0.12, G: 0.22, B: 0.55}},
			{transform: [2][3]float64{{1.0, 0, -0.15}, {0, 0.8, 0.45}}, innerA: 0.10, innerRGB: codegen.RGB{R: 0.05, G: 0.28, B: 0.30}},
		}
	case "center":
		return []glowRadial{
			{transform: [2][3]float64{{2.0, 0, -0.1}, {0, 1.8, -0.05}}, innerA: 0.18, innerRGB: codegen.RGB{R: 0.12, G: 0.22, B: 0.50}, midPos: 0.8},
		}
	case "cta":
		return []glowRadial{
			{transform: [2][3]float64{{1.5, 0, 0.0}, {0, 1.2, -0.05}}, innerA: 0.22, innerRGB: codegen.RGB{R: 0.14, G: 0.26, B: 0.58}, midPos: 0.7},
			{transform: [2][3]float64{{1.0, 0, 0.1}, {0, 0.8, 0.3}}, innerA: 0.12, innerRGB: codegen.RGB{R: 0.05, G: 0.28, B: 0.30}},
		}
	default: // subtle
		return []glowRadial{
			{transform: [2][3]float64{{2.0, 0, 0.2}, {0, 1.6, -0.1}}, innerA: 0.10, innerRGB: codegen.RGB{R: 0.10, G: 0.18, B: 0.40}},
		}
	}
}

func newFXGlowCmd() *cobra.Command {
	var (
		position  string
		intensity float64
		colorHex  string
	)
	cmd := &cobra.Command{
		Use:   "glow <nodeId>",
		Short: "Add radial gradient glow fills to a frame",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pos := position
			switch pos {
			case "topRight", "center", "subtle", "cta":
			default:
				return fmt.Errorf("position must be topRight, center, subtle, or cta")
			}
			var tint codegen.RGB
			var useCustom bool
			if strings.TrimSpace(colorHex) != "" {
				c, err := codegen.HexToRGB(colorHex)
				if err != nil {
					return err
				}
				tint = c
				useCustom = true
			}
			layers := glowPresetLayers(pos)

			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('fills' in node)) throw new Error('Node has no fills');")
			b.Line("function __fxFirstSolid(n) {")
			b.Line("  if (!n.fills || n.fills.length === 0) return null;")
			b.Line("  const s = n.fills.find(f => f.type === 'SOLID');")
			b.Line("  return s ? s.color : null;")
			b.Line("}")
			b.Line("const __bg = __fxFirstSolid(node) || { r: 0.06, g: 0.07, b: 0.11 };")
			b.Linef("const __k = %s;", codegen.FmtFloat(intensity))
			b.Line("const __fills = [{ type: 'SOLID', color: __bg }];")

			for _, layer := range layers {
				t := layer.transform
				ir, ig, ib := layer.innerRGB.R, layer.innerRGB.G, layer.innerRGB.B
				ia := layer.innerA
				if useCustom {
					ir, ig, ib = tint.R, tint.G, tint.B
				}
				a0 := ia * intensity
				if a0 > 1 {
					a0 = 1
				}
				b.Linef(`__fills.push({
  type: 'GRADIENT_RADIAL',
  gradientTransform: [[%s,%s,%s],[%s,%s,%s]],
  gradientStops: [`,
					codegen.FmtFloat(t[0][0]), codegen.FmtFloat(t[0][1]), codegen.FmtFloat(t[0][2]),
					codegen.FmtFloat(t[1][0]), codegen.FmtFloat(t[1][1]), codegen.FmtFloat(t[1][2]))
				b.Linef("    { position: 0, color: %s },", codegen.FormatRGBA(codegen.RGB{R: ir, G: ig, B: ib}, a0))
				if layer.midPos > 0 {
					b.Linef("    { position: %s, color: { ...__bg, a: 0 } }", codegen.FmtFloat(layer.midPos))
				} else {
					b.Line("    { position: 1, color: { ...__bg, a: 0 } }")
				}
				b.Line("  ],")
				b.Line("});")
			}

			b.Line("node.fills = __fills;")
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&position, "position", "subtle", "Glow preset: topRight, center, subtle, cta")
	cmd.Flags().Float64Var(&intensity, "intensity", 1, "Multiplier for glow gradient alphas")
	cmd.Flags().StringVar(&colorHex, "color", "", "Optional hex tint for glow (#RRGGBB)")
	return cmd
}

// --- mesh ---

func resolvePaletteColor(token string) (codegen.RGB, error) {
	token = strings.TrimSpace(strings.ToLower(token))
	named := map[string]string{
		"blue": "#2563eb", "teal": "#14b8a6", "purple": "#8b5cf6",
		"red": "#ef4444", "green": "#22c55e", "orange": "#f97316",
		"pink": "#ec4899", "cyan": "#06b6d4", "yellow": "#eab308",
	}
	if hex, ok := named[token]; ok {
		return codegen.HexToRGB(hex)
	}
	return codegen.HexToRGB(token)
}

func newFXMeshCmd() *cobra.Command {
	var (
		pointsN int
		palette string
	)
	cmd := &cobra.Command{
		Use:   "mesh <nodeId>",
		Short: "Multi-point mesh gradient (layered radial fills)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if pointsN < 2 {
				return fmt.Errorf("--points must be at least 2")
			}
			tokens := splitCommaTrim(palette)
			if len(tokens) < 1 {
				return fmt.Errorf("palette must list at least one color")
			}
			var colors []codegen.RGB
			for _, t := range tokens {
				c, err := resolvePaletteColor(t)
				if err != nil {
					return fmt.Errorf("palette token %q: %w", t, err)
				}
				colors = append(colors, c)
			}

			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('fills' in node)) throw new Error('Node has no fills');")
			b.Line("if (!('width' in node) || !('height' in node)) throw new Error('Node needs width/height');")
			b.Linef("const __nw = %d;", pointsN)
			b.Line("const __cols = [")
			for i, c := range colors {
				comma := ","
				if i == len(colors)-1 {
					comma = ""
				}
				b.Linef("  %s%s", codegen.FormatRGB(c), comma)
			}
			b.Line("];")
			b.Raw(`const __fills = [{ type: 'SOLID', color: __cols[0] }];
const golden = 2.399963229728653;
for (let i = 0; i < __nw; i++) {
  const t = i / Math.max(1, __nw - 1);
  const ang = t * golden * Math.PI * 2;
  const cx = 0.35 + 0.35 * Math.cos(ang);
  const cy = 0.35 + 0.35 * Math.sin(ang);
  const col = __cols[i % __cols.length];
  const col2 = __cols[(i + 1) % __cols.length];
  __fills.push({
    type: 'GRADIENT_RADIAL',
    gradientTransform: [[1.4, 0, cx], [0, 1.1, cy]],
    gradientStops: [
      { position: 0, color: { r: col.r, g: col.g, b: col.b, a: 0.45 } },
      { position: 0.55, color: { r: col2.r, g: col2.g, b: col2.b, a: 0.12 } },
      { position: 1, color: { r: col.r, g: col.g, b: col.b, a: 0 } },
    ],
  });
}
node.fills = __fills;
`)
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&pointsN, "points", 5, "Number of radial mesh points")
	cmd.Flags().StringVar(&palette, "palette", "#2563eb,#14b8a6,#8b5cf6", "Comma-separated hex or named colors (blue,teal,...)")
	return cmd
}

func splitCommaTrim(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// --- noise ---

func newFXNoiseCmd() *cobra.Command {
	var opacity float64
	cmd := &cobra.Command{
		Use:   "noise <nodeId>",
		Short: "Subtle noise overlay via gradient dithering",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('appendChild' in node)) throw new Error('Expected frame-like parent');")
			b.Line("if (!('width' in node) || !('height' in node)) throw new Error('Parent needs width/height');")
			b.Linef("const o = %s;", codegen.FmtFloat(opacity))
			b.Raw(`const n = figma.createRectangle();
n.name = 'FX Noise';
n.resize(node.width, node.height);
n.x = 0;
n.y = 0;
n.fills = [];
const stops = [];
const steps = 48;
for (let i = 0; i <= steps; i++) {
  const p = i / steps;
  const flick = (i % 2) * 0.04;
  const base = 0.48 + flick;
  stops.push({
    position: p,
    color: { r: base, g: base - 0.02, b: base + 0.02, a: o * 0.22 },
  });
}
n.fills = [{
  type: 'GRADIENT_LINEAR',
  gradientTransform: [[0.02, 1.2, 0], [-1.1, 0.05, 1]],
  gradientStops: stops,
}];
n.blendMode = 'OVERLAY';
node.appendChild(n);
`)
			b.ReturnIDs("node.id", "n.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().Float64Var(&opacity, "opacity", 0.35, "Noise overlay strength (0–1)")
	return cmd
}

// --- vignette ---

func newFXVignetteCmd() *cobra.Command {
	var strength float64
	cmd := &cobra.Command{
		Use:   "vignette <nodeId>",
		Short: "Edge vignette darkening overlay",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('appendChild' in node)) throw new Error('Expected frame-like parent');")
			b.Line("if (!('width' in node) || !('height' in node)) throw new Error('Parent needs width/height');")
			b.Linef("const s = %s;", codegen.FmtFloat(strength))
			b.Raw(`const v = figma.createRectangle();
v.name = 'FX Vignette';
v.resize(node.width, node.height);
v.x = 0;
v.y = 0;
const edgeA = Math.min(0.85, 0.25 + s * 0.45);
v.fills = [{
  type: 'GRADIENT_RADIAL',
  gradientTransform: [[1, 0, 0.5], [0, 1, 0.45]],
  gradientStops: [
    { position: 0, color: { r: 0, g: 0, b: 0, a: 0 } },
    { position: 0.55, color: { r: 0, g: 0, b: 0, a: edgeA * 0.35 } },
    { position: 1, color: { r: 0, g: 0, b: 0, a: edgeA } },
  ],
}];
v.blendMode = 'MULTIPLY';
node.appendChild(v);
`)
			b.ReturnIDs("node.id", "v.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().Float64Var(&strength, "strength", 0.6, "Vignette strength")
	return cmd
}

// --- grain ---

func newFXGrainCmd() *cobra.Command {
	var amount string
	cmd := &cobra.Command{
		Use:   "grain <nodeId>",
		Short: "Film grain overlay (dithered gradient)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var steps int
			var amp float64
			switch strings.ToLower(amount) {
			case "light":
				steps, amp = 32, 0.14
			case "heavy":
				steps, amp = 72, 0.32
			default: // medium
				steps, amp = 52, 0.22
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('appendChild' in node)) throw new Error('Expected frame-like parent');")
			b.Line("if (!('width' in node) || !('height' in node)) throw new Error('Parent needs width/height');")
			b.Linef("const __steps = %d;", steps)
			b.Linef("const __amp = %s;", codegen.FmtFloat(amp))
			b.Raw(`const g = figma.createRectangle();
g.name = 'FX Grain';
g.resize(node.width, node.height);
g.x = 0;
g.y = 0;
const stops = [];
for (let i = 0; i <= __steps; i++) {
  const p = i / __steps;
  const j = (i % 3) * __amp * 0.5;
  const v = 0.5 + j - __amp * 0.25;
  stops.push({
    position: p,
    color: { r: v, g: v - 0.03, b: v + 0.04, a: __amp * 0.55 },
  });
}
g.fills = [{
  type: 'GRADIENT_LINEAR',
  gradientTransform: [[1.1, 0.15, 0], [0.12, 0.9, 0]],
  gradientStops: stops,
}];
g.blendMode = 'SOFT_LIGHT';
node.appendChild(g);
`)
			b.ReturnIDs("node.id", "g.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&amount, "amount", "medium", "Grain density: light, medium, heavy")
	return cmd
}

// --- blur-bg ---

func parseRGBAString(s string) (codegen.RGB, float64, error) {
	s = strings.TrimSpace(s)
	low := strings.ToLower(s)
	if !strings.HasPrefix(low, "rgba(") || !strings.HasSuffix(s, ")") {
		return codegen.RGB{}, 0, fmt.Errorf("tint must look like rgba(r,g,b,a)")
	}
	inner := s[len("rgba(") : len(s)-1]
	parts := splitCommaTrim(inner)
	if len(parts) != 4 {
		return codegen.RGB{}, 0, fmt.Errorf("rgba() needs four comma-separated values")
	}
	parseCh := func(tok string) (float64, error) {
		v, err := strconv.ParseFloat(strings.TrimSpace(tok), 64)
		if err != nil {
			return 0, err
		}
		if v > 1 {
			v /= 255
		}
		if v < 0 {
			return 0, nil
		}
		if v > 1 {
			return 1, nil
		}
		return v, nil
	}
	rf, err := parseCh(parts[0])
	if err != nil {
		return codegen.RGB{}, 0, fmt.Errorf("rgba r: %w", err)
	}
	gf, err := parseCh(parts[1])
	if err != nil {
		return codegen.RGB{}, 0, fmt.Errorf("rgba g: %w", err)
	}
	bf, err := parseCh(parts[2])
	if err != nil {
		return codegen.RGB{}, 0, fmt.Errorf("rgba b: %w", err)
	}
	a, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
	if err != nil {
		return codegen.RGB{}, 0, fmt.Errorf("rgba a: %w", err)
	}
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	return codegen.RGB{R: rf, G: gf, B: bf}, a, nil
}

func newFXBlurBgCmd() *cobra.Command {
	var (
		radius int
		tint   string
	)
	cmd := &cobra.Command{
		Use:   "blur-bg <nodeId>",
		Short: "Frosted-glass overlay (fill + background blur)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fillLit := "{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.35 }"
			if strings.TrimSpace(tint) != "" {
				c, a, err := parseRGBAString(tint)
				if err != nil {
					return err
				}
				fillLit = fmt.Sprintf("{ type: 'SOLID', color: %s, opacity: %s }",
					codegen.FormatRGB(c), codegen.FmtFloat(a))
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('appendChild' in node)) throw new Error('Expected frame-like parent');")
			b.Line("if (!('width' in node) || !('height' in node)) throw new Error('Parent needs width/height');")
			b.Linef("const __fill = %s;", fillLit)
			b.Linef("const __r = %d;", radius)
			b.Raw(`const glass = figma.createRectangle();
glass.name = 'FX Blur BG';
glass.resize(node.width, node.height);
glass.x = 0;
glass.y = 0;
glass.fills = [__fill];
glass.effects = [{ type: 'BACKGROUND_BLUR', radius: __r, visible: true }];
node.appendChild(glass);
`)
			b.ReturnIDs("node.id", "glass.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&radius, "radius", 24, "Background blur radius")
	cmd.Flags().StringVar(&tint, "tint", "", `Optional tint as rgba(255,255,255,0.35)`)
	return cmd
}

// --- accent bar ---

func newFXAccentBarCmd() *cobra.Command {
	var (
		fromHex string
		toHex   string
		w, h    int
		x, y    int
	)
	cmd := &cobra.Command{
		Use:   "accent-bar <parentId>",
		Short: "Gradient accent bar (linear fill)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c0, err := codegen.HexToRGB(fromHex)
			if err != nil {
				return err
			}
			c1, err := codegen.HexToRGB(toHex)
			if err != nil {
				return err
			}
			if w <= 0 {
				return fmt.Errorf("--w must be positive")
			}
			if h <= 0 {
				h = 4
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const par = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!par) throw new Error('Parent not found');")
			b.Line("if (!('appendChild' in par)) throw new Error('Expected frame-like parent');")
			b.Line("const bar = figma.createRectangle();")
			b.Linef("bar.name = %q;", "Accent bar")
			b.Linef("bar.resize(%d, %d);", w, h)
			b.Linef("bar.x = %d;", x)
			b.Linef("bar.y = %d;", y)
			b.Linef("bar.cornerRadius = %s;", codegen.FmtFloat(float64(h)/2))
			b.Line("bar.fills = [{")
			b.Line("  type: 'GRADIENT_LINEAR',")
			b.Line("  gradientTransform: [[1, 0, 0], [0, 1, 0]],")
			b.Linef("  gradientStops: [{ position: 0, color: %s }, { position: 1, color: %s }],",
				codegen.FormatRGBA(c0, 1), codegen.FormatRGBA(c1, 1))
			b.Line("}];")
			b.Line("par.appendChild(bar);")
			b.ReturnIDs("bar.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&fromHex, "from", "#3366FF", "Start color (hex)")
	cmd.Flags().StringVar(&toHex, "to", "#14B8A6", "End color (hex)")
	cmd.Flags().IntVar(&w, "w", 240, "Bar width")
	cmd.Flags().IntVar(&h, "h", 4, "Bar height")
	cmd.Flags().IntVar(&x, "x", 0, "X offset inside parent")
	cmd.Flags().IntVar(&y, "y", 0, "Y offset inside parent")
	return cmd
}

// --- shadow ---

type shadowPreset struct {
	typ            string // DROP_SHADOW | INNER_SHADOW
	r, g, b, a     float64
	offX, offY     float64
	radius, spread float64
}

func shadowPresetByName(name string) (shadowPreset, error) {
	switch strings.ToLower(name) {
	case "sm":
		return shadowPreset{"DROP_SHADOW", 0, 0, 0, 0.18, 0, 2, 4, 0}, nil
	case "md":
		return shadowPreset{"DROP_SHADOW", 0, 0, 0, 0.22, 0, 4, 10, 0}, nil
	case "lg":
		return shadowPreset{"DROP_SHADOW", 0, 0, 0, 0.28, 0, 8, 24, 0}, nil
	case "xl":
		return shadowPreset{"DROP_SHADOW", 0, 0, 0, 0.32, 0, 12, 40, 0}, nil
	case "glow":
		return shadowPreset{"DROP_SHADOW", 0.35, 0.45, 0.95, 0.55, 0, 0, 28, 2}, nil
	case "inner":
		return shadowPreset{"INNER_SHADOW", 0, 0, 0, 0.35, 0, 4, 12, 0}, nil
	default:
		return shadowPreset{}, fmt.Errorf("preset must be sm, md, lg, xl, glow, or inner")
	}
}

func newFXShadowCmd() *cobra.Command {
	var preset string
	cmd := &cobra.Command{
		Use:   "shadow <nodeId>",
		Short: "Apply a drop/inner shadow preset",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := shadowPresetByName(preset)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("if (!('effects' in node)) throw new Error('Node has no effects');")
			b.Linef(`node.effects = [{
  type: %q,
  color: { r: %s, g: %s, b: %s, a: %s },
  offset: { x: %s, y: %s },
  radius: %s,
  spread: %s,
  visible: true,
  blendMode: 'NORMAL',
}];`,
				p.typ,
				codegen.FmtFloat(p.r), codegen.FmtFloat(p.g), codegen.FmtFloat(p.b), codegen.FmtFloat(p.a),
				codegen.FmtFloat(p.offX), codegen.FmtFloat(p.offY),
				codegen.FmtFloat(p.radius), codegen.FmtFloat(p.spread))
			b.ReturnIDs("node.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&preset, "preset", "md", "sm | md | lg | xl | glow | inner")
	return cmd
}

// --- parallax ---

func newFXParallaxLayerCmd() *cobra.Command {
	var layersN int
	cmd := &cobra.Command{
		Use:   "parallax-layer <parentId>",
		Short: "Layered depth stack (offset frames for parallax-style composition)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if layersN < 2 {
				return fmt.Errorf("--layers must be at least 2")
			}
			b := codegen.New()
			b.PageSetup(resolvePage())
			b.Linef("const parent = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!parent) throw new Error('Parent not found');")
			b.Line("if (!('appendChild' in parent)) throw new Error('Expected frame-like parent');")
			b.Line("if (!('width' in parent) || !('height' in parent)) throw new Error('Parent needs width/height');")
			b.Linef("const __n = %d;", layersN)
			b.Raw(`const ids = [];
for (let i = 0; i < __n; i++) {
  const layer = figma.createFrame();
  layer.name = 'Parallax ' + (i + 1);
  const inset = i * 10;
  layer.resize(Math.max(1, parent.width - inset * 2), Math.max(1, parent.height - inset * 2));
  layer.x = inset;
  layer.y = inset + i * 6;
  layer.fills = [{ type: 'SOLID', color: { r: 0.12 + i * 0.04, g: 0.14 + i * 0.03, b: 0.18 + i * 0.02 }, opacity: 0.25 + (__n - i) * 0.08 }];
  layer.effects = i === 0 ? [] : [{
    type: 'DROP_SHADOW',
    color: { r: 0, g: 0, b: 0, a: 0.15 + i * 0.04 },
    offset: { x: 0, y: 4 + i * 2 },
    radius: 8 + i * 4,
    spread: 0,
    visible: true,
    blendMode: 'NORMAL',
  }];
  parent.appendChild(layer);
  ids.push(layer.id);
}
`)
			b.ReturnExpr("{ createdNodeIds: ids }")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&layersN, "layers", 3, "Number of depth layers")
	return cmd
}
