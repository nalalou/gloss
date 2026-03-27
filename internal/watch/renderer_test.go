package watch

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderPanelToBuffer(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build"}, 0)
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
	r.DrawPanel([]string{"─── gloss ───", "  ○ Build"}, 0)
	buf.Reset()
	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build"}, 2)
	out := buf.String()
	if !strings.Contains(out, "\033[2A") {
		t.Error("expected cursor-up for redraw")
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

func TestWriteScrollWithPanelNoPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.WriteScrollWithPanel([]string{"line1", "line2"}, nil, 0)
	out := buf.String()
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Errorf("missing scroll lines: %q", out)
	}
}

func TestClearPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 60, false)
	r.ClearPanel(3)
	out := buf.String()
	if !strings.Contains(out, "\033[3A") {
		t.Error("expected cursor-up in clear")
	}
}
