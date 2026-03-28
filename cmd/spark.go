package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var sparkCmd = &cobra.Command{
	Use:   "spark <values>",
	Short: "Render a sparkline",
	Long: `Renders a Unicode sparkline from numeric values. Values can be passed as
comma-separated (1,4,7), space-separated arguments, or piped via stdin.
Each value maps to a bar height character.`,
	Example: `  gloss spark 1,4,7,3,9,2,5
  gloss spark 1 4 7 3 9 2 5
  echo "1,4,7,3,9" | gloss spark`,
	Args: cobra.ArbitraryArgs,
	RunE: runSpark,
}

func init() {
	rootCmd.AddCommand(sparkCmd)
}

func runSpark(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("no values provided; see 'gloss spark --help'")
	}

	values, err := render.ParseSparkValues(input)
	if err != nil {
		return err
	}

	result := render.RenderSpark(values)
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := render.ColorizeLines([]string{result}, gradient, flagColor, noColor)
	fmt.Println(lines[0])
	return nil
}
