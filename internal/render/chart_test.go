package render

import (
	"strings"
	"testing"
)

func TestRenderChartBasic(t *testing.T) {
	lines := RenderChart([]float64{10, 5, 8, 3}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "█") {
		t.Error("top row should contain █ for max value column")
	}
}

func TestRenderChartSingleValue(t *testing.T) {
	lines := RenderChart([]float64{42}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	for _, line := range lines {
		if !strings.Contains(line, "█") {
			t.Error("single value should fill all rows")
		}
	}
}

func TestRenderChartAllSame(t *testing.T) {
	lines := RenderChart([]float64{5, 5, 5}, 4)
	for _, line := range lines {
		if !strings.Contains(line, "█") {
			t.Error("all-same values should fill all rows")
		}
	}
}

func TestRenderChartAllZeros(t *testing.T) {
	lines := RenderChart([]float64{0, 0, 0}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[3], "▄") {
		t.Error("all-zeros should have ▄ at bottom row")
	}
}

func TestRenderChartEmpty(t *testing.T) {
	lines := RenderChart([]float64{}, 4)
	if len(lines) != 0 {
		t.Errorf("empty input should return empty slice, got %d lines", len(lines))
	}
}

func TestRenderChartNegativesClamped(t *testing.T) {
	lines := RenderChart([]float64{-5, 10, -3}, 4)
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestRenderChartHeight(t *testing.T) {
	lines := RenderChart([]float64{10, 5}, 8)
	if len(lines) != 8 {
		t.Errorf("expected 8 lines for height=8, got %d", len(lines))
	}
}
