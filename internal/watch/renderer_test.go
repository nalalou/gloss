package watch

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderPanelToBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build"})
	if !strings.Contains(buf.String(), "gloss") {
		t.Errorf("missing panel: %q", buf.String())
	}
}

func TestRenderScrollLine(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.WriteScroll("Hello world")
	if !strings.Contains(buf.String(), "Hello world") {
		t.Errorf("missing scroll: %q", buf.String())
	}
}

func TestRenderPanelRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"─── gloss ───", "  ○ Build"})
	buf.Reset()
	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build"})
	out := buf.String()
	// Should move up 2 lines to erase old panel
	if !strings.Contains(out, "\033[2A") {
		t.Error("expected cursor-up(2) for 2-line panel redraw")
	}
}

func TestRenderHideCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.HideCursor()
	if !strings.Contains(buf.String(), "\033[?25l") {
		t.Error("missing hide cursor")
	}
}

func TestRenderShowCursor(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.ShowCursor()
	if !strings.Contains(buf.String(), "\033[?25h") {
		t.Error("missing show cursor")
	}
}

func TestRenderScrollNoPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.Render([]string{"line1", "line2"}, nil)
	out := buf.String()
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Errorf("missing scroll lines: %q", out)
	}
}

func TestClearPanelFunc(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	// First draw a panel
	r.DrawPanel([]string{"a", "b", "c"})
	buf.Reset()
	r.ClearPanel()
	out := buf.String()
	if !strings.Contains(out, "\033[3A") {
		t.Error("expected cursor-up(3) for 3-line clear")
	}
}
