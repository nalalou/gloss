package render

import (
	"fmt"
	"strings"
)

type rgbColor struct {
	R, G, B uint8
}

func parseHex(hex string) (rgbColor, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return rgbColor{}, fmt.Errorf("invalid hex color: %q (want 6 hex digits)", hex)
	}
	var r, g, b uint8
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return rgbColor{}, fmt.Errorf("parse hex %q: %w", hex, err)
	}
	return rgbColor{r, g, b}, nil
}

func interpolate(a, b rgbColor, t float64) rgbColor {
	lerp := func(x, y uint8) uint8 {
		return uint8(float64(x)*(1-t) + float64(y)*t)
	}
	return rgbColor{lerp(a.R, b.R), lerp(a.G, b.G), lerp(a.B, b.B)}
}

func ansiColor(c rgbColor, s string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, s)
}

var gradientPresets = map[string][2]string{
	"fire":   {"#FF4500", "#FFD700"},
	"ocean":  {"#006994", "#00CED1"},
	"mono":   {"#FFFFFF", "#888888"},
	"neon":   {"#FF6B9D", "#C0F0A0"},
	"aurora": {"#00FA9A", "#9370DB"},
}

// GradientPreset returns the start and end colors for a named preset.
func GradientPreset(name string) (start, end rgbColor, ok bool) {
	pair, ok := gradientPresets[name]
	if !ok {
		return rgbColor{}, rgbColor{}, false
	}
	s, _ := parseHex(pair[0])
	e, _ := parseHex(pair[1])
	return s, e, true
}

// ApplyGradient colors lines with a gradient from start to end.
func ApplyGradient(lines []string, start, end rgbColor, direction string, noColor bool) []string {
	if noColor {
		return lines
	}
	result := make([]string, len(lines))

	if direction == "vertical" {
		for i, line := range lines {
			t := 0.0
			if len(lines) > 1 {
				t = float64(i) / float64(len(lines)-1)
			}
			result[i] = ansiColor(interpolate(start, end, t), line)
		}
		return result
	}

	// Horizontal: color each rune individually
	for i, line := range lines {
		runes := []rune(line)
		n := len(runes)
		var sb strings.Builder
		for j, ch := range runes {
			t := 0.0
			if n > 1 {
				t = float64(j) / float64(n-1)
			}
			c := interpolate(start, end, t)
			sb.WriteString(ansiColor(c, string(ch)))
		}
		result[i] = sb.String()
	}
	return result
}

// ApplySolidColor colors all lines with a single hex color.
func ApplySolidColor(lines []string, hex string, noColor bool) []string {
	if noColor || hex == "" {
		return lines
	}
	c, err := parseHex(hex)
	if err != nil {
		return lines
	}
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = ansiColor(c, line)
	}
	return result
}
