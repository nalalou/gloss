package env

import (
	"os"

	"github.com/muesli/termenv"
	"golang.org/x/term"
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
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))

	output := termenv.NewOutput(os.Stdout)
	profile := output.ColorProfile()

	// Auto-disable color when not a TTY (piped output)
	if !isTTY || noColor {
		profile = termenv.Ascii
		noColor = true
	}

	return Info{
		NoColor: noColor,
		CI:      ci,
		IsTTY:   isTTY,
		Profile: profile,
	}
}
