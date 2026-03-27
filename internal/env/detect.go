package env

import (
	"os"

	"github.com/muesli/termenv"
)

// Info holds detected terminal environment state.
type Info struct {
	NoColor bool // true if output should be plain text
	CI      bool // true if running in CI (disables animation, keeps color)
	IsTTY   bool // true if stdout is a TTY
	Profile termenv.Profile
}

// Detect reads environment variables and terminal capabilities.
func Detect() Info {
	noColor := os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"
	ci := os.Getenv("CI") == "true"

	output := termenv.NewOutput(os.Stdout)
	profile := output.ColorProfile()
	isTTY := output.HasDarkBackground() || profile != termenv.Ascii

	if noColor {
		profile = termenv.Ascii
	}

	return Info{
		NoColor: noColor,
		CI:      ci,
		IsTTY:   isTTY,
		Profile: profile,
	}
}
