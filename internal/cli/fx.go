package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dop-amine/figma-kit/internal/codegen"
)

func newFXCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fx",
		Short: "Visual effects (glow, mesh, noise, vignette, grain, blur, shadow, ...)",
		Example: `  # "Add a colorful mesh gradient to the hero background"
  figma-kit fx mesh <frameId> -t noir

  # "Add a subtle glow effect behind the card"
  figma-kit fx glow <frameId> -t noir

  # "Apply film grain for texture"
  figma-kit fx grain <frameId>`,
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
	cmd.AddCommand(newFXAuroraCmd())
	cmd.AddCommand(newFXMorphCmd())
	cmd.AddCommand(newFXGradientBorderCmd())
	cmd.AddCommand(newFXSpotlightCmd())
	cmd.AddCommand(newFXPatternCmd())
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

// --- aurora ---

func newFXAuroraCmd() *cobra.Command {
	var palette string
	cmd := &cobra.Command{
		Use:     "aurora <nodeId>",
		Short:   "Apply aurora borealis gradient overlay to a frame",
		Example: `  figma-kit fx aurora <frameId> --palette sunset`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")

			var colors string
			switch palette {
			case "sunset":
				colors = "{r:0.95,g:0.3,b:0.2},{r:0.85,g:0.15,b:0.5},{r:0.6,g:0.2,b:0.8},{r:0.3,g:0.2,b:0.7}"
			case "custom":
				colors = "{r:0.2,g:0.8,b:0.6},{r:0.4,g:0.3,b:0.9},{r:0.8,g:0.2,b:0.5},{r:0.2,g:0.5,b:0.9}"
			default: // northern
				colors = "{r:0.1,g:0.8,b:0.4},{r:0.2,g:0.5,b:0.9},{r:0.5,g:0.2,b:0.8},{r:0.1,g:0.6,b:0.7}"
			}

			b.Linef("const auroraColors = [%s];", colors)
			b.Line("const existing = node.fills ? [...node.fills] : [];")
			b.Line("const newFills = auroraColors.map((c, i) => ({")
			b.Line("  type: 'GRADIENT_RADIAL',")
			b.Line("  gradientStops: [")
			b.Line("    {position: 0, color: {...c, a: 0.6 - i * 0.1}},")
			b.Line("    {position: 1, color: {...c, a: 0}}")
			b.Line("  ],")
			b.Line("  gradientTransform: [")
			b.Line("    [1.5 + i * 0.3, 0, -0.2 + i * 0.15],")
			b.Line("    [0, 1.2 + i * 0.2, -0.1 + i * 0.1]")
			b.Line("  ],")
			b.Line("  blendMode: i % 2 === 0 ? 'SCREEN' : 'LINEAR_DODGE',")
			b.Line("  opacity: 0.7,")
			b.Line("  visible: true")
			b.Line("}));")
			b.Line("node.fills = [...existing, ...newFills];")
			b.ReturnExpr("'Applied aurora effect with ' + auroraColors.length + ' layers'")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&palette, "palette", "northern", "Color palette: northern, sunset, custom")
	return cmd
}

// --- morph ---

func newFXMorphCmd() *cobra.Command {
	var (
		count   int
		spread  int
		palette string
	)
	cmd := &cobra.Command{
		Use:     "morph <nodeId>",
		Short:   "Add organic blob shapes as background elements",
		Example: `  figma-kit fx morph <frameId> --count 5 --spread 200`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")

			var colors string
			switch palette {
			case "warm":
				colors = "[{r:0.95,g:0.3,b:0.3},{r:0.95,g:0.6,b:0.2},{r:0.9,g:0.4,b:0.5}]"
			case "cool":
				colors = "[{r:0.2,g:0.5,b:0.9},{r:0.4,g:0.8,b:0.7},{r:0.3,g:0.3,b:0.8}]"
			default:
				colors = "[{r:0.5,g:0.3,b:0.9},{r:0.3,g:0.7,b:0.9},{r:0.9,g:0.4,b:0.6}]"
			}

			b.Linef("const blobColors = %s;", colors)
			b.Linef("const spread = %d;", spread)
			b.Linef("const count = %d;", count)
			b.Line("const w = node.width || 400; const h = node.height || 300;")
			b.Line("for (let i = 0; i < count; i++) {")
			b.Line("  const blob = figma.createEllipse();")
			b.Line("  const size = 80 + (i * 37) % 120;")
			b.Line("  blob.resize(size, size * (0.7 + (i * 0.13) % 0.6));")
			b.Line("  blob.x = (i * 97 + 30) % w;")
			b.Line("  blob.y = (i * 73 + 20) % h;")
			b.Line("  blob.name = 'blob-' + i;")
			b.Line("  const c = blobColors[i % blobColors.length];")
			b.Line("  blob.fills = [{type:'SOLID', color:c, opacity:0.15 + (i % 3) * 0.05}];")
			b.Line("  blob.effects = [{type:'LAYER_BLUR', radius:30 + i * 5, visible:true}];")
			b.Line("  if ('appendChild' in node) node.appendChild(blob);")
			b.Line("  else node.parent.appendChild(blob);")
			b.Line("}")
			b.ReturnExpr("'Added ' + count + ' blob shapes'")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().IntVar(&count, "count", 4, "Number of blobs")
	cmd.Flags().IntVar(&spread, "spread", 150, "Spread area in pixels")
	cmd.Flags().StringVar(&palette, "palette", "default", "Color palette: default, warm, cool")
	return cmd
}

// --- gradient-border ---

func newFXGradientBorderCmd() *cobra.Command {
	var (
		fromHex string
		toHex   string
		width   int
	)
	cmd := &cobra.Command{
		Use:     "gradient-border <nodeId>",
		Short:   "Simulate a gradient border around a node",
		Example: `  figma-kit fx gradient-border <frameId> --from "#3B82F6" --to "#8B5CF6"`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			fromRGB, e1 := codegen.HexToRGB(fromHex)
			if e1 != nil {
				fromRGB = codegen.RGB{R: 0.23, G: 0.51, B: 0.96}
			}
			toRGB, e2 := codegen.HexToRGB(toHex)
			if e2 != nil {
				toRGB = codegen.RGB{R: 0.55, G: 0.36, B: 0.96}
			}

			b := codegen.New()
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const w = node.width; const h = node.height;")
			b.Line("const parent = node.parent;")
			b.Linef("const bw = %d;", width)
			b.Line("const outer = figma.createFrame();")
			b.Line("outer.name = node.name + ' / gradient-border';")
			b.Line("outer.resize(w + bw * 2, h + bw * 2);")
			b.Line("outer.x = ('x' in node ? node.x : 0) - bw;")
			b.Line("outer.y = ('y' in node ? node.y : 0) - bw;")
			b.Line("outer.cornerRadius = ('cornerRadius' in node ? node.cornerRadius + bw : bw);")
			b.Linef("outer.fills = [{type:'GRADIENT_LINEAR', gradientStops:[")
			b.Linef("  {position:0, color:{r:%.3f,g:%.3f,b:%.3f,a:1}},", fromRGB.R, fromRGB.G, fromRGB.B)
			b.Linef("  {position:1, color:{r:%.3f,g:%.3f,b:%.3f,a:1}}", toRGB.R, toRGB.G, toRGB.B)
			b.Line("], gradientTransform:[[1,0,0],[0,1,0]]}];")
			b.Line("outer.clipsContent = true;")
			b.Line("if (parent && 'appendChild' in parent) {")
			b.Line("  const idx = parent.children.indexOf(node);")
			b.Line("  parent.insertChild(idx, outer);")
			b.Line("}")
			b.Line("node.x = bw; node.y = bw;")
			b.Line("outer.appendChild(node);")
			b.ReturnIDs("outer.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&fromHex, "from", "#3B82F6", "Start gradient color (hex)")
	cmd.Flags().StringVar(&toHex, "to", "#8B5CF6", "End gradient color (hex)")
	cmd.Flags().IntVar(&width, "width", 2, "Border width in pixels")
	return cmd
}

// --- spotlight ---

func newFXSpotlightCmd() *cobra.Command {
	var (
		x, y      float64
		radius    int
		intensity float64
	)
	cmd := &cobra.Command{
		Use:     "spotlight <nodeId>",
		Short:   "Add a circular spotlight/highlight effect",
		Example: `  figma-kit fx spotlight <frameId> --x 0.5 --y 0.3 --radius 300`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			b := codegen.New()
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const existing = node.fills ? [...node.fills] : [];")
			b.Linef("const spotlight = {")
			b.Line("  type: 'GRADIENT_RADIAL',")
			b.Line("  gradientStops: [")
			b.Linef("    {position: 0, color: {r:1, g:1, b:1, a:%.2f}},", intensity)
			b.Line("    {position: 0.6, color: {r:1, g:1, b:1, a:0.02}},")
			b.Line("    {position: 1, color: {r:0, g:0, b:0, a:0.15}}")
			b.Line("  ],")
			b.Linef("  gradientTransform: [[%.2f, 0, %.2f], [0, %.2f, %.2f]],",
				float64(radius)/200.0, x-float64(radius)/400.0,
				float64(radius)/200.0, y-float64(radius)/400.0)
			b.Line("  visible: true")
			b.Line("};")
			b.Line("node.fills = [...existing, spotlight];")
			b.ReturnExpr("'Applied spotlight effect'")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().Float64Var(&x, "x", 0.5, "Spotlight center X (0-1)")
	cmd.Flags().Float64Var(&y, "y", 0.3, "Spotlight center Y (0-1)")
	cmd.Flags().IntVar(&radius, "radius", 200, "Spotlight radius in pixels")
	cmd.Flags().Float64Var(&intensity, "intensity", 0.25, "Spotlight intensity (0-1)")
	return cmd
}

// --- pattern ---

func newFXPatternCmd() *cobra.Command {
	var (
		patternType string
		spacing     int
		size        int
		colorHex    string
	)
	cmd := &cobra.Command{
		Use:     "pattern <nodeId>",
		Short:   "Add a repeating geometric pattern to a frame",
		Example: `  figma-kit fx pattern <frameId> --type dots --spacing 24 --size 3`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := resolveTheme(cmd)
			if err != nil {
				return err
			}
			rgb, cErr := codegen.HexToRGB(colorHex)
			if cErr != nil {
				rgb = codegen.RGB{R: 0.4, G: 0.4, B: 0.45}
			}

			b := codegen.New()
			b.Linef("const node = await figma.getNodeByIdAsync(%q);", args[0])
			b.Line("if (!node) throw new Error('Node not found');")
			b.Line("const w = node.width || 400; const h = node.height || 300;")
			b.Linef("const spacing = %d;", spacing)
			b.Linef("const sz = %d;", size)
			b.Linef("const pColor = {r:%.3f,g:%.3f,b:%.3f};", rgb.R, rgb.G, rgb.B)
			b.Line("const container = figma.createFrame();")
			b.Line("container.name = 'pattern-' + " + fmt.Sprintf("%q", patternType) + ";")
			b.Line("container.resize(w, h);")
			b.Line("container.fills = [];")
			b.Line("container.clipsContent = true;")

			switch patternType {
			case "lines":
				b.Line("for (let x = 0; x < w; x += spacing) {")
				b.Line("  const line = figma.createRectangle();")
				b.Line("  line.resize(1, h); line.x = x; line.y = 0;")
				b.Line("  line.fills = [{type:'SOLID', color:pColor, opacity:0.2}];")
				b.Line("  container.appendChild(line);")
				b.Line("}")
			case "crosses":
				b.Line("for (let y = spacing; y < h; y += spacing) {")
				b.Line("  for (let x = spacing; x < w; x += spacing) {")
				b.Line("    const h1 = figma.createRectangle(); h1.resize(sz*2, 1); h1.x = x-sz; h1.y = y;")
				b.Line("    h1.fills = [{type:'SOLID', color:pColor, opacity:0.2}]; container.appendChild(h1);")
				b.Line("    const v1 = figma.createRectangle(); v1.resize(1, sz*2); v1.x = x; v1.y = y-sz;")
				b.Line("    v1.fills = [{type:'SOLID', color:pColor, opacity:0.2}]; container.appendChild(v1);")
				b.Line("  }")
				b.Line("}")
			case "diagonal":
				b.Line("for (let i = -Math.max(w,h); i < Math.max(w,h)*2; i += spacing) {")
				b.Line("  const line = figma.createRectangle();")
				b.Line("  line.resize(1, Math.max(w,h)*2); line.x = i; line.y = -h/2;")
				b.Line("  line.rotation = 45;")
				b.Line("  line.fills = [{type:'SOLID', color:pColor, opacity:0.15}];")
				b.Line("  container.appendChild(line);")
				b.Line("}")
			case "grid":
				b.Line("for (let x = 0; x < w; x += spacing) {")
				b.Line("  const vl = figma.createRectangle(); vl.resize(1, h); vl.x = x; vl.y = 0;")
				b.Line("  vl.fills = [{type:'SOLID', color:pColor, opacity:0.1}]; container.appendChild(vl);")
				b.Line("}")
				b.Line("for (let y = 0; y < h; y += spacing) {")
				b.Line("  const hl = figma.createRectangle(); hl.resize(w, 1); hl.x = 0; hl.y = y;")
				b.Line("  hl.fills = [{type:'SOLID', color:pColor, opacity:0.1}]; container.appendChild(hl);")
				b.Line("}")
			default: // dots
				b.Line("for (let y = spacing; y < h; y += spacing) {")
				b.Line("  for (let x = spacing; x < w; x += spacing) {")
				b.Line("    const dot = figma.createEllipse();")
				b.Line("    dot.resize(sz, sz); dot.x = x - sz/2; dot.y = y - sz/2;")
				b.Line("    dot.fills = [{type:'SOLID', color:pColor, opacity:0.25}];")
				b.Line("    container.appendChild(dot);")
				b.Line("  }")
				b.Line("}")
			}

			b.Line("if ('appendChild' in node) node.appendChild(container);")
			b.Line("else node.parent.appendChild(container);")
			b.ReturnIDs("container.id")
			output(b.String())
			return nil
		},
	}
	cmd.Flags().StringVar(&patternType, "type", "dots", "Pattern type: dots, lines, crosses, diagonal, grid")
	cmd.Flags().IntVar(&spacing, "spacing", 24, "Space between pattern elements")
	cmd.Flags().IntVar(&size, "size", 3, "Size of pattern elements")
	cmd.Flags().StringVar(&colorHex, "color", "#6B7280", "Pattern color (hex)")
	return cmd
}
