package render

import (
	"strings"

	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/theme"
)

func Render(text string, f font.Font, opts theme.Options) string {
	lines := f.Render(text)

	if len(opts.Gradient) >= 2 && !opts.NoColor {
		colors, err := ParseGradientColors(opts.Gradient)
		if err == nil {
			lines = ApplyGradient(lines, colors, "horizontal", opts.NoColor)
		}
	} else if opts.Color != "" {
		lines = ApplySolidColor(lines, opts.Color, opts.NoColor)
	}

	joined := strings.Join(lines, "\n")

	if opts.Shadow && !opts.NoColor {
		joined = ApplyShadow(joined)
	}
	if opts.Border != "none" && opts.Border != "" {
		joined = ApplyBorder(joined, opts.Border)
	}

	joined = ApplyLayout(joined, opts.Align, opts.Width, opts.NoColor)
	return joined
}
