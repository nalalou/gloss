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
// Cursor invariant: after every render call, the cursor is ON the last
// panel line (no trailing newline). If there's no panel, cursor is at
// the start of a new line after the last scroll output.
type Renderer struct {
	out     io.Writer
	width   int
	noColor bool
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

// WriteScrollWithPanel prints scroll content then redraws the panel.
// prevPanelHeight is the number of lines the panel occupied last time.
// Cursor must be ON the last panel line (or at a fresh line if prevPanelHeight=0).
func (r *Renderer) WriteScrollWithPanel(scrollLines []string, panelLines []string, prevPanelHeight int) {
	// No panel at all — just print scroll lines
	if prevPanelHeight == 0 && len(panelLines) == 0 {
		for _, line := range scrollLines {
			fmt.Fprintf(r.out, "%s\n", line)
		}
		return
	}

	fmt.Fprint(r.out, syncBegin)

	// Move cursor to the FIRST line of the old panel.
	// Cursor is ON the last panel line, so up (prevHeight - 1) reaches first.
	if prevPanelHeight > 1 {
		fmt.Fprintf(r.out, "\033[%dA", prevPanelHeight-1)
	}

	// Write scroll lines (overwriting old panel lines, pushing content up).
	for _, line := range scrollLines {
		fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
	}

	// Write panel lines. Last line has NO trailing \n.
	for i, line := range panelLines {
		if i < len(panelLines)-1 {
			fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
		} else {
			fmt.Fprintf(r.out, "\r%s%s", clearLineSeq, line)
		}
	}

	// Clear orphaned lines if panel shrank.
	// After writing, cursor is on the new last panel line.
	// Old panel had more lines below — they still have stale content.
	totalWritten := len(scrollLines) + len(panelLines)
	if totalWritten < prevPanelHeight {
		orphans := prevPanelHeight - totalWritten
		for i := 0; i < orphans; i++ {
			fmt.Fprintf(r.out, "\n\r%s", clearLineSeq)
		}
		// Move back up to the actual last panel line
		fmt.Fprintf(r.out, "\033[%dA", orphans)
	}

	fmt.Fprint(r.out, syncEnd)
}

// DrawPanel redraws the panel in place.
// prevHeight is the number of lines the panel occupied last time.
func (r *Renderer) DrawPanel(lines []string, prevHeight int) {
	if len(lines) == 0 {
		return
	}

	fmt.Fprint(r.out, syncBegin)

	// Move to first panel line. Cursor is ON last line, so up (prev - 1).
	if prevHeight > 1 {
		fmt.Fprintf(r.out, "\033[%dA", prevHeight-1)
	}

	// Write all panel lines. Last line has no trailing \n.
	for i, line := range lines {
		if i < len(lines)-1 {
			fmt.Fprintf(r.out, "\r%s%s\n", clearLineSeq, line)
		} else {
			fmt.Fprintf(r.out, "\r%s%s", clearLineSeq, line)
		}
	}

	// Clear orphaned lines if panel shrank.
	if len(lines) < prevHeight {
		orphans := prevHeight - len(lines)
		for i := 0; i < orphans; i++ {
			fmt.Fprintf(r.out, "\n\r%s", clearLineSeq)
		}
		fmt.Fprintf(r.out, "\033[%dA", orphans)
	}

	fmt.Fprint(r.out, syncEnd)
}

// ClearPanel erases the panel area and repositions cursor to where the panel started.
func (r *Renderer) ClearPanel(height int) {
	if height == 0 {
		return
	}
	// Move to first panel line
	if height > 1 {
		fmt.Fprintf(r.out, "\033[%dA", height-1)
	}
	// Clear each line
	for i := 0; i < height; i++ {
		fmt.Fprintf(r.out, "\r%s\n", clearLineSeq)
	}
	// Move back up to where panel started
	fmt.Fprintf(r.out, "\033[%dA", height)
}

func (r *Renderer) SetWidth(width int) { r.width = width }
