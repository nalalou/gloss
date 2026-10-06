package protocol

import (
	"bufio"
	"io"
	"strings"
)

// ReadLines calls fn for each line of r, with no limit on line length.
// Trailing "\n" and "\r\n" are stripped.
func ReadLines(r io.Reader, fn func(line string)) error {
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if line != "" {
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
			fn(line)
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// Masker hides values registered with "::add-mask::value" in every later line,
// the same way GitHub Actions does.
type Masker struct {
	secrets []string
}

// Apply registers a new secret if line is an add-mask command, and returns
// line with all known secrets replaced by "***".
func (m *Masker) Apply(line string) string {
	if name, _, value, ok := ParseGitHub(line); ok && name == "add-mask" {
		if value = strings.TrimSpace(value); value != "" {
			m.secrets = append(m.secrets, value)
		}
		return line
	}
	for _, s := range m.secrets {
		line = strings.ReplaceAll(line, s, "***")
	}
	return line
}
