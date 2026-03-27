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
	if lines > 30 {
		lines = 30
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
			pos:   -(rng.Intn(lines + 4)), // staggered start
			speed: 1 + rng.Intn(3),
			trail: 2 + rng.Intn(4),
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

	// Green shades for trail: bright head → dim tail
	greenShades := []string{
		"\033[38;2;0;255;0m", // bright green (head)
		"\033[38;2;0;200;0m",
		"\033[38;2;0;150;0m",
		"\033[38;2;0;100;0m",
		"\033[38;2;0;60;0m", // dark green (tail)
	}
	reset := "\033[0m"

	if flagNoColor {
		greenShades = []string{"", "", "", "", ""}
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
	frameDelay := 50 * time.Millisecond

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
					cols[i].pos = -(rng.Intn(lines + 6))
					cols[i].speed = 1 + rng.Intn(3)
					cols[i].trail = 2 + rng.Intn(4)
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
				if shadeIdx >= len(greenShades) {
					shadeIdx = len(greenShades) - 1
				}

				sb.WriteString(greenShades[shadeIdx])
				sb.WriteRune(ch)
				if greenShades[shadeIdx] != "" {
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
