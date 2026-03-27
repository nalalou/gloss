package render

import (
	"strings"
	"testing"
)

func TestResolveColorNamed(t *testing.T) {
	hex, ok := ResolveColor("red")
	if !ok || hex != "#FF0000" {
		t.Errorf("expected #FF0000, got %q ok=%v", hex, ok)
	}
}

func TestResolveColorHex(t *testing.T) {
	hex, ok := ResolveColor("#FF6B9D")
	if !ok || hex != "#FF6B9D" {
		t.Errorf("expected #FF6B9D, got %q ok=%v", hex, ok)
	}
}

func TestResolveColorHexNoHash(t *testing.T) {
	hex, ok := ResolveColor("FF6B9D")
	if !ok || hex != "#FF6B9D" {
		t.Errorf("expected #FF6B9D, got %q ok=%v", hex, ok)
	}
}

func TestResolveColorInvalid(t *testing.T) {
	_, ok := ResolveColor("notacolor")
	if ok {
		t.Error("expected false for invalid color")
	}
}

func TestRenderStyledFgOnly(t *testing.T) {
	result := RenderStyled("hello", "#FF0000", false, false)
	if !strings.Contains(result, "\033[38;2;255;0;0m") {
		t.Error("should contain red ANSI code")
	}
	if !strings.Contains(result, "hello") {
		t.Error("should contain text")
	}
	if !strings.HasSuffix(result, "\033[0m") {
		t.Error("should end with reset")
	}
}

func TestRenderStyledBold(t *testing.T) {
	result := RenderStyled("hello", "", true, false)
	if !strings.Contains(result, "\033[1m") {
		t.Error("should contain bold code")
	}
}

func TestRenderStyledDim(t *testing.T) {
	result := RenderStyled("hello", "", false, true)
	if !strings.Contains(result, "\033[2m") {
		t.Error("should contain dim code")
	}
}

func TestRenderStyledCombined(t *testing.T) {
	result := RenderStyled("hello", "#FF0000", true, false)
	if !strings.Contains(result, "\033[38;2;") {
		t.Error("should contain color code")
	}
	if !strings.Contains(result, "\033[1m") {
		t.Error("should contain bold code")
	}
}

func TestRenderStyledNoStyle(t *testing.T) {
	result := RenderStyled("hello", "", false, false)
	if result != "hello" {
		t.Errorf("no style should return plain text, got %q", result)
	}
}

func TestAllNamedColors(t *testing.T) {
	names := []string{"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white", "gray", "grey", "orange", "pink", "purple"}
	for _, name := range names {
		hex, ok := ResolveColor(name)
		if !ok {
			t.Errorf("named color %q not found", name)
			continue
		}
		if hex == "" {
			t.Errorf("named color %q returned empty hex", name)
		}
	}
}
