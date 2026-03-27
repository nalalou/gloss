package theme

import "testing"

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
