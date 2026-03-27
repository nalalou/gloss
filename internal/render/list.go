package render

import (
	"fmt"
	"strings"
)

var listIcons = map[string]string{
	"bullet": "•", "arrow": "→", "dash": "—", "star": "★", "check": "✓",
}
var statusIcons = map[string]string{
	"done": "✓", "pending": "○", "fail": "✗",
}

func RenderList(items []string, style string, statusMode bool) string {
	if len(items) == 0 { return "" }
	var lines []string
	for i, item := range items {
		if statusMode {
			lastColon := strings.LastIndex(item, ":")
			text, status := item, ""
			if lastColon != -1 {
				text = item[:lastColon]
				status = item[lastColon+1:]
			}
			icon := "•"
			if ic, ok := statusIcons[status]; ok { icon = ic }
			lines = append(lines, fmt.Sprintf("%s %s", icon, text))
		} else if style == "numbered" {
			lines = append(lines, fmt.Sprintf("%d. %s", i+1, item))
		} else {
			icon := "•"
			if ic, ok := listIcons[style]; ok { icon = ic }
			lines = append(lines, fmt.Sprintf("%s %s", icon, item))
		}
	}
	return strings.Join(lines, "\n")
}
