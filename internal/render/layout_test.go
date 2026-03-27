package render

import (
	"strings"
	"testing"
)

func TestAlignCenter(t *testing.T) {
	result := ApplyLayout("Hello", "center", 20, false)
	if !strings.Contains(result, "Hello") {
		t.Error("centered output should contain original text")
	}
}

func TestAlignLeft(t *testing.T) {
	result := ApplyLayout("Hello", "left", 20, false)
	if !strings.Contains(result, "Hello") {
		t.Error("left-aligned output should contain original text")
	}
}

func TestWidthClamp(t *testing.T) {
	long := strings.Repeat("A", 100)
	result := ApplyLayout(long, "left", 20, false)
	_ = result // Just verify no panic
}
