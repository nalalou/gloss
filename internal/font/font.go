package font

// Font renders text as multi-line ASCII art.
// Each element in the returned slice is one row of the rendered output.
type Font interface {
	Render(text string) []string
	Height() int
}
