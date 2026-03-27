package theme

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	opts := Defaults()
	if opts.Font != "block" {
		t.Errorf("expected font=block, got %q", opts.Font)
	}
	if opts.Align != "left" {
		t.Errorf("expected align=left, got %q", opts.Align)
	}
	if opts.Border != "none" {
		t.Errorf("expected border=none, got %q", opts.Border)
	}
	if opts.Shadow {
		t.Error("expected shadow=false by default")
	}
	if opts.Animate {
		t.Error("expected animate=false by default")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	tomlPath := filepath.Join(dir, "gloss.toml")
	os.WriteFile(tomlPath, []byte(`
font = "outline"
shadow = true
align = "center"
`), 0644)

	opts, err := LoadFile(tomlPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Font != "outline" {
		t.Errorf("expected font=outline, got %q", opts.Font)
	}
	if !opts.Shadow {
		t.Error("expected shadow=true")
	}
	if opts.Align != "center" {
		t.Errorf("expected align=center, got %q", opts.Align)
	}
	// Unset keys should keep defaults
	if opts.Border != "none" {
		t.Errorf("expected border=none (default), got %q", opts.Border)
	}
}

func TestMerge(t *testing.T) {
	base := Defaults()
	base.Font = "outline"

	override := Options{Font: "chrome", Align: "right"}
	merged := Merge(base, override)

	if merged.Font != "chrome" {
		t.Errorf("expected font=chrome after merge, got %q", merged.Font)
	}
	if merged.Align != "right" {
		t.Errorf("expected align=right after merge, got %q", merged.Align)
	}
	if merged.Border != "none" {
		t.Errorf("expected border=none from base, got %q", merged.Border)
	}
}

func TestResolveNoFile(t *testing.T) {
	dir := t.TempDir()
	os.Unsetenv("GLOSS_THEME")
	opts := resolveFromDir(dir)
	if opts.Font != "block" {
		t.Errorf("expected default font=block, got %q", opts.Font)
	}
}
