package render

import "fmt"

var namedColors = map[string]string{
	"black":   "#000000",
	"red":     "#FF0000",
	"green":   "#00FF00",
	"yellow":  "#FFFF00",
	"blue":    "#0000FF",
	"magenta": "#FF00FF",
	"cyan":    "#00FFFF",
	"white":   "#FFFFFF",
	"gray":    "#888888",
	"grey":    "#888888",
	"orange":  "#FF8800",
	"pink":    "#FF6B9D",
	"purple":  "#8800FF",
}

// ResolveColor converts a named color or hex string to a hex string.
func ResolveColor(name string) (string, bool) {
	if hex, ok := namedColors[name]; ok {
		return hex, true
	}
	// Check if it's already a hex color
	if len(name) > 0 && (name[0] == '#' || len(name) == 6) {
		_, err := parseHex(name)
		if err == nil {
			if name[0] != '#' {
				name = "#" + name
			}
			return name, true
		}
	}
	return "", false
}

// RenderStyled applies foreground color, bold, and dim to text.
func RenderStyled(text string, fg string, bold bool, dim bool) string {
	if fg == "" && !bold && !dim {
		return text
	}

	var prefix string

	if fg != "" {
		c, err := parseHex(fg)
		if err == nil {
			prefix += fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
		}
	}
	if bold {
		prefix += "\033[1m"
	}
	if dim {
		prefix += "\033[2m"
	}

	if prefix == "" {
		return text
	}
	return prefix + text + "\033[0m"
}
