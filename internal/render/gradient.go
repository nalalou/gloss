package render

import (
	"fmt"
	"math"
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

func interpolateMulti(colors []rgbColor, t float64) rgbColor {
	if len(colors) == 0 {
		return rgbColor{}
	}
	if len(colors) == 1 {
		return colors[0]
	}
	t = math.Max(0, math.Min(1, t))
	n := len(colors) - 1
	segment := t * float64(n)
	idx := int(segment)
	if idx >= n {
		return colors[n]
	}
	localT := segment - float64(idx)
	return interpolate(colors[idx], colors[idx+1], localT)
}

func ansiColor(c rgbColor, s string) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", c.R, c.G, c.B, s)
}

var gradientPresets = map[string][]string{
	"fire":      {"#FF4500", "#FFD700"},
	"ocean":     {"#006994", "#00CED1"},
	"mono":      {"#FFFFFF", "#888888"},
	"neon":      {"#FF6B9D", "#C0F0A0"},
	"aurora":    {"#00FA9A", "#9370DB"},
	"sunset":    {"#FF6B35", "#F7931E", "#FF355E"},
	"synthwave": {"#FF00FF", "#00FFFF", "#FF00FF"},
	"matrix":    {"#00FF00", "#003300"},
	"cyberpunk": {"#FFFF00", "#FF00FF"},
	"pastel":    {"#FFB3BA", "#BAE1FF"},
	"lavender":  {"#9B59B6", "#3498DB"},
	"ice":       {"#FFFFFF", "#00BFFF", "#0000FF"},
	"autumn":    {"#8B4513", "#FF8C00", "#FFD700"},
	"mint":      {"#00FF7F", "#20B2AA"},
	"rainbow":   {"#FF0000", "#FF8800", "#FFFF00", "#00FF00", "#0088FF", "#8800FF"},
}

func GradientPresetNames() []string {
	return []string{"fire", "ocean", "mono", "neon", "aurora", "sunset", "synthwave", "matrix", "cyberpunk", "pastel", "lavender", "ice", "autumn", "mint", "rainbow"}
}

func GradientPreset(name string) ([]rgbColor, bool) {
	hexes, ok := gradientPresets[name]
	if !ok {
		return nil, false
	}
	colors := make([]rgbColor, len(hexes))
	for i, h := range hexes {
		c, err := parseHex(h)
		if err != nil {
			return nil, false
		}
		colors[i] = c
	}
	return colors, true
}

func ParseGradientColors(hexes []string) ([]rgbColor, error) {
	colors := make([]rgbColor, len(hexes))
	for i, h := range hexes {
		c, err := parseHex(h)
		if err != nil {
			return nil, fmt.Errorf("gradient color %d: %w", i+1, err)
		}
		colors[i] = c
	}
	return colors, nil
}

func ApplyGradient(lines []string, colors []rgbColor, direction string, noColor bool) []string {
	if noColor || len(colors) == 0 {
		return lines
	}
	result := make([]string, len(lines))

	if direction == "vertical" {
		for i, line := range lines {
			t := 0.0
			if len(lines) > 1 {
				t = float64(i) / float64(len(lines)-1)
			}
			result[i] = ansiColor(interpolateMulti(colors, t), line)
		}
		return result
	}

	for i, line := range lines {
		runes := []rune(line)
		n := len(runes)
		var sb strings.Builder
		for j, ch := range runes {
			t := 0.0
			if n > 1 {
				t = float64(j) / float64(n-1)
			}
			c := interpolateMulti(colors, t)
			sb.WriteString(ansiColor(c, string(ch)))
		}
		result[i] = sb.String()
	}
	return result
}

func ApplyGradientWithOffset(lines []string, colors []rgbColor, direction string, offset float64, noColor bool) []string {
	if noColor || len(colors) == 0 {
		return lines
	}
	result := make([]string, len(lines))

	if direction == "vertical" {
		for i, line := range lines {
			t := 0.0
			if len(lines) > 1 {
				t = float64(i) / float64(len(lines)-1)
			}
			t = math.Mod(t+offset, 1.0)
			result[i] = ansiColor(interpolateMulti(colors, t), line)
		}
		return result
	}

	for i, line := range lines {
		runes := []rune(line)
		n := len(runes)
		var sb strings.Builder
		for j, ch := range runes {
			t := 0.0
			if n > 1 {
				t = float64(j) / float64(n-1)
			}
			t = math.Mod(t+offset, 1.0)
			c := interpolateMulti(colors, t)
			sb.WriteString(ansiColor(c, string(ch)))
		}
		result[i] = sb.String()
	}
	return result
}

func ColorizeLines(lines []string, gradient []string, color string, noColor bool) []string {
	if noColor {
		return lines
	}
	if len(gradient) >= 2 {
		colors, err := ParseGradientColors(gradient)
		if err == nil {
			return ApplyGradient(lines, colors, "horizontal", false)
		}
	}
	if color != "" {
		return ApplySolidColor(lines, color, false)
	}
	return lines
}

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
