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

var flagDividerStyle string

var dividerCmd = &cobra.Command{
	Use:   "divider [label]",
	Short: "Render a horizontal divider",
	Example: `  gloss divider
  gloss divider "Section Title"
  gloss divider --style=double`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDivider,
}

func init() {
	dividerCmd.Flags().StringVar(&flagDividerStyle, "style", "heavy", "divider style: heavy, light, double, dashed, dots, ascii")
	rootCmd.AddCommand(dividerCmd)
}

func runDivider(cmd *cobra.Command, args []string) error {
	label := ""
	if len(args) > 0 {
		label = args[0]
	}

	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	validStyles := []string{"heavy", "light", "double", "dashed", "dots", "ascii"}
	valid := false
	for _, s := range validStyles {
		if flagDividerStyle == s {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid --style %q; expected: %s", flagDividerStyle, strings.Join(validStyles, ", "))
	}

	result := render.RenderDivider(label, width, flagDividerStyle)

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := render.ColorizeLines([]string{result}, gradient, flagColor, noColor)
	fmt.Println(lines[0])
	return nil
}
