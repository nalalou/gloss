package cmd

import (
	"fmt"
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

	if !envInfo.IsTTY {
		return runWatchStateless(noColor)
	}

	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGWINCH)
	defer signal.Stop(sigCh)

	panel := watch.NewPanel(width)
	renderer := watch.NewRenderer(os.Stdout, width, noColor)

	renderer.HideCursor()
	defer renderer.ShowCursor()

	lines := make(chan string, 256)
	readErr := make(chan error, 1)
	go func() {
		var masker protocol.Masker
		readErr <- protocol.ReadLines(os.Stdin, func(line string) {
			lines <- masker.Apply(line)
		})
		close(lines)
	}()

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
					for _, subline := range strings.Split(rendered, "\n") {
						scrollLines = append(scrollLines, subline)
					}
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
	summaryLines := panel.RenderLines()
	for _, line := range summaryLines {
		fmt.Println(line)
	}
	select {
	case err := <-readErr:
		return err
	default: // interrupted before input ended
		return nil
	}
}

func runWatchStateless(noColor bool) error {
	width := flagWidth
	if width == 0 {
		width = 80
	}
	return formatStream(width, noColor)
}
