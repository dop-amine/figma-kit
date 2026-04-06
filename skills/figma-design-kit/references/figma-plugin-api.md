# Figma Plugin API — Gotchas & Patterns

Reference for common pitfalls when writing `use_figma` code. Every issue below has been encountered and debugged in real production work.

## Critical Gotchas

### Font Style Naming
Inter uses **spaced** style names. This fails silently if wrong.
```
CORRECT: { family: 'Inter', style: 'Semi Bold' }
WRONG:   { family: 'Inter', style: 'SemiBold' }
```
Same for `Extra Bold`, `Extra Light`, etc.

### Setting Current Page
Never assign `figma.currentPage`. Use the async method:
```javascript
// CORRECT
await figma.setCurrentPageAsync(page);

// WRONG — will throw
figma.currentPage = page;
```

### DROP_SHADOW Requires blendMode
Every `DROP_SHADOW` and `INNER_SHADOW` effect MUST include `blendMode`:
```javascript
{
  type: 'DROP_SHADOW',
  color: { r: 0.2, g: 0.4, b: 1, a: 0.12 },
  offset: { x: 0, y: 4 },
  radius: 24,
  spread: 0,
  visible: true,
  blendMode: 'NORMAL'  // REQUIRED — omitting causes validation error
}
```

### Stroke Color Has No Alpha
Stroke paint colors use `{r, g, b}` only. For semi-transparent strokes, use `opacity` on the paint object:
```javascript
// CORRECT
fr.strokes = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.08 }];

// WRONG — 'a' field causes validation error
fr.strokes = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1, a: 0.08 } }];
```

### Fill Paint Opacity
Fill paints support an `opacity` field on the Paint object (not in the color):
```javascript
fr.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.04 }];
```

### appendChild Is Required
Newly created nodes are NOT automatically parented. Always call:
```javascript
parent.appendChild(child);
```
Without this, the node exists in memory but is invisible.

### Page Children Reporting
After creating frames via `use_figma`, querying `page.children` in a *separate* call may report the page as empty. This is a known quirk of the MCP bridge. Workarounds:
- Consolidate all creation and querying into a single `use_figma` call
- Use `get_screenshot` to visually confirm frames exist
- Re-set the current page before querying: `await figma.setCurrentPageAsync(pg)`

## Effect Types

| Type | Required Fields |
|------|----------------|
| `DROP_SHADOW` | color (RGBA), offset, radius, spread, visible, blendMode |
| `INNER_SHADOW` | color (RGBA), offset, radius, spread, visible, blendMode |
| `LAYER_BLUR` | radius, visible |
| `BACKGROUND_BLUR` | radius, visible |

`BACKGROUND_BLUR` blurs content *behind* the frame in z-order, including the parent frame's fills. This is how glassmorphism works.

## Gradient Transforms

The `gradientTransform` is a 2x3 affine matrix `[[a,b,tx],[c,d,ty]]`:
- Identity: `[[1,0,0],[0,1,0]]` — default horizontal gradient
- Scale X 2x, offset right: `[[2,0,0.3],[0,1,0]]`
- Centered radial: `[[2,0,-0.1],[0,1.8,-0.05]]`

For radial gradients, the transform controls the ellipse's position and scale. Larger values = wider/taller falloff. Negative tx/ty shifts the center left/up.

## Text Properties

```javascript
const t = figma.createText();
t.fontName = { family: 'Inter', style: 'Regular' };  // Set BEFORE characters
t.characters = 'Hello';
t.fontSize = 16;
t.fills = [{ type: 'SOLID', color: { r: 0, g: 0, b: 0 } }];
t.lineHeight = { value: 24, unit: 'PIXELS' };  // or { unit: 'AUTO' }
t.letterSpacing = { value: 5, unit: 'PERCENT' };
t.textAlignHorizontal = 'CENTER';  // 'LEFT' | 'CENTER' | 'RIGHT' | 'JUSTIFIED'
t.textAutoResize = 'HEIGHT';  // 'NONE' | 'WIDTH_AND_HEIGHT' | 'HEIGHT'
```

## Frame vs Rectangle

- Use `figma.createFrame()` when you need children (containers, cards, layouts)
- Use `figma.createRectangle()` for pure decoration (dividers, accent bars, backgrounds)
- Frames support `clipsContent`, `layoutMode`, `itemSpacing`, `padding*`
- Rectangles are lighter weight and render faster

## Auto Layout

```javascript
frame.layoutMode = 'HORIZONTAL';  // or 'VERTICAL'
frame.itemSpacing = 16;
frame.paddingLeft = 24;
frame.paddingRight = 24;
frame.paddingTop = 16;
frame.paddingBottom = 16;
frame.counterAxisAlignItems = 'CENTER';  // cross-axis alignment
frame.primaryAxisAlignItems = 'SPACE_BETWEEN';  // main-axis distribution
```

## Output Limits
- `use_figma` return values have a 20KB limit
- Keep return strings concise — return IDs and names, not full node trees
- For large operations, split into multiple `use_figma` calls
