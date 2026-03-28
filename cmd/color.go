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
	flagColorFg   string
	flagColorBold bool
	flagColorDim  bool
)

var colorCmd = &cobra.Command{
	Use:   "color <text>",
	Short: "Style text with colors, bold, or dim",
	Long: `Applies ANSI color and style to text for inline use in scripts. Outputs
without a trailing newline when piped, so it composes inside echo/printf.
When stdout is a terminal, a newline is appended for readability.`,
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
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		text = strings.TrimRight(t, "\n")
	}
	if text == "" {
		return fmt.Errorf("no text provided; see 'gloss color --help'")
	}

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	if noColor {
		fmt.Print(text)
		if isTTY {
			fmt.Println()
		}
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

	// Print WITHOUT newline when piped (enables inline use).
	// Append newline when stdout is a terminal for readability.
	fmt.Print(result)
	if isTTY {
		fmt.Println()
	}
	return nil
}
