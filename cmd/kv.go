package cmd

import (
	"fmt"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/render"
	"github.com/spf13/cobra"
)

var flagKVSeparator string

var kvCmd = &cobra.Command{
	Use:   "kv [key=value...]",
	Short: "Render aligned key-value pairs",
	Long: `Renders key-value pairs with aligned columns. Pass pairs as "key=value"
arguments or pipe newline-delimited "key=value" lines via stdin. The
--separator flag controls the display separator between keys and values.`,
	Example: `  gloss kv "Status=Running" "Pods=3/3" "Image=nginx:1.25"
  gloss kv "Name=Alice" "Role=Engineer" --separator="→"
  echo -e "Status=Running\nPods=3/3" | gloss kv`,
	Args: cobra.ArbitraryArgs,
	RunE: runKV,
}

func init() {
	kvCmd.Flags().StringVar(&flagKVSeparator, "separator", ":", "separator between key and value")
	rootCmd.AddCommand(kvCmd)
}

func runKV(cmd *cobra.Command, args []string) error {
	var pairs [][]string

	if len(args) > 0 {
		pairs = render.ParseKVPairs(args)
	} else {
		t, err := readStdinText(int64(maxInputSize))
		if err != nil {
			return err
		}
		text := strings.TrimSpace(t)
		if text != "" {
			lines := strings.Split(text, "\n")
			pairs = render.ParseKVPairs(lines)
		}
	}

	if len(pairs) == 0 {
		return fmt.Errorf("no data provided; see 'gloss kv --help'")
	}

	result := render.RenderKV(pairs, flagKVSeparator)

	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	gradient := resolveGradientFlag()
	lines := strings.Split(result, "\n")
	lines = render.ColorizeLines(lines, gradient, flagColor, noColor)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
