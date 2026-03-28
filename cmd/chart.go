package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagChartHeight int
	flagChartLabel  string
)

var chartCmd = &cobra.Command{
	Use:   "chart <values>",
	Short: "Render a vertical bar chart",
	Long: `Renders a tall vertical bar chart from numeric values. Values can be
comma-separated (220,195,180) or space-separated. Pipe via stdin or pass as args.
Each value becomes a column scaled relative to the maximum.`,
	Example: `  gloss chart 220,195,180,140,95,72,45
  gloss chart 10 5 8 3 --height=8
  gloss chart 220,195,180 --label="p99 (ms)" --gradient=fire
  echo "10,20,30,40,50" | gloss chart`,
	Args: cobra.ArbitraryArgs,
	RunE: runChart,
}

func init() {
	chartCmd.Flags().IntVar(&flagChartHeight, "height", 6, "chart height in rows")
	chartCmd.Flags().StringVar(&flagChartLabel, "label", "", "label above the chart")
	rootCmd.AddCommand(chartCmd)
}

func runChart(cmd *cobra.Command, args []string) error {
	input := ""
	if len(args) > 0 {
		input = strings.Join(args, " ")
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		input = strings.TrimSpace(t)
	}
	if input == "" {
		return fmt.Errorf("no data provided; see 'gloss chart --help'")
	}

	values, err := render.ParseSparkValues(input)
	if err != nil {
		return err
	}

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()

	height := flagChartHeight
	if height < 1 {
		height = 1
	}
	if height > 30 {
		height = 30
	}

	if flagChartLabel != "" {
		fmt.Printf("  \033[2m%s\033[0m\n", flagChartLabel)
	}

	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if isTTY && !noColor {
		animateChart(values, height, gradient, flagColor, noColor)
	} else {
		lines := render.RenderChart(values, height)
		lines = render.ColorizeLines(lines, gradient, flagColor, noColor)
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	return nil
}

func animateChart(values []float64, height int, gradient []string, color string, noColor bool) {
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	for i := 0; i < height; i++ {
		fmt.Println()
	}

	for col := 1; col <= len(values); col++ {
		lines := render.RenderChartPartial(values, height, col)
		lines = render.ColorizeLines(lines, gradient, color, noColor)

		fmt.Printf("\033[%dA", height)
		for _, line := range lines {
			fmt.Printf("\033[2K%s\n", line)
		}

		delay := 120 + (col * 10)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
