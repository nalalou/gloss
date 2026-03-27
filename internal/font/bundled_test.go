package font

import "testing"

func TestLoadBundledBlock(t *testing.T) {
	f, err := Load("block")
	if err != nil {
		t.Fatalf("Load(block): %v", err)
	}
	lines := f.Render("Hi")
	if len(lines) == 0 {
		t.Error("expected non-empty render")
	}
}

func TestLoadBundledAllFonts(t *testing.T) {
	names := []string{"block", "outline", "round", "thin", "3d", "chrome"}
	for _, name := range names {
		f, err := Load(name)
		if err != nil {
			t.Errorf("Load(%q): %v", name, err)
			continue
		}
		lines := f.Render("A")
		if len(lines) == 0 {
			t.Errorf("font %q: empty render", name)
		}
	}
}

func TestLoadUnknownFontReturnsError(t *testing.T) {
	_, err := Load("doesnotexist")
	if err == nil {
		t.Error("expected error for unknown font name")
	}
}
