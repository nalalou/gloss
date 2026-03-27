package render

import "testing"

func TestRenderSparkBasic(t *testing.T) {
	result := RenderSpark([]float64{1, 4, 7, 3, 9, 2, 5})
	if len([]rune(result)) != 7 {
		t.Errorf("expected 7 chars, got %d", len([]rune(result)))
	}
	if string([]rune(result)[4]) != "█" {
		t.Errorf("max should be █, got %q", string([]rune(result)[4]))
	}
}

func TestRenderSparkAllSame(t *testing.T) {
	result := RenderSpark([]float64{5, 5, 5, 5})
	runes := []rune(result)
	for i := range runes {
		if runes[i] != runes[0] {
			t.Errorf("all same should produce same char at pos %d", i)
		}
	}
}

func TestRenderSparkSingleValue(t *testing.T) {
	if len([]rune(RenderSpark([]float64{5}))) != 1 {
		t.Error("single value should produce 1 char")
	}
}

func TestRenderSparkEmpty(t *testing.T) {
	if RenderSpark([]float64{}) != "" {
		t.Error("empty should return empty")
	}
}

func TestParseSparkValues(t *testing.T) {
	vals, err := ParseSparkValues("1,4,7,3,9")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 5 {
		t.Errorf("expected 5, got %d", len(vals))
	}
}

func TestParseSparkValuesSpaceSep(t *testing.T) {
	vals, err := ParseSparkValues("1 4 7 3 9")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 5 {
		t.Errorf("expected 5, got %d", len(vals))
	}
}

func TestParseSparkValuesInvalid(t *testing.T) {
	if _, err := ParseSparkValues("abc"); err == nil {
		t.Error("expected error")
	}
}
