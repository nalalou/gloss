package render

import (
	"strings"
	"testing"
)

func TestRenderCalloutWarning(t *testing.T) {
	r := RenderCallout("Approval needed", "warning")
	if !strings.Contains(r, "⚠") || !strings.Contains(r, "Warning") || !strings.Contains(r, "Approval needed") {
		t.Error("wrong")
	}
}

func TestRenderCalloutSuccess(t *testing.T) {
	r := RenderCallout("All passing", "success")
	if !strings.Contains(r, "✓") || !strings.Contains(r, "Success") {
		t.Error("wrong")
	}
}

func TestRenderCalloutError(t *testing.T) {
	if !strings.Contains(RenderCallout("Deleted", "error"), "✗") {
		t.Error("missing ✗")
	}
}

func TestRenderCalloutInfo(t *testing.T) {
	if !strings.Contains(RenderCallout("Docs", "info"), "ℹ") {
		t.Error("missing ℹ")
	}
}
