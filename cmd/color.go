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

var (
	flagColorFg   string
	flagColorBold bool
	flagColorDim  bool
)

var colorCmd = &cobra.Command{
	Use:   "color <text>",
	Short: "Style text with colors, bold, or dim",
	Example: `  gloss color "FAILED" --fg=red
  gloss color "OK" --fg=green --bold
  gloss color "note" --dim
  gloss color "custom" --fg="#FF6B9D"
  echo "Status: $(gloss color OK --fg=green)"`,
	Args: cobra.MaximumNArgs(1),
	RunE: runColor,
}

func init() {
	colorCmd.Flags().StringVar(&flagColorFg, "fg", "", "foreground color: name (red, green, blue...) or hex (#RRGGBB)")
	colorCmd.Flags().BoolVar(&flagColorBold, "bold", false, "bold text")
	colorCmd.Flags().BoolVar(&flagColorDim, "dim", false, "dim/faint text")
	rootCmd.AddCommand(colorCmd)
}

func runColor(cmd *cobra.Command, args []string) error {
	text := ""
	if len(args) > 0 {
		text = args[0]
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(io.LimitReader(os.Stdin, int64(maxInputSize)))
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			text = strings.TrimRight(string(data), "\n")
		}
	}
	if text == "" {
		return fmt.Errorf("no text provided; usage: gloss color <text> --fg=red")
	}

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	if noColor {
		fmt.Print(text)
		return nil
	}

	// Resolve the foreground color
	fg := ""
	if flagColorFg != "" {
		hex, ok := render.ResolveColor(flagColorFg)
		if !ok {
			return fmt.Errorf("unknown color %q; use a name (red, green, blue, yellow, cyan, magenta, white, gray, orange, pink, purple) or hex (#RRGGBB)", flagColorFg)
		}
		fg = hex
	}

	result := render.RenderStyled(text, fg, flagColorBold, flagColorDim)

	// Print WITHOUT newline — this enables inline use: echo "$(gloss color OK --fg=green)"
	fmt.Print(result)
	return nil
}
