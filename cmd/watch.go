package cmd

import (
	"bufio"
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
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	spinnerTicker := time.NewTicker(80 * time.Millisecond)
	defer spinnerTicker.Stop()
	spinnerFrame := 0
	prevPanelHeight := 0

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
			panelChanged := false

			for _, bline := range batch {
				dir, id, dargs := protocol.ParseDirective(bline)

				if dir == "remove" && id != "" {
					panel.Remove(id)
					panelChanged = true
				} else if dir != "" && id != "" {
					panel.Set(id, dir, dargs, noColor)
					panelChanged = true
				} else {
					rendered := protocol.RenderLine(bline, width, noColor)
					// Split multi-line rendered content (callouts, tables)
					// into separate scroll lines for correct cursor tracking
					for _, subline := range strings.Split(rendered, "\n") {
						scrollLines = append(scrollLines, subline)
					}
				}
			}

			panelLines := panel.RenderLines()
			newPanelHeight := len(panelLines)

			if len(scrollLines) > 0 || panelChanged {
				renderer.WriteScrollWithPanel(scrollLines, panelLines, prevPanelHeight)
			}
			prevPanelHeight = newPanelHeight

		case <-spinnerTicker.C:
			if panel.HasRunning() {
				spinnerFrame++
				panel.UpdateSpinnerFrame(spinnerFrame, noColor)
				panelLines := panel.RenderLines()
				renderer.DrawPanel(panelLines, prevPanelHeight)
				prevPanelHeight = len(panelLines)
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
	if prevPanelHeight > 0 {
		renderer.ClearPanel(prevPanelHeight)
	}
	summaryLines := panel.RenderLines()
	for _, line := range summaryLines {
		fmt.Println(line)
	}
	return nil
}

func runWatchStateless(noColor bool) error {
	width := 80
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		dir, _, args := protocol.ParseDirective(line)
		if dir != "" {
			rendered := protocol.RenderLine("::"+dir+" "+args, width, noColor)
			fmt.Println(rendered)
		} else {
			fmt.Println(line)
		}
	}
	return scanner.Err()
}
