package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ApplyBorder wraps text in a lipgloss border.
func ApplyBorder(text string, borderStyle string) string {
	var b lipgloss.Border
	switch borderStyle {
	case "single":
		b = lipgloss.NormalBorder()
	case "double":
		b = lipgloss.DoubleBorder()
	case "rounded":
		b = lipgloss.RoundedBorder()
	case "thick":
		b = lipgloss.ThickBorder()
	default:
		return text
	}
	return lipgloss.NewStyle().Border(b).Padding(0, 1).Render(text)
}

// ApplyShadow adds a dim gray drop-shadow below the text.
func ApplyShadow(text string) string {
	lines := strings.Split(text, "\n")
	shadowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))

	shadowLines := make([]string, len(lines))
	for i, line := range lines {
		shadowLines[i] = " " + shadowStyle.Render(stripANSI(line))
	}

	result := make([]string, 0, len(lines)+1)
	result = append(result, lines...)
	result = append(result, shadowLines[len(shadowLines)-1])

	return strings.Join(result, "\n")
}

// stripANSI removes ANSI escape codes from s.
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, ch := range s {
		if ch == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(ch)
	}
	return result.String()
}
