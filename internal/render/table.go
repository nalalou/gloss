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

	// Determine available width for wrapping: terminal minus border chrome (4) minus separators
	tw := termWidth()
	borderChrome := 0
	if borderStyle != "none" && borderStyle != "" {
		borderChrome = 4 // 2 border + 2 padding
	}
	separatorWidth := (cols - 1) * 3 // " │ " between columns
	availForCells := tw - borderChrome - separatorWidth
	maxCellWidth := 0
	if cols > 0 && availForCells > 0 {
		maxCellWidth = availForCells / cols
	}

	// Wrap cell content if it exceeds max cell width
	if maxCellWidth > 10 {
		for r, row := range rows {
			for c, cell := range row {
				if displayWidth(cell) > maxCellWidth {
					rows[r][c] = wrapText(cell, maxCellWidth)
				}
			}
		}
	}

	widths := make([]int, cols)
	for _, row := range rows {
		for i, cell := range row {
			// For wrapped cells, use the widest line
			for _, line := range strings.Split(cell, "\n") {
				if displayWidth(line) > widths[i] { widths[i] = displayWidth(line) }
			}
		}
	}
	var sb strings.Builder
	for _, row := range rows {
		var cells []string
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(row) { cell = row[i] }
			cells = append(cells, cell+strings.Repeat(" ", widths[i]-displayWidth(cell)))
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
