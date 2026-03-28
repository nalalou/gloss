package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var spinCmd = &cobra.Command{
	Use:   "spin <message> [-- command args...]",
	Short: "Show a spinner while a command runs",
	Long: `Displays an animated spinner with a label while an external command runs.
If no command is given after --, runs a 2-second demo. The completion badge
and label text respect --color and --gradient flags.`,
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

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	if noColor {
		frames = []string{".", "..", "...", "....", "....."}
	}

	// Colorize the label text
	colorizeLabel := func(s string) string {
		if noColor {
			return s
		}
		gradient := resolveGradientFlag()
		color := flagColor
		lines := render.ColorizeLines([]string{s}, gradient, color, noColor)
		return lines[0]
	}

	coloredMessage := colorizeLabel(message)

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
				badge := colorizeLabel("✗ " + message)
				fmt.Fprintf(os.Stderr, "%s\n", badge)
				return fmt.Errorf("command failed: %w", err)
			}
			badge := colorizeLabel("✓ " + message)
			fmt.Fprintf(os.Stderr, "%s\n", badge)
			return nil
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "\r\033[2K")
			cancelled := colorizeLabel("✗ " + message + " (cancelled)")
			fmt.Fprintf(os.Stderr, "%s\n", cancelled)
			return nil
		case <-ticker.C:
			frame := frames[frameIdx%len(frames)]
			fmt.Fprintf(os.Stderr, "\r\033[2K  %s %s", frame, coloredMessage)
			frameIdx++
		}
	}
}
