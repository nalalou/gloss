# `gloss chart` Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `gloss chart` command that renders tall vertical bar charts in the terminal with gradient coloring and optional TTY animation.

**Architecture:** Pure render function in `internal/render/chart.go` returns lines of text. Cobra command in `cmd/chart.go` handles input parsing, colorization, and TTY animation. Follows the exact same pattern as `spark.go` / `cmd/spark.go`.

**Tech Stack:** Go, cobra, existing render/ColorizeLines pipeline, ANSI cursor movement for animation.

**Spec:** `docs/superpowers/specs/2026-03-28-chart-command-design.md`

---

### Task 1: Render function and tests

**Files:**
- Create: `internal/render/chart.go`
- Create: `internal/render/chart_test.go`

- [ ] **Step 1: Write failing tests**

```go
package render

import (
	"strings"
	"testing"
)

func TestRenderChartBasic(t *testing.T) {
	lines := RenderChart([]float64{10, 5, 8, 3}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	// Max value (10) should have full block in top row
	if !strings.Contains(lines[0], "█") {
		t.Error("top row should contain █ for max value column")
	}
}

func TestRenderChartSingleValue(t *testing.T) {
	lines := RenderChart([]float64{42}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	// Single value should be full height
	for _, line := range lines {
		if !strings.Contains(line, "█") {
			t.Error("single value should fill all rows")
		}
	}
}

func TestRenderChartAllSame(t *testing.T) {
	lines := RenderChart([]float64{5, 5, 5}, 4)
	// All bars should be the same — all full height
	for _, line := range lines {
		if !strings.Contains(line, "█") {
			t.Error("all-same values should fill all rows")
		}
	}
}

func TestRenderChartAllZeros(t *testing.T) {
	lines := RenderChart([]float64{0, 0, 0}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	// Bottom row should have half-blocks
	if !strings.Contains(lines[3], "▄") {
		t.Error("all-zeros should have ▄ at bottom row")
	}
}

func TestRenderChartEmpty(t *testing.T) {
	lines := RenderChart([]float64{}, 4)
	if len(lines) != 0 {
		t.Errorf("empty input should return empty slice, got %d lines", len(lines))
	}
}

func TestRenderChartNegativesClamped(t *testing.T) {
	lines := RenderChart([]float64{-5, 10, -3}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	// Should not panic, negative values treated as 0
}

func TestRenderChartHeight(t *testing.T) {
	lines := RenderChart([]float64{10, 5}, 8)
	if len(lines) != 8 {
		t.Errorf("expected 8 lines for height=8, got %d", len(lines))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/render/ -run TestRenderChart -v`
Expected: FAIL — `RenderChart` not defined

- [ ] **Step 3: Write the render function**

```go
package render

import "math"

// RenderChart renders a vertical bar chart as a slice of strings (top to bottom).
// Each bar is 2 chars wide with 1 char gap. Uses █ for full cells, ▄ for half.
// 2-space left indent on each line. No color — caller applies ColorizeLines.
func RenderChart(values []float64, height int) []string {
	if len(values) == 0 {
		return nil
	}

	// Clamp negatives to 0
	clamped := make([]float64, len(values))
	for i, v := range values {
		if v < 0 {
			v = 0
		}
		clamped[i] = v
	}

	// Find max
	max := 0.0
	for _, v := range clamped {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	barWidth := 2
	gap := 1
	colWidth := barWidth + gap
	totalWidth := len(clamped)*colWidth - gap

	// Build grid
	grid := make([][]rune, height)
	for r := range grid {
		grid[r] = make([]rune, totalWidth)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	for col, v := range clamped {
		barH := v / max * float64(height)
		fullCells := int(barH)
		hasHalf := (barH - float64(fullCells)) >= 0.5

		// All zeros: show a half-block at bottom
		if v == 0 && max > 0 {
			hasHalf = true
			fullCells = 0
		}

		x := col * colWidth
		for r := 0; r < height; r++ {
			rowFromBottom := height - 1 - r
			for b := 0; b < barWidth; b++ {
				if x+b < totalWidth {
					if rowFromBottom < fullCells {
						grid[r][x+b] = '█'
					} else if rowFromBottom == fullCells && hasHalf {
						grid[r][x+b] = '▄'
					}
				}
			}
		}
	}

	// Convert grid to strings with indent
	lines := make([]string, height)
	for r := 0; r < height; r++ {
		lines[r] = "  " + string(grid[r])
	}
	return lines
}

// ChartColumnCount returns how many columns a chart with n values would have.
// Useful for animation: iterate 0..n-1 and render partial charts.
func ChartColumnCount(n int) int {
	if n == 0 {
		return 0
	}
	return n*3 - 1 // 2 chars per bar + 1 gap, minus trailing gap
}

// RenderChartPartial renders only the first numCols values, useful for animation.
func RenderChartPartial(values []float64, height int, numCols int) []string {
	if numCols <= 0 || numCols > len(values) {
		numCols = len(values)
	}
	// Render the subset but use the FULL dataset's max for consistent scaling
	max := 0.0
	for _, v := range values {
		if v > 0 && v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}

	subset := values[:numCols]
	barWidth := 2
	gap := 1
	colWidth := barWidth + gap
	// Total width is based on FULL dataset so the chart doesn't shift during animation
	fullWidth := len(values)*colWidth - gap

	grid := make([][]rune, height)
	for r := range grid {
		grid[r] = make([]rune, fullWidth)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	for col, v := range subset {
		if v < 0 {
			v = 0
		}
		barH := v / max * float64(height)
		fullCells := int(barH)
		hasHalf := (barH - math.Floor(barH)) >= 0.5

		if v == 0 {
			hasHalf = true
			fullCells = 0
		}

		x := col * colWidth
		for r := 0; r < height; r++ {
			rowFromBottom := height - 1 - r
			for b := 0; b < barWidth; b++ {
				if x+b < fullWidth {
					if rowFromBottom < fullCells {
						grid[r][x+b] = '█'
					} else if rowFromBottom == fullCells && hasHalf {
						grid[r][x+b] = '▄'
					}
				}
			}
		}
	}

	lines := make([]string, height)
	for r := 0; r < height; r++ {
		lines[r] = "  " + string(grid[r])
	}
	return lines
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/render/ -run TestRenderChart -v`
Expected: all 7 tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/render/chart.go internal/render/chart_test.go
git commit -m "feat: add RenderChart function for tall bar charts"
```

---

### Task 2: Cobra command

**Files:**
- Create: `cmd/chart.go`

- [ ] **Step 1: Write the command**

```go
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagChartHeight int
	flagChartLabel  string
)

var chartCmd = &cobra.Command{
	Use:   "chart <values>",
	Short: "Render a vertical bar chart",
	Long: `Renders a tall vertical bar chart from numeric values. Values can be
comma-separated (220,195,180) or space-separated. Pipe via stdin or pass as args.
Each value becomes a column scaled relative to the maximum.`,
	Example: `  gloss chart 220,195,180,140,95,72,45
  gloss chart 10 5 8 3 --height=8
  gloss chart 220,195,180 --label="p99 (ms)" --gradient=fire
  echo "10,20,30,40,50" | gloss chart`,
	Args: cobra.ArbitraryArgs,
	RunE: runChart,
}

func init() {
	chartCmd.Flags().IntVar(&flagChartHeight, "height", 6, "chart height in rows")
	chartCmd.Flags().StringVar(&flagChartLabel, "label", "", "label above the chart")
	rootCmd.AddCommand(chartCmd)
}

func runChart(cmd *cobra.Command, args []string) error {
	// Parse input — reuse spark's value parsing
	input := ""
	if len(args) > 0 {
		input = strings.Join(args, " ")
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		input = strings.TrimSpace(t)
	}
	if input == "" {
		return fmt.Errorf("no data provided; see 'gloss chart --help'")
	}

	values, err := render.ParseSparkValues(input)
	if err != nil {
		return err
	}

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()

	height := flagChartHeight
	if height < 1 {
		height = 1
	}
	if height > 30 {
		height = 30
	}

	// Print label if set
	if flagChartLabel != "" {
		fmt.Printf("  \033[2m%s\033[0m\n", flagChartLabel)
	}

	// Animate if TTY, static if piped
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if isTTY && !noColor {
		animateChart(values, height, gradient, flagColor, noColor)
	} else {
		lines := render.RenderChart(values, height)
		lines = render.ColorizeLines(lines, gradient, flagColor, noColor)
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	return nil
}

func animateChart(values []float64, height int, gradient []string, color string, noColor bool) {
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	// Print empty lines to reserve space
	for i := 0; i < height; i++ {
		fmt.Println()
	}

	for col := 1; col <= len(values); col++ {
		lines := render.RenderChartPartial(values, height, col)
		lines = render.ColorizeLines(lines, gradient, color, noColor)

		// Move cursor up and redraw
		fmt.Printf("\033[%dA", height)
		for _, line := range lines {
			fmt.Printf("\033[2K%s\n", line)
		}

		delay := 120 + (col * 10)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
```

- [ ] **Step 2: Build and verify**

Run: `go build -o gloss . && echo "OK"`
Expected: OK, no errors

- [ ] **Step 3: Manual test — static**

Run: `echo "220,195,180,140,95,72,45" | ./gloss chart --height=6`
Expected: 6-row bar chart printed, no animation

- [ ] **Step 4: Manual test — animated**

Run: `./gloss chart 220,195,180,140,95,72,45 --height=6 --gradient=fire --label="p99 (ms)"`
Expected: dim label, then bars appear one at a time in fire gradient

- [ ] **Step 5: Manual test — edge cases**

Run: `./gloss chart 42` (single value)
Run: `./gloss chart 5,5,5,5` (all same)
Run: `./gloss chart` (no input — should error)
Expected: single bar full height, all bars full height, error message

- [ ] **Step 6: Commit**

```bash
git add cmd/chart.go
git commit -m "feat: add gloss chart command with TTY animation"
```

---

### Task 3: Smoke test

**Files:**
- Modify: `cmd/smoke_test.go`

- [ ] **Step 1: Add smoke tests**

Add these tests to `cmd/smoke_test.go`:

```go
func TestE2EChartBasic(t *testing.T) {
	out := captureOutput(t, func() {
		chartCmd.SetArgs([]string{"10,5,8,3"})
		flagChartHeight = 4
		flagChartLabel = ""
		err := chartCmd.Execute()
		if err != nil {
			t.Fatal(err)
		}
		flagChartHeight = 6 // reset default
	})
	if len(out) == 0 {
		t.Error("expected chart output")
	}
	if !strings.Contains(out, "█") {
		t.Error("expected block chars in output")
	}
}

func TestE2EChartEmpty(t *testing.T) {
	chartCmd.SetArgs([]string{})
	// Simulate no stdin
	err := runChart(chartCmd, []string{})
	if err == nil {
		t.Error("expected error for empty input")
	}
}
```

- [ ] **Step 2: Run all tests**

Run: `go test ./cmd/... -v`
Expected: all tests pass including new chart tests

- [ ] **Step 3: Run full test suite**

Run: `go test ./...`
Expected: all packages pass

- [ ] **Step 4: Commit**

```bash
git add cmd/smoke_test.go
git commit -m "test: add smoke tests for gloss chart"
```

---

### Task 4: Update demo to use render function

**Files:**
- Modify: `cmd/demo.go`

- [ ] **Step 1: Replace `animateTallBars` calls with chart render + animation**

In `cmd/demo.go`, replace the two `animateTallBars` calls with the new `render.RenderChartPartial` approach. The `animateTallBars` function in demo.go can be removed since `cmd/chart.go` now has `animateChart`.

Update the demo's Act 3 to call:

```go
// p99 latency dropping
fmt.Printf("  \033[2mp99 (ms)\033[0m\n")
animateChart(
    []float64{220, 195, 180, 140, 95, 72, 45, 38, 22, 18},
    6, []string{"#FF4444", "#44FF88"}, "", false,
)
time.Sleep(400 * time.Millisecond)

// rps climbing
fmt.Printf("  \033[2mrps\033[0m\n")
animateChart(
    []float64{12, 18, 31, 45, 64, 87, 93, 120, 142, 158},
    6, []string{"#006994", "#00CED1"}, "", false,
)
```

- [ ] **Step 2: Remove `animateTallBars` and its helpers (`lerpHex`, `hexComponents`, `hexByte`, `hexToRGB`) from demo.go**

These are no longer needed since `animateChart` in chart.go uses `render.ColorizeLines` instead of manual hex interpolation.

- [ ] **Step 3: Build and test**

Run: `go build -o gloss . && ./gloss demo`
Expected: demo still works, charts animate the same way

Run: `go test ./...`
Expected: all tests pass

- [ ] **Step 4: Commit**

```bash
git add cmd/demo.go cmd/chart.go
git commit -m "refactor: demo uses chart render function instead of inline animation"
```
