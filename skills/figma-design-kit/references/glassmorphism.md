# Glassmorphism Recipes

Standardized frosted-glass card patterns for the Figma Design Kit.

## Anatomy of a Glass Card

A glassmorphism card in Figma consists of four layers:

1. **Fill** — Semi-transparent white overlay (2-8% opacity)
2. **Stroke** — Subtle white border (6-12% opacity, 1px)
3. **Background Blur** — Frosted effect (16-28px radius)
4. **Glow Shadow** — Blue halo around the card (optional)

## Presets

### Subtle Glass
Low-contrast, minimal presence. For secondary containers and pills.
```javascript
G(parent, x, y, w, h, { r: 16, f: 0.03, s: 0.06, ga: 0.04, bl: 20 })
```

### Default Glass
Standard holding shape. For content cards and sections.
```javascript
G(parent, x, y, w, h, { r: 20, f: 0.04, s: 0.08, ga: 0.06, bl: 24 })
```

### Strong Glass
High presence, hero cards. For headlines and featured content.
```javascript
G(parent, x, y, w, h, { r: 24, f: 0.06, s: 0.12, ga: 0.12, bl: 24 })
```

### Glass Pill
Rounded chip/tag. For CTAs, feature labels, metadata.
```javascript
G(parent, x, y, w, h, { r: h/2, f: 0.06, s: 0.10, ga: 0.06, bl: 20 })
```

### No-Glow Glass
Clean card without blue halo. For neutral containers.
```javascript
G(parent, x, y, w, h, { r: 20, f: 0.04, s: 0.08, glow: false, bl: 24 })
```

## Making Glass Visible

Glass cards need something behind them for the blur to affect. Options:

### 1. Gradient Background Fills (Preferred)
Add radial gradient fills to the parent frame. The glass card blurs these.
```javascript
slide.fills = [
  { type: 'SOLID', color: BG },
  { type: 'GRADIENT_RADIAL',
    gradientTransform: [[1.8, 0, 0.35], [0, 1.5, -0.1]],
    gradientStops: [
      { position: 0, color: { r: 0.12, g: 0.22, b: 0.55, a: 0.14 } },
      { position: 1, color: { ...BG, a: 0 } }
    ]}
];
```
Use the `glowFills()` helper for standard presets.

### 2. Drop Shadow Glow
The `ga` parameter adds a blue `DROP_SHADOW` at offset (0,0) which radiates evenly around the card, creating a glow halo.

### 3. Layered Cards
Place glass cards over other elements (images, shapes, other cards) for natural depth.

## Light Mode Glass

For light themes, increase fill opacity significantly:
```javascript
G(parent, x, y, w, h, { r: 20, f: 0.60, s: 0.15, ga: 0.04, bl: 16 })
```
The white overlay is 60% instead of 4%, creating an opaque-ish frosted panel.

## Content Inside Glass Cards

When placing text inside a glass card, position relative to the card's origin:
```javascript
const card = G(slide, 80, 220, 920, 400, { r: 24, f: 0.05, s: 0.10 });
T(card, 'Headline', 28, 28, 864, 72, 'Bold', WT, 86);  // relative to card
T(card, 'Subtitle', 28, 300, 864, 20, 'Regular', BD);
accentBar(card, 28, 340, 60, 4);
```

## Composition Rules

1. Every distinct content group should be wrapped in a glass card
2. Use stronger glass for primary content, subtle for secondary
3. Maintain 20-28px inner padding within glass cards
4. Keep glass card cornerRadius between 16-28px (larger = more premium)
5. Don't nest glass cards (blur-on-blur creates visual noise)
6. Use glass pills for small interactive elements (CTAs, chips, tags)
