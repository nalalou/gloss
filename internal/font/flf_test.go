package font

import (
	"strings"
	"testing"
)

// buildMiniFont creates a minimal FLF font string with only space, A, B.
// Fills chars 33-64 with blank rows so A lands at position 65.
func buildMiniFont() string {
	var b strings.Builder
	// Header
	b.WriteString("flf2a$ 3 2 5 0 0\n")
	// Char 32 (space): 3 rows
	b.WriteString("  @\n")
	b.WriteString("  @\n")
	b.WriteString("  @@\n")
	// Chars 33–64: 32 filler characters (each 3 rows of empty)
	for ch := 33; ch <= 64; ch++ {
		b.WriteString(" @\n")
		b.WriteString(" @\n")
		b.WriteString(" @@\n")
	}
	// Char 65 (A)
	b.WriteString("/\\ @\n")
	b.WriteString("/--\\@\n")
	b.WriteString("/  \\@@\n")
	// Char 66 (B)
	b.WriteString("|--\\@\n")
	b.WriteString("|  |@\n")
	b.WriteString("|--/@@\n")
	return b.String()
}

func TestParseMiniFont(t *testing.T) {
	font, err := ParseFLF(strings.NewReader(buildMiniFont()))
	if err != nil {
		t.Fatalf("ParseFLF: %v", err)
	}
	if font.Height() != 3 {
		t.Errorf("expected height=3, got %d", font.Height())
	}
}

func TestRenderSingleChar(t *testing.T) {
	font, _ := ParseFLF(strings.NewReader(buildMiniFont()))

	lines := font.Render("A")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines for 'A', got %d", len(lines))
	}
	if lines[0] != "/\\ " {
		t.Errorf("expected line 0 = '/\\ ', got %q", lines[0])
	}
}

func TestRenderSpace(t *testing.T) {
	font, _ := ParseFLF(strings.NewReader(buildMiniFont()))

	lines := font.Render(" ")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines for space, got %d", len(lines))
	}
	for _, line := range lines {
		if line != "  " {
			t.Errorf("expected 2 spaces per row for space char, got %q", line)
		}
	}
}

func TestRenderMultiChar(t *testing.T) {
	font, _ := ParseFLF(strings.NewReader(buildMiniFont()))

	lines := font.Render("AB")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "/\\ |--\\" {
		t.Errorf("unexpected row 0: %q", lines[0])
	}
}
