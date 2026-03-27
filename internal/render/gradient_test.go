package render

import (
	"strings"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	c, err := parseHex("#FF6B9D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.R != 0xFF || c.G != 0x6B || c.B != 0x9D {
		t.Errorf("wrong color: got R=%d G=%d B=%d", c.R, c.G, c.B)
	}
}

func TestParseHexWithoutHash(t *testing.T) {
	_, err := parseHex("FF6B9D")
	if err != nil {
		t.Fatalf("should accept hex without #: %v", err)
	}
}

func TestInterpolateColors(t *testing.T) {
	black := rgbColor{0, 0, 0}
	white := rgbColor{255, 255, 255}
	mid := interpolate(black, white, 0.5)
	if mid.R != 127 && mid.R != 128 {
		t.Errorf("expected ~127 for midpoint R, got %d", mid.R)
	}
}

func TestApplyGradientVertical(t *testing.T) {
	lines := []string{"AAAA", "BBBB", "CCCC"}
	start, _ := parseHex("#000000")
	end, _ := parseHex("#FFFFFF")
	result := ApplyGradient(lines, start, end, "vertical", false)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	for _, line := range result {
		if !strings.Contains(line, "\033[") {
			t.Errorf("expected ANSI codes in output, got %q", line)
		}
	}
}

func TestApplyGradientNoColorPassthrough(t *testing.T) {
	lines := []string{"AAAA", "BBBB"}
	start, _ := parseHex("#FF0000")
	end, _ := parseHex("#0000FF")
	result := ApplyGradient(lines, start, end, "horizontal", true)
	for i, line := range result {
		if strings.Contains(line, "\033[") {
			t.Errorf("expected no ANSI codes when noColor=true, line %d: %q", i, line)
		}
	}
}

func TestGradientPresetFire(t *testing.T) {
	start, end, ok := GradientPreset("fire")
	if !ok {
		t.Fatal("preset 'fire' not found")
	}
	if start.R == 0 && start.G == 0 && start.B == 0 {
		t.Error("expected non-zero start color for fire preset")
	}
	_ = end
}
