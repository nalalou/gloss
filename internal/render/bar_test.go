package render

import (
	"strings"
	"testing"
)

func TestRenderBarBasic(t *testing.T) {
	result := RenderBar(50, 20, "blocks", true)
	if !strings.Contains(result, "█") { t.Error("should contain █") }
	if !strings.Contains(result, "░") { t.Error("should contain ░") }
	if !strings.Contains(result, "50%") { t.Error("should contain 50%") }
}

func TestRenderBarFull(t *testing.T) {
	result := RenderBar(100, 20, "blocks", true)
	if strings.Contains(result, "░") { t.Error("100% should have no empty") }
}

func TestRenderBarEmpty(t *testing.T) {
	result := RenderBar(0, 20, "blocks", true)
	if strings.Contains(result, "█") { t.Error("0% should have no filled") }
}

func TestRenderBarNoPercent(t *testing.T) {
	result := RenderBar(50, 20, "blocks", false)
	if strings.Contains(result, "%") { t.Error("should not show percent") }
}

func TestRenderBarStyles(t *testing.T) {
	for style, expected := range map[string]string{"blocks": "█", "dots": "▮", "ascii": "=", "thin": "━"} {
		result := RenderBar(50, 20, style, true)
		if !strings.Contains(result, expected) { t.Errorf("style %q: missing %q", style, expected) }
	}
}

func TestRenderBarClamp(t *testing.T) {
	if !strings.Contains(RenderBar(150, 20, "blocks", true), "100%") { t.Error(">100 should clamp") }
	if !strings.Contains(RenderBar(-10, 20, "blocks", true), "0%") { t.Error("<0 should clamp") }
}

func TestParseBarValue(t *testing.T) {
	tests := []struct{ input string; expected float64 }{
		{"75", 75}, {"0.75", 75}, {"3/4", 75}, {"100", 100}, {"0", 0}, {"1.0", 100},
	}
	for _, tt := range tests {
		v, err := ParseBarValue(tt.input)
		if err != nil { t.Errorf("ParseBarValue(%q): %v", tt.input, err); continue }
		if v != tt.expected { t.Errorf("ParseBarValue(%q) = %v, want %v", tt.input, v, tt.expected) }
	}
}

func TestParseBarValueInvalid(t *testing.T) {
	if _, err := ParseBarValue("abc"); err == nil { t.Error("expected error") }
}
