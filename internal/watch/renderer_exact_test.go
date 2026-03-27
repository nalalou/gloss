package watch

import (
	"bytes"
	"testing"
)

const (
	tclr     = "\033[2K"
	tsyncOn  = "\033[?2026h"
	tsyncOff = "\033[?2026l"
)

func tup(n int) string {
	s := "\033["
	if n >= 10 {
		s += string(rune('0'+n/10)) + string(rune('0'+n%10))
	} else {
		s += string(rune('0' + n))
	}
	return s + "A"
}

// New model: cursor ends AFTER the last panel line (every line has \n).
// To erase: move up N, clear N lines, move up N again.
// Then write new content.

func TestExact_FirstPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"div", "status"})

	// No previous panel → no erase. Just write 2 lines.
	expected := tsyncOn +
		"\r" + tclr + "div\n" +
		"\r" + tclr + "status\n" +
		tsyncOff

	if got := buf.String(); got != expected {
		t.Errorf("first panel:\nexpected: %q\n     got: %q", expected, got)
	}
	if r.LinesOnScreen() != 2 {
		t.Errorf("linesOnScreen: %d", r.LinesOnScreen())
	}
}

func TestExact_PanelRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	// First draw
	r.DrawPanel([]string{"div", "old"})
	buf.Reset()

	// Redraw (same size)
	r.DrawPanel([]string{"div", "new"})

	// Erase old 2 lines: up(2), clear+\n x2, up(2). Then write new 2 lines.
	expected := tsyncOn +
		tup(2) +
		"\r" + tclr + "\n" +
		"\r" + tclr + "\n" +
		tup(2) +
		"\r" + tclr + "div\n" +
		"\r" + tclr + "new\n" +
		tsyncOff

	if got := buf.String(); got != expected {
		t.Errorf("panel redraw:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_ScrollWithPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	// Establish a 2-line panel
	r.DrawPanel([]string{"div", "status"})
	buf.Reset()

	// Scroll lines + panel redraw
	r.Render([]string{"output1"}, []string{"div", "status2"})

	// Erase old 2: up(2), clear x2, up(2). Print scroll. Print panel.
	expected := tsyncOn +
		tup(2) +
		"\r" + tclr + "\n" +
		"\r" + tclr + "\n" +
		tup(2) +
		"\r" + tclr + "output1\n" +
		"\r" + tclr + "div\n" +
		"\r" + tclr + "status2\n" +
		tsyncOff

	if got := buf.String(); got != expected {
		t.Errorf("scroll+panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_ScrollNoPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.Render([]string{"hello", "world"}, nil)
	expected := "hello\nworld\n"
	if got := buf.String(); got != expected {
		t.Errorf("scroll no panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_ClearPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"a", "b", "c"})
	buf.Reset()
	r.ClearPanel()

	expected := tup(3) +
		"\r" + tclr + "\n" +
		"\r" + tclr + "\n" +
		"\r" + tclr + "\n" +
		tup(3)

	if got := buf.String(); got != expected {
		t.Errorf("clear:\nexpected: %q\n     got: %q", expected, got)
	}
	if r.LinesOnScreen() != 0 {
		t.Errorf("linesOnScreen after clear: %d", r.LinesOnScreen())
	}
}
