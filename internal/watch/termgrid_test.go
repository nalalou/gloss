package watch

import (
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// TermGrid simulates a fixed-size terminal screen. Implements io.Writer.
// Interprets the ANSI sequences our Renderer emits.
type TermGrid struct {
	rows int
	cols int
	grid [][]rune
	crow int
	ccol int
}

func NewTermGrid(rows, cols int) *TermGrid {
	g := &TermGrid{rows: rows, cols: cols}
	g.grid = make([][]rune, rows)
	for r := 0; r < rows; r++ {
		g.grid[r] = makeBlankRow(cols)
	}
	return g
}

func makeBlankRow(cols int) []rune {
	row := make([]rune, cols)
	for i := range row {
		row[i] = ' '
	}
	return row
}

func (g *TermGrid) Write(p []byte) (int, error) {
	total := len(p)
	i := 0
	for i < len(p) {
		b := p[i]
		if b == '\r' {
			g.ccol = 0
			i++
			continue
		}
		if b == '\n' {
			if g.crow == g.rows-1 {
				g.scrollUp()
			} else {
				g.crow++
			}
			g.ccol = 0
			i++
			continue
		}
		if b == 0x1b && i+1 < len(p) && p[i+1] == '[' {
			i += 2
			var param []byte
			for i < len(p) && (p[i] >= '0' && p[i] <= '9' || p[i] == '?' || p[i] == ';') {
				param = append(param, p[i])
				i++
			}
			if i < len(p) {
				cmd := p[i]
				i++
				g.dispatchCSI(string(param), cmd)
			}
			continue
		}
		r, size := utf8.DecodeRune(p[i:])
		if r == utf8.RuneError && size <= 1 {
			i++
			continue
		}
		i += size
		if g.ccol < g.cols {
			g.grid[g.crow][g.ccol] = r
			g.ccol++
		}
	}
	return total, nil
}

func (g *TermGrid) dispatchCSI(param string, cmd byte) {
	switch cmd {
	case 'A':
		n := 1
		if param != "" {
			if v, err := strconv.Atoi(param); err == nil {
				n = v
			}
		}
		g.crow -= n
		if g.crow < 0 {
			g.crow = 0
		}
	case 'B': // Cursor Down
		n := 1
		if param != "" {
			if v, err := strconv.Atoi(param); err == nil {
				n = v
			}
		}
		g.crow += n
		if g.crow >= g.rows {
			overshoot := g.crow - (g.rows - 1)
			for i := 0; i < overshoot; i++ {
				g.scrollUp()
			}
			g.crow = g.rows - 1
		}
	case 'G': // Cursor Horizontal Absolute
		col := 0
		if param != "" {
			if v, err := strconv.Atoi(param); err == nil {
				col = v
			}
		}
		if col > 0 {
			col--
		}
		g.ccol = col
	case 'K':
		for c := 0; c < g.cols; c++ {
			g.grid[g.crow][c] = ' '
		}
	case 'h', 'l':
		// ignore private modes
	}
}

func (g *TermGrid) scrollUp() {
	for r := 0; r < g.rows-1; r++ {
		g.grid[r] = g.grid[r+1]
	}
	g.grid[g.rows-1] = makeBlankRow(g.cols)
}

func (g *TermGrid) Row(n int) string {
	if n < 0 || n >= g.rows {
		return ""
	}
	return strings.TrimRight(string(g.grid[n]), " ")
}

func (g *TermGrid) String() string {
	last := -1
	for r := g.rows - 1; r >= 0; r-- {
		if g.Row(r) != "" {
			last = r
			break
		}
	}
	if last < 0 {
		return ""
	}
	var sb strings.Builder
	for r := 0; r <= last; r++ {
		if r > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(g.Row(r))
	}
	return sb.String()
}

func (g *TermGrid) CursorRow() int { return g.crow }
func (g *TermGrid) CursorCol() int { return g.ccol }

// --- Tests using TermGrid ---

func TestTermGrid_DrawPanel(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.DrawPanel([]string{"--- gloss ---", "  ok Build"})
	if !strings.Contains(grid.Row(0), "gloss") {
		t.Errorf("row 0: %q", grid.Row(0))
	}
	if !strings.Contains(grid.Row(1), "Build") {
		t.Errorf("row 1: %q", grid.Row(1))
	}
	if grid.Row(2) != "" {
		t.Errorf("row 2 should be empty: %q", grid.Row(2))
	}
}

func TestTermGrid_PanelRedraw(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.DrawPanel([]string{"--- gloss ---", "  A Build"})
	r.DrawPanel([]string{"--- gloss ---", "  B Build"})

	if strings.Count(grid.String(), "gloss") != 1 {
		t.Errorf("divider should appear once:\n%s", grid.String())
	}
	if !strings.Contains(grid.Row(1), "B Build") {
		t.Errorf("row 1 should have updated content: %q", grid.Row(1))
	}
}

func TestTermGrid_ScrollWithPanel(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.Render(nil, []string{"--- gloss ---", "  spin Build"})
	r.Render([]string{"compiling main.go"}, []string{"--- gloss ---", "  ok Build"})

	if !strings.Contains(grid.Row(0), "compiling main.go") {
		t.Errorf("row 0: %q", grid.Row(0))
	}
	if !strings.Contains(grid.Row(1), "gloss") {
		t.Errorf("row 1: %q", grid.Row(1))
	}
	if !strings.Contains(grid.Row(2), "ok Build") {
		t.Errorf("row 2: %q", grid.Row(2))
	}
	if grid.Row(3) != "" {
		t.Errorf("row 3 should be empty: %q", grid.Row(3))
	}
}

func TestTermGrid_PanelGrowth(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.DrawPanel([]string{"--- gloss ---", "  ok Build"})
	r.DrawPanel([]string{"--- gloss ---", "  ok Build", "  ok Test"})

	if !strings.Contains(grid.Row(2), "Test") {
		t.Errorf("row 2: %q", grid.Row(2))
	}
	if grid.Row(3) != "" {
		t.Errorf("row 3 should be empty: %q", grid.Row(3))
	}
}

func TestTermGrid_MultipleSpinnerTicks(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	for _, frame := range []string{"A", "B", "C", "D"} {
		r.DrawPanel([]string{"--- gloss ---", "  " + frame + " Build"})
	}

	if !strings.Contains(grid.Row(1), "D Build") {
		t.Errorf("row 1 should have last frame: %q", grid.Row(1))
	}
	if strings.Count(grid.String(), "gloss") != 1 {
		t.Errorf("divider should appear once:\n%s", grid.String())
	}
}

func TestTermGrid_PanelShrink(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.DrawPanel([]string{"--- gloss ---", "  ok Build", "  ok Test"})
	r.DrawPanel([]string{"--- gloss ---", "  ok Build"})

	if strings.Contains(grid.Row(2), "Test") {
		t.Errorf("row 2 should be cleared: %q", grid.Row(2))
	}
}

func TestTermGrid_ClearPanel(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.DrawPanel([]string{"--- gloss ---", "  ok Build"})
	r.ClearPanel()

	if grid.Row(0) != "" {
		t.Errorf("row 0 should be empty: %q", grid.Row(0))
	}
	if grid.Row(1) != "" {
		t.Errorf("row 1 should be empty: %q", grid.Row(1))
	}
}

func TestTermGrid_MultipleScrollBatches(t *testing.T) {
	grid := NewTermGrid(24, 80)
	r := NewRenderer(grid, 80, true)
	r.Render(nil, []string{"divider", "  status"})
	r.Render([]string{"line-1"}, []string{"divider", "  status"})
	r.Render([]string{"line-2"}, []string{"divider", "  status"})

	if !strings.Contains(grid.Row(0), "line-1") {
		t.Errorf("row 0: %q", grid.Row(0))
	}
	if !strings.Contains(grid.Row(1), "line-2") {
		t.Errorf("row 1: %q", grid.Row(1))
	}
	if !strings.Contains(grid.Row(2), "divider") {
		t.Errorf("row 2: %q", grid.Row(2))
	}
	if !strings.Contains(grid.Row(3), "status") {
		t.Errorf("row 3: %q", grid.Row(3))
	}
}

func TestTermGrid_RawEscapes(t *testing.T) {
	grid := NewTermGrid(5, 20)
	grid.Write([]byte("hello\n"))
	grid.Write([]byte("world\n"))
	if grid.Row(0) != "hello" {
		t.Errorf("row 0: %q", grid.Row(0))
	}
	if grid.Row(1) != "world" {
		t.Errorf("row 1: %q", grid.Row(1))
	}
	grid.Write([]byte("\033[2A"))
	if grid.CursorRow() != 0 {
		t.Errorf("cursor should be at row 0: %d", grid.CursorRow())
	}
	grid.Write([]byte("\r\033[2K"))
	if grid.Row(0) != "" {
		t.Errorf("row 0 should be cleared: %q", grid.Row(0))
	}
	grid.Write([]byte("replaced"))
	if grid.Row(0) != "replaced" {
		t.Errorf("row 0: %q", grid.Row(0))
	}
}

// Test: Terminal nearly full — panel at bottom, spinner ticks after many scroll lines
func TestTermGrid_PanelAtBottomOfTerminal(t *testing.T) {
	grid := NewTermGrid(10, 40) // small terminal
	r := NewRenderer(grid, 40, true)

	// Push 8 scroll lines + 2-line panel (fills 10-row terminal)
	for i := 0; i < 8; i++ {
		r.Render([]string{"scroll " + strconv.Itoa(i)}, []string{"div", "status"})
	}

	// Now do 5 spinner redraws at the bottom of a full terminal
	for i := 0; i < 5; i++ {
		r.DrawPanel([]string{"div", "frame " + strconv.Itoa(i)})
	}

	screenshot := grid.String()

	// Divider should appear exactly once
	if strings.Count(screenshot, "div") != 1 {
		t.Errorf("divider should appear once, got:\n%s", screenshot)
	}

	// Last spinner frame should be visible
	if !strings.Contains(grid.Row(grid.CursorRow()-1), "frame 4") {
		// Check the row before cursor
		t.Errorf("expected last frame, screenshot:\n%s\ncursor at row %d", screenshot, grid.CursorRow())
	}
}
