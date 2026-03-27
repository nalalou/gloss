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

var flagKVSeparator string

var kvCmd = &cobra.Command{
	Use:   "kv [key=value...]",
	Short: "Render aligned key-value pairs",
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
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err := io.ReadAll(io.LimitReader(os.Stdin, int64(maxInputSize)))
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			text := strings.TrimSpace(string(data))
			if text != "" {
				lines := strings.Split(text, "\n")
				pairs = render.ParseKVPairs(lines)
			}
		}
	}

	if len(pairs) == 0 {
		return fmt.Errorf("no data provided; usage: gloss kv key=value [key=value...]")
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
