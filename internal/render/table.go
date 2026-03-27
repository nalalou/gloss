package render

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func ParseKVArgs(args []string) [][]string {
	rows := make([][]string, 0, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 { rows = append(rows, parts) } else { rows = append(rows, []string{arg, ""}) }
	}
	return rows
}

func RenderTable(rows [][]string, borderStyle string) string {
	if len(rows) == 0 { return "" }
	cols := 0
	for _, row := range rows { if len(row) > cols { cols = len(row) } }
	widths := make([]int, cols)
	for _, row := range rows {
		for i, cell := range row { if len(cell) > widths[i] { widths[i] = len(cell) } }
	}
	var sb strings.Builder
	for _, row := range rows {
		var cells []string
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) { cell = row[i] }
			cells = append(cells, cell+strings.Repeat(" ", widths[i]-len(cell)))
		}
		sb.WriteString(strings.Join(cells, " │ ") + "\n")
	}
	content := strings.TrimRight(sb.String(), "\n")
	if borderStyle == "none" || borderStyle == "" { return content }
	var b lipgloss.Border
	switch borderStyle {
	case "single": b = lipgloss.NormalBorder()
	case "double": b = lipgloss.DoubleBorder()
	case "rounded": b = lipgloss.RoundedBorder()
	case "thick": b = lipgloss.ThickBorder()
	default: b = lipgloss.RoundedBorder()
	}
	return lipgloss.NewStyle().Border(b).Padding(0, 1).Render(content)
}
