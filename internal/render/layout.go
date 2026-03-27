package render

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// ApplyLayout applies alignment and width constraints using lipgloss.
func ApplyLayout(text string, align string, width int, noColor bool) string {
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	var lipAlign lipgloss.Position
	switch align {
	case "center":
		lipAlign = lipgloss.Center
	case "right":
		lipAlign = lipgloss.Right
	default:
		lipAlign = lipgloss.Left
	}

	style := lipgloss.NewStyle().
		Width(width).
		Align(lipAlign)

	return style.Render(text)
}
