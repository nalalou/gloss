package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// readStdinText reads up to maxSize bytes from stdin if piped.
// It trims trailing newlines only, preserving leading whitespace.
// Returns empty string if stdin is a terminal (not piped).
func readStdinText(maxSize int64) (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", nil
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	data, err := io.ReadAll(io.LimitReader(os.Stdin, maxSize+1))
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}
	if int64(len(data)) > maxSize {
		return "", fmt.Errorf("input too large (max %dKB)", maxSize/1024)
	}
	return strings.TrimRight(string(data), "\n"), nil
}
