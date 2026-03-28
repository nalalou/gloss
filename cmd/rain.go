package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagRainLines    int
	flagRainDuration float64
	flagRainChars    string
)

var rainCmd = &cobra.Command{
	Use:   "rain",
	Short: "Brief matrix-style rain transition",
	Example: `  gloss rain
  gloss rain --lines=3 --duration=1
  gloss rain --lines=10 --gradient=fire
  gloss rain --chars="01"`,
	Args: cobra.NoArgs,
	RunE: runRain,
}

func init() {
	rainCmd.Flags().IntVar(&flagRainLines, "lines", 6, "number of lines of rain")
	rainCmd.Flags().Float64Var(&flagRainDuration, "duration", 2.0, "duration in seconds")
	rainCmd.Flags().StringVar(&flagRainChars, "chars", "", "characters to use (default: matrix-style)")
	rootCmd.AddCommand(rainCmd)
}

// parseHexRGB parses a "#RRGGBB" string into r, g, b components.
func parseHexRGB(hex string) (r, g, b uint8, ok bool) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b, err == nil
}

// rainShades returns a slice of ANSI color strings for the rain trail.
// If a gradient flag is set, it derives shades from the gradient colors.
// Otherwise it falls back to the default green shades.
func rainShades() []string {
	gradient := resolveGradientFlag()
	if len(gradient) >= 2 {
		type rgb struct{ r, g, b uint8 }
		colors := make([]rgb, 0, len(gradient))
		for _, h := range gradient {
			r, g, b, ok := parseHexRGB(h)
			if ok {
				colors = append(colors, rgb{r, g, b})
			}
		}
		if len(colors) >= 2 {
			// White head + 9 shades that dim through the gradient
			shades := make([]string, 10)
			shades[0] = "\033[38;2;255;255;255m" // white head
			for i := 1; i < 10; i++ {
				t := float64(i-1) / 8.0
				idx := t * float64(len(colors)-1)
				lo := int(idx)
				hi := lo + 1
				if hi >= len(colors) {
					hi = len(colors) - 1
				}
				frac := idx - float64(lo)
				cr := float64(colors[lo].r)*(1-frac) + float64(colors[hi].r)*frac
				cg := float64(colors[lo].g)*(1-frac) + float64(colors[hi].g)*frac
				cb := float64(colors[lo].b)*(1-frac) + float64(colors[hi].b)*frac
				dim := 1.0 - float64(i)*0.09
				if dim < 0.15 {
					dim = 0.15
				}
				shades[i] = fmt.Sprintf("\033[38;2;%d;%d;%dm",
					int(cr*dim), int(cg*dim), int(cb*dim))
			}
			return shades
		}
	}
	// Default green shades — white head like cmatrix, long gradient tail
	return []string{
		"\033[38;2;255;255;255m", // white head
		"\033[38;2;180;255;180m", // near-white green
		"\033[38;2;0;255;0m",    // bright green
		"\033[38;2;0;220;0m",
		"\033[38;2;0;180;0m",
		"\033[38;2;0;140;0m",
		"\033[38;2;0;100;0m",
		"\033[38;2;0;70;0m",
		"\033[38;2;0;45;0m",
		"\033[38;2;0;25;0m", // barely visible tail
	}
}

func runRain(cmd *cobra.Command, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	width := flagWidth
	if width == 0 {
		w, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil || w <= 0 {
			w = 80
		}
		width = w
	}

	lines := flagRainLines
	if lines < 1 {
		lines = 1
	}
	if lines > 200 {
		lines = 200
	}

	// Character set
	chars := []rune("ｱｲｳｴｵｶｷｸｹｺ0123456789ﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉ")
	if flagRainChars != "" {
		chars = []rune(flagRainChars)
	}
	if flagNoColor {
		chars = []rune(".:+*#")
	}

	// Column state: each column has a "drop" position (-1 = inactive)
	type column struct {
		pos   int // current head position (row index), -1 = waiting
		speed int // frames between drops (1=fast, 3=slow)
		tick  int // frame counter for this column
		trail int // trail length
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cols := make([]column, width)
	for i := range cols {
		cols[i] = column{
			pos:   -(rng.Intn(lines/2 + 3)),
			speed: 1 + rng.Intn(2),
			trail: 4 + rng.Intn(8),
		}
	}

	// Grid
	grid := make([][]rune, lines)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	shades := rainShades()
	reset := "\033[0m"

	if flagNoColor {
		shades = []string{"", "", "", "", ""}
		reset = ""
	}

	// Hide cursor, print initial empty lines
	fmt.Fprint(os.Stderr, "\033[?25l")
	defer fmt.Fprint(os.Stderr, "\033[?25h")

	for i := 0; i < lines; i++ {
		fmt.Fprintln(os.Stderr)
	}

	duration := time.Duration(flagRainDuration * float64(time.Second))
	deadline := time.Now().Add(duration)
	frameDelay := 65 * time.Millisecond

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			clearRain(lines)
			return nil
		default:
		}

		// Update columns
		for i := range cols {
			cols[i].tick++
			if cols[i].tick >= cols[i].speed {
				cols[i].tick = 0
				cols[i].pos++

				// Place new character at head position if in bounds
				if cols[i].pos >= 0 && cols[i].pos < lines {
					grid[cols[i].pos][i] = chars[rng.Intn(len(chars))]
				}

				// Clear trail tail
				tailPos := cols[i].pos - cols[i].trail
				if tailPos >= 0 && tailPos < lines {
					grid[tailPos][i] = ' '
				}

				// Reset column when fully off screen
				if cols[i].pos-cols[i].trail >= lines {
					cols[i].pos = -(rng.Intn(lines/3 + 3))
					cols[i].speed = 1 + rng.Intn(2)
					cols[i].trail = 4 + rng.Intn(8)
				}
			}
		}

		// Render grid
		fmt.Fprintf(os.Stderr, "\033[%dA", lines) // move cursor up
		for row := 0; row < lines; row++ {
			var sb strings.Builder
			for col := 0; col < width; col++ {
				ch := grid[row][col]
				if ch == ' ' {
					sb.WriteRune(' ')
					continue
				}

				// Find distance from nearest column head for color
				dist := 0
				if cols[col].pos >= row {
					dist = cols[col].pos - row
				}
				shadeIdx := dist
				if shadeIdx >= len(shades) {
					shadeIdx = len(shades) - 1
				}

				sb.WriteString(shades[shadeIdx])
				sb.WriteRune(ch)
				if shades[shadeIdx] != "" {
					sb.WriteString(reset)
				}
			}
			fmt.Fprintf(os.Stderr, "\033[2K%s\n", sb.String())
		}

		time.Sleep(frameDelay)
	}

	// Clean up: clear the rain area
	clearRain(lines)
	return nil
}

func clearRain(lines int) {
	fmt.Fprintf(os.Stderr, "\033[%dA", lines)
	for i := 0; i < lines; i++ {
		fmt.Fprintf(os.Stderr, "\033[2K\n")
	}
	fmt.Fprintf(os.Stderr, "\033[%dA", lines)
}
