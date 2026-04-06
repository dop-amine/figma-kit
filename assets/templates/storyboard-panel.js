// Template: Storyboard Styleframe Panel
// Format: 1920x1080 (16:9 horizontal)
// Structure: Full-bleed visual area + glassmorphism annotation strip at bottom
// Requires: helpers.js functions (T, GM, R, G, annot)
//
// Each panel represents one scene in a teaser/video storyboard.
// The annotation strip provides production-ready notes for the animation team.

/**
 * Create a storyboard styleframe panel.
 *
 * @param {PageNode} page    - Figma page
 * @param {object}   scene   - Scene configuration:
 *   @param {number}  scene.index      - Scene index (0-based, controls x position)
 *   @param {number}  scene.number     - Scene number (1-based, for labels)
 *   @param {string}  scene.name       - Frame name
 *   @param {string}  scene.title      - Scene title (in quotes)
 *   @param {string}  scene.timecode   - Timecode range (e.g., "0:00 \u2014 0:02.5")
 *   @param {string}  scene.duration   - Duration string (e.g., "2.5s")
 *   @param {string}  scene.camera     - Camera/motion direction
 *   @param {string}  scene.sound      - Sound/music description
 *   @param {Paint[]} scene.fills      - Background fills for the panel
 *   @param {string}  [scene.text]     - On-screen text overlay (null for no text)
 *   @param {string}  [scene.textStyle]- Text style: 'Bold'|'Medium'|'Light'
 *   @param {number}  [scene.textSize] - Text font size (default 56)
 *   @param {string}  [scene.textAlign]- 'CENTER'|'LEFT' (default 'CENTER')
 *   @param {string}  [scene.textNote] - Small annotation below text (e.g., "FADE-IN")
 * @param {object}   theme   - Theme config
 * @returns {FrameNode}
 */
function createStyleframe(page, scene, theme) {
  const FW = theme.spacing.frame16.width;
  const FH = theme.spacing.frame16.height;
  const GAP = 120;

  const sf = figma.createFrame();
  sf.name = 'SF ' + String(scene.number).padStart(2, '0') + ' \u2014 ' + scene.name;
  sf.resize(FW, FH);
  sf.x = scene.index * (FW + GAP);
  sf.y = 0;
  sf.fills = scene.fills || [{ type: 'SOLID', color: theme.colors.BG }];
  sf.clipsContent = true;
  page.appendChild(sf);

  // On-screen text overlay
  if (scene.text) {
    T(sf, scene.text, 0, FH / 2 - 60, FW,
      scene.textSize || 56,
      scene.textStyle || 'Medium',
      { r: 1, g: 1, b: 1 },
      null,
      scene.textAlign || 'CENTER');

    if (scene.textNote) {
      GM(sf, scene.textNote, FW / 2 - 80, FH / 2 + 20, 9,
        { r: 0.35, g: 0.38, b: 0.45 });
    }
  }

  // Production annotation strip
  annot(sf, scene.number, scene.title, scene.timecode,
    scene.duration, scene.camera, scene.sound);

  return sf;
}

/**
 * Create a panoramic storyboard overview strip.
 *
 * @param {PageNode} page     - Figma page
 * @param {Array}    scenes   - Array of scene summary objects:
 *   [{num, title, timecode, duration, description, bgColor, text, camera}]
 * @param {object}   meta     - Project metadata:
 *   @param {string}  meta.title    - Project title
 *   @param {string}  meta.duration - Total duration
 *   @param {string}  meta.format   - Format description
 *   @param {string}  meta.tone     - Tone description
 * @param {object}   theme    - Theme config
 * @param {object}   [opts]   - Options:
 *   @param {number}  opts.y - Y position (default 1280)
 * @returns {FrameNode}
 */
function createStoryboardPanoramic(page, scenes, meta, theme, opts) {
  opts = opts || {};
  const panelW = 960;
  const panelH = 540;
  const panelGap = 100;
  const PW = scenes.length * (panelW + panelGap) + 120;
  const PH = 1100;
  const BG = theme.colors.BG;

  const pan = figma.createFrame();
  pan.name = 'Storyboard Panor\u00e1mico \u2014 ' + meta.title;
  pan.resize(PW, PH);
  pan.x = 0;
  pan.y = opts.y || 1280;
  pan.fills = [
    { type: 'SOLID', color: BG },
    { type: 'GRADIENT_RADIAL',
      gradientTransform: [[2, 0, 0], [0, 2, 0]],
      gradientStops: [
        { position: 0, color: { r: 0.08, g: 0.12, b: 0.25, a: 0.10 } },
        { position: 1, color: { ...BG, a: 0 } },
      ]},
  ];
  pan.clipsContent = true;
  page.appendChild(pan);

  // Header
  const hdr = G(pan, 40, 30, PW - 80, 80, theme.effects.glass.default);
  GM(hdr, 'STORYBOARD', 20, 16, 14, theme.colors.BL);
  T(hdr, meta.title, 220, 14, null, 18, 'Semi Bold', theme.colors.WT);
  GM(hdr, meta.duration, 220, 44, 11, theme.colors.MT);
  GM(hdr, meta.format, 420, 44, 11, theme.colors.MT);
  if (meta.tone) GM(hdr, meta.tone, 680, 44, 11, { r: 0.35, g: 0.38, b: 0.45 });

  // Scene panels
  const SY = 160;
  scenes.forEach((sc, i) => {
    const sx = 60 + i * (panelW + panelGap);
    GM(pan, 'ESCENA ' + String(sc.num).padStart(2, '0'), sx, SY - 28, 11, theme.colors.BL);
    GM(pan, sc.timecode, sx + panelW - 60, SY - 28, 11, theme.colors.TL);

    const pnl = figma.createFrame();
    pnl.x = sx; pnl.y = SY; pnl.resize(panelW, panelH);
    pnl.cornerRadius = 16; pnl.clipsContent = true;
    pnl.fills = [{
      type: 'GRADIENT_RADIAL',
      gradientTransform: [[1.5, 0, 0], [0, 1.5, 0]],
      gradientStops: [
        { position: 0, color: { ...sc.bgColor, a: 1 } },
        { position: 1, color: { ...BG, a: 1 } },
      ],
    }];
    pnl.strokes = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.06 }];
    pnl.strokeWeight = 1;
    pan.appendChild(pnl);

    T(pnl, sc.text, 0, panelH / 2 - 30, panelW, 32, 'Bold',
      { r: 1, g: 1, b: 1 }, 40, 'CENTER');

    // Duration badge
    const dbg = figma.createFrame();
    dbg.x = panelW - 80; dbg.y = 12; dbg.resize(64, 28); dbg.cornerRadius = 14;
    dbg.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.08 }];
    pnl.appendChild(dbg);
    T(dbg, sc.duration, 0, 6, 64, 12, 'Medium', theme.colors.TL, null, 'CENTER');

    // Camera badge
    const cbg = figma.createFrame();
    cbg.x = 12; cbg.y = 12; cbg.resize(140, 28); cbg.cornerRadius = 14;
    cbg.fills = [{ type: 'SOLID', color: { r: 1, g: 1, b: 1 }, opacity: 0.08 }];
    pnl.appendChild(cbg);
    T(cbg, sc.camera, 0, 6, 140, 11, 'Medium',
      { r: 0.75, g: 0.78, b: 0.85 }, null, 'CENTER');

    // Title + description below
    T(pan, '"' + sc.title + '"', sx, SY + panelH + 16, panelW,
      18, 'Semi Bold', theme.colors.WT);
    T(pan, sc.description, sx, SY + panelH + 44, panelW,
      13, 'Regular', theme.colors.BD, 20);

    // Arrow to next
    if (i < scenes.length - 1) {
      const ax = sx + panelW + 10;
      const arrowLine = figma.createRectangle();
      arrowLine.x = ax; arrowLine.y = SY + panelH / 2 - 1;
      arrowLine.resize(panelGap - 20, 2);
      arrowLine.fills = [theme.gradients.accent];
      pan.appendChild(arrowLine);
      T(pan, '\u25B6', ax + panelGap - 36, SY + panelH / 2 - 10,
        null, 16, 'Regular', theme.colors.BL);
    }
  });

  // Timeline bar
  const tlY = SY + panelH + 120;
  R(pan, 60, tlY, PW - 120, 4, { r: 0.10, g: 0.12, b: 0.18 }, 2);
  const prog = figma.createRectangle();
  prog.x = 60; prog.y = tlY; prog.resize(PW - 120, 4); prog.cornerRadius = 2;
  prog.fills = [theme.gradients.accent];
  prog.opacity = 0.4;
  pan.appendChild(prog);

  return pan;
}
