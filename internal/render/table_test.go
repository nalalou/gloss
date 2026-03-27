package render

import (
	"strings"
	"testing"
)

func TestRenderTableBasic(t *testing.T) {
	rows := [][]string{{"Name", "Alice"}, {"Role", "Engineer"}}
	result := RenderTable(rows, "rounded")
	if !strings.Contains(result, "Name") || !strings.Contains(result, "Alice") { t.Error("missing data") }
	if !strings.Contains(result, "╭") { t.Error("missing border") }
}

func TestRenderTableNoBorder(t *testing.T) {
	result := RenderTable([][]string{{"A", "1"}}, "none")
	if !strings.Contains(result, "A") { t.Error("missing data") }
	if strings.Contains(result, "╭") { t.Error("none should have no border") }
}

func TestRenderTableEmpty(t *testing.T) {
	if RenderTable([][]string{}, "rounded") != "" { t.Error("empty should return empty") }
}

func TestParseKVArgs(t *testing.T) {
	rows := ParseKVArgs([]string{"Name=Alice", "Role=Engineer"})
	if len(rows) != 2 { t.Fatalf("expected 2 rows, got %d", len(rows)) }
	if rows[0][0] != "Name" || rows[0][1] != "Alice" { t.Errorf("wrong: %v", rows[0]) }
}
