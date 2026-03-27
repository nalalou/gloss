package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	hideCursor = "\033[?25l"
	showCursor = "\033[?25h"
)

// animate prints text with a typewriter effect at 30ms per character.
// Restores cursor on SIGINT/SIGTERM.
func animate(text string) error {
	fmt.Fprint(os.Stdout, hideCursor)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})

	go func() {
		select {
		case <-sigs:
			fmt.Fprint(os.Stdout, showCursor)
			os.Exit(0)
		case <-done:
		}
	}()

	defer func() {
		close(done)
		signal.Stop(sigs)
		fmt.Fprint(os.Stdout, showCursor)
	}()

	for _, ch := range text {
		fmt.Fprint(os.Stdout, string(ch))
		if ch != '\n' {
			time.Sleep(30 * time.Millisecond)
		}
	}
	fmt.Println()
	return nil
}
