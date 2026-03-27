package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var flagCalloutType string

var calloutCmd = &cobra.Command{
	Use:   "callout <text>",
	Short: "Render a callout box",
	Example: `  gloss callout "Deploy requires approval" --type=warning
  gloss callout "All tests passing" --type=success
  echo "message" | gloss callout --type=error`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCallout,
}

func init() {
	calloutCmd.Flags().StringVar(&flagCalloutType, "type", "info", "callout type: success, error, warning, info")
	rootCmd.AddCommand(calloutCmd)
}

func runCallout(cmd *cobra.Command, args []string) error {
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
			text = strings.TrimSpace(string(data))
		}
	}
	if text == "" {
		return fmt.Errorf("no text provided; usage: gloss callout <text> --type=warning")
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

	result := render.RenderCallout(text, flagCalloutType)
	fmt.Println(result)
	return nil
}
