package codegen

import (
	"strings"

	"github.com/amine/figma-kit/assets"
)

// helperBoundaries maps function names to their source text.
var helperBoundaries map[string]string

func init() {
	helperBoundaries = parseHelperFunctions(assets.HelpersJS)
}

// AllHelpers returns the full helpers.js source.
func AllHelpers() string {
	return assets.HelpersJS
}

// HelperFunc returns the source of a single named helper function.
func HelperFunc(name string) string {
	return helperBoundaries[name]
}

// HelperFuncs returns the source of multiple helper functions concatenated.
func HelperFuncs(names ...string) string {
	var sb strings.Builder
	for i, name := range names {
		src := helperBoundaries[name]
		if src == "" {
			continue
		}
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(src)
	}
	return sb.String()
}

// AvailableHelpers returns the names of all parsed helper functions.
func AvailableHelpers() []string {
	names := make([]string, 0, len(helperBoundaries))
	for name := range helperBoundaries {
		names = append(names, name)
	}
	return names
}

func parseHelperFunctions(src string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(src, "\n")

	var currentName string
	var currentStart int

	flush := func(endExclusive int) {
		if currentName == "" {
			return
		}
		block := strings.Join(lines[currentStart:endExclusive], "\n")
		result[currentName] = strings.TrimRight(block, "\n") + "\n"
		currentName = ""
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "function ") || strings.HasPrefix(trimmed, "async function ") {
			flush(i)
			s := trimmed
			s = strings.TrimPrefix(s, "async ")
			s = strings.TrimPrefix(s, "function ")
			if paren := strings.IndexByte(s, '('); paren > 0 {
				currentName = s[:paren]
				currentStart = i
			}
			continue
		}

		if strings.HasPrefix(trimmed, "// ───") && currentName != "" {
			flush(i)
		}
	}

	flush(len(lines))
	return result
}
