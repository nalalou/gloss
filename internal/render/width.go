package render

import (
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// displayWidth returns the number of terminal columns needed to display s,
// after stripping any ANSI escape sequences.
func displayWidth(s string) int {
	return runewidth.StringWidth(stripANSI(s))
}

// termWidth returns the current terminal width, defaulting to 80.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// wrapText wraps text so that no line exceeds maxWidth display columns.
// It splits on whitespace boundaries. Lines already shorter than maxWidth
// are left unchanged.
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	var result []string
	for _, line := range strings.Split(text, "\n") {
		if displayWidth(line) <= maxWidth {
			result = append(result, line)
			continue
		}
		words := strings.Fields(line)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}
		current := words[0]
		for _, word := range words[1:] {
			trial := current + " " + word
			if displayWidth(trial) <= maxWidth {
				current = trial
			} else {
				result = append(result, current)
				current = word
			}
		}
		result = append(result, current)
	}
	return strings.Join(result, "\n")
}
