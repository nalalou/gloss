package render

import (
	"strings"
	"testing"
)

func TestRenderListBullet(t *testing.T) {
	r := RenderList([]string{"Build", "Test", "Deploy"}, "bullet", false)
	lines := strings.Split(r, "\n")
	if len(lines) != 3 { t.Fatalf("expected 3 lines, got %d", len(lines)) }
	if !strings.Contains(lines[0], "•") || !strings.Contains(lines[0], "Build") { t.Error("wrong") }
}
func TestRenderListNumbered(t *testing.T) {
	r := RenderList([]string{"A", "B", "C"}, "numbered", false)
	if !strings.Contains(r, "1.") || !strings.Contains(r, "3.") { t.Error("missing numbers") }
}
func TestRenderListArrow(t *testing.T) {
	if !strings.Contains(RenderList([]string{"Item"}, "arrow", false), "→") { t.Error("missing →") }
}
func TestRenderListStatus(t *testing.T) {
	r := RenderList([]string{"Build:done", "Test:done", "Deploy:pending", "Rollback:fail"}, "bullet", true)
	if !strings.Contains(r, "✓") || !strings.Contains(r, "○") || !strings.Contains(r, "✗") { t.Error("missing status icons") }
}
func TestRenderListEmpty(t *testing.T) {
	if RenderList([]string{}, "bullet", false) != "" { t.Error("empty should return empty") }
}
