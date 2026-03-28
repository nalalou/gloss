package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/nalalou/gloss/internal/font"
	"github.com/nalalou/gloss/internal/render"
	"github.com/nalalou/gloss/internal/theme"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "See what gloss does",
	Args:  cobra.NoArgs,
	RunE:  runDemo,
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

// ── animation helpers ───────────────────────────────────────

const clearLine = "\033[2K\r"

var spinFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// animateSpin shows a spinner with text, then resolves to a badge.
func animateSpin(label string, duration time.Duration, badgeText, badgeType string) {
	dim := "\033[2m"
	reset := "\033[0m"
	cyan := "\033[36m"

	frames := int(duration / (80 * time.Millisecond))
	for i := 0; i < frames; i++ {
		frame := spinFrames[i%len(spinFrames)]
		fmt.Printf("%s%s%s %s%s%s%s", clearLine, cyan, frame, dim, label, reset, reset)
		time.Sleep(80 * time.Millisecond)
	}

	badge := render.RenderBadge(badgeText, badgeType)
	color := render.BadgeDefaultColor(badgeType)
	colored := render.ColorizeLines([]string{badge}, nil, color, false)[0]
	fmt.Printf("%s%s\n", clearLine, colored)
}

// animateBar fills a progress bar from 0 to target over duration.
func animateBar(target int, width int, duration time.Duration, gradient []string) {
	steps := 30
	delay := duration / time.Duration(steps)
	for i := 1; i <= steps; i++ {
		pct := float64(target) * float64(i) / float64(steps)
		bar := render.RenderBar(pct, width, "blocks", true)
		colored := render.ColorizeLines([]string{bar}, gradient, "", false)[0]
		fmt.Printf("%s%s", clearLine, colored)
		time.Sleep(delay)
	}
	fmt.Println()
}

// animateSparkline reveals a sparkline one bar at a time.
func animateSparkline(values []float64, gradient []string) {
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	max := 0.0
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	built := ""
	for i, v := range values {
		idx := int(math.Round(v / max * float64(len(blocks)-1)))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(blocks) {
			idx = len(blocks) - 1
		}
		built += string(blocks[idx])
		colored := render.ColorizeLines([]string{built}, gradient, "", false)[0]
		fmt.Printf("%s%s", clearLine, colored)
		delay := 80 + (i * 15) // slightly accelerating rhythm
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
	fmt.Println()
}

// animateDivider draws a divider from center outward.
func animateDivider(label string, width int, style string, gradient []string) {
	char := '━'
	switch style {
	case "light":
		char = '─'
	case "double":
		char = '═'
	case "dashed":
		char = '╌'
	}

	padded := ""
	if label != "" {
		padded = " " + label + " "
	}
	labelLen := len([]rune(padded))
	sideTotal := width - labelLen
	if sideTotal < 2 {
		sideTotal = 2
	}
	halfSteps := sideTotal / 2

	// Grow from center outward
	for i := 1; i <= halfSteps; i++ {
		left := strings.Repeat(string(char), i)
		right := strings.Repeat(string(char), i)
		line := left + padded + right
		colored := render.ColorizeLines([]string{line}, gradient, "", false)[0]
		fmt.Printf("%s%s", clearLine, colored)
		time.Sleep(15 * time.Millisecond)
	}
	// Final full-width render
	full := render.RenderDivider(label, width, style)
	colored := render.ColorizeLines([]string{full}, gradient, "", false)[0]
	fmt.Printf("%s%s\n", clearLine, colored)
}

// animateTypewriter prints text character by character with gradient.
func animateTypewriter(text string, gradient []string, delay time.Duration) {
	runes := []rune(text)
	for i := 1; i <= len(runes); i++ {
		partial := string(runes[:i])
		colored := render.ColorizeLines([]string{partial}, gradient, "", false)[0]
		fmt.Printf("%s%s", clearLine, colored)
		time.Sleep(delay)
	}
	fmt.Println()
}

// animateBannerReveal renders a FIGlet banner line by line.
func animateBannerReveal(text string, opts theme.Options, delay time.Duration) error {
	f, err := font.Load(opts.Font)
	if err != nil {
		return err
	}
	banner := render.Render(text, f, opts)
	lines := strings.Split(banner, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Println(line)
			time.Sleep(delay)
		}
	}
	return nil
}

// ── the demo ────────────────────────────────────────────────

func runDemo(cmd *cobra.Command, args []string) error {
	dim := "\033[2m"
	reset := "\033[0m"

	// Get terminal width
	width := 60
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		width = w
		if width > 80 {
			width = 80
		}
	}
	barWidth := width - 6
	if barWidth < 30 {
		barWidth = 30
	}

	// Hide cursor for clean animations
	fmt.Print(hideCursor)
	defer fmt.Print(showCursor)

	// ═══════════════════════════════════════════════════════════
	// ACT 1: YOUR FRIDAY AFTERNOON
	// ═══════════════════════════════════════════════════════════

	noise := []string{
		"$ npm run build",
		"",
		"> app@2.4.1 build",
		"> next build",
		"",
		"   Creating an optimized production build...",
	}
	for _, line := range noise {
		fmt.Printf("%s%s%s\n", dim, line, reset)
		time.Sleep(60 * time.Millisecond)
	}
	time.Sleep(700 * time.Millisecond)

	slow := []string{
		"   Compiled successfully.",
		"   Collecting page data...",
		"   Generating static pages (0/7)...",
		"   Generating static pages (7/7)",
		"   Finalizing page optimization...",
	}
	for _, line := range slow {
		fmt.Printf("%s%s%s\n", dim, line, reset)
		time.Sleep(140 * time.Millisecond)
	}
	fmt.Println()
	time.Sleep(100 * time.Millisecond)

	routes := []string{
		"Route (app)                    Size  First Load JS",
		"┌ ○ /                          5.2 kB     89.3 kB",
		"├ ○ /about                     1.8 kB     85.9 kB",
		"├ ○ /api/health                0 B        79.1 kB",
		"├ λ /api/users                 0 B        79.1 kB",
		"├ ○ /dashboard                 12.4 kB    96.5 kB",
		"├ ○ /login                     3.1 kB     87.2 kB",
		"└ ○ /settings                  4.7 kB     88.8 kB",
	}
	for _, line := range routes {
		fmt.Printf("%s%s%s\n", dim, line, reset)
		time.Sleep(40 * time.Millisecond)
	}
	time.Sleep(250 * time.Millisecond)

	docker := []string{
		"",
		"$ docker build -t app:v2.4.1 .",
		"[+] Building 24.3s (12/12) FINISHED",
		" => [1/8] FROM node:20-alpine",
		" => [2/8] WORKDIR /app",
		" => [3/8] COPY package*.json ./",
		" => [4/8] RUN npm ci --production",
		" => [5/8] COPY . .",
		" => [6/8] RUN npm run build",
		" => [7/8] RUN rm -rf node_modules/.cache",
		" => [8/8] EXPOSE 3000",
		" => exporting to image",
	}
	for _, line := range docker {
		fmt.Printf("%s%s%s\n", dim, line, reset)
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(350 * time.Millisecond)

	deploy := []string{
		"",
		"$ kubectl rollout status deployment/app -n prod",
		"Waiting for deployment rollout to finish: 0 of 3 updated...",
		"Waiting for deployment rollout to finish: 1 of 3 updated...",
		"Waiting for deployment rollout to finish: 2 of 3 updated...",
		"deployment \"app\" successfully rolled out",
	}
	for _, line := range deploy {
		fmt.Printf("%s%s%s\n", dim, line, reset)
		time.Sleep(180 * time.Millisecond)
	}

	time.Sleep(1000 * time.Millisecond)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// ACT 2: PIPE IT THROUGH
	// ═══════════════════════════════════════════════════════════

	animateDivider("", width, "heavy", []string{"#FF4500", "#FFD700"})
	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	// Spinners resolve to badges — feels like watching a live deploy
	animateSpin("next build", 800*time.Millisecond, "Build (12s)", "success")
	animateSpin("jest --ci", 1200*time.Millisecond, "14 tests", "success")
	animateSpin("docker push", 700*time.Millisecond, "app:v2.4.1 → ECR", "success")
	animateSpin("kubectl rollout", 900*time.Millisecond, "3/3 pods", "success")

	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// ACT 3: NOW YOU SEE IT
	// ═══════════════════════════════════════════════════════════

	// Bar fills — the satisfying part
	animateBar(100, barWidth, 1800*time.Millisecond, []string{"#FF4500", "#FFD700"})
	time.Sleep(400 * time.Millisecond)

	fmt.Println()

	// Two tall bar charts — the visual centerpiece
	// p99 latency dropping
	fmt.Printf("  \033[2mp99 (ms)\033[0m\n")
	animateChart(
		[]float64{220, 195, 180, 140, 95, 72, 45, 38, 22, 18},
		6, []string{"#FF4444", "#44FF88"}, "", false,
	)
	time.Sleep(400 * time.Millisecond)

	// rps climbing
	fmt.Printf("  \033[2mrps\033[0m\n")
	animateChart(
		[]float64{12, 18, 31, 45, 64, 87, 93, 120, 142, 158},
		6, []string{"#006994", "#00CED1"}, "", false,
	)
	time.Sleep(500 * time.Millisecond)

	fmt.Println()

	// The deploy receipt
	pairs := [][]string{
		{"app", "v2.4.1"},
		{"bundle", "412 KB"},
		{"build", "12s"},
		{"rollout", "3s"},
	}
	kv := render.RenderKV(pairs, "·")
	kvLines := strings.Split(kv, "\n")
	kvLines = render.ColorizeLines(kvLines, []string{"#006994", "#00CED1"}, "", false)
	for _, line := range kvLines {
		if line != "" {
			fmt.Println(line)
			time.Sleep(100 * time.Millisecond)
		}
	}

	time.Sleep(800 * time.Millisecond)
	fmt.Println()

	// ═══════════════════════════════════════════════════════════
	// ACT 4: THE DROP
	// ═══════════════════════════════════════════════════════════

	opts := theme.Options{
		Font:     "doom",
		Gradient: []string{"#FF4500", "#FFD700", "#FF4500"},
		Border:   "none",
	}
	if err := animateBannerReveal("shipped.", opts, 70*time.Millisecond); err != nil {
		return err
	}

	time.Sleep(1500 * time.Millisecond)

	animateTypewriter(
		"$ ./deploy.sh | gloss watch",
		[]string{"#888888", "#FFFFFF"},
		45*time.Millisecond,
	)

	time.Sleep(800 * time.Millisecond)
	fmt.Println()

	callout := render.RenderCallout("The script didn't change. The output did.", "success")
	fmt.Println(callout)
	fmt.Println()

	return nil
}
