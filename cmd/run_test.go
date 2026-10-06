package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCaptured runs gloss with stdout redirected to a pipe (so not a TTY).
func runCaptured(t *testing.T, args ...string) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	_, runErr := run(args...)
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), runErr
}

func writeScript(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "child.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunPassesExitCodeAndStderr(t *testing.T) {
	script := writeScript(t, "echo '::ok fine'\necho oops >&2\necho \"GLOSS=$GLOSS\"\nexit 3\n")
	out, err := runCaptured(t, "run", "--no-color", script)

	var exit exitCodeError
	if !errors.As(err, &exit) || exit.code != 3 {
		t.Fatalf("want exit code 3, got %v", err)
	}
	for _, want := range []string{"✓ fine", "oops", "GLOSS=1", "✗ child.sh (exit 3)"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestRunSuccess(t *testing.T) {
	script := writeScript(t, "echo hello\n")
	out, err := runCaptured(t, "run", "--no-color", "--", script, "--not-a-gloss-flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "hello") || !strings.Contains(out, "✓ child.sh --not-a-gloss-flag") {
		t.Errorf("output:\n%s", out)
	}
}

func TestRunMissingCommand(t *testing.T) {
	_, err := runCaptured(t, "run", "gloss-no-such-command")
	var exit exitCodeError
	if !errors.As(err, &exit) || exit.code != 127 {
		t.Fatalf("want exit code 127, got %v", err)
	}
}

func TestRunDoesNotWaitForBackgroundChildren(t *testing.T) {
	script := writeScript(t, "sleep 30 &\necho started\n")
	out, err := runCaptured(t, "run", "--no-color", script)
	if err != nil || !strings.Contains(out, "started") {
		t.Fatalf("err=%v out=%s", err, out)
	}
}
