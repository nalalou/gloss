package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/render"
	"github.com/nalalou/gloss/internal/theme"
	"github.com/spf13/cobra"
)

var (
	flagFont     string
	flagGradient string
	flagColor    string
	flagShadow   bool
	flagAlign    string
	flagBorder   string
	flagAnimate  bool
	flagWidth    int
	flagNoColor  bool
)

var rootCmd = &cobra.Command{
	Use:   "gloss [text]",
	Short: "Render beautiful styled text banners",
	Long:  `gloss renders text as styled ASCII art banners. Reads gloss.toml if present.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runRoot,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagFont, "font", "", "font name or path to .flf file (default: block)")
	rootCmd.PersistentFlags().StringVar(&flagGradient, "gradient", "", `gradient colors "#hex1,#hex2" or preset name (fire, ocean, mono, neon, aurora)`)
	rootCmd.PersistentFlags().StringVar(&flagColor, "color", "", "solid text color as hex (#RRGGBB)")
	rootCmd.PersistentFlags().BoolVar(&flagShadow, "shadow", false, "add drop shadow")
	rootCmd.PersistentFlags().StringVar(&flagAlign, "align", "", "text alignment: left, center, right")
	rootCmd.PersistentFlags().StringVar(&flagBorder, "border", "", "border style: none, single, double, rounded, thick")
	rootCmd.PersistentFlags().BoolVar(&flagAnimate, "animate", false, "typewriter animation (TTY only)")
	rootCmd.PersistentFlags().IntVar(&flagWidth, "width", 0, "max output width in columns (0 = terminal width)")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "force plain text output")
}

func runRoot(cmd *cobra.Command, args []string) error {
	text, err := resolveText(args)
	if err != nil {
		return err
	}
	if strings.TrimSpace(text) == "" {
		return cmd.Usage()
	}

	envInfo := env.Detect()
	opts := theme.Resolve()

	flagOpts := flagsToOptions()
	opts = theme.Merge(opts, flagOpts)

	if envInfo.NoColor || flagNoColor {
		opts.NoColor = true
	}
	if envInfo.CI {
		opts.Animate = false
	}

	f, err := font.Load(opts.Font)
	if err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	output := render.Render(text, f, opts)

	if opts.Animate && envInfo.IsTTY && !envInfo.CI {
		return animate(output)
	}
	fmt.Println(output)
	return nil
}

func resolveText(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return strings.TrimRight(string(data), "\n"), nil
	}
	return "", nil
}

func flagsToOptions() theme.Options {
	opts := theme.Options{}
	opts.Font = flagFont
	opts.Color = flagColor
	opts.Shadow = flagShadow
	opts.Align = flagAlign
	opts.Border = flagBorder
	opts.Animate = flagAnimate
	opts.Width = flagWidth
	opts.NoColor = flagNoColor

	if flagGradient != "" {
		if start, end, ok := render.GradientPreset(flagGradient); ok {
			opts.Gradient = []string{
				fmt.Sprintf("#%02X%02X%02X", start.R, start.G, start.B),
				fmt.Sprintf("#%02X%02X%02X", end.R, end.G, end.B),
			}
		} else {
			parts := strings.Split(flagGradient, ",")
			if len(parts) == 2 {
				opts.Gradient = []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])}
			}
		}
	}

	return opts
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
