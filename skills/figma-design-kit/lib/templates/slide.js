// Template: LinkedIn Carousel Slide
// Format: 1080x1350 (4:5 portrait), optimized for mobile LinkedIn feed
// Requires: helpers.js functions (T, GM, R, G, brand, mkSlide, glowFills, accentBar)
//
// Usage: Call createSlide() with a slide config object.
// This template handles layout, branding, and glow backgrounds.
// Customize content by passing different config objects.

/**
 * Create a single LinkedIn carousel slide.
 *
 * @param {PageNode} page    - Figma page to append to
 * @param {object}   config  - Slide configuration:
 *   @param {number}  config.index       - Slide index (0-based)
 *   @param {string}  config.name        - Frame name
 *   @param {number}  config.total       - Total slides (for counter)
 *   @param {string}  config.glowType    - Glow preset: 'topRight'|'center'|'subtle'|'cta'
 *   @param {string}  config.headline    - Main headline text
 *   @param {string}  [config.subtitle]  - Subtitle text
 *   @param {boolean} [config.centered]  - Center-align text (default: left)
 *   @param {object}  [config.cta]       - CTA button config: { text, centered }
 *   @param {string[]} [config.chips]    - Feature chip labels
 * @param {object}   theme   - Theme config object
 * @returns {FrameNode}
 */
function createSlide(page, config, theme) {
  const SW = theme.spacing.slide.width;
  const SH = theme.spacing.slide.height;
  const PAD = theme.spacing.slide.padding;
  const CW = SW - PAD * 2;
  const BG = theme.colors.BG;
  const al = config.centered ? 'CENTER' : 'LEFT';

  // Create slide with glow background
  const fills = glowFills(BG, config.glowType || 'subtle');
  const slide = mkSlide(page, config.index, config.name, fills, {
    sw: SW, sh: SH, gap: 60,
  });

  // Branding
  brand(slide, config.index + 1, { total: config.total || 7 });

  // Headline glass card
  const headlineLines = (config.headline.match(/\n/g) || []).length + 1;
  const headlineH = headlineLines * theme.type.h1.lineHeight + 80;
  const gc = G(slide, PAD, 220, CW, headlineH, theme.effects.glass.strong);
  T(gc, config.headline, 28, 28, CW - 56,
    theme.type.h1.fontSize, theme.type.h1.style, theme.colors.WT,
    theme.type.h1.lineHeight, al);

  // Subtitle
  if (config.subtitle) {
    T(gc, config.subtitle, 28, headlineH - 50, CW - 56,
      theme.type.body.fontSize, theme.type.body.style, theme.colors.BD,
      null, al);
  }

  // Accent bar
  const barX = config.centered ? CW / 2 - 30 : 28;
  accentBar(gc, barX, headlineH - 20, 60, 4);

  // CTA button
  if (config.cta) {
    const ctaW = 360;
    const ctaX = config.cta.centered ? SW / 2 - ctaW / 2 : PAD;
    const cta = figma.createFrame();
    cta.x = ctaX; cta.y = SH - 180; cta.resize(ctaW, 56);
    cta.cornerRadius = 28;
    cta.fills = [theme.gradients.cta];
    cta.effects = [theme.effects.shadow.glow];
    slide.appendChild(cta);
    T(cta, config.cta.text, 0, 17, ctaW, 16, 'Semi Bold',
      { r: 1, g: 1, b: 1 }, null, 'CENTER');
  }

  // Feature chips
  if (config.chips && config.chips.length > 0) {
    const chipRow = G(slide, PAD - 10, SH - 136, CW + 20, 48,
      theme.effects.glass.pill);
    const chipW = Math.floor((CW + 20) / config.chips.length);
    config.chips.forEach((ch, i) => {
      T(chipRow, ch, i * chipW, 14, chipW, 13, 'Medium',
        theme.colors.BD, null, 'CENTER');
    });
  }

  return slide;
}

// Example usage:
//
// createSlide(page, {
//   index: 0,
//   name: 'Slide 1 - Cover',
//   total: 7,
//   glowType: 'topRight',
//   headline: '\u00BFCu\u00E1ntas\ndecisiones\nperdiste hoy?',
//   subtitle: 'Por no tener los datos a la mano.',
//   centered: false,
// }, theme);
