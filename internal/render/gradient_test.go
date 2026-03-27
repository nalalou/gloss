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

func TestInterpolateMulti3Colors(t *testing.T) {
	red := rgbColor{255, 0, 0}
	green := rgbColor{0, 255, 0}
	blue := rgbColor{0, 0, 255}
	colors := []rgbColor{red, green, blue}

	c := interpolateMulti(colors, 0.0)
	if c.R != 255 || c.G != 0 || c.B != 0 {
		t.Errorf("t=0.0: expected red, got R=%d G=%d B=%d", c.R, c.G, c.B)
	}
	c = interpolateMulti(colors, 0.5)
	if c.R != 0 || c.G != 255 || c.B != 0 {
		t.Errorf("t=0.5: expected green, got R=%d G=%d B=%d", c.R, c.G, c.B)
	}
	c = interpolateMulti(colors, 1.0)
	if c.R != 0 || c.G != 0 || c.B != 255 {
		t.Errorf("t=1.0: expected blue, got R=%d G=%d B=%d", c.R, c.G, c.B)
	}
	c = interpolateMulti(colors, 0.25)
	if c.R < 120 || c.R > 135 || c.G < 120 || c.G > 135 {
		t.Errorf("t=0.25: expected ~halfway, got R=%d G=%d B=%d", c.R, c.G, c.B)
	}
}

func TestInterpolateMulti1Color(t *testing.T) {
	c := interpolateMulti([]rgbColor{{255, 0, 0}}, 0.5)
	if c.R != 255 {
		t.Errorf("single color: expected red, got R=%d", c.R)
	}
}

func TestInterpolateMulti2Colors(t *testing.T) {
	c := interpolateMulti([]rgbColor{{0, 0, 0}, {255, 255, 255}}, 0.5)
	if c.R < 127 || c.R > 128 {
		t.Errorf("2 colors: expected ~127, got R=%d", c.R)
	}
}

func TestApplyGradientVertical(t *testing.T) {
	lines := []string{"AAAA", "BBBB", "CCCC"}
	black, _ := parseHex("#000000")
	white, _ := parseHex("#FFFFFF")
	result := ApplyGradient(lines, []rgbColor{black, white}, "vertical", false)
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	for _, line := range result {
		if !strings.Contains(line, "\033[") {
			t.Errorf("expected ANSI codes, got %q", line)
		}
	}
}

func TestApplyGradientNoColorPassthrough(t *testing.T) {
	lines := []string{"AAAA", "BBBB"}
	red, _ := parseHex("#FF0000")
	blue, _ := parseHex("#0000FF")
	result := ApplyGradient(lines, []rgbColor{red, blue}, "horizontal", true)
	for i, line := range result {
		if strings.Contains(line, "\033[") {
			t.Errorf("noColor line %d has ANSI: %q", i, line)
		}
	}
}

func TestGradientPresetReturnsSlice(t *testing.T) {
	colors, ok := GradientPreset("fire")
	if !ok {
		t.Fatal("fire not found")
	}
	if len(colors) < 2 {
		t.Errorf("expected >=2 colors, got %d", len(colors))
	}
}

func TestGradientPresetRainbow(t *testing.T) {
	colors, ok := GradientPreset("rainbow")
	if !ok {
		t.Fatal("rainbow not found")
	}
	if len(colors) != 6 {
		t.Errorf("expected 6 colors, got %d", len(colors))
	}
}

func TestApplyGradientMultiStop(t *testing.T) {
	lines := []string{"ABCDEFGHIJ"}
	red, _ := parseHex("#FF0000")
	green, _ := parseHex("#00FF00")
	blue, _ := parseHex("#0000FF")
	result := ApplyGradient(lines, []rgbColor{red, green, blue}, "horizontal", false)
	if !strings.Contains(result[0], "\033[") {
		t.Error("expected ANSI codes in multi-stop output")
	}
}

func TestAllPresetsLoad(t *testing.T) {
	for _, name := range GradientPresetNames() {
		colors, ok := GradientPreset(name)
		if !ok {
			t.Errorf("preset %q not found", name)
			continue
		}
		if len(colors) < 2 {
			t.Errorf("preset %q: expected >=2 colors, got %d", name, len(colors))
		}
	}
}

func TestApplyGradientWithOffset(t *testing.T) {
	lines := []string{"AAAA"}
	red, _ := parseHex("#FF0000")
	blue, _ := parseHex("#0000FF")
	colors := []rgbColor{red, blue}
	r0 := ApplyGradientWithOffset(lines, colors, "horizontal", 0.0, false)
	r5 := ApplyGradientWithOffset(lines, colors, "horizontal", 0.5, false)
	if r0[0] == r5[0] {
		t.Error("different offsets should produce different output")
	}
}
