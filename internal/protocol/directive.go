package protocol

import "strings"

// ParseDirective splits a :: line into (directive, id, args).
// Returns ("", "", "") for non-directive lines.
func ParseDirective(line string) (directive, id, args string) {
	if !strings.HasPrefix(line, "::") {
		return "", "", ""
	}
	rest := strings.TrimPrefix(line, "::")
	if rest == "" {
		return "", "", ""
	}

	tokens := strings.SplitN(rest, " ", 3)
	directive = strings.ToLower(tokens[0])

	if len(tokens) == 1 {
		return directive, "", ""
	}

	if strings.HasPrefix(tokens[1], "id=") {
		id = strings.TrimPrefix(tokens[1], "id=")
		if len(tokens) == 3 {
			args = tokens[2]
		}
		return directive, id, args
	}

	if len(tokens) == 2 {
		args = tokens[1]
	} else {
		args = tokens[1] + " " + tokens[2]
	}
	return directive, "", args
}
