package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var flagCalloutType string

var calloutCmd = &cobra.Command{
	Use:   "callout <text>",
	Short: "Render a callout box",
	Long: `Renders a bordered callout box with an icon and header. Types (success,
error, warning, info) set the border color and icon. Use --gradient or --color
to override the default type color. Respects NO_COLOR.`,
	Example: `  gloss callout "Deploy requires approval" --type=warning
  gloss callout "All tests passing" --type=success
  echo "message" | gloss callout --type=error`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCallout,
}

func init() {
	calloutCmd.Flags().StringVar(&flagCalloutType, "type", "info", "callout type: success, error, warning, info")
	calloutCmd.Flags().StringVar(&flagCalloutType, "style", "info", "alias for --type")
	calloutCmd.Flags().MarkHidden("style")
	rootCmd.AddCommand(calloutCmd)
}

func runCallout(cmd *cobra.Command, args []string) error {
	text := ""
	if len(args) > 0 {
		text = args[0]
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		text = strings.TrimSpace(t)
	}
	if text == "" {
		return fmt.Errorf("no text provided; see 'gloss callout --help'")
	}

	validTypes := []string{"success", "error", "warning", "info"}
	valid := false
	for _, v := range validTypes {
		if flagCalloutType == v {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid --type %q; expected: %s", flagCalloutType, strings.Join(validTypes, ", "))
	}

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	result := render.RenderCallout(text, flagCalloutType)

	gradient := resolveGradientFlag()
	color := flagColor
	if color == "" && gradient == nil && !noColor {
		color = render.CalloutDefaultColor(flagCalloutType)
	}
	lines := strings.Split(result, "\n")
	lines = render.ColorizeLines(lines, gradient, color, noColor)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
