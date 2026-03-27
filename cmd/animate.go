package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nalalou/gloss/internal/render"
)

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

// animate prints text with a typewriter effect at 30ms per visible character.
// Skips delay for ANSI escape sequences. Restores cursor on SIGINT/SIGTERM.
func animate(text string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Fprint(os.Stdout, hideCursor)
	defer fmt.Fprint(os.Stdout, showCursor)

	runes := []rune(text)
	i := 0
	for i < len(runes) {
		// Check for cancellation
		select {
		case <-ctx.Done():
			fmt.Println()
			return nil
		default:
		}

		ch := runes[i]

		// If this is the start of an ANSI escape, print the whole sequence without delay
		if ch == '\033' && i+1 < len(runes) && runes[i+1] == '[' {
			// Find the end of the CSI sequence (letter terminates it)
			j := i + 2
			for j < len(runes) && !isCSITerminator(runes[j]) {
				j++
			}
			if j < len(runes) {
				j++ // include the terminator
			}
			// Print the entire escape sequence at once
			for _, r := range runes[i:j] {
				fmt.Fprint(os.Stdout, string(r))
			}
			i = j
			continue
		}

		fmt.Fprint(os.Stdout, string(ch))
		i++

		// Only delay on visible, non-newline characters
		if ch != '\n' && ch != '\r' {
			time.Sleep(30 * time.Millisecond)
		}
	}
	fmt.Println()
	return nil
}

// animateGradient cycles gradient colors across the rendered text.
// Runs for 90 frames (~3s at 33ms/frame), then shows final static frame.
func animateGradient(fontLines []string, gradientHexes []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Fprint(os.Stdout, hideCursor)
	defer fmt.Fprint(os.Stdout, showCursor)

	const (
		frames  = 90
		frameMs = 33
	)

	numLines := len(fontLines)

	for frame := 0; frame < frames; frame++ {
		select {
		case <-ctx.Done():
			fmt.Println()
			return nil
		default:
		}

		offset := float64(frame) / float64(frames)
		colored := render.RenderGradientFrame(fontLines, gradientHexes, offset)

		if frame > 0 {
			fmt.Fprintf(os.Stdout, "\033[%dA", numLines)
		}
		for _, line := range colored {
			fmt.Fprintf(os.Stdout, "\033[2K%s\n", line)
		}

		time.Sleep(frameMs * time.Millisecond)
	}

	// Final static frame
	fmt.Fprintf(os.Stdout, "\033[%dA", numLines)
	colored := render.RenderGradientFrame(fontLines, gradientHexes, 0)
	for _, line := range colored {
		fmt.Fprintf(os.Stdout, "\033[2K%s\n", line)
	}

	return nil
}

func isCSITerminator(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
