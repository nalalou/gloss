package protocol

import (
	"strings"
	"testing"
)

func TestOk(t *testing.T) {
	r := RenderLine("::ok Tests passing", 80, true)
	if !strings.Contains(r, "✓") || !strings.Contains(r, "Tests passing") {
		t.Errorf("got %q", r)
	}
}

func TestErr(t *testing.T) {
	r := RenderLine("::err Build failed", 80, true)
	if !strings.Contains(r, "✗") {
		t.Errorf("got %q", r)
	}
}

func TestErrorAlias(t *testing.T) {
	r := RenderLine("::error Something broke", 80, true)
	if !strings.Contains(r, "✗") {
		t.Errorf("got %q", r)
	}
}

func TestWarn(t *testing.T) {
	r := RenderLine("::warn Slow query", 80, true)
	if !strings.Contains(r, "⚠") {
		t.Errorf("got %q", r)
	}
}

func TestWarningAlias(t *testing.T) {
	r := RenderLine("::warning Watch out", 80, true)
	if !strings.Contains(r, "⚠") {
		t.Errorf("got %q", r)
	}
}

func TestInfo(t *testing.T) {
	r := RenderLine("::info Using cache", 80, true)
	if !strings.Contains(r, "ℹ") {
		t.Errorf("got %q", r)
	}
}

func TestBar(t *testing.T) {
	r := RenderLine("::bar 75 Coverage", 40, true)
	if !strings.Contains(r, "75%") || !strings.Contains(r, "Coverage") || !strings.Contains(r, "█") {
		t.Errorf("got %q", r)
	}
}

func TestBarNoLabel(t *testing.T) {
	r := RenderLine("::bar 50", 30, true)
	if !strings.Contains(r, "50%") {
		t.Errorf("got %q", r)
	}
}

func TestDivider(t *testing.T) {
	r := RenderLine("::divider Testing", 40, true)
	if !strings.Contains(r, "Testing") || !strings.Contains(r, "─") {
		t.Errorf("got %q", r)
	}
}

func TestDivAlias(t *testing.T) {
	r := RenderLine("::div Section", 40, true)
	if !strings.Contains(r, "Section") {
		t.Errorf("got %q", r)
	}
}

func TestDividerEmpty(t *testing.T) {
	r := RenderLine("::divider", 40, true)
	if !strings.Contains(r, "─") {
		t.Errorf("got %q", r)
	}
}

func TestCallout(t *testing.T) {
	r := RenderLine("::callout error Root cause found", 80, true)
	if !strings.Contains(r, "Root cause found") || !strings.Contains(r, "Error") {
		t.Errorf("got %q", r)
	}
}

func TestPlainText(t *testing.T) {
	r := RenderLine("Just regular text", 80, true)
	if r != "Just regular text" {
		t.Errorf("got %q", r)
	}
}

func TestUnknownDirective(t *testing.T) {
	r := RenderLine("::foobar something", 80, true)
	if r != "::foobar something" {
		t.Errorf("got %q", r)
	}
}

func TestEmptyDirective(t *testing.T) {
	r := RenderLine("::", 80, true)
	if r != "::" {
		t.Errorf("got %q", r)
	}
}

func TestList(t *testing.T) {
	r := RenderLine("::list done Build | done Test | fail Deploy", 80, true)
	if !strings.Contains(r, "✓") || !strings.Contains(r, "✗") {
		t.Errorf("got %q", r)
	}
}

func TestTable(t *testing.T) {
	r := RenderLine("::table Name=Alice | Role=Engineer", 80, true)
	if !strings.Contains(r, "Alice") {
		t.Errorf("got %q", r)
	}
}

func TestKV(t *testing.T) {
	r := RenderLine("::kv Status=Running | Pods=3/3", 80, true)
	if !strings.Contains(r, "Status") || !strings.Contains(r, "Running") {
		t.Errorf("got %q", r)
	}
}

func TestSpark(t *testing.T) {
	r := RenderLine("::spark 1,4,7,3,9", 80, true)
	if len([]rune(r)) != 5 {
		t.Errorf("expected 5 chars, got %q", r)
	}
}

func TestColorWithColor(t *testing.T) {
	r := RenderLine("::color red FAILED", 80, false)
	if !strings.Contains(r, "FAILED") || !strings.Contains(r, "\033[") {
		t.Errorf("got %q", r)
	}
}

func TestColorNoColor(t *testing.T) {
	r := RenderLine("::color red FAILED", 80, true)
	if r != "FAILED" {
		t.Errorf("got %q", r)
	}
}

func TestRain(t *testing.T) {
	r := RenderLine("::rain", 80, true)
	if r != "" {
		t.Errorf("rain should be empty in pipe, got %q", r)
	}
}
