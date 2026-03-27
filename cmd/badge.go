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

var flagBadgeType string

var badgeCmd = &cobra.Command{
	Use:   "badge <text>",
	Short: "Render a status badge",
	Example: `  gloss badge "Tests passing" --type=success
  gloss badge "Build failed" --type=error`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBadge,
}

func init() {
	badgeCmd.Flags().StringVar(&flagBadgeType, "type", "plain", "badge type: success, error, warning, info, plain")
	rootCmd.AddCommand(badgeCmd)
}

func runBadge(cmd *cobra.Command, args []string) error {
	text := ""
	if len(args) > 0 {
		text = args[0]
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(io.LimitReader(os.Stdin, 4096))
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			text = strings.TrimSpace(string(data))
		}
	}
	if text == "" {
		return fmt.Errorf("no text provided; usage: gloss badge <text> --type=success")
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
