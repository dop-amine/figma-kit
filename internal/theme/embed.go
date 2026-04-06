package theme

import (
	"github.com/dop-amine/figma-kit/assets"
)

// embeddedThemes maps theme names to their raw JSON data.
var embeddedThemes = map[string][]byte{
	"default": assets.ThemeDefault,
	"light":   assets.ThemeLight,
	"arkham":  assets.ThemeArkham,
}
