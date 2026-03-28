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

// ApplyShadow adds a 2D drop-shadow offset right by 1 and down by 1.
// Shadow characters are drawn in dim gray behind the original content.
func ApplyShadow(text string) string {
	lines := strings.Split(text, "\n")
	shadowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))

	// Build shadow lines: each original line becomes a shadow shifted right by 1.
	shadowLines := make([]string, len(lines))
	for i, line := range lines {
		shadowLines[i] = " " + shadowStyle.Render(stripANSI(line))
	}

	// Composite: for each output row, overlay the original content on top of
	// the shadow from the line above (shifted down by 1).
	// Row 0: original line 0 (no shadow behind it — shadow is offset down)
	// Row i (1..n-1): original line i overlaid on shadow of line i-1
	// Row n: shadow-only from line n-1 (the bottom shadow edge)
	result := make([]string, 0, len(lines)+1)
	result = append(result, lines[0])
	for i := 1; i < len(lines); i++ {
		// The shadow from line i-1 appears behind line i.
		// We need to composite: take the original line i, and where it ends,
		// show the remaining shadow from line i-1.
		origWidth := displayWidth(stripANSI(lines[i]))
		shadowPlain := stripANSI(lines[i-1])
		shadowVisWidth := displayWidth(shadowPlain) + 1 // +1 for the leading space offset

		if shadowVisWidth > origWidth {
			// There's shadow peeking out to the right of the original content.
			// Render the original line, then append the right portion of shadow.
			rightShadowStart := origWidth - 1 // -1 because shadow is offset right by 1
			if rightShadowStart < 0 {
				rightShadowStart = 0
			}
			shadowRunes := []rune(shadowPlain)
			rightPortion := ""
			if rightShadowStart < len(shadowRunes) {
				rightPortion = string(shadowRunes[rightShadowStart:])
			}
			if rightPortion != "" {
				result = append(result, lines[i]+shadowStyle.Render(rightPortion))
			} else {
				result = append(result, lines[i])
			}
		} else {
			result = append(result, lines[i])
		}
	}
	// Bottom shadow row: shadow of the last line
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
