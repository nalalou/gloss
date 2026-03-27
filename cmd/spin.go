package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var spinCmd = &cobra.Command{
	Use:   "spin <message> [-- command args...]",
	Short: "Show a spinner while a command runs",
	Example: `  gloss spin "Installing..." -- npm install
  gloss spin "Building..." -- make build
  gloss spin "Loading..."`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE:               runSpin,
}

func init() {
	rootCmd.AddCommand(spinCmd)
}

func runSpin(cmd *cobra.Command, args []string) error {
	// Parse: everything before "--" is our args, everything after is the command
	message := args[0]
	var cmdArgs []string

	// Look for "--" in os.Args to find the external command
	osArgs := os.Args
	dashIdx := -1
	for i, a := range osArgs {
		if a == "--" {
			dashIdx = i
			break
		}
	}
	if dashIdx != -1 && dashIdx+1 < len(osArgs) {
		cmdArgs = osArgs[dashIdx+1:]
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	if flagNoColor {
		frames = []string{".", "..", "...", "....", "....."}
	}

	// Hide cursor
	fmt.Fprint(os.Stderr, "\033[?25l")
	defer fmt.Fprint(os.Stderr, "\033[?25h")

	done := make(chan error, 1)

	if len(cmdArgs) > 0 {
		// Run the command in background
		c := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		go func() {
			done <- c.Run()
		}()
	} else {
		// Demo mode: spin for 2 seconds
		go func() {
			time.Sleep(2 * time.Second)
			done <- nil
		}()
	}

	// Animate spinner
	frameIdx := 0
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			// Clear spinner line
			fmt.Fprintf(os.Stderr, "\r\033[2K")
			if err != nil {
				fmt.Fprintf(os.Stderr, "✗ %s\n", message)
				return fmt.Errorf("command failed: %w", err)
			}
			fmt.Fprintf(os.Stderr, "✓ %s\n", message)
			return nil
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "\r\033[2K")
			fmt.Fprintf(os.Stderr, "✗ %s (cancelled)\n", message)
			return nil
		case <-ticker.C:
			frame := frames[frameIdx%len(frames)]
			fmt.Fprintf(os.Stderr, "\r\033[2K  %s %s", frame, message)
			frameIdx++
		}
	}
}
