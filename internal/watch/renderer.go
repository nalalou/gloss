package watch

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	escCursorUp1  = "\033[1A"
	escCursorDn1  = "\033[1B"
	escCursorCol0 = "\033[0G"
	escHideCursor = "\033[?25l"
	escShowCursor = "\033[?25h"
)

// Renderer uses the Docker BuildKit pattern with single-syscall buffered writes.
type Renderer struct {
	out           io.Writer
	width         int
	noColor       bool
	linesOnScreen int
	repeated      bool
}

func NewRenderer(out io.Writer, width int, noColor bool) *Renderer {
	return &Renderer{out: out, width: width, noColor: noColor}
}

func (r *Renderer) HideCursor() { fmt.Fprint(r.out, escHideCursor) }
func (r *Renderer) ShowCursor() { fmt.Fprint(r.out, escShowCursor) }

func (r *Renderer) WriteScroll(line string) {
	fmt.Fprintf(r.out, "%s\n", line)
}

func (r *Renderer) Render(scrollLines []string, panelLines []string) {
	if r.linesOnScreen == 0 && len(panelLines) == 0 && !r.repeated {
		for _, line := range scrollLines {
			fmt.Fprintf(r.out, "%s\n", line)
		}
		return
	}

	var buf bytes.Buffer

	buf.WriteString(escHideCursor)
	writeUpN(&buf, r.linesOnScreen)

	if !r.repeated && len(panelLines) > 0 {
		buf.WriteString(escCursorDn1)
		writeUpN(&buf, 1)
	}

	for _, line := range scrollLines {
		writePadded(&buf, line, r.width)
	}
	for _, line := range panelLines {
		writePadded(&buf, line, r.width)
	}

	written := len(scrollLines) + len(panelLines)
	if leftover := r.linesOnScreen - written; leftover > 0 {
		for i := 0; i < leftover; i++ {
			writePadded(&buf, "", r.width)
		}
		writeUpN(&buf, leftover)
	}

	r.linesOnScreen = len(panelLines)
	if len(panelLines) > 0 {
		r.repeated = true
	}
	buf.WriteString(escShowCursor)

	r.out.Write(buf.Bytes()) // single syscall
}

func (r *Renderer) DrawPanel(panelLines []string) {
	if len(panelLines) == 0 && r.linesOnScreen == 0 {
		return
	}

	var buf bytes.Buffer

	buf.WriteString(escHideCursor)
	writeUpN(&buf, r.linesOnScreen)

	if !r.repeated && len(panelLines) > 0 {
		buf.WriteString(escCursorDn1)
		writeUpN(&buf, 1)
	}

	for _, line := range panelLines {
		writePadded(&buf, line, r.width)
	}

	if leftover := r.linesOnScreen - len(panelLines); leftover > 0 {
		for i := 0; i < leftover; i++ {
			writePadded(&buf, "", r.width)
		}
		writeUpN(&buf, leftover)
	}

	r.linesOnScreen = len(panelLines)
	if len(panelLines) > 0 {
		r.repeated = true
	}
	buf.WriteString(escShowCursor)

	r.out.Write(buf.Bytes()) // single syscall
}

func (r *Renderer) ClearPanel() {
	if r.linesOnScreen == 0 {
		return
	}
	var buf bytes.Buffer
	buf.WriteString(escHideCursor)
	writeUpN(&buf, r.linesOnScreen)
	for i := 0; i < r.linesOnScreen; i++ {
		writePadded(&buf, "", r.width)
	}
	writeUpN(&buf, r.linesOnScreen)
	r.linesOnScreen = 0
	r.repeated = false
	buf.WriteString(escShowCursor)
	r.out.Write(buf.Bytes())
}

func (r *Renderer) SetWidth(width int) { r.width = width }
func (r *Renderer) LinesOnScreen() int { return r.linesOnScreen }

func writeUpN(buf *bytes.Buffer, n int) {
	for i := 0; i < n; i++ {
		buf.WriteString(escCursorUp1)
	}
}

func writePadded(buf *bytes.Buffer, line string, width int) {
	vis := visibleLen(line)
	maxCols := width - 1 // never fill last column — avoids Terminal.app pending-wrap double advance
	pad := maxCols - vis
	if pad < 0 {
		pad = 0
	}
	buf.WriteString(escCursorCol0)
	buf.WriteString(line)
	buf.WriteString(strings.Repeat(" ", pad))
	buf.WriteByte('\n')
}

func visibleLen(s string) int {
	n := 0
	i := 0
	for i < len(s) {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && s[i] >= 0x20 && s[i] <= 0x3F {
				i++
			}
			if i < len(s) {
				i++
			}
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		n++
	}
	return n
}
