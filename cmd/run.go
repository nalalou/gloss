package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/nalalou/gloss/internal/env"
	"github.com/nalalou/gloss/internal/watch"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [--] command [args...]",
	Short: "Run a command and render its :: output",
	Long: `Runs a command and renders its output like gloss watch, including
error output. Shows a timer while it runs and whether it passed at the end.
gloss exits with the command's own exit code.

The command sees GLOSS=1 in its environment. When gloss is drawing to a
terminal, FORCE_COLOR=1 and CLICOLOR_FORCE=1 are also set, so tools keep
their colors even though their output goes through gloss.`,
	Example: `  gloss run npm test
  gloss run -- ./deploy.sh --env=prod
  gloss run python agent.py`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRun,
}

func init() {
	// Flags after the command name belong to the command, not to gloss.
	runCmd.Flags().SetInterspersed(false)
	rootCmd.AddCommand(runCmd)
}

// exitCodeError makes gloss exit with a specific code and no message.
type exitCodeError struct{ code int }

func (e exitCodeError) Error() string { return fmt.Sprintf("exit status %d", e.code) }

const runStatusID = "gloss-run"

// How long to keep reading after the command exits, in case a background
// process it started still holds its output open.
const runDrainGrace = 250 * time.Millisecond

func runRun(cmd *cobra.Command, args []string) error {
	envInfo := env.Detect()
	noColor := envInfo.NoColor || flagNoColor
	label := strings.Join(append([]string{filepath.Base(args[0])}, args[1:]...), " ")

	// stdout and stderr share one pipe so their lines stay in order.
	pr, pw, err := os.Pipe()
	if err != nil {
		return err
	}
	child := exec.Command(args[0], args[1:]...)
	child.Stdin = os.Stdin
	child.Stdout = pw
	child.Stderr = pw
	child.Env = childEnv(envInfo.IsTTY && !noColor)

	// Ctrl+C reaches the command directly (same terminal), so gloss stays up to
	// report how it ended. A second Ctrl+C kills it.
	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	if err := child.Start(); err != nil {
		pw.Close()
		pr.Close()
		fmt.Fprintf(os.Stderr, "gloss: can't run %s: %v\n", args[0], err)
		return exitCodeError{127}
	}
	pw.Close() // the child has its own copy; EOF arrives when it's done
	started := time.Now()

	exited := make(chan error, 1)
	go func() {
		err := child.Wait()
		pr.SetReadDeadline(time.Now().Add(runDrainGrace))
		exited <- err
	}()

	go func() {
		interrupts := 0
		for sig := range sigCh {
			interrupts++
			switch {
			case interrupts > 1:
				child.Process.Kill()
			case sig == syscall.SIGTERM:
				child.Process.Signal(syscall.SIGTERM) // not sent by the terminal, so pass it on
			}
		}
	}()

	output, _ := readInput(pr) // a read error just means the output ended early
	lines := make(chan string, 256)
	resultCode := make(chan int, 1)
	go func() {
		if envInfo.IsTTY {
			lines <- "::status id=" + runStatusID + " running " + label
		}
		for line := range output {
			lines <- line
		}
		code := exitCode(<-exited)
		lines <- runSummary(label, code, time.Since(started), envInfo.IsTTY)
		close(lines)
		resultCode <- code
	}()

	if envInfo.IsTTY {
		livePanel(lines, noColor, false)
	} else {
		formatLines(lines, plainWidth(), noColor)
	}
	pr.Close()

	if code := <-resultCode; code != 0 {
		return exitCodeError{code}
	}
	return nil
}

// runSummary is the final pass/fail line. In the live panel it updates the
// running status, which adds the duration itself.
func runSummary(label string, code int, took time.Duration, live bool) string {
	state, text := "done", label
	if code != 0 {
		state, text = "error", fmt.Sprintf("%s (exit %d)", label, code)
	}
	if live {
		return "::status id=" + runStatusID + " " + state + " " + text
	}
	if took >= time.Second {
		text += " in " + watch.FormatDuration(took)
	}
	return "::status " + state + " " + text
}

func childEnv(forceColor bool) []string {
	vars := append(os.Environ(), "GLOSS=1")
	if forceColor {
		for _, kv := range []string{"FORCE_COLOR=1", "CLICOLOR_FORCE=1"} {
			name, _, _ := strings.Cut(kv, "=")
			if _, set := os.LookupEnv(name); !set {
				vars = append(vars, kv)
			}
		}
	}
	return vars
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return 128 + int(status.Signal()) // shell convention, e.g. 130 for Ctrl+C
		}
		return exitErr.ExitCode()
	}
	return 1
}
