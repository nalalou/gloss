package render

import "testing"

func TestDisplayWidthASCII(t *testing.T) {
	if w := displayWidth("hello"); w != 5 {
		t.Errorf("ASCII: expected 5, got %d", w)
	}
}

func TestDisplayWidthCJK(t *testing.T) {
	// Each CJK character is 2 columns wide
	if w := displayWidth("你好"); w != 4 {
		t.Errorf("CJK: expected 4, got %d", w)
	}
}

func TestDisplayWidthEmoji(t *testing.T) {
	if w := displayWidth("🎉"); w != 2 {
		t.Errorf("emoji: expected 2, got %d", w)
	}
}

func TestDisplayWidthANSI(t *testing.T) {
	colored := "\033[38;2;255;0;0mhello\033[0m"
	if w := displayWidth(colored); w != 5 {
		t.Errorf("ANSI colored: expected 5, got %d", w)
	}
}

func TestDisplayWidthMixed(t *testing.T) {
	// "A" (1) + "你" (2) + "🎉" (2) = 5, with ANSI wrapping
	mixed := "\033[1mA你🎉\033[0m"
	if w := displayWidth(mixed); w != 5 {
		t.Errorf("mixed: expected 5, got %d", w)
	}
}

func TestDisplayWidthEmpty(t *testing.T) {
	if w := displayWidth(""); w != 0 {
		t.Errorf("empty: expected 0, got %d", w)
	}
}
