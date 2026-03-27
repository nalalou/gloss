package render

import (
	"fmt"
	"strings"
)

// RenderKV renders key-value pairs in aligned columns.
func RenderKV(pairs [][]string, separator string) string {
	if len(pairs) == 0 {
		return ""
	}

	// Find max key width for alignment
	maxKey := 0
	for _, pair := range pairs {
		if len(pair[0]) > maxKey {
			maxKey = len(pair[0])
		}
	}

	var lines []string
	for _, pair := range pairs {
		key := pair[0]
		val := ""
		if len(pair) > 1 {
			val = pair[1]
		}
		padding := strings.Repeat(" ", maxKey-len(key))
		lines = append(lines, fmt.Sprintf("  %s%s %s %s", key, padding, separator, val))
	}
	return strings.Join(lines, "\n")
}

// ParseKVPairs parses "key=value" strings into pairs.
func ParseKVPairs(args []string) [][]string {
	pairs := make([][]string, 0, len(args))
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			pairs = append(pairs, parts)
		} else {
			pairs = append(pairs, []string{arg, ""})
		}
	}
	return pairs
}
