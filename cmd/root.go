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

var version = "dev"

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
	Use:     "gloss [text]",
	Short:   "Render beautiful styled text banners",
	Long:    `gloss renders text as styled ASCII art banners. Reads gloss.toml if present.`,
	Version: version,
	Example: `  gloss "Hello World"
  gloss "Deploy v2.0" --gradient=fire --border=rounded --shadow
  echo "OK" | gloss --font=outline
  gloss "Title" --color="#FF5500" --align=center
  gloss fonts                         # preview all bundled fonts
  gloss init                          # scaffold a gloss.toml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runRoot,
}

func init() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVar(&flagFont, "font", "", "font name or path to .flf file (default: block)")
	rootCmd.PersistentFlags().StringVar(&flagGradient, "gradient", "", `gradient preset or hex colors "#hex1,#hex2[,#hex3...]"`)
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
		return showShowcase()
	}
	if err := validateFlags(); err != nil {
		return err
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

	// Animated gradient cycling
	if opts.Animate && envInfo.IsTTY && !envInfo.CI && len(opts.Gradient) >= 2 && !opts.NoColor {
		fontLines := f.Render(text)
		return animateGradient(fontLines, opts.Gradient)
	}

	output := render.Render(text, f, opts)

	if opts.Animate && envInfo.IsTTY && !envInfo.CI {
		return animate(output)
	}
	fmt.Println(output)
	return nil
}

func validateFlags() error {
	if flagAlign != "" {
		switch flagAlign {
		case "left", "center", "right":
		default:
			return fmt.Errorf("invalid --align %q; expected: left, center, right", flagAlign)
		}
	}
	if flagBorder != "" {
		switch flagBorder {
		case "none", "single", "double", "rounded", "thick":
		default:
			return fmt.Errorf("invalid --border %q; expected: none, single, double, rounded, thick", flagBorder)
		}
	}
	if flagWidth < 0 {
		return fmt.Errorf("invalid --width %d; must be non-negative", flagWidth)
	}
	if flagGradient != "" {
		if _, ok := render.GradientPreset(flagGradient); !ok {
			parts := strings.Split(flagGradient, ",")
			if len(parts) < 2 {
				names := strings.Join(render.GradientPresetNames(), ", ")
				return fmt.Errorf("invalid --gradient %q; use a preset name (%s) or 2+ hex colors: \"#hex1,#hex2[,#hex3...]\"", flagGradient, names)
			}
		}
	}
	return nil
}

const maxInputSize = 64 * 1024 // 64KB

func resolveText(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", nil
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(io.LimitReader(os.Stdin, maxInputSize+1))
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		if len(data) > maxInputSize {
			return "", fmt.Errorf("input too large (max %dKB); gloss is for short banner text", maxInputSize/1024)
		}
		return strings.TrimRight(string(data), "\n"), nil
	}
	return "", nil
}

func resolveGradientFlag() []string {
	if flagGradient == "" {
		return nil
	}
	if colors, ok := render.GradientPreset(flagGradient); ok {
		hexes := make([]string, len(colors))
		for i, c := range colors {
			hexes[i] = fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
		}
		return hexes
	}
	parts := strings.Split(flagGradient, ",")
	if len(parts) >= 2 {
		trimmed := make([]string, len(parts))
		for i, p := range parts {
			trimmed[i] = strings.TrimSpace(p)
		}
		return trimmed
	}
	return nil
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
	opts.Gradient = resolveGradientFlag()
	return opts
}

func showShowcase() error {
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor

	// Render "gloss" banner
	opts := theme.Defaults()
	opts.NoColor = noColor
	opts.Align = "left"
	if !noColor {
		opts.Gradient = []string{"#FF6B9D", "#6B9DFF"}
	}
	f, err := font.Load("small")
	if err != nil {
		f, _ = font.Load("")
	}
	if f != nil {
		banner := render.Render("gloss", f, opts)
		fmt.Println(banner)
	}

	fmt.Println("  Terminal styling toolkit — banners, badges, bars, and more.")
	fmt.Println()
	fmt.Println("  Try these:")
	fmt.Println(`    gloss "Hello World" --gradient=fire`)
	fmt.Println(`    gloss badge "OK" --type=success`)
	fmt.Println(`    gloss bar 75 --label="Progress"`)
	fmt.Println(`    gloss divider "Section"`)
	fmt.Println()
	fmt.Println("  Run 'gloss --help' for all commands and options.")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
