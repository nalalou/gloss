package env

import (
	"os"
	"testing"
)

func TestNoColorEnv(t *testing.T) {
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")
	d := Detect()
	if !d.NoColor {
		t.Error("expected NoColor=true when NO_COLOR=1")
	}
}

func TestTermDumb(t *testing.T) {
	os.Setenv("TERM", "dumb")
	defer os.Unsetenv("TERM")
	d := Detect()
	if !d.NoColor {
		t.Error("expected NoColor=true when TERM=dumb")
	}
}

func TestCIDisablesAnimation(t *testing.T) {
	os.Setenv("CI", "true")
	defer os.Unsetenv("CI")
	d := Detect()
	if !d.CI {
		t.Error("expected CI=true when CI=true")
	}
}
