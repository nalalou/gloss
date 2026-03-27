package render

import (
	"strings"
	"testing"
)

func TestRenderDividerPlain(t *testing.T) {
	result := RenderDivider("", 20, "heavy")
	if len([]rune(result)) != 20 {
		t.Errorf("expected 20 runes, got %d", len([]rune(result)))
	}
	if !strings.Contains(result, "━") {
		t.Error("heavy should use ━")
	}
}

func TestRenderDividerWithLabel(t *testing.T) {
	result := RenderDivider("Title", 30, "heavy")
	if !strings.Contains(result, "Title") {
		t.Error("should contain label")
	}
	if !strings.Contains(result, "━") {
		t.Error("should contain ━")
	}
}

func TestRenderDividerStyles(t *testing.T) {
	styles := map[string]string{
		"heavy": "━", "light": "─", "double": "═",
		"dashed": "╌", "dots": "·", "ascii": "-",
	}
	for style, expected := range styles {
		result := RenderDivider("", 10, style)
		if !strings.Contains(result, expected) {
			t.Errorf("style %q: expected %q", style, expected)
		}
	}
}

func TestRenderDividerLabelLongerThanWidth(t *testing.T) {
	result := RenderDivider("Very Long Label", 10, "heavy")
	if !strings.Contains(result, "Very Long Label") {
		t.Error("label should be preserved")
	}
}
