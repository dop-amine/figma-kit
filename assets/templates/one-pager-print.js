// Template: One-Pager (Print, Letter Size)
// Format: 1224x1584 (US Letter proportions at 2x)
// Structure: Problem -> Solution -> Proof -> CTA (B2B sales format)
// Requires: helpers.js functions (T, GM, R, G, accentBar)
//
// This template creates a light-mode, print-friendly one-pager.
// Use light.json theme for print, or default.json for digital dark version.

/**
 * Create a B2B sales one-pager.
 *
 * @param {PageNode} page    - Figma page
 * @param {object}   content - Content config:
 *   @param {string}  content.badge      - Badge text (e.g., "NEW  •  Q2 2026")
 *   @param {string}  content.headline   - Main headline
 *   @param {string}  content.subhead    - Supporting description
 *   @param {string}  content.tags       - Industry tags (monospace)
 *   @param {Array}   content.metrics    - [{value, label, delta}] (3 items)
 *   @param {Array}   content.capabilities - [{num, title, desc}] (4 items)
 *   @param {string}  content.clients    - Client names joined by " • "
 *   @param {object}  content.testimonial - {quote, author}
 *   @param {object}  content.cta        - {headline, button, compliance}
 * @param {object}   theme   - Theme config (use light.json for print)
 * @param {object}   [opts]  - Options:
 *   @param {number}  opts.x - X position (default 0)
 *   @param {number}  opts.y - Y position (default 0)
 * @returns {FrameNode}
 */
function createOnePager(page, content, theme, opts) {
  opts = opts || {};
  const W = theme.spacing.letter.width;
  const H = theme.spacing.letter.height;
  const MG = theme.spacing.letter.margin;
  const CW = W - MG * 2;
  const BG = theme.colors.BG;
  const WT = theme.colors.WT || theme.colors.CARD;
  const TH = theme.colors.TH || theme.colors.WT;
  const BD = theme.colors.BD;
  const MT = theme.colors.MT;
  const BL = theme.colors.BL;
  const TL = theme.colors.TL;
  const STK = theme.colors.STK;
  const SH = theme.effects.shadow.card;

  const f = figma.createFrame();
  f.name = 'One-Pager Print';
  f.resize(W, H);
  f.x = opts.x || 0;
  f.y = opts.y || 0;
  f.fills = [{ type: 'SOLID', color: BG }];
  f.clipsContent = true;
  page.appendChild(f);

  // Top accent gradient line
  accentBar(f, 0, 0, W, 4);

  // --- HEADER ---
  GM(f, 'ARKHAM', MG, 24, 11, TH);
  T(f, 'arkham.tech', W - MG - 80, 24, null, 11, 'Medium', BL);
  R(f, MG, 48, CW, 1, STK);

  // --- HERO ---
  if (content.badge) {
    const bdg = figma.createFrame();
    bdg.x = MG; bdg.y = 64; bdg.resize(152, 26); bdg.cornerRadius = 13;
    bdg.fills = [{ type: 'SOLID', color: BG }];
    bdg.strokes = [{ type: 'SOLID', color: BL }]; bdg.strokeWeight = 1;
    f.appendChild(bdg);
    T(bdg, content.badge, 16, 5, null, 10, 'Semi Bold', BL);
  }

  T(f, content.headline, MG, 104, 620,
    theme.type.h1.fontSize, theme.type.h1.style, TH, theme.type.h1.lineHeight);
  T(f, content.subhead, MG, 276, 550,
    theme.type.body.fontSize, theme.type.body.style, BD, theme.type.body.lineHeight);
  if (content.tags) GM(f, content.tags, MG, 380, 9, MT);

  // --- METRICS ---
  R(f, MG, 508, CW, 1, STK);
  GM(f, 'KEY OUTCOMES', MG, 524, 10, MT);

  if (content.metrics) {
    const mw = Math.floor((CW - 32) / content.metrics.length);
    content.metrics.forEach((m, i) => {
      const mx = MG + i * (mw + 16);
      const mc = figma.createFrame();
      mc.x = mx; mc.y = 548; mc.resize(mw, 90); mc.cornerRadius = 12;
      mc.fills = [{ type: 'SOLID', color: WT }];
      mc.strokes = [{ type: 'SOLID', color: STK }]; mc.strokeWeight = 1;
      mc.effects = [SH];
      f.appendChild(mc);
      T(mc, m.value, 16, 14, null, 28, 'Bold', BL);
      T(mc, m.label, 16, 52, mw - 32, 12, 'Regular', BD);
    });
  }

  // --- CAPABILITIES ---
  R(f, MG, 658, CW, 1, STK);
  GM(f, 'WHAT YOU GET', MG, 674, 10, MT);

  if (content.capabilities) {
    const cw2 = Math.floor((CW - 16) / 2);
    content.capabilities.forEach((cap, i) => {
      const row = Math.floor(i / 2), col = i % 2;
      const cc = figma.createFrame();
      cc.x = MG + col * (cw2 + 16);
      cc.y = 698 + row * (160 + 16);
      cc.resize(cw2, 160); cc.cornerRadius = 12;
      cc.fills = [{ type: 'SOLID', color: WT }];
      cc.strokes = [{ type: 'SOLID', color: STK }]; cc.strokeWeight = 1;
      cc.effects = [SH];
      f.appendChild(cc);
      R(cc, 0, 0, cw2, 4, i < 2 ? BL : TL, 0);
      GM(cc, cap.num, 16, 20, 10, MT);
      T(cc, cap.title, 16, 44, cw2 - 32, 15, 'Semi Bold', TH);
      T(cc, cap.desc, 16, 72, cw2 - 32, 12, 'Regular', BD, 18);
    });
  }

  // --- SOCIAL PROOF ---
  R(f, MG, 1050, CW, 1, STK);
  GM(f, 'TRUSTED BY INDUSTRY LEADERS', MG, 1066, 10, MT);
  if (content.clients) T(f, content.clients, MG, 1092, CW, 12, 'Semi Bold', MT);
  R(f, MG, 1124, CW, 1, STK);

  if (content.testimonial) {
    const qc = figma.createFrame();
    qc.x = MG; qc.y = 1140; qc.resize(CW, 120); qc.cornerRadius = 12;
    qc.fills = [{ type: 'SOLID', color: WT }];
    qc.effects = [SH];
    f.appendChild(qc);
    T(qc, '\u201C', 16, -10, null, 64, 'Bold', BL);
    T(qc, content.testimonial.quote, 16, 50, CW - 32, 14, 'Medium', TH, 22);
    T(qc, content.testimonial.author, 16, 90, null, 11, 'Regular', MT);
  }

  // --- CTA ---
  if (content.cta) {
    const ctaFrame = figma.createFrame();
    ctaFrame.x = MG; ctaFrame.y = 1280; ctaFrame.resize(CW, 146); ctaFrame.cornerRadius = 16;
    ctaFrame.fills = [theme.gradients.cta];
    ctaFrame.effects = [theme.effects.shadow.strong];
    f.appendChild(ctaFrame);
    T(ctaFrame, content.cta.headline, 32, 24, null, 22, 'Bold', { r: 1, g: 1, b: 1 });
    if (content.cta.button) {
      const btn = figma.createFrame();
      btn.x = 32; btn.y = 70; btn.resize(220, 44); btn.cornerRadius = 22;
      btn.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 } }];
      ctaFrame.appendChild(btn);
      T(btn, content.cta.button, 0, 12, 220, 14, 'Semi Bold', BL, null, 'CENTER');
    }
  }

  // --- FOOTER ---
  GM(f, 'ARKHAM \u2014 THE DATA & AI PLATFORM', MG, 1448, 9, MT);
  T(f, '\u00A9 2026 Arkham Technologies Inc.', W - MG - 190, 1448, null, 10, 'Regular', MT);

  return f;
}
