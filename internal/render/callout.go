package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var calloutTypes = map[string]struct{ Icon, Header, Color string }{
	"success": {"✓", "Success", "#44FF88"},
	"error":   {"✗", "Error", "#FF4444"},
	"warning": {"⚠", "Warning", "#FFD700"},
	"info":    {"ℹ", "Info", "#4488FF"},
}

func RenderCallout(text string, calloutType string) string {
	ct, ok := calloutTypes[calloutType]
	if !ok {
		ct = calloutTypes["info"]
	}
	header := fmt.Sprintf("%s %s", ct.Icon, ct.Header)
	// Wrap text to fit within the border box (border=2 + padding=2 = 4 cols of chrome)
	availWidth := termWidth() - 4
	if availWidth > 0 {
		text = wrapText(text, availWidth)
	}
	content := header + "\n" + text
	border := lipgloss.RoundedBorder()
	style := lipgloss.NewStyle().Border(border).Padding(0, 1).BorderForeground(lipgloss.Color(ct.Color))
	return style.Render(content)
}

func CalloutDefaultColor(calloutType string) string {
	ct, ok := calloutTypes[calloutType]
	if !ok {
		return ""
	}
	return ct.Color
}
