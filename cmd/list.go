package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var (
	flagListStyle  string
	flagListStatus bool
)

var listCmd = &cobra.Command{
	Use:   "list <items...>",
	Short: "Render a styled list",
	Long: `Renders a styled list from arguments or newline-delimited stdin. Supports
bullet, arrow, dash, star, check, and numbered styles. Use --status to parse
items as "text:done|pending|fail" for status indicators.`,
	Example: `  gloss list "Build" "Test" "Deploy"
  gloss list --style=numbered "Step 1" "Step 2"
  gloss list "Build:done" "Test:pending" --status`,
	Args: cobra.ArbitraryArgs,
	RunE: runList,
}

func init() {
	listCmd.Flags().StringVar(&flagListStyle, "style", "bullet", "list style: bullet, arrow, dash, star, check, numbered")
	listCmd.Flags().BoolVar(&flagListStatus, "status", false, "status mode: items as 'text:done|pending|fail'")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	var items []string
	if len(args) > 0 {
		items = args
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		text := strings.TrimSpace(t)
		if text != "" {
			items = strings.Split(text, "\n")
		}
	}
	if len(items) == 0 {
		return fmt.Errorf("no items provided; see 'gloss list --help'")
	}

	result := render.RenderList(items, flagListStyle, flagListStatus)
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := strings.Split(result, "\n")
	lines = render.ColorizeLines(lines, gradient, flagColor, noColor)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
