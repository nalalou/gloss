package render

import (
	"strings"
	"testing"
)

func TestRenderBadgeSuccess(t *testing.T) {
	r := RenderBadge("Tests passing", "success")
	if !strings.Contains(r, "✓") || !strings.Contains(r, "Tests passing") {
		t.Error("wrong output")
	}
}

func TestRenderBadgeError(t *testing.T) {
	if !strings.Contains(RenderBadge("Fail", "error"), "✗") {
		t.Error("missing ✗")
	}
}

func TestRenderBadgeWarning(t *testing.T) {
	if !strings.Contains(RenderBadge("Warn", "warning"), "⚠") {
		t.Error("missing ⚠")
	}
}

func TestRenderBadgeInfo(t *testing.T) {
	if !strings.Contains(RenderBadge("Info", "info"), "ℹ") {
		t.Error("missing ℹ")
	}
}

func TestRenderBadgePlain(t *testing.T) {
	r := RenderBadge("Hello", "plain")
	if strings.ContainsAny(r, "✓✗⚠ℹ") {
		t.Error("plain should have no icon")
	}
}

func TestBadgeDefaultColor(t *testing.T) {
	if BadgeDefaultColor("success") != "#44FF88" {
		t.Error("wrong color")
	}
}
