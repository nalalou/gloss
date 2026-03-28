package render

import "strings"

var dividerChars = map[string]string{
	"heavy":  "━",
	"light":  "─",
	"double": "═",
	"dashed": "╌",
	"dots":   "·",
	"ascii":  "-",
}

func RenderDivider(label string, width int, style string) string {
	ch, ok := dividerChars[style]
	if !ok {
		ch = "━"
	}
	if label == "" {
		return strings.Repeat(ch, width)
	}
	padded := " " + label + " "
	remaining := width - displayWidth(padded)
	if remaining <= 0 {
		return padded
	}
	left := remaining / 2
	right := remaining - left
	return strings.Repeat(ch, left) + padded + strings.Repeat(ch, right)
}
