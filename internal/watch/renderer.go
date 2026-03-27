package watch

import (
	"fmt"
	"io"
)

const (
	hideCursorSeq = "\033[?25l"
	showCursorSeq = "\033[?25h"
	clearLineSeq  = "\033[2K"
	syncBegin     = "\033[?2026h"
	syncEnd       = "\033[?2026l"
)

// Renderer handles terminal output with cursor control for the panel.
//
// Model: we track how many physical lines the panel currently occupies.
// To redraw, we move up that many lines, clear them, and write new content.
// The cursor always ends AFTER the last line we wrote (on a new blank line).
type Renderer struct {
	out           io.Writer
	width         int
	noColor       bool
	linesOnScreen int // how many lines we've written that we need to erase on next redraw
}

func NewRenderer(out io.Writer, width int, noColor bool) *Renderer {
	return &Renderer{out: out, width: width, noColor: noColor}
}

func (r *Renderer) HideCursor() { fmt.Fprint(r.out, hideCursorSeq) }
func (r *Renderer) ShowCursor() { fmt.Fprint(r.out, showCursorSeq) }

// WriteScroll prints a line that scrolls normally (no panel).
func (r *Renderer) WriteScroll(line string) {
	fmt.Fprintf(r.out, "%s\n", line)
}

// Render does a full redraw cycle:
// 1. Erase the previous panel (move up, clear each line)
// 2. Print scroll lines (these become permanent scroll output)
// 3. Print new panel lines
//
// After this call, linesOnScreen reflects the new panel height.
func (r *Renderer) Render(scrollLines []string, panelLines []string) {
	// If nothing to do and no previous panel, just print scroll lines
	if r.linesOnScreen == 0 && len(panelLines) == 0 {
		for _, line := range scrollLines {
			fmt.Fprintf(r.out, "%s\n", line)
		}
		return
	}

	fmt.Fprint(r.out, syncBegin)

	// Step 1: Move up to erase the old panel
	if r.linesOnScreen > 0 {
		fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
		for i := 0; i < r.linesOnScreen; i++ {
			fmt.Fprintf(r.out, "\r%s\n", clearLineSeq)
		}
		// Move back to where we started erasing
		fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
	}

	// Step 2: Print scroll lines (permanent, won't be erased next time)
	for _, line := range scrollLines {
		fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
	}

	// Step 3: Print new panel lines
	for _, line := range panelLines {
		fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
	}

	r.linesOnScreen = len(panelLines)

	fmt.Fprint(r.out, syncEnd)
}

// DrawPanel redraws only the panel (no scroll content). Used for spinner ticks.
func (r *Renderer) DrawPanel(panelLines []string) {
	if len(panelLines) == 0 && r.linesOnScreen == 0 {
		return
	}

	fmt.Fprint(r.out, syncBegin)

	// Erase old panel
	if r.linesOnScreen > 0 {
		fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
		for i := 0; i < r.linesOnScreen; i++ {
			fmt.Fprintf(r.out, "\r%s\n", clearLineSeq)
		}
		fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
	}

	// Write new panel
	for _, line := range panelLines {
		fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
	}

	r.linesOnScreen = len(panelLines)

	fmt.Fprint(r.out, syncEnd)
}

// ClearPanel erases the panel. Used at cleanup/exit.
func (r *Renderer) ClearPanel() {
	if r.linesOnScreen == 0 {
		return
	}
	fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
	for i := 0; i < r.linesOnScreen; i++ {
		fmt.Fprintf(r.out, "\r%s\n", clearLineSeq)
	}
	fmt.Fprintf(r.out, "\033[%dA", r.linesOnScreen)
	r.linesOnScreen = 0
}

func (r *Renderer) SetWidth(width int)  { r.width = width }
func (r *Renderer) LinesOnScreen() int  { return r.linesOnScreen }
