package assets

import (
	"embed"
)

//go:embed helpers.js
var HelpersJS string

//go:embed themes/default.json
var ThemeDefault []byte

//go:embed themes/light.json
var ThemeLight []byte

//go:embed themes/noir.json
var ThemeNoir []byte

//go:embed templates/slide.js
var TemplateSlide string

//go:embed templates/one-pager-print.js
var TemplateOnePager string

//go:embed templates/storyboard-panel.js
var TemplateStoryboard string

//go:embed themes/community
var CommunityThemesFS embed.FS

//go:embed examples
var ExamplesFS embed.FS

//go:embed cookbook.md
var CookbookMD string
