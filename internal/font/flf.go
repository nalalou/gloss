package font

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// FLFFont implements Font by parsing the FIGlet .flf format.
type FLFFont struct {
	height    int
	baseline  int
	hardblank string
	chars     map[rune][]string
}

// ParseFLF reads a .flf font file and returns a Font ready to render text.
func ParseFLF(r io.Reader) (*FLFFont, error) {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, fmt.Errorf("empty font file")
	}
	parts := strings.Fields(scanner.Text())
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid header: need at least 6 fields")
	}
	prefix := parts[0]
	hardblank := string(prefix[len(prefix)-1])

	height, err := strconv.Atoi(parts[1])
	if err != nil || height <= 0 {
		return nil, fmt.Errorf("invalid height %q", parts[1])
	}
	baseline, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid baseline %q", parts[2])
	}
	commentLines, err := strconv.Atoi(parts[5])
	if err != nil {
		return nil, fmt.Errorf("invalid comment line count %q", parts[5])
	}

	for i := 0; i < commentLines; i++ {
		scanner.Scan()
	}

	f := &FLFFont{
		height:    height,
		baseline:  baseline,
		hardblank: hardblank,
		chars:     make(map[rune][]string),
	}

	// Parse printable ASCII 32–126, stop gracefully at EOF
	for ch := rune(32); ch <= 126; ch++ {
		lines := make([]string, 0, height)
		complete := true
		for i := 0; i < height; i++ {
			if !scanner.Scan() {
				complete = false
				break
			}
			line := scanner.Text()
			line = strings.TrimRight(line, "@")
			line = strings.ReplaceAll(line, hardblank, " ")
			lines = append(lines, line)
		}
		if !complete {
			break
		}
		f.chars[ch] = lines
	}

	return f, nil
}

func (f *FLFFont) Height() int { return f.height }

func (f *FLFFont) Render(text string) []string {
	rows := make([]string, f.height)

	for _, ch := range text {
		charLines, ok := f.chars[ch]
		if !ok {
			charLines = f.chars[' ']
		}
		for i := 0; i < f.height; i++ {
			if i < len(charLines) {
				rows[i] += charLines[i]
			}
		}
	}

	return rows
}
