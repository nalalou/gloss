package theme

// Options holds all rendering configuration. Populated from defaults →
// gloss.toml → CLI flags (each layer overrides the previous).
type Options struct {
	Font     string   // "block", "outline", etc. or path to .flf file
	Gradient []string // ["#hex1", "#hex2"] or nil for solid color
	Color    string   // hex color or "" for terminal default
	Shadow   bool
	Align    string // "left", "center", "right"
	Border   string // "none", "single", "double", "rounded", "thick"
	Animate  bool
	Width    int  // 0 = use terminal width
	NoColor  bool // force plain text output
}

// Defaults returns the baseline Options before any config or flags are applied.
func Defaults() Options {
	return Options{
		Font:   "block",
		Align:  "left",
		Border: "none",
	}
}
