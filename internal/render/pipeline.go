package render

import (
	"strings"

	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/theme"
)

// Render runs text through the full 4-stage pipeline.
func Render(text string, f font.Font, opts theme.Options) string {
	// Stage 1: font render
	lines := f.Render(text)

	// Stage 2: color/gradient
	if len(opts.Gradient) == 2 && !opts.NoColor {
		start, err1 := parseHex(opts.Gradient[0])
		end, err2 := parseHex(opts.Gradient[1])
		if err1 == nil && err2 == nil {
			lines = ApplyGradient(lines, start, end, "horizontal", opts.NoColor)
		}
	} else if opts.Color != "" {
		lines = ApplySolidColor(lines, opts.Color, opts.NoColor)
	}

	// Join lines for stages 3-4
	joined := strings.Join(lines, "\n")

	// Stage 3: effects
	if opts.Shadow && !opts.NoColor {
		joined = ApplyShadow(joined)
	}
	if opts.Border != "none" && opts.Border != "" {
		joined = ApplyBorder(joined, opts.Border)
	}

	// Stage 4: layout
	joined = ApplyLayout(joined, opts.Align, opts.Width, opts.NoColor)

	return joined
}
