package theme

import (
	"path/filepath"
	"strings"

	"github.com/dop-amine/figma-kit/assets"
)

var embeddedThemes = map[string][]byte{
	"default": assets.ThemeDefault,
	"light":   assets.ThemeLight,
	"noir":    assets.ThemeNoir,
}

var communityThemes = map[string][]byte{}

func init() {
	entries, err := assets.CommunityThemesFS.ReadDir("themes/community")
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := assets.CommunityThemesFS.ReadFile(filepath.Join("themes", "community", e.Name()))
		if err != nil {
			continue
		}
		key := strings.TrimSuffix(e.Name(), ".json")
		communityThemes[key] = data
	}
}
