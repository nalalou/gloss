package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/watch"
)

func main() {
	// Simulate what happens in the watch loop

	panel := watch.NewPanel(88) // terminal width 88

	// Step 1: Set a status element (like the watch loop does on first ::status line)
	panel.Set("build", "status", "running Building", false)

	fmt.Println("=== After Set (initial render via renderElement) ===")
	lines := panel.RenderLines()
	fmt.Printf("Line count: %d\n", len(lines))
	for i, line := range lines {
		fmt.Printf("Line %d: %q\n", i, line)
		fmt.Printf("  Hex: %s\n", hex.EncodeToString([]byte(line)))
		fmt.Printf("  Bytes: %d, Contains \\n: %v\n", len(line), strings.Contains(line, "\n"))
	}

	// Step 2: Update spinner frame (like the spinner ticker does)
	panel.UpdateSpinnerFrame(1, false)

	fmt.Println("\n=== After UpdateSpinnerFrame ===")
	lines2 := panel.RenderLines()
	fmt.Printf("Line count: %d\n", len(lines2))
	for i, line := range lines2 {
		fmt.Printf("Line %d: %q\n", i, line)
		fmt.Printf("  Hex: %s\n", hex.EncodeToString([]byte(line)))
		fmt.Printf("  Bytes: %d, Contains \\n: %v\n", len(line), strings.Contains(line, "\n"))
	}

	// Check if line count changed!
	if len(lines) != len(lines2) {
		fmt.Printf("\n*** LINE COUNT CHANGED: %d -> %d ***\n", len(lines), len(lines2))
		fmt.Println("THIS IS THE BUG — linesOnScreen will be wrong!")
	} else {
		fmt.Println("\nLine count is stable.")
	}

	// Step 3: Add more elements like the demo script does
	panel.Set("lint", "status", "pending Lint", false)
	panel.Set("test", "status", "pending Test", false)
	panel.Set("deploy", "status", "pending Deploy", false)
	panel.Set("prog", "bar", "0 Progress", false)

	fmt.Println("\n=== Full panel with 4 elements + bar ===")
	lines3 := panel.RenderLines()
	fmt.Printf("Line count: %d\n", len(lines3))
	for i, line := range lines3 {
		fmt.Printf("Line %d (%d bytes): %q\n", i, len(line), line)
	}
}
