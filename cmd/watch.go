package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/protocol"
	"github.com/nalalou/gloss/internal/watch"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Live-updating panel for :: protocol streams",
	Long: `Reads stdin and renders a persistent status panel at the bottom
of the terminal. Directives with id= update in place in the panel.
Everything else scrolls normally above.`,
	Example: `  my-agent | gloss watch
  ./deploy.sh | gloss watch`,
	Args: cobra.NoArgs,
	RunE: runWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}

func runWatch(cmd *cobra.Command, args []string) error {
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	lines, readErr := readInput(os.Stdin)
	if envInfo.IsTTY {
		livePanel(lines, noColor, true)
	} else {
		formatLines(lines, plainWidth(), noColor)
	}

	select {
	case err := <-readErr:
		return err
	default: // interrupted before input ended
		return nil
	}
}

// readInput streams r's lines, with ::add-mask:: secrets hidden, on the returned
// channel. The error channel receives ReadLines' result just before lines closes.
func readInput(r io.Reader) (<-chan string, <-chan error) {
	lines := make(chan string, 256)
	readErr := make(chan error, 1)
	go func() {
		var masker protocol.Masker
		readErr <- protocol.ReadLines(r, func(line string) {
			lines <- masker.Apply(line)
		})
		close(lines)
	}()
	return lines, readErr
}

// plainWidth is the width for non-TTY output: --width, else 80.
func plainWidth() int {
	if flagWidth > 0 {
		return flagWidth
	}
	return 80
}

// formatLines renders each line as styled text, without cursor tricks.
func formatLines(lines <-chan string, width int, noColor bool) {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for line := range lines {
		rendered := protocol.RenderLine(line, width, noColor)
		if rendered == "" && line != "" {
			continue // hidden directive, e.g. ::endgroup:: or ::remove
		}
		fmt.Fprintln(out, rendered)
		if len(lines) == 0 {
			out.Flush() // nothing queued, so show what we have
		}
	}
}

// livePanel scrolls plain lines and keeps id= directives in a panel at the
// bottom of the terminal until lines closes. If stopOnInterrupt is true,
// Ctrl+C ends it early; otherwise the caller handles interrupts.
func livePanel(lines <-chan string, noColor bool, stopOnInterrupt bool) {
	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	if stopOnInterrupt {
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	}
	defer signal.Stop(sigCh)

	panel := watch.NewPanel(width)
	renderer := watch.NewRenderer(os.Stdout, width, noColor)

	renderer.HideCursor()
	defer renderer.ShowCursor()

	spinnerTicker := time.NewTicker(80 * time.Millisecond)
	defer spinnerTicker.Stop()
	spinnerFrame := 0

	for {
		select {
		case line, ok := <-lines:
			if !ok {
				goto cleanup
			}

			batch := []string{line}
		drain:
			for {
				select {
				case l, ok := <-lines:
					if !ok {
						break drain
					}
					batch = append(batch, l)
				default:
					break drain
				}
			}

			var scrollLines []string
			for _, bline := range batch {
				dir, id, dargs := protocol.ParseDirective(bline)

				if dir == "remove" && id != "" {
					panel.Remove(id)
				} else if dir != "" && id != "" {
					panel.Set(id, dir, dargs, noColor)
				} else {
					rendered := protocol.RenderLine(bline, width, noColor)
					if rendered == "" && bline != "" {
						continue // hidden directive, e.g. ::endgroup::
					}
					scrollLines = append(scrollLines, strings.Split(rendered, "\n")...)
				}
			}

			renderer.Render(scrollLines, panel.RenderLines())

		case <-spinnerTicker.C:
			if panel.HasRunning() {
				spinnerFrame++
				panel.UpdateSpinnerFrame(spinnerFrame, noColor)
				renderer.DrawPanel(panel.RenderLines())
			}

		case sig := <-sigCh:
			switch sig {
			case syscall.SIGWINCH:
				w, _, err := term.GetSize(int(os.Stdout.Fd()))
				if err == nil && w > 0 {
					width = w
					panel.SetWidth(width)
					renderer.SetWidth(width)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				goto cleanup
			}
		}
	}

cleanup:
	renderer.ClearPanel()
	for _, line := range panel.RenderLines() {
		fmt.Println(line)
	}
}
