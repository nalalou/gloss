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

type Renderer struct {
	out     io.Writer
	width   int
	noColor bool
}

func NewRenderer(out io.Writer, width int, noColor bool) *Renderer {
	return &Renderer{out: out, width: width, noColor: noColor}
}

func (r *Renderer) HideCursor() {
	fmt.Fprint(r.out, hideCursorSeq)
}

func (r *Renderer) ShowCursor() {
	fmt.Fprint(r.out, showCursorSeq)
}

func (r *Renderer) WriteScroll(line string) {
	fmt.Fprintf(r.out, "%s\n", line)
}

func (r *Renderer) WriteScrollWithPanel(scrollLines []string, panelLines []string, prevPanelHeight int) {
	if prevPanelHeight == 0 && len(panelLines) == 0 {
		for _, line := range scrollLines {
			fmt.Fprintf(r.out, "%s\n", line)
		}
		return
	}

	fmt.Fprint(r.out, syncBegin)

	if prevPanelHeight > 0 {
		fmt.Fprintf(r.out, "\033[%dA", prevPanelHeight)
	}

	for _, line := range scrollLines {
		fmt.Fprintf(r.out, "%s%s\n", clearLineSeq, line)
	}

	for i, line := range panelLines {
		fmt.Fprintf(r.out, "%s%s", clearLineSeq, line)
		if i < len(panelLines)-1 {
			fmt.Fprint(r.out, "\n")
		}
	}

	fmt.Fprint(r.out, syncEnd)
}

func (r *Renderer) DrawPanel(lines []string, prevHeight int) {
	if len(lines) == 0 {
		return
	}

	fmt.Fprint(r.out, syncBegin)

	if prevHeight > 0 {
		fmt.Fprintf(r.out, "\033[%dA", prevHeight)
	}

	for i, line := range lines {
		fmt.Fprintf(r.out, "%s%s", clearLineSeq, line)
		if i < len(lines)-1 {
			fmt.Fprint(r.out, "\n")
		}
	}

	fmt.Fprint(r.out, syncEnd)
}

func (r *Renderer) ClearPanel(height int) {
	if height == 0 {
		return
	}
	fmt.Fprintf(r.out, "\033[%dA", height)
	for i := 0; i < height; i++ {
		fmt.Fprintf(r.out, "%s\n", clearLineSeq)
	}
	fmt.Fprintf(r.out, "\033[%dA", height)
}

func (r *Renderer) SetWidth(width int) {
	r.width = width
}
