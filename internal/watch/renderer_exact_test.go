package watch

import (
	"bytes"
	"testing"
)

const (
	clr     = "\033[2K"
	syncOn  = "\033[?2026h"
	syncOff = "\033[?2026l"
)

func upN(n int) string {
	s := "\033["
	if n >= 10 {
		s += string(rune('0'+n/10)) + string(rune('0'+n%10))
	} else {
		s += string(rune('0' + n))
	}
	return s + "A"
}

// Scenario 1: First panel render, no previous panel.
// Cursor starts at a blank line. No cursor-up needed.
// Cursor ends ON last panel line (no trailing \n).
func TestExact_FirstPanelRender(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build", "  ✓ Test"}, 0)

	expected := syncOn +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ✓ Build" + "\n" +
		"\r" + clr + "  ✓ Test" +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("first panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 2: Panel-only update, same size.
// Cursor is ON last panel line. up(prevHeight-1) to reach first line.
func TestExact_PanelOnlyUpdate(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.DrawPanel([]string{"─── gloss ───", "  ⠙ Build", "  ✓ Test"}, 3)

	expected := syncOn +
		upN(2) +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ⠙ Build" + "\n" +
		"\r" + clr + "  ✓ Test" +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("panel update:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 3: Scroll lines + panel redraw.
// Cursor ON last panel line. up(prevHeight-1) to first panel line.
// Scroll lines overwrite old panel lines, panel redraws below.
func TestExact_ScrollWithPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.WriteScrollWithPanel(
		[]string{"output 1", "output 2"},
		[]string{"─── gloss ───", "  ✓ Build", "  ✓ Test"},
		3,
	)

	expected := syncOn +
		upN(2) +
		"\r" + clr + "output 1" + "\n" +
		"\r" + clr + "output 2" + "\n" +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ✓ Build" + "\n" +
		"\r" + clr + "  ✓ Test" +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("scroll+panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 4: Panel grows (3 → 4).
// up(2) to first line, write 4 lines. Extra line scrolls terminal.
func TestExact_PanelGrows(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build", "  ✓ Test", "  ✓ Deploy"}, 3)

	expected := syncOn +
		upN(2) +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ✓ Build" + "\n" +
		"\r" + clr + "  ✓ Test" + "\n" +
		"\r" + clr + "  ✓ Deploy" +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("panel grows:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 5: Panel shrinks (4 → 3).
// up(3) to first line, write 3 lines, clear orphan line, move back.
func TestExact_PanelShrinks(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.DrawPanel([]string{"─── gloss ───", "  ✓ Build", "  ✓ Test"}, 4)

	expected := syncOn +
		upN(3) +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ✓ Build" + "\n" +
		"\r" + clr + "  ✓ Test" + "\n" +
		"\r" + clr +
		upN(1) +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("panel shrinks:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 6: Pure scroll, no panel at all.
func TestExact_ScrollNoPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.WriteScrollWithPanel([]string{"hello", "world"}, nil, 0)

	expected := "hello\nworld\n"

	if got := buf.String(); got != expected {
		t.Errorf("scroll no panel:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 7: Panel appears mid-stream (prevPanelHeight=0, panel lines present).
func TestExact_PanelAppearsMidStream(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.WriteScrollWithPanel(
		[]string{"compiling main.go"},
		[]string{"─── gloss ───", "  ⠋ Build", "  ○ Test"},
		0,
	)

	expected := syncOn +
		"\r" + clr + "compiling main.go" + "\n" +
		"\r" + clr + "─── gloss ───" + "\n" +
		"\r" + clr + "  ⠋ Build" + "\n" +
		"\r" + clr + "  ○ Test" +
		syncOff

	if got := buf.String(); got != expected {
		t.Errorf("panel mid-stream:\nexpected: %q\n     got: %q", expected, got)
	}
}

// Scenario 8: ClearPanel.
// Cursor ON last panel line. up(height-1) to first, clear all, up(height) back.
func TestExact_ClearPanel(t *testing.T) {
	var buf bytes.Buffer
	r := NewRenderer(&buf, 80, false)

	r.ClearPanel(3)

	expected := upN(2) +
		"\r" + clr + "\n" +
		"\r" + clr + "\n" +
		"\r" + clr + "\n" +
		upN(3)

	if got := buf.String(); got != expected {
		t.Errorf("clear panel:\nexpected: %q\n     got: %q", expected, got)
	}
}
