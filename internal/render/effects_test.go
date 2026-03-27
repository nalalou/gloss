package render

import (
	"strings"
	"testing"
)

func TestApplyBorderRounded(t *testing.T) {
	result := ApplyBorder("Hello", "rounded")
	if !strings.Contains(result, "Hello") {
		t.Error("border output should contain original text")
	}
	if !strings.Contains(result, "╭") {
		t.Error("rounded border should contain ╭")
	}
}

func TestApplyBorderNone(t *testing.T) {
	result := ApplyBorder("Hello", "none")
	if result != "Hello" {
		t.Errorf("border=none should return text unchanged, got %q", result)
	}
}

func TestApplyShadowAddsRows(t *testing.T) {
	lines := []string{"AAA", "BBB", "CCC"}
	result := ApplyShadow(strings.Join(lines, "\n"))
	resultLines := strings.Split(result, "\n")
	if len(resultLines) <= len(lines) {
		t.Errorf("shadow should add rows: input %d, output %d", len(lines), len(resultLines))
	}
}
