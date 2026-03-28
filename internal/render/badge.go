package render

var badgeTypes = map[string]struct{ Icon, Color string }{
	"success": {"✓", "#44FF88"},
	"error":   {"✗", "#FF4444"},
	"warning": {"⚠", "#FFD700"},
	"info":    {"ℹ", "#4488FF"},
	"plain":   {"", ""},
}

func RenderBadge(text string, badgeType string) string {
	bt, ok := badgeTypes[badgeType]
	if !ok {
		bt = badgeTypes["plain"]
	}
	if bt.Icon == "" {
		return text
	}
	return bt.Icon + " " + text
}

func BadgeDefaultColor(badgeType string) string {
	bt, ok := badgeTypes[badgeType]
	if !ok {
		return ""
	}
	return bt.Color
}
