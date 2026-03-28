package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagBarStyle     string
	flagBarLabel     string
	flagBarNoPercent bool
)

var barCmd = &cobra.Command{
	Use:   "bar <value>",
	Short: "Render a progress bar",
	Long: `Renders a terminal progress bar. Accepts a percentage (75), a decimal (0.75),
or a fraction (3/4). Use --label to add a prefix and --style to pick a visual
style (blocks, dots, ascii, thin).`,
	Example: `  gloss bar 75
  gloss bar 0.75
  gloss bar 3/4
  gloss bar 75 --style=dots --label="Tests"
  echo "0.75" | gloss bar`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBar,
}

func init() {
	barCmd.Flags().StringVar(&flagBarStyle, "style", "blocks", "bar style: blocks, dots, ascii, thin")
	barCmd.Flags().StringVar(&flagBarLabel, "label", "", "label prefix")
	barCmd.Flags().BoolVar(&flagBarNoPercent, "no-percent", false, "hide percentage")
	rootCmd.AddCommand(barCmd)
}

func runBar(cmd *cobra.Command, args []string) error {
	input := ""
	if len(args) > 0 {
		input = args[0]
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		input = strings.TrimSpace(t)
	}
	if input == "" {
		return fmt.Errorf("no value provided; see 'gloss bar --help'")
	}

	percent, err := render.ParseBarValue(input)
	if err != nil { return err }

	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 { w = 80 }
		width = w
	}

	// Subtract label width before rendering so the bar fits in one line.
	if flagBarLabel != "" {
		width -= len(flagBarLabel) + 1
		if width < 10 {
			width = 10
		}
	}

	result := render.RenderBar(percent, width, flagBarStyle, !flagBarNoPercent)
	if flagBarLabel != "" { result = flagBarLabel + " " + result }

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := render.ColorizeLines([]string{result}, gradient, flagColor, noColor)
	fmt.Println(lines[0])
	return nil
}
