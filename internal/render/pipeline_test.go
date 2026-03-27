package render

import (
	"strings"
	"testing"

	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/theme"
)

func loadTestFont(t *testing.T) font.Font {
	t.Helper()
	f, err := font.Load("block")
	if err != nil {
		t.Skipf("skipping: block font not available: %v", err)
	}
	return f
}

func TestRenderBasic(t *testing.T) {
	f := loadTestFont(t)
	opts := theme.Defaults()
	opts.NoColor = true

	result := Render("Hi", f, opts)
	if result == "" {
		t.Error("expected non-empty render output")
	}
	if !strings.Contains(result, "\n") {
		t.Error("expected multi-line output")
	}
}

func TestRenderWithGradient(t *testing.T) {
	f := loadTestFont(t)
	opts := theme.Defaults()
	opts.Gradient = []string{"#FF0000", "#0000FF"}

	result := Render("Hi", f, opts)
	if !strings.Contains(result, "\033[") {
		t.Error("expected ANSI codes in gradient output")
	}
}

func TestRenderNoColor(t *testing.T) {
	f := loadTestFont(t)
	opts := theme.Defaults()
	opts.NoColor = true
	opts.Gradient = []string{"#FF0000", "#0000FF"}

	result := Render("Hi", f, opts)
	if strings.Contains(result, "\033[") {
		t.Error("expected no ANSI codes when noColor=true")
	}
}

func TestRenderWithBorder(t *testing.T) {
	f := loadTestFont(t)
	opts := theme.Defaults()
	opts.Border = "rounded"
	opts.NoColor = true

	result := Render("Hi", f, opts)
	if !strings.Contains(result, "╭") {
		t.Error("expected rounded border character ╭ in output")
	}
}
