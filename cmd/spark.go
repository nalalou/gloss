package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var sparkCmd = &cobra.Command{
	Use:   "spark <values>",
	Short: "Render a sparkline",
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
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(io.LimitReader(os.Stdin, 4096))
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			input = strings.TrimSpace(string(data))
		}
	}
	if input == "" {
		return fmt.Errorf("no values provided; usage: gloss spark <values>")
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
