package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var calloutTypes = map[string]struct{ Icon, Header, Color string }{
	"success": {"✓", "Success", "#00FF00"},
	"error":   {"✗", "Error", "#FF0000"},
	"warning": {"⚠", "Warning", "#FFD700"},
	"info":    {"ℹ", "Info", "#00BFFF"},
}

func RenderCallout(text string, calloutType string) string {
	ct, ok := calloutTypes[calloutType]
	if !ok {
		ct = calloutTypes["info"]
	}
	header := fmt.Sprintf("%s %s", ct.Icon, ct.Header)
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
