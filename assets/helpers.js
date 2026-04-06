// Figma Design Kit — Core Helper Functions
// Inject this preamble into every use_figma call.
// All functions expect a loaded theme object with colors, type, effects, spacing.
//
// Usage: Copy the functions you need into your use_figma code block.
// The AI should inline these automatically when the figma-design-kit skill is active.

// ─────────────────────────────────────────────
// FONT LOADING
// ─────────────────────────────────────────────

async function loadFonts() {
  const fonts = [
    { family: 'Inter', style: 'Bold' },
    { family: 'Inter', style: 'Semi Bold' },
    { family: 'Inter', style: 'Medium' },
    { family: 'Inter', style: 'Regular' },
    { family: 'Inter', style: 'Light' },
    { family: 'Geist Mono', style: 'Regular' },
    { family: 'Geist Mono', style: 'Medium' },
  ];
  for (const fn of fonts) await figma.loadFontAsync(fn);
}

// ─────────────────────────────────────────────
// TEXT
// ─────────────────────────────────────────────

/**
 * Create a text node with Inter font.
 * @param {FrameNode} par   - Parent node
 * @param {string}    s     - Text content
 * @param {number}    x     - X position
 * @param {number}    y     - Y position
 * @param {number|null} w   - Width (null for auto-width)
 * @param {number}    sz    - Font size (default 14)
 * @param {string}    st    - Font style: 'Bold'|'Semi Bold'|'Medium'|'Regular'|'Light'
 * @param {RGB}       col   - Fill color {r,g,b}
 * @param {number|null} lh  - Line height in px (null for auto)
 * @param {string|null} al  - Horizontal align: 'LEFT'|'CENTER'|'RIGHT'
 * @returns {TextNode}
 */
function T(par, s, x, y, w, sz, st, col, lh, al) {
  const t = figma.createText();
  t.fontName = { family: 'Inter', style: st || 'Regular' };
  t.characters = s;
  t.fontSize = sz || 14;
  t.fills = [{ type: 'SOLID', color: col || { r: 0.96, g: 0.97, b: 0.98 } }];
  if (lh) t.lineHeight = { value: lh, unit: 'PIXELS' };
  if (al) t.textAlignHorizontal = al;
  t.x = x;
  t.y = y;
  if (w) {
    t.resize(w, t.height);
    t.textAutoResize = 'HEIGHT';
  }
  par.appendChild(t);
  return t;
}

/**
 * Create a Geist Mono label (uppercase tracking).
 * @param {FrameNode} par - Parent node
 * @param {string}    s   - Text content
 * @param {number}    x   - X position
 * @param {number}    y   - Y position
 * @param {number}    sz  - Font size (default 11)
 * @param {RGB}       col - Fill color
 * @returns {TextNode}
 */
function GM(par, s, x, y, sz, col) {
  const t = figma.createText();
  t.fontName = { family: 'Geist Mono', style: 'Medium' };
  t.characters = s;
  t.fontSize = sz || 11;
  t.fills = [{ type: 'SOLID', color: col || { r: 0.45, g: 0.48, b: 0.55 } }];
  t.letterSpacing = { value: 10, unit: 'PERCENT' };
  t.x = x;
  t.y = y;
  par.appendChild(t);
  return t;
}

// ─────────────────────────────────────────────
// SHAPES
// ─────────────────────────────────────────────

/**
 * Create a rectangle.
 * @param {FrameNode} par  - Parent node
 * @param {number}    x    - X position
 * @param {number}    y    - Y position
 * @param {number}    w    - Width
 * @param {number}    h    - Height
 * @param {RGB}       fill - Fill color (default: subtle divider)
 * @param {number}    cr   - Corner radius (default 0)
 * @returns {RectangleNode}
 */
function R(par, x, y, w, h, fill, cr) {
  const r = figma.createRectangle();
  r.x = x;
  r.y = y;
  r.resize(w, h);
  if (cr) r.cornerRadius = cr;
  r.fills = [{ type: 'SOLID', color: fill || { r: 0.14, g: 0.16, b: 0.22 } }];
  par.appendChild(r);
  return r;
}

// ─────────────────────────────────────────────
// GLASSMORPHISM CARD
// ─────────────────────────────────────────────

/**
 * Create a glassmorphism card with frosted glass effect.
 * @param {FrameNode} par - Parent node
 * @param {number}    x   - X position
 * @param {number}    y   - Y position
 * @param {number}    w   - Width
 * @param {number}    h   - Height
 * @param {object}    o   - Options:
 *   @param {number}  o.r    - Corner radius (default 20)
 *   @param {number}  o.f    - Fill opacity, white overlay (default 0.04)
 *   @param {number}  o.s    - Stroke opacity, white border (default 0.08)
 *   @param {number}  o.ga   - Glow alpha, blue halo (default 0.06)
 *   @param {number}  o.bl   - Background blur radius (default 24)
 *   @param {boolean} o.glow - Enable blue glow shadow (default true)
 *   @param {boolean} o.clip - Clip content (default true)
 * @returns {FrameNode}
 */
function G(par, x, y, w, h, o) {
  o = o || {};
  const fr = figma.createFrame();
  fr.x = x;
  fr.y = y;
  fr.resize(w, h);
  fr.cornerRadius = o.r || 20;
  fr.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: o.f || 0.04 }];
  fr.strokes = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: o.s || 0.08 }];
  fr.strokeWeight = 1;
  const fx = [{ type: 'BACKGROUND_BLUR', radius: o.bl || 24, visible: true }];
  if (o.glow !== false) {
    fx.push({
      type: 'DROP_SHADOW',
      color: { r: 0.2, g: 0.4, b: 1, a: o.ga || 0.06 },
      offset: { x: 0, y: 0 },
      radius: 24,
      spread: 0,
      visible: true,
      blendMode: 'NORMAL',
    });
  }
  fr.effects = fx;
  fr.clipsContent = o.clip !== false;
  par.appendChild(fr);
  return fr;
}

// ─────────────────────────────────────────────
// BRANDING
// ─────────────────────────────────────────────

/**
 * Add branded header (gradient accent line) and footer (logo + slide counter).
 * @param {FrameNode} slide      - The slide frame
 * @param {number}    slideNum   - Current slide number (1-based)
 * @param {object}    [opts]     - Options:
 *   @param {number}  opts.total - Total slides (default 7)
 *   @param {number}  opts.sw    - Slide width (default from slide)
 *   @param {number}  opts.sh    - Slide height (default from slide)
 *   @param {number}  opts.pad   - Horizontal padding (default 80)
 */
function brand(slide, slideNum, opts) {
  opts = opts || {};
  const sw = opts.sw || slide.width;
  const sh = opts.sh || slide.height;
  const pad = opts.pad || 80;
  const total = opts.total || 7;

  // Top gradient accent line
  const accent = figma.createRectangle();
  accent.x = 0;
  accent.y = 0;
  accent.resize(sw, 4);
  accent.fills = [{
    type: 'GRADIENT_LINEAR',
    gradientTransform: [[1, 0, 0], [0, 1, 0]],
    gradientStops: [
      { position: 0, color: { r: 0.2, g: 0.4, b: 1, a: 1 } },
      { position: 1, color: { r: 0.08, g: 0.72, b: 0.65, a: 1 } },
    ],
  }];
  slide.appendChild(accent);

  // Bottom logo
  GM(slide, '[\u039B] ARKHAM', pad, sh - 56, 11, { r: 0.35, g: 0.38, b: 0.45 });

  // Slide counter
  const num = String(slideNum).padStart(2, '0');
  const tot = String(total).padStart(2, '0');
  GM(slide, num + '/' + tot, sw - pad - 42, sh - 56, 11, { r: 0.45, g: 0.48, b: 0.55 });
}

// ─────────────────────────────────────────────
// SLIDE FACTORY
// ─────────────────────────────────────────────

/**
 * Create a positioned slide frame with optional gradient glow background.
 * @param {PageNode}  parent   - Page to append to
 * @param {number}    idx      - Slide index (0-based, controls x position)
 * @param {string}    name     - Frame name
 * @param {Paint[]}   [fills]  - Background fills (default: solid dark)
 * @param {object}    [opts]   - Options:
 *   @param {number}  opts.sw     - Width (default 1080)
 *   @param {number}  opts.sh     - Height (default 1350)
 *   @param {number}  opts.gap    - Gap between slides (default 60)
 *   @param {number}  opts.startX - X offset for first slide (default 0)
 * @returns {FrameNode}
 */
function mkSlide(parent, idx, name, fills, opts) {
  opts = opts || {};
  const sw = opts.sw || 1080;
  const sh = opts.sh || 1350;
  const gap = opts.gap || 60;
  const startX = opts.startX || 0;

  const s = figma.createFrame();
  s.name = name;
  s.resize(sw, sh);
  s.x = startX + idx * (sw + gap);
  s.y = 0;
  s.fills = fills || [{ type: 'SOLID', color: { r: 0.05, g: 0.06, b: 0.09 } }];
  s.clipsContent = true;
  parent.appendChild(s);
  return s;
}

// ─────────────────────────────────────────────
// ANNOTATION STRIP (Storyboard)
// ─────────────────────────────────────────────

/**
 * Create a glassmorphism annotation strip for storyboard styleframes.
 * @param {FrameNode} par   - Parent styleframe (1920x1080)
 * @param {number}    num   - Scene number (1-5)
 * @param {string}    title - Scene title in quotes
 * @param {string}    time  - Timecode range (e.g., "0:00 — 0:02.5")
 * @param {string}    dur   - Duration (e.g., "2.5s")
 * @param {string}    cam   - Camera/motion direction
 * @param {string}    sound - Sound/music description
 * @returns {FrameNode}
 */
function annot(par, num, title, time, dur, cam, sound) {
  const fw = par.width;
  const fh = par.height;
  const s = G(par, 40, fh - 160, fw - 80, 120, { r: 16, f: 0.06, s: 0.10, ga: 0.08 });
  GM(s, 'ESCENA 0' + num, 24, 16, 12, { r: 0.2, g: 0.4, b: 1.0 });
  T(s, title, 160, 14, null, 14, 'Bold', { r: 0.96, g: 0.97, b: 0.98 });
  GM(s, time, fw - 380, 16, 11, { r: 0.45, g: 0.48, b: 0.55 });
  GM(s, dur, fw - 180, 16, 11, { r: 0.08, g: 0.72, b: 0.65 });
  GM(s, 'C\u00C1MARA:', 24, 46, 9, { r: 0.35, g: 0.38, b: 0.45 });
  T(s, cam, 24, 62, Math.floor((fw - 80) / 2) - 40, 12, 'Regular', { r: 0.78, g: 0.81, b: 0.85 }, 18);
  GM(s, 'SONIDO:', Math.floor((fw - 80) / 2), 46, 9, { r: 0.35, g: 0.38, b: 0.45 });
  T(s, sound, Math.floor((fw - 80) / 2), 62, Math.floor((fw - 80) / 2) - 40, 12, 'Regular', { r: 0.55, g: 0.58, b: 0.65 }, 18);
  return s;
}

// ─────────────────────────────────────────────
// GRADIENT GLOW PRESETS
// ─────────────────────────────────────────────

/**
 * Generate background fill arrays with subtle radial gradient glows.
 * Use these as the `fills` parameter for mkSlide() or any frame.
 * @param {RGB}    bg       - Base background color
 * @param {string} position - Glow position: 'topRight'|'center'|'subtle'|'cta'
 * @returns {Paint[]}
 */
function glowFills(bg, position) {
  const base = { type: 'SOLID', color: bg };
  const presets = {
    topRight: [
      base,
      { type: 'GRADIENT_RADIAL', gradientTransform: [[1.8, 0, 0.35], [0, 1.5, -0.1]],
        gradientStops: [
          { position: 0, color: { r: 0.12, g: 0.22, b: 0.55, a: 0.14 } },
          { position: 1, color: { ...bg, a: 0 } },
        ]},
      { type: 'GRADIENT_RADIAL', gradientTransform: [[1.0, 0, -0.15], [0, 0.8, 0.45]],
        gradientStops: [
          { position: 0, color: { r: 0.05, g: 0.28, b: 0.30, a: 0.10 } },
          { position: 1, color: { ...bg, a: 0 } },
        ]},
    ],
    center: [
      base,
      { type: 'GRADIENT_RADIAL', gradientTransform: [[2.0, 0, -0.1], [0, 1.8, -0.05]],
        gradientStops: [
          { position: 0, color: { r: 0.12, g: 0.22, b: 0.50, a: 0.18 } },
          { position: 0.8, color: { ...bg, a: 0 } },
        ]},
    ],
    subtle: [
      base,
      { type: 'GRADIENT_RADIAL', gradientTransform: [[2.0, 0, 0.2], [0, 1.6, -0.1]],
        gradientStops: [
          { position: 0, color: { r: 0.10, g: 0.18, b: 0.40, a: 0.10 } },
          { position: 1, color: { ...bg, a: 0 } },
        ]},
    ],
    cta: [
      base,
      { type: 'GRADIENT_RADIAL', gradientTransform: [[1.5, 0, 0.0], [0, 1.2, -0.05]],
        gradientStops: [
          { position: 0, color: { r: 0.14, g: 0.26, b: 0.58, a: 0.22 } },
          { position: 0.7, color: { ...bg, a: 0 } },
        ]},
      { type: 'GRADIENT_RADIAL', gradientTransform: [[1.0, 0, 0.1], [0, 0.8, 0.3]],
        gradientStops: [
          { position: 0, color: { r: 0.05, g: 0.28, b: 0.30, a: 0.12 } },
          { position: 1, color: { ...bg, a: 0 } },
        ]},
    ],
  };
  return presets[position] || presets.subtle;
}

// ─────────────────────────────────────────────
// GRADIENT ACCENT BAR
// ─────────────────────────────────────────────

/**
 * Create a horizontal gradient accent bar (blue -> teal).
 * @param {FrameNode} par - Parent node
 * @param {number}    x   - X position
 * @param {number}    y   - Y position
 * @param {number}    w   - Width
 * @param {number}    h   - Height (default 4)
 * @returns {RectangleNode}
 */
function accentBar(par, x, y, w, h) {
  const bar = figma.createRectangle();
  bar.x = x;
  bar.y = y;
  bar.resize(w, h || 4);
  bar.cornerRadius = (h || 4) / 2;
  bar.fills = [{
    type: 'GRADIENT_LINEAR',
    gradientTransform: [[1, 0, 0], [0, 1, 0]],
    gradientStops: [
      { position: 0, color: { r: 0.2, g: 0.4, b: 1, a: 1 } },
      { position: 1, color: { r: 0.08, g: 0.72, b: 0.65, a: 1 } },
    ],
  }];
  par.appendChild(bar);
  return bar;
}
