package watch

import (
	"bytes"
	"strings"
	"testing"
)

func up1(n int) string { return strings.Repeat("\033[1A", n) }

const (
	col0 = "\033[0G"
	hide = "\033[?25l"
	show = "\033[?25h"
	dn1  = "\033[1B"
)

func padLine(content string, width int) string {
	vis := len(content) // test lines have no ANSI
	maxCols := width - 1
	pad := maxCols - vis
	if pad < 0 {
		pad = 0
	}
	return col0 + content + strings.Repeat(" ", pad) + "\n"
}

func TestExact_FirstPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"div", "status"})

	expected := hide + dn1 + up1(1) +
		padLine("div", 80) + padLine("status", 80) +
		show

	if got := buf.String(); got != expected {
		t.Errorf("first panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_PanelRedraw(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"div", "old"})
	buf.Reset()
	r.DrawPanel([]string{"div", "new"})

	expected := hide + up1(2) +
		padLine("div", 80) + padLine("new", 80) +
		show

	if got := buf.String(); got != expected {
		t.Errorf("redraw:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_ScrollWithPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"div", "status"})
	buf.Reset()
	r.Render([]string{"output1"}, []string{"div", "status2"})

	expected := hide + up1(2) +
		padLine("output1", 80) + padLine("div", 80) + padLine("status2", 80) +
		show

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
		t.Errorf("no panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_ClearPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"a", "b", "c"})
	buf.Reset()
	r.ClearPanel()

	expected := hide + up1(3) +
		padLine("", 80) + padLine("", 80) + padLine("", 80) +
		up1(3) + show

	if got := buf.String(); got != expected {
		t.Errorf("clear:\nexpected: %q\n     got: %q", expected, got)
	}
}

func TestExact_PanelShrink(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)
	r.DrawPanel([]string{"a", "b", "c"})
	buf.Reset()
	r.DrawPanel([]string{"a", "b"})

	expected := hide + up1(3) +
		padLine("a", 80) + padLine("b", 80) +
		padLine("", 80) + up1(1) +
		show

	if got := buf.String(); got != expected {
		t.Errorf("shrink:\nexpected: %q\n     got: %q", expected, got)
	}
}
