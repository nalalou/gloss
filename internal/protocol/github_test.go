package protocol

import (
	"strings"
	"testing"
)

func TestParseGitHubPlain(t *testing.T) {
	name, params, msg, ok := ParseGitHub("::warning::careful")
	if !ok || name != "warning" || msg != "careful" || len(params) != 0 {
		t.Errorf("got name=%q params=%v msg=%q ok=%v", name, params, msg, ok)
	}
}

func TestParseGitHubParams(t *testing.T) {
	name, params, msg, ok := ParseGitHub("::error file=a.go,line=3,title=Bad::boom%0Anext")
	if !ok || name != "error" || msg != "boom\nnext" {
		t.Fatalf("got name=%q msg=%q ok=%v", name, msg, ok)
	}
	if params["file"] != "a.go" || params["line"] != "3" || params["title"] != "Bad" {
		t.Errorf("params: %v", params)
	}
}

func TestParseGitHubRejectsGlossStyle(t *testing.T) {
	for _, line := range []string{"::warn Slow::query", "::warning Slow query::x", "::ok done", "::status id=b done Build"} {
		if _, _, _, ok := ParseGitHub(line); ok {
			t.Errorf("%q should not parse as a GitHub command", line)
		}
	}
}

func TestRenderLineGitHub(t *testing.T) {
	cases := map[string]string{
		"::error file=a.go,line=3::boom": "✗ a.go:3 boom",
		"::warning::careful":             "⚠ careful",
		"::notice title=FYI::read me":    "ℹ FYI: read me",
		"::endgroup::":                   "",
		"::debug::noise":                 "",
		"::add-mask::secret":             "",
	}
	for in, want := range cases {
		if got := RenderLine(in, 40, true); got != want {
			t.Errorf("RenderLine(%q) = %q, want %q", in, got, want)
		}
	}
	if got := RenderLine("::group::Build", 40, true); !strings.Contains(got, " Build ") || !strings.Contains(got, "─") {
		t.Errorf("group: %q", got)
	}
}

func TestRenderLineIgnoresID(t *testing.T) {
	cases := map[string]string{
		"::status id=b done Build": "✓ Build",
		"::status id=b fail Build": "✗ Build",
		"::ok id=x All green":      "✓ All green",
	}
	for in, want := range cases {
		if got := RenderLine(in, 40, true); got != want {
			t.Errorf("RenderLine(%q) = %q, want %q", in, got, want)
		}
	}
	if got := RenderLine("::bar id=p 50 Prog", 40, true); !strings.HasPrefix(got, "Prog ") || strings.Contains(got, "id=") {
		t.Errorf("bar with id: %q", got)
	}
}
