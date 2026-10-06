package watch

import (
	"strings"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	cases := map[time.Duration]string{
		1400 * time.Millisecond: "1.4s",
		38 * time.Second:        "38s",
		125 * time.Second:       "2m05s",
	}
	for d, want := range cases {
		if got := FormatDuration(d); got != want {
			t.Errorf("FormatDuration(%v) = %q, want %q", d, got, want)
		}
	}
}

func TestPanelStepTiming(t *testing.T) {
	clock := time.Unix(0, 0)
	p := NewPanel(60)
	p.now = func() time.Time { return clock }

	p.Set("build", "status", "running Build", true)
	clock = clock.Add(3 * time.Second)
	p.UpdateSpinnerFrame(1, true)
	if got := p.elements["build"].Rendered; !strings.HasSuffix(got, "Build 3.0s") {
		t.Errorf("live timer: %q", got)
	}

	clock = clock.Add(1200 * time.Millisecond)
	p.Set("build", "status", "done Build", true)
	if got := p.elements["build"].Rendered; got != "✓ Build 4.2s" {
		t.Errorf("finished: %q", got)
	}
}

func TestPanelQuickStepHasNoTiming(t *testing.T) {
	clock := time.Unix(0, 0)
	p := NewPanel(60)
	p.now = func() time.Time { return clock }
	p.Set("x", "status", "running Lint", true)
	clock = clock.Add(300 * time.Millisecond)
	p.Set("x", "status", "done Lint", true)
	if got := p.elements["x"].Rendered; got != "✓ Lint" {
		t.Errorf("got %q", got)
	}
}

func TestPanelSpinAnimates(t *testing.T) {
	p := NewPanel(60)
	p.Set("d", "spin", "Deploying now", true)
	if !p.HasRunning() {
		t.Fatal("::spin should count as running")
	}
	p.UpdateSpinnerFrame(1, true)
	if got := p.elements["d"].Rendered; !strings.HasSuffix(got, " Deploying now") {
		t.Errorf("got %q", got)
	}
}

func TestPanelTruncatesLongLines(t *testing.T) {
	p := NewPanel(20)
	p.Set("x", "ok", strings.Repeat("long ", 20), true)
	for _, l := range p.RenderLines() {
		if visibleLen(l) > 19 {
			t.Errorf("line wider than panel: %q", l)
		}
	}
}
