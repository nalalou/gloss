package theme

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

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

// fileConfig mirrors Options but with pointer fields so unset keys are nil.
type fileConfig struct {
	Font     *string  `toml:"font"`
	Gradient []string `toml:"gradient"`
	Color    *string  `toml:"color"`
	Shadow   *bool    `toml:"shadow"`
	Align    *string  `toml:"align"`
	Border   *string  `toml:"border"`
	Animate  *bool    `toml:"animate"`
	Width    *int     `toml:"width"`
}

// LoadFile parses a gloss.toml file and merges it with Defaults().
func LoadFile(path string) (Options, error) {
	opts := Defaults()
	var fc fileConfig
	if _, err := toml.DecodeFile(path, &fc); err != nil {
		return opts, err
	}
	if fc.Font != nil {
		opts.Font = *fc.Font
	}
	if fc.Gradient != nil {
		opts.Gradient = fc.Gradient
	}
	if fc.Color != nil {
		opts.Color = *fc.Color
	}
	if fc.Shadow != nil {
		opts.Shadow = *fc.Shadow
	}
	if fc.Align != nil {
		opts.Align = *fc.Align
	}
	if fc.Border != nil {
		opts.Border = *fc.Border
	}
	if fc.Animate != nil {
		opts.Animate = *fc.Animate
	}
	if fc.Width != nil {
		opts.Width = *fc.Width
	}
	return opts, nil
}

// Merge overlays non-zero fields from override onto base.
func Merge(base, override Options) Options {
	if override.Font != "" {
		base.Font = override.Font
	}
	if override.Gradient != nil {
		base.Gradient = override.Gradient
	}
	if override.Color != "" {
		base.Color = override.Color
	}
	if override.Shadow {
		base.Shadow = true
	}
	if override.Align != "" {
		base.Align = override.Align
	}
	if override.Border != "" {
		base.Border = override.Border
	}
	if override.Animate {
		base.Animate = true
	}
	if override.Width != 0 {
		base.Width = override.Width
	}
	if override.NoColor {
		base.NoColor = true
	}
	return base
}

// Resolve loads Options using the priority chain:
// GLOSS_THEME env → ./gloss.toml → ~/.config/gloss/gloss.toml → defaults
func Resolve() Options {
	dir, _ := os.Getwd()
	return resolveFromDir(dir)
}

func resolveFromDir(dir string) Options {
	if path := os.Getenv("GLOSS_THEME"); path != "" {
		if opts, err := LoadFile(path); err == nil {
			return opts
		}
	}
	projectPath := filepath.Join(dir, "gloss.toml")
	if _, err := os.Stat(projectPath); err == nil {
		if opts, err := LoadFile(projectPath); err == nil {
			return opts
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		userPath := filepath.Join(home, ".config", "gloss", "gloss.toml")
		if _, err := os.Stat(userPath); err == nil {
			if opts, err := LoadFile(userPath); err == nil {
				return opts
			}
		}
	}
	return Defaults()
}
