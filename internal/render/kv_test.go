package render

import (
	"strings"
	"testing"
)

func TestRenderKVBasic(t *testing.T) {
	pairs := [][]string{{"Status", "Running"}, {"Pods", "3/3"}}
	result := RenderKV(pairs, ":")
	if !strings.Contains(result, "Status") || !strings.Contains(result, "Running") {
		t.Error("should contain key and value")
	}
	if !strings.Contains(result, ":") {
		t.Error("should contain separator")
	}
}

func TestRenderKVAlignment(t *testing.T) {
	pairs := [][]string{{"A", "1"}, {"LongKey", "2"}}
	result := RenderKV(pairs, ":")
	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	// Both colons should be at the same column
	idx1 := strings.Index(lines[0], ":")
	idx2 := strings.Index(lines[1], ":")
	if idx1 != idx2 {
		t.Errorf("separators not aligned: line1=%d line2=%d", idx1, idx2)
	}
}

func TestRenderKVCustomSeparator(t *testing.T) {
	pairs := [][]string{{"Name", "Alice"}}
	result := RenderKV(pairs, "→")
	if !strings.Contains(result, "→") {
		t.Error("should use custom separator")
	}
}

func TestRenderKVEmpty(t *testing.T) {
	result := RenderKV([][]string{}, ":")
	if result != "" {
		t.Errorf("empty should return empty, got %q", result)
	}
}

func TestParseKVPairs(t *testing.T) {
	pairs := ParseKVPairs([]string{"Name=Alice", "Role=Engineer"})
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0][0] != "Name" || pairs[0][1] != "Alice" {
		t.Errorf("wrong pair: %v", pairs[0])
	}
}

func TestParseKVPairsValueWithEquals(t *testing.T) {
	pairs := ParseKVPairs([]string{"URL=https://example.com?a=1"})
	if pairs[0][1] != "https://example.com?a=1" {
		t.Errorf("should preserve = in value, got %q", pairs[0][1])
	}
}
