package watch

import (
	"strings"
	"testing"
)

func TestPanelAddElement(t *testing.T) {
	p := NewPanel(60)
	p.Set("build", "ok", "Build complete", true)
	if p.Len() != 1 { t.Errorf("len: %d", p.Len()) }
	found := false
	for _, l := range p.RenderLines() {
		if strings.Contains(l, "Build complete") { found = true }
	}
	if !found { t.Error("missing Build complete") }
}

func TestPanelUpdateElement(t *testing.T) {
	p := NewPanel(60)
	p.Set("build", "ok", "Building", true)
	p.Set("build", "ok", "Build complete", true)
	if p.Len() != 1 { t.Errorf("len after update: %d", p.Len()) }
}

func TestPanelInsertionOrder(t *testing.T) {
	p := NewPanel(60)
	p.Set("a", "ok", "First", true)
	p.Set("b", "ok", "Second", true)
	p.Set("c", "ok", "Third", true)
	lines := p.RenderLines()
	posA, posB, posC := -1, -1, -1
	for i, l := range lines {
		if strings.Contains(l, "First") { posA = i }
		if strings.Contains(l, "Second") { posB = i }
		if strings.Contains(l, "Third") { posC = i }
	}
	if posA >= posB || posB >= posC { t.Errorf("order: a=%d b=%d c=%d", posA, posB, posC) }
}

func TestPanelRemove(t *testing.T) {
	p := NewPanel(60)
	p.Set("a", "ok", "Test", true)
	p.Remove("a")
	if p.Len() != 0 { t.Errorf("len after remove: %d", p.Len()) }
}

func TestPanelEmpty(t *testing.T) {
	p := NewPanel(60)
	if len(p.RenderLines()) != 0 { t.Error("empty panel should have 0 lines") }
	if p.Height() != 0 { t.Error("empty height should be 0") }
}

func TestPanelHeight(t *testing.T) {
	p := NewPanel(60)
	p.Set("a", "ok", "Test", true)
	if p.Height() < 2 { t.Errorf("height: %d", p.Height()) }
}

func TestPanelHasRunning(t *testing.T) {
	p := NewPanel(60)
	p.Set("db", "status", "running Migrations", true)
	if !p.HasRunning() { t.Error("should have running") }
	p.Set("db", "status", "done Migrations", true)
	if p.HasRunning() { t.Error("should not have running") }
}

func TestPanelBar(t *testing.T) {
	p := NewPanel(40)
	p.Set("prog", "bar", "75 Coverage", true)
	found := false
	for _, l := range p.RenderLines() {
		if strings.Contains(l, "█") && strings.Contains(l, "75%") { found = true }
	}
	if !found { t.Error("missing bar") }
}
