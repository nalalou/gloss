package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var flagBadgeType string

var badgeCmd = &cobra.Command{
	Use:   "badge <text>",
	Short: "Render a status badge",
	Long: `Renders a colored status badge with an icon and label text. Badge types
(success, error, warning, info, plain) determine the icon and default color.
Use --gradient or --color to override the default badge color.`,
	Example: `  gloss badge "Tests passing" --type=success
  gloss badge "Build failed" --type=error`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBadge,
}

func init() {
	badgeCmd.Flags().StringVar(&flagBadgeType, "type", "plain", "badge type: success, error, warning, info, plain")
	badgeCmd.Flags().StringVar(&flagBadgeType, "style", "plain", "alias for --type")
	badgeCmd.Flags().MarkHidden("style")
	rootCmd.AddCommand(badgeCmd)
}

func runBadge(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("no text provided; see 'gloss badge --help'")
	}

	validTypes := []string{"success", "error", "warning", "info", "plain"}
	valid := false
	for _, v := range validTypes {
		if flagBadgeType == v {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid --type %q; expected: %s", flagBadgeType, strings.Join(validTypes, ", "))
	}

	result := render.RenderBadge(text, flagBadgeType)
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	color := flagColor
	if color == "" && gradient == nil && !noColor {
		color = render.BadgeDefaultColor(flagBadgeType)
	}
	lines := render.ColorizeLines([]string{result}, gradient, color, noColor)
	fmt.Println(lines[0])
	return nil
}
