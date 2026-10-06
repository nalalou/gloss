package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/protocol"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format :: protocol lines from stdin",
	Long: `Reads stdin line by line. Lines starting with :: are rendered as
visual elements. All other lines pass through unchanged.

Agent protocol:
  ::ok Text         ✓ success badge
  ::err Text        ✗ error badge
  ::warn Text       ⚠ warning badge
  ::info Text       ℹ info badge
  ::bar 75 Label    progress bar
  ::divider Label   horizontal rule
  ::callout type Text   boxed alert`,
	Example: `  echo "::ok Tests passing" | gloss fmt
  cat <<'EOF' | gloss fmt
  ::divider Results
  All tests completed.
  ::ok 3/4 suites passing
  EOF`,
	Args: cobra.NoArgs,
	RunE: runFmt,
}

func init() {
	rootCmd.AddCommand(fmtCmd)
}

func runFmt(cmd *cobra.Command, args []string) error {
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	return formatStream(width, noColor)
}

// formatStream renders stdin line by line, without cursor tricks.
func formatStream(width int, noColor bool) error {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var masker protocol.Masker
	return protocol.ReadLines(os.Stdin, func(line string) {
		line = masker.Apply(line)
		rendered := protocol.RenderLine(line, width, noColor)
		if rendered == "" && line != "" {
			return // hidden directive, e.g. ::endgroup:: or ::remove
		}
		fmt.Fprintln(out, rendered)
		if !strings.HasPrefix(line, "::") {
			out.Flush() // keep plain output live; batch only directive bursts
		}
	})
}
