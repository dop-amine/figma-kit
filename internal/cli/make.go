package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/dop-amine/figma-kit/internal/codegen"
	"github.com/dop-amine/figma-kit/internal/theme"
)

func loadContentFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading content file: %w", err)
	}
	var content map[string]any
	if err := yaml.Unmarshal(data, &content); err != nil {
		return nil, fmt.Errorf("parsing content file: %w", err)
	}
	return content, nil
}

func emitThemeJSON(b *codegen.Builder, t *theme.Theme) error {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("serializing theme: %w", err)
	}
	b.Comment("--- Theme object (for templates) ---")
	b.Linef("const theme = JSON.parse(%s);", strconv.Quote(string(data)))
	b.Blank()
	return nil
}

func jsStringLiteral(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		return strconv.Quote(s)
	}
	return string(b)
}

func marshalJSONForJS(v any) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return strconv.Quote(string(data)), nil
}

func startDeliverable(cmd *cobra.Command) (*theme.Theme, int, *codegen.Builder, error) {
	t, err := resolveTheme(cmd)
	if err != nil {
		return nil, 0, nil, err
	}
	page := resolvePage()
	b := codegen.New()
	codegen.PreambleWithPage(b, t, page)
	b.Comment("--- Helpers ---")
	b.Raw(codegen.AllHelpers())
	return t, page, b, nil
}

func finishDeliverable(b *codegen.Builder) {
	b.Blank()
	b.ReturnDone()
	output(b.String())
}

// --- Marketing & Social ---

func newMakeCarouselCmd() *cobra.Command {
	var contentPath string
	var slidesN int

	cmd := &cobra.Command{
		Use:   "carousel",
		Short: "LinkedIn-style carousel slides from YAML (1080×1350 via theme.slide)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			if err := emitThemeJSON(b, t); err != nil {
				return err
			}
			b.Comment("--- Template: slide ---")
			b.Raw(embeddedTemplates["slide"])

			content, err := loadContentFile(contentPath)
			if err != nil {
				return err
			}
			raw, ok := content["slides"]
			if !ok {
				return fmt.Errorf("content file must contain top-level 'slides' array")
			}
			arr, ok := raw.([]any)
			if !ok {
				return fmt.Errorf("'slides' must be an array")
			}
			if slidesN > 0 && len(arr) > slidesN {
				arr = arr[:slidesN]
			}
			slidesJSON, err := json.Marshal(arr)
			if err != nil {
				return err
			}
			quoted, err := marshalJSONForJS(json.RawMessage(slidesJSON))
			if err != nil {
				return err
			}
			b.Blank()
			b.Comment("--- Deliverable: carousel ---")
			b.Linef("const slideConfigs = JSON.parse(%s);", quoted)
			b.Line("for (let i = 0; i < slideConfigs.length; i++) {")
			b.Line("  const cfg = slideConfigs[i];")
			b.Line("  createSlide(pg, Object.assign({}, cfg, { index: i, total: slideConfigs.length }), theme);")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&contentPath, "content", "", "YAML with slides[] (slide configs for createSlide)")
	cmd.Flags().IntVar(&slidesN, "slides", 0, "Max slides (0 = use all from file)")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newMakeInstagramPostCmd() *cobra.Command {
	var typ, content string

	cmd := &cobra.Command{
		Use:   "instagram-post",
		Short: "Instagram feed post frame (1080×1080)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: instagram-post ---")
			b.Line("const W = 1080, H = 1080;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Instagram Post (%s)';", typ)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("f.clipsContent = true;")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, W, 8, BL);")
			b.Linef("GM(f, 'INSTAGRAM POST • %s', 48, 48, 10, MT);", strings.ToUpper(typ))
			qt, _ := marshalJSONForJS(content)
			b.Linef("T(f, JSON.parse(%s), 48, 120, W - 96, 36, 'Bold', WT, 44, 'LEFT');", qt)
			b.Line("R(f, 48, W - 120, W - 96, 4, STK);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "image", "quote or image")
	cmd.Flags().StringVar(&content, "content", "", "Headline or caption text")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newMakeInstagramStoryCmd() *cobra.Command {
	var contentPath string

	cmd := &cobra.Command{
		Use:   "instagram-story",
		Short: "Instagram story frame from YAML (1080×1920)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			content, err := loadContentFile(contentPath)
			if err != nil {
				return err
			}
			title := stringField(content, "title", "Story")
			sub := stringField(content, "subtitle", "")
			quotedTitle, _ := marshalJSONForJS(title)
			quotedSub, _ := marshalJSONForJS(sub)

			b.Comment("--- Deliverable: instagram-story ---")
			b.Line("const W = 1080, H = 1920;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Instagram Story';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("f.clipsContent = true;")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, H - 220, W, 220, BG);")
			b.Linef("T(f, JSON.parse(%s), 64, H - 180, W - 128, 40, 'Bold', WT, 48, 'LEFT');", quotedTitle)
			if sub != "" {
				b.Linef("T(f, JSON.parse(%s), 64, H - 100, W - 128, 22, 'Regular', BD, 32, 'LEFT');", quotedSub)
			}
			b.Line("R(f, W / 2 - 2, 80, 4, 640, BL);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&contentPath, "content", "", "YAML (title, subtitle, …)")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func stringField(m map[string]any, key, def string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func newMakeTwitterCardCmd() *cobra.Command {
	var headline, imageStyle string

	cmd := &cobra.Command{
		Use:   "twitter-card",
		Short: "Twitter / X large card (1200×675)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qh, _ := marshalJSONForJS(headline)
			b.Comment("--- Deliverable: twitter-card ---")
			b.Line("const W = 1200, H = 675;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Twitter Card (%s)';", imageStyle)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			if imageStyle == "minimal" {
				b.Line("R(f, W - 400, 0, 400, H, CARD);")
			} else {
				b.Line("R(f, 0, 0, W, H, CARD);")
				b.Line("R(f, 0, 0, W, H, { type: 'SOLID', color: { r: 0.1, g: 0.15, b: 0.25 }, opacity: 0.5 });")
			}
			b.Linef("T(f, JSON.parse(%s), 64, 200, 640, 52, 'Bold', WT, 62, 'LEFT');", qh)
			b.Line("GM(f, 'TWITTER CARD', 64, 64, 11, MT);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&headline, "headline", "", "Card headline")
	cmd.Flags().StringVar(&imageStyle, "image", "hero", "hero or minimal")
	_ = cmd.MarkFlagRequired("headline")
	return cmd
}

func newMakeFacebookCoverCmd() *cobra.Command {
	var scheme string

	cmd := &cobra.Command{
		Use:   "facebook-cover",
		Short: "Facebook cover photo (820×312)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			bg := "BG"
			if scheme == "light" {
				bg = "WT"
			}
			b.Comment("--- Deliverable: facebook-cover ---")
			b.Line("const W = 820, H = 312;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Facebook Cover (%s)';", scheme)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Linef("f.fills = [{ type: 'SOLID', color: %s }];", bg)
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, H - 6, W, 6, BL);")
			txtColor := "WT"
			if scheme == "light" {
				txtColor = "BG"
			}
			b.Linef("GM(f, 'COVER • %s', 32, 24, 10, MT);", strings.ToUpper(scheme))
			b.Linef("T(f, 'Your Page Name', 32, 80, W - 64, 28, 'Bold', %s, 34, 'LEFT');", txtColor)
			finishDeliverable(b)
			return nil
		},
	}
	// Named `scheme` to avoid clashing with the root persistent `--theme` (Figma kit theme).
	cmd.Flags().StringVar(&scheme, "scheme", "dark", "Cover color scheme: dark or light (same intent as dark|light cover theme)")
	return cmd
}

func newMakeYouTubeThumbCmd() *cobra.Command {
	var title, face string

	cmd := &cobra.Command{
		Use:   "youtube-thumb",
		Short: "YouTube thumbnail (1280×720)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qt, _ := marshalJSONForJS(title)
			var faceX string
			switch face {
			case "left":
				faceX = "120"
			case "right":
				faceX = "1280 - 120 - 200"
			default:
				faceX = "640 - 100"
			}
			b.Comment("--- Deliverable: youtube-thumb ---")
			b.Line("const W = 1280, H = 720;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'YouTube Thumbnail';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Linef("const faceX = %s;", faceX)
			b.Line("const facePl = figma.createEllipse();")
			b.Line("facePl.resize(200, 200); facePl.x = faceX; facePl.y = H / 2 - 100;")
			b.Line("facePl.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("f.appendChild(facePl);")
			b.Line("GM(f, 'FACE • " + face + "', faceX, H / 2 + 120, 10, MT);")
			b.Linef("T(f, JSON.parse(%s), 64, 64, 720, 56, 'Bold', WT, 66, 'LEFT');", qt)
			b.Line("R(f, 0, H - 12, W, 12, BL);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Video title on thumb")
	cmd.Flags().StringVar(&face, "face", "center", "left, right, or center (face placeholder)")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newMakeOGImageCmd() *cobra.Command {
	var title, description string

	cmd := &cobra.Command{
		Use:   "og-image",
		Short: "Open Graph share image (1200×630)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qt, _ := marshalJSONForJS(title)
			qd, _ := marshalJSONForJS(description)
			b.Comment("--- Deliverable: og-image ---")
			b.Line("const W = 1200, H = 630;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'OG Image';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, 8, H, BL);")
			b.Linef("T(f, JSON.parse(%s), 80, 160, W - 160, 48, 'Bold', WT, 58, 'LEFT');", qt)
			b.Linef("T(f, JSON.parse(%s), 80, 280, W - 160, 22, 'Regular', BD, 32, 'LEFT');", qd)
			b.Line("GM(f, 'OPEN GRAPH', 80, 80, 10, MT);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "OG title")
	cmd.Flags().StringVar(&description, "description", "", "OG description")
	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("description")
	return cmd
}

func newMakeBannerCmd() *cobra.Command {
	var sizes, contentPath string

	cmd := &cobra.Command{
		Use:   "banner",
		Short: "IAB-style display banners from YAML + size list",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			_, err = loadContentFile(contentPath)
			if err != nil {
				return err
			}
			parts := strings.Split(strings.ReplaceAll(sizes, " ", ""), ",")
			sizeMap := map[string][2]int{
				"leaderboard": {728, 90},
				"mrec":        {300, 250},
				"skyscraper":  {160, 600},
				"billboard":   {970, 250},
				"mobile":      {320, 50},
			}
			smJSON, err := json.Marshal(sizeMap)
			if err != nil {
				return err
			}
			keysJSON, err := json.Marshal(parts)
			if err != nil {
				return err
			}
			qsm, _ := marshalJSONForJS(json.RawMessage(smJSON))
			qk, _ := marshalJSONForJS(json.RawMessage(keysJSON))

			b.Comment("--- Deliverable: banner ---")
			b.Linef("const sizeMap = JSON.parse(%s);", qsm)
			b.Linef("const keys = JSON.parse(%s);", qk)
			b.Line("let x = 0;")
			b.Line("for (const k of keys) {")
			b.Line("  const dim = sizeMap[k];")
			b.Line("  if (!dim) continue;")
			b.Line("  const [bw, bh] = dim;")
			b.Line("  const bf = figma.createFrame();")
			b.Line("  bf.name = 'Banner ' + k;")
			b.Line("  bf.resize(bw, bh); bf.x = x; bf.y = 0;")
			b.Line("  bf.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("  pg.appendChild(bf);")
			b.Line("  GM(bf, k.toUpperCase(), 8, 8, 9, MT);")
			b.Line("  T(bf, 'Ad creative', 8, bh / 2 - 10, bw - 16, 14, 'Semi Bold', WT, null, 'LEFT');")
			b.Line("  x += bw + 40;")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&sizes, "sizes", "leaderboard,mrec", "Comma sizes: leaderboard,mrec,skyscraper,…")
	cmd.Flags().StringVar(&contentPath, "content", "", "YAML ad copy (loaded; extend generated JS as needed)")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newMakeEmailHeaderCmd() *cobra.Command {
	var width int

	cmd := &cobra.Command{
		Use:   "email-header",
		Short: "Email header strip (default width 600px)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			h := 120
			b.Comment("--- Deliverable: email-header ---")
			b.Linef("const W = %d, H = %d;", width, h)
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Email Header';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, W, 4, BL);")
			b.Line("GM(f, 'EMAIL HEADER', 24, 20, 10, MT);")
			b.Line("T(f, 'Brand • Newsletter', 24, 52, W - 48, 20, 'Bold', BG, null, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&width, "width", 600, "Header width in px")
	return cmd
}

func newMakeAdSetCmd() *cobra.Command {
	var contentPath, platforms string

	cmd := &cobra.Command{
		Use:   "ad-set",
		Short: "Multi-platform ad frames from campaign YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			if _, err := loadContentFile(contentPath); err != nil {
				return err
			}
			pl := strings.Split(strings.ReplaceAll(platforms, " ", ""), ",")
			var keys []string
			for _, p := range pl {
				p = strings.TrimSpace(strings.ToLower(p))
				if p != "" {
					keys = append(keys, p)
				}
			}
			pj, err := json.Marshal(keys)
			if err != nil {
				return err
			}
			qpl, err := marshalJSONForJS(json.RawMessage(pj))
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: ad-set ---")
			b.Line("const specs = { linkedin: [1200, 627], instagram: [1080, 1080], twitter: [1200, 675], facebook: [1200, 628] };")
			b.Linef("const platformList = JSON.parse(%s);", qpl)
			b.Line("let x = 0;")
			b.Line("for (const pl of platformList) {")
			b.Line("  const dim = specs[pl] || [1080, 1080];")
			b.Line("  const af = figma.createFrame();")
			b.Line("  af.name = 'Ad • ' + pl;")
			b.Line("  af.resize(dim[0], dim[1]); af.x = x; af.y = 0;")
			b.Line("  af.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("  pg.appendChild(af);")
			b.Line("  GM(af, pl.toUpperCase(), 32, 32, 11, MT);")
			b.Line("  T(af, 'Campaign placement', 32, 80, dim[0] - 64, 28, 'Bold', WT, 34, 'LEFT');")
			b.Line("  x += dim[0] + 80;")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&contentPath, "content", "", "campaign.yml")
	cmd.Flags().StringVar(&platforms, "platforms", "linkedin,instagram,twitter", "Comma platforms")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

// --- Sales & Business ---

func newMakeOnePagerCmd() *cobra.Command {
	var format, mode, contentPath string

	cmd := &cobra.Command{
		Use:   "one-pager",
		Short: "B2B one-pager from YAML (embedded one-pager-print template)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			if err := emitThemeJSON(b, t); err != nil {
				return err
			}
			b.Comment("--- Template: one-pager-print ---")
			b.Raw(embeddedTemplates["one-pager-print"])

			content, err := loadContentFile(contentPath)
			if err != nil {
				return err
			}
			inner := content
			if c, ok := content["content"].(map[string]any); ok {
				inner = c
			}
			payload, err := json.Marshal(inner)
			if err != nil {
				return err
			}
			quoted, err := marshalJSONForJS(json.RawMessage(payload))
			if err != nil {
				return err
			}

			b.Blank()
			b.Comment("--- Deliverable: one-pager ---")
			b.Linef("// format=%s mode=%s", format, mode)
			b.Linef("const onePagerContent = JSON.parse(%s);", quoted)
			b.Line("createOnePager(pg, onePagerContent, theme, { x: 0, y: 0 });")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "letter", "letter or A4 (metadata; template uses theme letter spacing)")
	cmd.Flags().StringVar(&mode, "mode", "print", "print or digital")
	cmd.Flags().StringVar(&contentPath, "content", "", "one-pager.yml")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newMakePitchDeckCmd() *cobra.Command {
	var slides int
	var template string

	cmd := &cobra.Command{
		Use:   "pitch-deck",
		Short: "Pitch deck slide frames (1920×1080)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: pitch-deck ---")
			b.Line("const W = 1920, H = 1080, GAP = 80;")
			b.Linef("const tmpl = %s;", jsStringLiteral(template))
			for i := 0; i < slides; i++ {
				b.Linef("{ const i = %d;", i)
				b.Line("  const sf = figma.createFrame();")
				b.Linef("  sf.name = 'Slide ' + (i + 1) + ' • ' + tmpl;")
				b.Line("  sf.resize(W, H); sf.x = i * (W + GAP); sf.y = 0;")
				b.Line("  sf.fills = [{ type: 'SOLID', color: BG }];")
				b.Line("  sf.clipsContent = true;")
				b.Line("  pg.appendChild(sf);")
				b.Line("  GM(sf, tmpl.toUpperCase(), 64, 48, 11, MT);")
				b.Line("  T(sf, 'Pitch headline ' + (i + 1), 64, 120, W - 128, 56, 'Bold', WT, 66, 'LEFT');")
				b.Line("  R(sf, 64, 220, W - 128, 4, BL);")
				b.Line("}")
			}
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&slides, "slides", 10, "Number of slides")
	cmd.Flags().StringVar(&template, "template", "saas", "saas, agency, or startup")
	return cmd
}

func newMakeCaseStudyCmd() *cobra.Command {
	var sections string

	cmd := &cobra.Command{
		Use:   "case-study",
		Short: "Case study layout with section blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			names := strings.Split(strings.ReplaceAll(sections, " ", ""), ",")
			arr, _ := json.Marshal(names)
			qa, _ := marshalJSONForJS(json.RawMessage(arr))

			b.Comment("--- Deliverable: case-study ---")
			b.Line("const W = 1440; let y = 0;")
			b.Linef("const sectionNames = JSON.parse(%s);", qa)
			b.Line("for (let s = 0; s < sectionNames.length; s++) {")
			b.Line("  const name = sectionNames[s];")
			b.Line("  const sec = figma.createFrame();")
			b.Line("  sec.name = 'Case Study • ' + name;")
			b.Line("  sec.resize(W, 420); sec.x = 0; sec.y = y;")
			b.Line("  sec.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("  pg.appendChild(sec);")
			b.Line("  GM(sec, name.replace(/,/g, '').toUpperCase(), 48, 40, 10, MT);")
			b.Line("  T(sec, 'Content for ' + name + ' goes here.', 48, 100, W - 96, 22, 'Regular', BD, 30, 'LEFT');")
			b.Line("  y += 420 + 32;")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&sections, "sections", "overview,challenge,solution,results", "Comma section keys")
	return cmd
}

func newMakeProposalCmd() *cobra.Command {
	var client, scopePath string

	cmd := &cobra.Command{
		Use:   "proposal",
		Short: "Client proposal cover + scope from YAML",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			scope, err := loadContentFile(scopePath)
			if err != nil {
				return err
			}
			summary := stringField(scope, "summary", "Scope of work")
			qc, _ := marshalJSONForJS(client)
			qs, _ := marshalJSONForJS(summary)

			b.Comment("--- Deliverable: proposal ---")
			b.Line("const W = 1224, H = 1584;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Proposal';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("GM(f, 'PROPOSAL', 72, 72, 11, MT);")
			b.Linef("T(f, JSON.parse(%s), 72, 140, W - 144, 48, 'Bold', WT, 56, 'LEFT');", qc)
			b.Linef("T(f, JSON.parse(%s), 72, 240, W - 144, 18, 'Regular', BD, 28, 'LEFT');", qs)
			b.Line("R(f, 72, 320, W - 144, 1, STK);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&client, "client", "", "Client name")
	cmd.Flags().StringVar(&scopePath, "scope", "", "scope.yml")
	_ = cmd.MarkFlagRequired("client")
	_ = cmd.MarkFlagRequired("scope")
	return cmd
}

func newMakeInvoiceCmd() *cobra.Command {
	var template string

	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "Invoice layout frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: invoice ---")
			b.Line("const W = 800, H = 1100;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Invoice (%s)';", template)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			headCol := "WT"
			bodyCol := "BD"
			if template == "minimal" {
				b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
				headCol = "BG"
				bodyCol = "MT"
			} else {
				b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			}
			b.Line("pg.appendChild(f);")
			b.Line("GM(f, 'INVOICE', 48, 48, 11, MT);")
			b.Linef("T(f, 'Invoice #0001', 48, 100, 400, 28, 'Bold', %s, null, 'LEFT');", headCol)
			b.Line("R(f, 48, 160, W - 96, 1, STK);")
			b.Linef("T(f, 'Line items…', 48, 200, W - 96, 14, 'Regular', %s, null, 'LEFT');", bodyCol)
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&template, "template", "modern", "modern or minimal")
	return cmd
}

func newMakeBusinessCardCmd() *cobra.Command {
	var name, jobTitle string

	cmd := &cobra.Command{
		Use:   "business-card",
		Short: "Business card (1050×600)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qn, _ := marshalJSONForJS(name)
			qt, _ := marshalJSONForJS(jobTitle)
			b.Comment("--- Deliverable: business-card ---")
			b.Line("const W = 1050, H = 600;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Business Card';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, 12, H, BL);")
			b.Linef("T(f, JSON.parse(%s), 72, 200, W - 144, 36, 'Bold', WT, 42, 'LEFT');", qn)
			b.Linef("T(f, JSON.parse(%s), 72, 280, W - 144, 18, 'Medium', BD, null, 'LEFT');", qt)
			b.Line("GM(f, 'BUSINESS CARD', 72, 72, 10, MT);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Name on card")
	cmd.Flags().StringVar(&jobTitle, "title", "", "Title / role")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newMakeLetterheadCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "letterhead",
		Short: "Letterhead document frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			w, h := 1224, 1584
			if format == "A4" {
				w, h = 2480, 3508
			}
			b.Comment("--- Deliverable: letterhead ---")
			b.Linef("const W = %d, H = %d;", w, h)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Letterhead (%s)';", format)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 72, 72, W - 144, 4, BL);")
			b.Line("GM(f, 'LETTERHEAD', 72, 100, 10, MT);")
			b.Line("T(f, 'Company Legal Name', 72, 140, 400, 14, 'Semi Bold', BG, null, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "letter", "A4 or letter")
	return cmd
}

func newMakeContractCmd() *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "contract",
		Short: "Contract title page frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qt, _ := marshalJSONForJS(title)
			b.Comment("--- Deliverable: contract ---")
			b.Line("const W = 1224, H = 1584;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Contract';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Linef("T(f, JSON.parse(%s), 120, 400, W - 240, 36, 'Bold', BG, 44, 'CENTER');", qt)
			b.Line("GM(f, 'AGREEMENT', 120, 320, 11, MT);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Contract title")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

// --- Motion & Storyboard ---

func newMakeStoryboardCmd() *cobra.Command {
	var scenesN int
	var contentPath string

	cmd := &cobra.Command{
		Use:   "storyboard",
		Short: "Storyboard styleframes from YAML (embedded storyboard-panel template)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			if err := emitThemeJSON(b, t); err != nil {
				return err
			}
			b.Comment("--- Template: storyboard-panel ---")
			b.Raw(embeddedTemplates["storyboard-panel"])

			content, err := loadContentFile(contentPath)
			if err != nil {
				return err
			}
			raw, ok := content["scenes"]
			if !ok {
				return fmt.Errorf("storyboard.yml must contain 'scenes' array")
			}
			arr, ok := raw.([]any)
			if !ok {
				return fmt.Errorf("'scenes' must be an array")
			}
			if scenesN > 0 && len(arr) > scenesN {
				arr = arr[:scenesN]
			}
			scenesJSON, err := json.Marshal(arr)
			if err != nil {
				return err
			}
			quoted, err := marshalJSONForJS(json.RawMessage(scenesJSON))
			if err != nil {
				return err
			}

			b.Blank()
			b.Comment("--- Deliverable: storyboard ---")
			b.Linef("const scenes = JSON.parse(%s);", quoted)
			b.Line("scenes.forEach((scene, i) => {")
			b.Line("  createStyleframe(pg, Object.assign({}, scene, { index: i }), theme);")
			b.Line("});")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&scenesN, "scenes", 0, "Max scenes (0 = all from file)")
	cmd.Flags().StringVar(&contentPath, "content", "", "storyboard.yml with scenes[]")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newMakeStyleframeCmd() *cobra.Command {
	var mood, scene string

	cmd := &cobra.Command{
		Use:   "styleframe",
		Short: "Single motion styleframe (1920×1080)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			var cr, cg, cb string
			switch mood {
			case "warm":
				cr, cg, cb = "0.25", "0.12", "0.08"
			case "cool":
				cr, cg, cb = "0.06", "0.12", "0.22"
			default:
				cr, cg, cb = "0.1", "0.1", "0.12"
			}
			qs, _ := marshalJSONForJS(scene)
			b.Comment("--- Deliverable: styleframe ---")
			b.Line("const W = 1920, H = 1080;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Styleframe';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Linef("f.fills = [{ type: 'SOLID', color: { r: %s, g: %s, b: %s } }];", cr, cg, cb)
			b.Line("f.clipsContent = true;")
			b.Line("pg.appendChild(f);")
			b.Linef("T(f, JSON.parse(%s), 80, H / 2 - 40, W - 160, 64, 'Bold', WT, 72, 'LEFT');", qs)
			b.Linef("GM(f, 'MOOD • %s', 80, 80, 11, MT);", strings.ToUpper(mood))
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&mood, "mood", "neutral", "warm, cool, or neutral")
	cmd.Flags().StringVar(&scene, "scene", "", "Scene description text")
	_ = cmd.MarkFlagRequired("scene")
	return cmd
}

func newMakeAnimaticCmd() *cobra.Command {
	var fps int
	var duration string

	cmd := &cobra.Command{
		Use:   "animatic",
		Short: "Animatic timeline overview frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qd, _ := marshalJSONForJS(duration)
			b.Comment("--- Deliverable: animatic ---")
			b.Line("const W = 1920, H = 400;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Animatic';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Linef("GM(f, 'ANIMATIC • %d FPS', 48, 40, 11, MT);", fps)
			b.Linef("T(f, JSON.parse(%s), 48, 80, 600, 22, 'Regular', BD, 30, 'LEFT');", qd)
			b.Line("for (let i = 0; i < 12; i++) {")
			b.Line("  const c = figma.createFrame();")
			b.Line("  c.resize(120, 200); c.x = 48 + i * 140; c.y = 140;")
			b.Line("  c.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("  f.appendChild(c);")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&fps, "fps", 24, "Frames per second")
	cmd.Flags().StringVar(&duration, "duration", "30s", "Total duration label")
	return cmd
}

func newMakeTransitionSpecCmd() *cobra.Command {
	var typ, easing string

	cmd := &cobra.Command{
		Use:   "transition-spec",
		Short: "UI transition specification card",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: transition-spec ---")
			b.Line("const W = 720, H = 480;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Transition Spec';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("pg.appendChild(f);")
			b.Linef("GM(f, 'TRANSITION', 40, 36, 11, MT);")
			b.Linef("T(f, 'Type: %s', 40, 80, W - 80, 22, 'Semi Bold', WT, null, 'LEFT');", typ)
			b.Linef("T(f, 'Easing: %s', 40, 120, W - 80, 18, 'Regular', BD, null, 'LEFT');", easing)
			b.Line("R(f, 40, 200, W - 80, 120, BG);")
			b.Line("T(f, 'Curve / keyframes', 56, 240, W - 112, 14, 'Regular', MT, null, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "fade", "page-enter, fade, or slide")
	cmd.Flags().StringVar(&easing, "easing", "ease-out", "CSS-style easing label")
	return cmd
}

// --- UI/UX Design ---

func newMakeWireframeCmd() *cobra.Command {
	var typ, breakpoint string

	cmd := &cobra.Command{
		Use:   "wireframe",
		Short: "Low-fidelity wireframe shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			w := 1440
			switch breakpoint {
			case "tablet":
				w = 834
			case "mobile":
				w = 390
			}
			b.Comment("--- Deliverable: wireframe ---")
			b.Linef("const W = %d, H = 900;", w)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Wireframe • %s • %s';", typ, breakpoint)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, W, 64, STK);")
			b.Line("R(f, 24, 88, W - 48, 200, STK);")
			b.Line("R(f, 24, 320, (W - 56) / 2, 400, STK);")
			b.Line("R(f, 32 + (W - 56) / 2, 320, (W - 56) / 2, 400, STK);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "landing", "landing, dashboard, or form")
	cmd.Flags().StringVar(&breakpoint, "breakpoint", "desktop", "desktop, tablet, or mobile")
	return cmd
}

func newMakeScreenCmd() *cobra.Command {
	var typ, sections string

	cmd := &cobra.Command{
		Use:   "screen",
		Short: "Marketing screen with section strips",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			names := strings.Split(strings.ReplaceAll(sections, " ", ""), ",")
			arr, _ := json.Marshal(names)
			qa, _ := marshalJSONForJS(json.RawMessage(arr))

			b.Comment("--- Deliverable: screen ---")
			b.Line("const W = 1440; let y = 0;")
			b.Linef("const secList = JSON.parse(%s);", qa)
			b.Linef("const screenType = %s;", jsStringLiteral(typ))
			b.Line("for (let i = 0; i < secList.length; i++) {")
			b.Line("  const blk = figma.createFrame();")
			b.Line("  blk.name = screenType + ' • ' + secList[i];")
			b.Line("  blk.resize(W, 360); blk.x = 0; blk.y = y;")
			b.Line("  blk.fills = [{ type: 'SOLID', color: i % 2 === 0 ? BG : CARD }];")
			b.Line("  pg.appendChild(blk);")
			b.Line("  GM(blk, secList[i].toUpperCase(), 48, 40, 10, MT);")
			b.Line("  y += 360;")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "landing", "landing, pricing, or features")
	cmd.Flags().StringVar(&sections, "sections", "hero,features,pricing,cta", "Comma section names")
	return cmd
}

func newMakeDashboardCmd() *cobra.Command {
	var widgets string
	var cols int

	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Dashboard grid of widget placeholders",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			wlist := strings.Split(strings.ReplaceAll(widgets, " ", ""), ",")
			arr, _ := json.Marshal(wlist)
			qa, _ := marshalJSONForJS(json.RawMessage(arr))

			b.Comment("--- Deliverable: dashboard ---")
			b.Line("const W = 1440, H = 900;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Dashboard';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, 0, 0, W, 72, CARD);")
			b.Linef("const widgets = JSON.parse(%s);", qa)
			b.Linef("const cols = %d;", cols)
			b.Line("const gw = (W - 80 - (cols - 1) * 24) / cols;")
			b.Line("widgets.forEach((w, i) => {")
			b.Line("  const col = i % cols, row = Math.floor(i / cols);")
			b.Line("  const wf = figma.createFrame();")
			b.Line("  wf.name = 'Widget: ' + w;")
			b.Line("  wf.resize(gw, 200);")
			b.Line("  wf.x = 40 + col * (gw + 24);")
			b.Line("  wf.y = 100 + row * 224;")
			b.Line("  wf.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("  f.appendChild(wf);")
			b.Line("  GM(wf, w, 16, 16, 10, MT);")
			b.Line("});")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&widgets, "widgets", "stat,chart,table,list", "Comma widget types")
	cmd.Flags().IntVar(&cols, "cols", 2, "Grid columns")
	return cmd
}

func newMakeFormCmd() *cobra.Command {
	var fieldsPath string

	cmd := &cobra.Command{
		Use:   "form",
		Short: "Form layout from JSON/YAML schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			content, err := loadContentFile(fieldsPath)
			if err != nil {
				return err
			}
			raw, ok := content["fields"]
			if !ok {
				return fmt.Errorf("schema must contain 'fields' array")
			}
			fields, ok := raw.([]any)
			if !ok {
				return fmt.Errorf("'fields' must be an array")
			}
			fj, err := json.Marshal(fields)
			if err != nil {
				return err
			}
			qf, err := marshalJSONForJS(json.RawMessage(fj))
			if err != nil {
				return err
			}

			b.Comment("--- Deliverable: form ---")
			b.Line("const W = 480;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Form';")
			b.Line("f.resize(W, 80); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("f.layoutMode = 'VERTICAL';")
			b.Line("f.itemSpacing = 20;")
			b.Line("f.paddingLeft = 24; f.paddingRight = 24; f.paddingTop = 24; f.paddingBottom = 24;")
			b.Linef("const fields = JSON.parse(%s);", qf)
			b.Line("let fh = 24;")
			b.Line("fields.forEach((fld) => {")
			b.Line("  const row = figma.createFrame();")
			b.Line("  row.resize(W - 48, 56);")
			b.Line("  row.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("  f.appendChild(row);")
			b.Line("  const label = (fld.label || fld.name || 'Field');")
			b.Line("  T(row, label, 12, 8, row.width - 24, 12, 'Medium', WT, null, 'LEFT');")
			b.Line("  fh += 76;")
			b.Line("});")
			b.Line("f.resize(W, fh);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&fieldsPath, "fields", "", "form-schema.json (or YAML) with fields[]")
	_ = cmd.MarkFlagRequired("fields")
	return cmd
}

func newMakeModalCmd() *cobra.Command {
	var size, typ string

	cmd := &cobra.Command{
		Use:   "modal",
		Short: "Modal overlay frame",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			mw := 520
			switch size {
			case "sm":
				mw = 400
			case "lg":
				mw = 640
			}
			b.Comment("--- Deliverable: modal ---")
			b.Linef("const MW = %d, MH = 360;", mw)
			b.Line("const overlay = figma.createFrame();")
			b.Linef("overlay.name = 'Modal • %s • %s';", typ, size)
			b.Line("overlay.resize(1440, 900); overlay.x = 0; overlay.y = 0;")
			b.Line("overlay.fills = [{ type: 'SOLID', color: { r: 0, g: 0, b: 0 }, opacity: 0.45 }];")
			b.Line("pg.appendChild(overlay);")
			b.Line("const m = figma.createFrame();")
			b.Line("m.resize(MW, MH); m.x = (1440 - MW) / 2; m.y = (900 - MH) / 2;")
			b.Line("m.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("m.cornerRadius = 16;")
			b.Line("overlay.appendChild(m);")
			b.Linef("GM(m, '%s', 24, 20, 10, MT);", strings.ToUpper(typ))
			b.Line("T(m, 'Modal title', 24, 52, MW - 48, 18, 'Semi Bold', WT, null, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&size, "size", "md", "sm, md, or lg")
	cmd.Flags().StringVar(&typ, "type", "confirmation", "confirmation, form, or alert")
	return cmd
}

func newMakeEmptyStateCmd() *cobra.Command {
	var message string

	cmd := &cobra.Command{
		Use:   "empty-state",
		Short: "Empty state illustration block",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			qm, _ := marshalJSONForJS(message)
			b.Comment("--- Deliverable: empty-state ---")
			b.Line("const W = 400, H = 320;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Empty State';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("pg.appendChild(f);")
			b.Line("R(f, W/2 - 40, 48, 80, 80, STK);")
			b.Linef("T(f, JSON.parse(%s), 40, 160, W - 80, 16, 'Regular', BD, 24, 'CENTER');", qm)
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&message, "message", "Nothing here yet.", "Empty state copy")
	return cmd
}

func newMakeErrorPageCmd() *cobra.Command {
	var typ string

	cmd := &cobra.Command{
		Use:   "error-page",
		Short: "Full-page error state",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			code := typ
			var msg string
			switch typ {
			case "404":
				msg = "Page not found."
			case "offline":
				msg = "You appear to be offline."
			default:
				msg = "Something went wrong."
			}
			qm, _ := marshalJSONForJS(msg)
			b.Comment("--- Deliverable: error-page ---")
			b.Line("const W = 1440, H = 900;")
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Error %s';", typ)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Linef("T(f, %s, W/2 - 80, H/2 - 80, 160, 72, 'Bold', ERR, null, 'CENTER');", jsStringLiteral(code))
			b.Linef("T(f, JSON.parse(%s), W/2 - 200, H/2 + 20, 400, 20, 'Regular', BD, null, 'CENTER');", qm)
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "404", "404, 500, or offline")
	return cmd
}

func newMakeOnboardingCmd() *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "onboarding",
		Short: "Onboarding step screens",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: onboarding ---")
			b.Line("const W = 390, H = 844, GAP = 48;")
			for i := 0; i < steps; i++ {
				b.Linef("{ const si = %d;", i)
				b.Line("  const s = figma.createFrame();")
				b.Line("  s.name = 'Onboarding Step ' + (si + 1);")
				b.Line("  s.resize(W, H); s.x = si * (W + GAP); s.y = 0;")
				b.Line("  s.fills = [{ type: 'SOLID', color: BG }];")
				b.Line("  pg.appendChild(s);")
				b.Line("  GM(s, 'STEP ' + (si + 1), 32, 60, 11, MT);")
				b.Line("  R(s, 32, 120, W - 64, 200, CARD);")
				b.Line("}")
			}
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&steps, "steps", 3, "Number of steps")
	return cmd
}

func newMakeSettingsCmd() *cobra.Command {
	var sections string

	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings app layout with section list",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			names := strings.Split(strings.ReplaceAll(sections, " ", ""), ",")
			arr, _ := json.Marshal(names)
			qa, _ := marshalJSONForJS(json.RawMessage(arr))

			b.Comment("--- Deliverable: settings ---")
			b.Line("const W = 1440, H = 900;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Settings';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("f.layoutMode = 'HORIZONTAL';")
			b.Line("pg.appendChild(f);")
			b.Line("const nav = figma.createFrame();")
			b.Line("nav.resize(280, H); nav.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("f.appendChild(nav);")
			b.Linef("const navItems = JSON.parse(%s);", qa)
			b.Line("navItems.forEach((name, i) => {")
			b.Line("  T(nav, name, 24, 32 + i * 44, 232, 15, 'Medium', WT, null, 'LEFT');")
			b.Line("});")
			b.Line("const panel = figma.createFrame();")
			b.Line("panel.resize(W - 280, H); panel.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("f.appendChild(panel);")
			b.Line("GM(panel, 'DETAIL PANEL', 40, 40, 10, MT);")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&sections, "sections", "profile,notifications,billing,security", "Comma section ids")
	return cmd
}

// --- Print ---

func newMakePosterCmd() *cobra.Command {
	var size, bleed string

	cmd := &cobra.Command{
		Use:   "poster",
		Short: "Print poster frame with bleed note",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			w, h := 4961, 7016 // A2 @300dpi approx
			switch size {
			case "A3":
				w, h = 3508, 4961
			case "24x36":
				w, h = 7200, 10800
			case "custom":
				w, h = 3600, 4800
			}
			b.Comment("--- Deliverable: poster ---")
			b.Linef("const W = %d, H = %d;", w, h)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Poster %s';", size)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Linef("GM(f, 'BLEED %s', 48, 48, 11, MT);", bleed)
			b.Line("T(f, 'Poster headline', 48, 120, W - 96, 96, 'Bold', BG, 100, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&size, "size", "A2", "A2, A3, 24x36, or custom")
	cmd.Flags().StringVar(&bleed, "bleed", "3mm", "Bleed label")
	return cmd
}

func newMakeBrochureCmd() *cobra.Command {
	var fold, format string

	cmd := &cobra.Command{
		Use:   "brochure",
		Short: "Brochure panel layout",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			w, h := 1224, 1584
			if format == "A4" {
				w = 2480
				h = 3508
			}
			panels := 3
			if fold == "bifold" {
				panels = 2
			}
			b.Comment("--- Deliverable: brochure ---")
			b.Linef("const W = %d, H = %d, PANELS = %d;", w, h, panels)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Brochure %s %s';", fold, format)
			b.Line("f.resize(W * PANELS + (PANELS - 1) * 32, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("pg.appendChild(f);")
			b.Line("for (let p = 0; p < PANELS; p++) {")
			b.Line("  const pan = figma.createFrame();")
			b.Line("  pan.resize(W, H); pan.x = p * (W + 32); pan.y = 0;")
			b.Line("  pan.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("  f.appendChild(pan);")
			b.Line("  GM(pan, 'PANEL ' + (p + 1), 48, 48, 10, MT);")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&fold, "fold", "trifold", "trifold or bifold")
	cmd.Flags().StringVar(&format, "format", "letter", "letter or A4")
	return cmd
}

func newMakePackagingCmd() *cobra.Command {
	var typ string
	var w, h, d float64

	cmd := &cobra.Command{
		Use:   "packaging",
		Short: "Packaging die-line style frame (flat)",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: packaging ---")
			b.Linef("const boxW = %.0f, boxH = %.0f, boxD = %.0f;", w, h, d)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Packaging %s';", typ)
			b.Line("f.resize(boxW * 4 + 200, boxH * 2 + boxD + 200); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Line("GM(f, 'DIE LINE • UNFOLDED', 40, 32, 10, MT);")
			b.Line("R(f, 60, 80, boxW, boxH, STK);")
			b.Line("R(f, 60 + boxW, 80, boxD, boxH, STK);")
			b.Line("T(f, 'Front', 60, 80 + boxH + 8, boxW, 12, 'Medium', MT, null, 'LEFT');")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "box", "box (flat layout)")
	cmd.Flags().Float64Var(&w, "w", 200, "Front width")
	cmd.Flags().Float64Var(&h, "h", 280, "Front height")
	cmd.Flags().Float64Var(&d, "d", 100, "Depth")
	return cmd
}

func newMakeSignageCmd() *cobra.Command {
	var wStr, hStr string

	cmd := &cobra.Command{
		Use:   "signage",
		Short: "Large-format signage board",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			b.Comment("--- Deliverable: signage ---")
			b.Linef("// Requested size: %s × %s (interpreted as px for canvas)", wStr, hStr)
			b.Line("const W = 1920, H = 960;")
			b.Line("const f = figma.createFrame();")
			b.Line("f.name = 'Signage';")
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: CARD }];")
			b.Line("pg.appendChild(f);")
			b.Linef("GM(f, 'SIGNAGE', 64, 48, 14, MT);")
			b.Linef("T(f, %s + ' × ' + %s, 64, 120, W - 128, 28, 'Bold', WT, null, 'LEFT');",
				jsStringLiteral(wStr), jsStringLiteral(hStr))
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().StringVar(&wStr, "w", "48in", "Width label")
	cmd.Flags().StringVar(&hStr, "h", "24in", "Height label")
	return cmd
}

func newMakeMenuCmd() *cobra.Command {
	var sections int
	var format string

	cmd := &cobra.Command{
		Use:   "menu",
		Short: "Restaurant menu columns",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, b, err := startDeliverable(cmd)
			if err != nil {
				return err
			}
			w, h := 1224, 1584
			if format == "A4" {
				w, h = 2480, 3508
			}
			if sections < 1 {
				sections = 1
			}
			colW := (w - 120) / sections
			b.Comment("--- Deliverable: menu ---")
			b.Linef("const W = %d, H = %d, SECTIONS = %d, COLW = %d;", w, h, sections, colW)
			b.Line("const f = figma.createFrame();")
			b.Linef("f.name = 'Menu %s';", format)
			b.Line("f.resize(W, H); f.x = 0; f.y = 0;")
			b.Line("f.fills = [{ type: 'SOLID', color: WT }];")
			b.Line("pg.appendChild(f);")
			b.Line("GM(f, 'MENU', 48, 40, 12, MT);")
			b.Line("for (let s = 0; s < SECTIONS; s++) {")
			b.Line("  const col = figma.createFrame();")
			b.Line("  col.resize(COLW, H - 160); col.x = 48 + s * (COLW + 24); col.y = 100;")
			b.Line("  col.fills = [{ type: 'SOLID', color: BG }];")
			b.Line("  f.appendChild(col);")
			b.Line("  T(col, 'Section ' + (s + 1), 16, 16, COLW - 32, 16, 'Semi Bold', WT, null, 'LEFT');")
			b.Line("}")
			finishDeliverable(b)
			return nil
		},
	}
	cmd.Flags().IntVar(&sections, "sections", 2, "Number of menu columns/sections")
	cmd.Flags().StringVar(&format, "format", "letter", "letter or A4")
	return cmd
}

// newMakeCmd is the parent for all deliverable generators.
func newMakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make",
		Short: "Generate complete design deliverables",
	}
	// Marketing & Social (10)
	cmd.AddCommand(newMakeCarouselCmd())
	cmd.AddCommand(newMakeInstagramPostCmd())
	cmd.AddCommand(newMakeInstagramStoryCmd())
	cmd.AddCommand(newMakeTwitterCardCmd())
	cmd.AddCommand(newMakeFacebookCoverCmd())
	cmd.AddCommand(newMakeYouTubeThumbCmd())
	cmd.AddCommand(newMakeOGImageCmd())
	cmd.AddCommand(newMakeBannerCmd())
	cmd.AddCommand(newMakeEmailHeaderCmd())
	cmd.AddCommand(newMakeAdSetCmd())
	// Sales & Business (8)
	cmd.AddCommand(newMakeOnePagerCmd())
	cmd.AddCommand(newMakePitchDeckCmd())
	cmd.AddCommand(newMakeCaseStudyCmd())
	cmd.AddCommand(newMakeProposalCmd())
	cmd.AddCommand(newMakeInvoiceCmd())
	cmd.AddCommand(newMakeBusinessCardCmd())
	cmd.AddCommand(newMakeLetterheadCmd())
	cmd.AddCommand(newMakeContractCmd())
	// Motion & Storyboard (4)
	cmd.AddCommand(newMakeStoryboardCmd())
	cmd.AddCommand(newMakeStyleframeCmd())
	cmd.AddCommand(newMakeAnimaticCmd())
	cmd.AddCommand(newMakeTransitionSpecCmd())
	// UI/UX Design (9)
	cmd.AddCommand(newMakeWireframeCmd())
	cmd.AddCommand(newMakeScreenCmd())
	cmd.AddCommand(newMakeDashboardCmd())
	cmd.AddCommand(newMakeFormCmd())
	cmd.AddCommand(newMakeModalCmd())
	cmd.AddCommand(newMakeEmptyStateCmd())
	cmd.AddCommand(newMakeErrorPageCmd())
	cmd.AddCommand(newMakeOnboardingCmd())
	cmd.AddCommand(newMakeSettingsCmd())
	// Print (5)
	cmd.AddCommand(newMakePosterCmd())
	cmd.AddCommand(newMakeBrochureCmd())
	cmd.AddCommand(newMakePackagingCmd())
	cmd.AddCommand(newMakeSignageCmd())
	cmd.AddCommand(newMakeMenuCmd())
	return cmd
}
