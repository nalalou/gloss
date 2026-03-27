package watch

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderPanelToBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"--- gloss ---", "  ok Build"})
	if !strings.Contains(buf.String(), "gloss") {
		t.Errorf("missing panel: %q", buf.String())
	}
}

func TestRenderScrollLine(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.WriteScroll("Hello world")
	if !strings.Contains(buf.String(), "Hello world") {
		t.Errorf("missing: %q", buf.String())
	}
}

func TestRenderPanelRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"--- gloss ---", "  old"})
	buf.Reset()
	r.DrawPanel([]string{"--- gloss ---", "  new"})
	if strings.Count(buf.String(), "\033[1A") < 2 {
		t.Error("expected 2x cursor-up(1)")
	}
}

func TestRenderHideCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.HideCursor()
	if !strings.Contains(buf.String(), "\033[?25l") {
		t.Error("missing hide")
	}
}

func TestRenderShowCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.ShowCursor()
	if !strings.Contains(buf.String(), "\033[?25h") {
		t.Error("missing show")
	}
}

func TestRenderScrollNoPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.Render([]string{"line1", "line2"}, nil)
	out := buf.String()
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Errorf("missing: %q", out)
	}
}

func TestClearPanelFunc(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"a", "b", "c"})
	buf.Reset()
	r.ClearPanel()
	if strings.Count(buf.String(), "\033[1A") < 3 {
		t.Error("expected 3+ cursor-up(1)")
	}
}
